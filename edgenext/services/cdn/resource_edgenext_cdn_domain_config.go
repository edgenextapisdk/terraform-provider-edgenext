package cdn

import (
	"fmt"
	"log"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define single Map configuration items that need to be converted to List (MaxItems: 1)
var singleMapConfigs = map[string]bool{
	"origin":            true,
	"origin_host":       true,
	"referer":           true,
	"ip_black_list":     true,
	"ip_white_list":     true,
	"add_response_head": true,
	"https":             true,
	"compress_response": true,
	"rate_limit":        true,
	"cache_share":       true,
	"head_control":      true,
	"timeout":           true,
	"connect_timeout":   true,
	"deny_url":          true,
}

// Define Map array configuration items (maintain array format)
var arrayMapConfigs = map[string]bool{
	"cache_rule":           true,
	"cache_rule_list":      true,
	"add_back_source_head": true,
	"speed_limit":          true,
}

// Define all configuration items
var allConfigs = helper.MergeStringBoolMaps(singleMapConfigs, arrayMapConfigs)

func isEmpty(value interface{}) bool {
	return !isNonEmpty(value)
}

// isNonEmpty checks if the value is non-empty
func isNonEmpty(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	case string:
		return v != ""
	case int:
		return v != 0
	case bool:
		return true // bool values are always meaningful
	default:
		return true
	}
}

// getNonEmptyConfigItemsFromResource extracts non-empty configuration item names from Terraform resource configuration
func getNonEmptyConfigItemsFromResource(d *schema.ResourceData) []string {
	var configItems []string
	if config, ok := d.GetOk("config.0"); ok {
		configMap := config.(map[string]interface{})

		// Check if each configuration item is non-empty
		for key, value := range configMap {
			if isNonEmpty(value) {
				configItems = append(configItems, key)
			}
		}
	}

	return configItems
}

// getConfigItemsFromResource extracts configuration item names from Terraform resource configuration
func getConfigItemsFromResource(d *schema.ResourceData) []string {
	var configItems []string

	if config, ok := d.GetOk("config.0"); ok {
		configMap := config.(map[string]interface{})
		for key, _ := range configMap {
			configItems = append(configItems, key)
		}
	}

	return configItems
}

// buildConfigFromResource builds configuration map from Terraform resource
func buildConfigFromResource(d *schema.ResourceData) map[string]interface{} {
	config := make(map[string]interface{})

	if resourceConfig, ok := d.GetOk("config.0"); ok {
		configMap := resourceConfig.(map[string]interface{})

		log.Printf("[DEBUG] Original configuration: %+v", configMap)

		// Only include non-empty configuration items
		for key, value := range configMap {
			convertedValue := convertListTypeConfig(key, value)
			if convertedValue != nil {
				config[key] = convertedValue
			}
		}
	}

	log.Printf("[DEBUG] Final built configuration: %+v", config)
	return config
}

// convertListTypeConfig converts List type configuration items
func convertListTypeConfig(key string, value interface{}) interface{} {
	if isEmpty(value) {
		return nil
	}
	log.Printf("DEBUG convertListTypeConfig key:%s,value:%+v", key, value)

	// Determine if it is a List type
	list, ok := value.([]interface{})
	if !ok || len(list) == 0 {
		return value // Non-List type returns directly
	}

	if singleMapConfigs[key] {
		// Convert to single Map (take the first element)
		if configMap, ok := list[0].(map[string]interface{}); ok {
			// Recursively process nested List configurations
			processedMap := make(map[string]interface{})
			for k, v := range configMap {
				processedMap[k] = convertListTypeConfig(k, v)
			}
			return processedMap
		}
	} else if arrayMapConfigs[key] {
		// Convert to Map array
		var configArray []map[string]interface{}
		for _, item := range list {
			if itemMap, ok := item.(map[string]interface{}); ok {
				// Recursively process nested List configurations
				processedMap := make(map[string]interface{})
				for k, v := range itemMap {
					processedMap[k] = convertListTypeConfig(k, v)
				}
				configArray = append(configArray, processedMap)
			}
		}
		if len(configArray) > 0 {
			return configArray
		}
	} else {
		// For other possible nested Lists (like list fields), convert to string array
		if isStringList(list) {
			var stringArray []string
			for _, item := range list {
				if str, ok := item.(string); ok {
					stringArray = append(stringArray, str)
				}
			}
			if len(stringArray) > 0 {
				return stringArray
			}
		}
		// If type is uncertain, keep as is
		return list
	}

	return nil
}

// isStringList checks if it is a string list
func isStringList(list []interface{}) bool {
	if len(list) == 0 {
		return true
	}
	// Check if the first element is a string
	_, isString := list[0].(string)
	return isString
}

// convertAPIConfigToTerraform converts API returned configuration to Terraform schema expected format
func convertAPIConfigToTerraform(apiConfig map[string]interface{}) map[string]interface{} {
	terraformConfig := make(map[string]interface{})

	for key, value := range apiConfig {
		if value == nil {
			continue
		}
		// If the configuration item is not defined, skip
		if !allConfigs[key] {
			continue
		}
		if singleMapConfigs[key] {
			// Convert Map to List format (MaxItems: 1)
			if configMap, ok := value.(map[string]interface{}); ok && len(configMap) > 0 {
				// Recursively process nested configuration
				processedMap := convertAPINestedConfig(configMap)
				terraformConfig[key] = []interface{}{processedMap}
			}
		} else if arrayMapConfigs[key] {
			// Map array configuration items maintain array format
			if configArray, ok := value.([]interface{}); ok && len(configArray) > 0 {
				var processedArray []interface{}
				for _, item := range configArray {
					if itemMap, ok := item.(map[string]interface{}); ok {
						processedMap := convertAPINestedConfig(itemMap)
						processedArray = append(processedArray, processedMap)
					}
				}
				if len(processedArray) > 0 {
					terraformConfig[key] = processedArray
				}
			}
		} else {
			// Other configuration items set directly
			terraformConfig[key] = value
		}
	}

	return terraformConfig
}

// convertAPINestedConfig recursively processes nested configuration
func convertAPINestedConfig(apiConfig map[string]interface{}) map[string]interface{} {
	processedMap := make(map[string]interface{})

	for k, v := range apiConfig {
		if v == nil {
			continue
		}

		// Process nested List configurations (like referer.list, ip_white_list.list, etc.)
		if k == "list" {
			// Ensure list field is in string array format
			if strArray, ok := v.([]string); ok {
				processedMap[k] = strArray
			} else if interfaceArray, ok := v.([]interface{}); ok {
				var stringArray []string
				for _, item := range interfaceArray {
					if str, ok := item.(string); ok {
						stringArray = append(stringArray, str)
					}
				}
				processedMap[k] = stringArray
			} else {
				processedMap[k] = v
			}
		} else {
			processedMap[k] = v
		}
	}

	return processedMap
}

// compareAndUpdateConfig compares current configuration and desired configuration, executes necessary updates
func compareAndUpdateConfig(service *CdnService, d *schema.ResourceData, domain string, desiredConfig map[string]interface{}) error {
	var toDelete []string
	toSet := make(map[string]interface{})
	configItem := getConfigItemsFromResource(d)
	for _, item := range configItem {
		key := fmt.Sprintf("config.0.%s", item)
		// desireValue, ok := d.GetOk(key)
		// log.Printf("DEBUG key:%s,item:%s,value:%+v,HasChange:%+v,ok:%+v,desireValue:%+v", key, item, d.Get(key), d.HasChange(key), ok, desireValue)
		if d.HasChange(key) {
			if _, ok := d.GetOk(key); ok {
				toSet[item] = desiredConfig[item]
			} else {
				toDelete = append(toDelete, item)
			}
		}
	}

	// 3. First delete unnecessary configuration items
	if len(toDelete) > 0 {
		log.Printf("[INFO] Deleting configuration items: %v", toDelete)
		deleteReq := DeleteDomainConfigRequest{
			Domains: domain,
			Config:  toDelete,
		}
		err := service.DeleteDomainConfig(deleteReq)
		if err != nil {
			return fmt.Errorf("failed to delete configuration items: %w", err)
		}
	}

	// 4. Then set the required configuration items
	if len(toSet) > 0 {
		log.Printf("[INFO] Setting configuration items: %+v", toSet)
		_, err := service.SetDomainConfig(domain, toSet)
		if err != nil {
			return fmt.Errorf("failed to set configuration items: %w", err)
		}
	}

	return nil
}

// compareAndUpdateConfig compares current configuration and desired configuration, executes necessary updates
// func compareAndUpdateConfig(service *CdnService, domain string, currentConfig, desiredConfig map[string]interface{}) error {
// 	// 1. Find configuration items that need to be deleted (in current config but not in desired config)
// 	var toDelete []string
// 	for key := range currentConfig {
// 		if _, exists := desiredConfig[key]; !exists {
// 			toDelete = append(toDelete, key)
// 		}
// 	}

// 	// 2. Find configuration items that need to be set/updated
// 	toSet := make(map[string]interface{})
// 	for key, desiredValue := range desiredConfig {
// 		currentValue, exists := currentConfig[key]

// 		// If configuration item doesn't exist or value has changed, it needs to be set
// 		if !exists || !reflect.DeepEqual(currentValue, desiredValue) {
// 			toSet[key] = desiredValue
// 		}
// 	}

// 	// 3. First delete unnecessary configuration items
// 	if len(toDelete) > 0 {
// 		log.Printf("[INFO] Deleting configuration items: %v", toDelete)
// 		deleteReq := DeleteDomainConfigRequest{
// 			Domains: domain,
// 			Config:  toDelete,
// 		}
// 		err := service.DeleteDomainConfig(deleteReq)
// 		if err != nil {
// 			return fmt.Errorf("failed to delete configuration items: %w", err)
// 		}
// 	}

// 	// 4. Then set the required configuration items
// 	if len(toSet) > 0 {
// 		log.Printf("[INFO] Setting configuration items: %v", toSet)
// 		_, err := service.SetDomainConfig(domain, toSet)
// 		if err != nil {
// 			return fmt.Errorf("failed to set configuration items: %w", err)
// 		}
// 	}

// 	return nil
// }

func ResourceEdgenextCdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainConfigCreate,
		Read:   resourceDomainConfigRead,
		Update: resourceDomainConfigUpdate,
		Delete: resourceDomainConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Domain change requires resource recreation
				Description: "Accelerated domain name for setting functions",
			},
			"area": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Acceleration area: mainland_china(China mainland), outside_mainland_china(outside China mainland), global(global)",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain type: page(web), download(download), video_demand(video on demand), dynamic(dynamic)",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain status",
			},
			"icp_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ICP filing number",
			},
			"icp_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ICP filing status",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME",
			},
			"https": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "HTTPS",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time",
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Domain configuration items",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Origin configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_master": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Primary origin address, can fill multiple IPs or one domain name.\n" +
											"Multiple IPs separated by comma(,); primary and backup origins cannot have same IP or domain name.",
									},
									"default_slave": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Backup origin address, can fill multiple IPs or one domain name.\n" +
											"Multiple IPs separated by comma(,); primary and backup origins cannot have same IP or domain name.",
									},
									"origin_mode": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Origin mode: \n" +
											"default: Origin with user request protocol and port \n" +
											"http: Origin with http protocol on port 80 \n" +
											"https: Origin with https protocol on port 443 \n" +
											"custom: Origin with custom protocol(ori_https) and port(port) \n" +
											"Default value is default when not specified.",
									},
									"ori_https": {
										Type:     schema.TypeString,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// 从k中提取origin_mode的路径
											// 例如: "config.0.origin.0.ori_https" -> "config.0.origin.0.origin_mode"
											originModePath := strings.Replace(k, ".ori_https", ".origin_mode", 1)
											originModeValue := d.Get(originModePath)
											log.Printf("DEBUG ori_https k:%s, DiffSuppressFunc: originModePath=%s, originModeValue=%s", k, originModePath, originModeValue)
											var originMode string
											if originModeValue != nil {
												originMode = originModeValue.(string)
											}
											// 当origin_mode不是custom时，忽略ori_https的变化
											if originMode == "http" || originMode == "https" || originMode == "default" {
												return true
											}
											return old == new
										},
										Description: "This value needs to be set when origin_mode=custom. \n" +
											"HTTPS protocol origin: \n" +
											"yes: Yes \n" +
											"no: No",
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
											v := val.(int)
											if v < 1 || v > 65535 {
												errs = append(errs, fmt.Errorf("%q must be between 1 and 65535, got: %d", key, v))
											}
											return
										},
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// 从k中提取origin_mode的路径
											// 例如: "config.0.origin.0.port" -> "config.0.origin.0.origin_mode"
											originModePath := strings.Replace(k, ".port", ".origin_mode", 1)
											originModeValue := d.Get(originModePath)
											log.Printf("DEBUG port k:%s, DiffSuppressFunc: originModePath=%s, originModeValue=%s", k, originModePath, originModeValue)

											var originMode string
											if originModeValue != nil {
												originMode = originModeValue.(string)
											}
											// 当origin_mode不是custom时，忽略port的变化
											if originMode == "http" || originMode == "https" || originMode == "default" {
												return true
											}
											return old == new
										},
										Description: "This value needs to be set when origin_mode=custom. \n" +
											"Origin port, valid value range (1-65535).",
									},
								},
							},
						},
						"origin_host": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Origin HOST",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Origin HOST",
									},
								},
							},
						},
						"cache_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Cache rule list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache type: 1(file extension), 2(directory), 3(full path matching), 4(regex)",
									},
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Cache rules, multiple separated by commas. For example: \n" +
											"When type=1: jpg,png,gif \n" +
											"When type=2: /product/index,/test/index,/user/index \n" +
											"When type=3: /index.html,/test/*.jpg,/user/get?index \n" +
											"When type=4: see below for examples, set the corresponding regex. After setting regex, the request URL is matched against the regex, and if matched, this cache rule is used. \n" +
											"Default is cached with parameters, ignoring expiration time, not ignoring no-cache headers",
									},
									"time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache time, used with timeunit, maximum time not exceeding 2 years, when time=0 no caching for specified pattern",
									},
									"timeunit": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "s",
										Description: "Cache time unit, default value is s, optional values are (Y year, M month, D day, h hour, i minute, s second)",
									},
									"ignore_no_cache": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "off",
										Description: "Valid when cache time is greater than 0, ignore origin server no-cache headers, default value is off, optional parameters: on,off",
									},
									"ignore_expired": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "on",
										Description: "Valid when cache time is greater than 0, ignore origin server expiration time, default value is on, optional parameters: on,off",
									},
									"ignore_query": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "off",
										Description: "Valid when cache time is greater than 0, ignore parameters for caching and ignore parameters for origin requests, default value is off, optional parameters: on,off",
									},
								},
							},
						},
						"cache_rule_list": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "New cache rules",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_method": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Cache type, ext: file extension, dir: directory, route: full path matching, regex: regular expression",
									},
									"pattern": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Cache rules, multiple separated by commas. For example: when type=ext: jpg,png,gif; when type=dir: /product/index,/test/index,/user/index; when type=route: /index.html,/test/*.jpg,/user/get?index; when type=regex: set the corresponding regex as described below. After setting regex, the request URL is matched against the regex, and if matched, this cache rule is used. Default is cached with parameters, ignoring expiration time, not ignoring no-cache headers",
									},
									"case_ignore": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to ignore case, yes to ignore, no to not ignore. Default is yes",
									},
									"expire": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache time, used with expire_unit, maximum time not exceeding 2 years. When time=0, no caching for the specified pattern (i.e., caching disabled)",
									},
									"expire_unit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cache time unit, default value is s, optional values are (Y year, M month, D day, h hour, i minute, s second)",
									},
									"ignore_no_cache_headers": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to ignore no-cache information in origin server response headers<br/>such as Cache-Control:no-cache, default value is no, optional parameters: no,yes",
									},
									"follow_expired": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to follow origin server cache time, default value is no, optional parameters: no,yes",
									},
									"query_params_op": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "do_nothing",
										Description: "Query parameter operation mode, default value is no, optional parameters: no,yes,customer",
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Sort value, lower priority value means higher priority, duplicates are not allowed",
									},
									"query_params_op_way": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Only takes effect when query_params_op=customer, optional values: keep: retain, remove: remove",
									},
									"query_params_op_when": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Only takes effect when query_params_op=customer, optional values: cache: only process during caching, cache_back_source: process both caching and origin requests",
									},
									"params": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Only takes effect when query_params_op=customer, parameter list",
									},
								},
							},
						},
						"referer": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Referer blacklist and whitelist",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeInt,
										Required: true,
										Description: "Anti-hotlinking type: \n" +
											"1: referer blacklist \n" +
											"2: referer whitelist",
									},
									"list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Referer list, maximum 200 entries, multiple separated by commas; regex not supported; for wildcard domains, start with *, e.g.: *.example2.com, including any matching host headers and empty host headers.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"allow_empty": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Whether to allow empty referer, default value is true, optional parameters: true,false",
									},
								},
							},
						},
						"ip_black_list": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "IP blacklist",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "IP blacklist. IP format supports /8, /16, /24 subnet formats, IPs between subnets cannot overlap; maximum 500 IP formats can be set, multiple IP formats separated by commas; IP blacklist cannot coexist with IP whitelist, setting IP blacklist will clear IP whitelist functionality.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									// "mode": {
									// 	Type:      schema.TypeString,
									// 	Optional:  true,
									// 	WriteOnly: true,
									// 	Description: "IP list mode: \n" +
									// 		"append: Append mode \n" +
									// 		"cover: Cover mode, default cover",
									// },
								},
							},
						},
						"ip_white_list": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "IP whitelist",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "IP whitelist. IP format supports /8, /16, /24 subnet formats, IPs between subnets cannot overlap; maximum 500 IP formats can be set, multiple IP formats separated by commas; IP whitelist cannot coexist with IP blacklist, setting IP whitelist will clear IP blacklist functionality.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									// "mode": {
									// 	Type:      schema.TypeString,
									// 	Optional:  true,
									// 	WriteOnly: true,
									// 	Description: "IP list mode: \n" +
									// 		"append: Append mode \n" +
									// 		"cover: Cover mode, default cover",
									// },
								},
							},
						},
						"add_response_head": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Add response headers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Response header setting mode: \n" +
											"reset: Reset to the response headers set this time \n" +
											"add: Append the response headers set this time. If the appended response header key exists, overwrite the original response header. \n" +
											"remove: Delete response headers. \n" +
											"Default value is reset when not specified.",
									},
									"list": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Response header name",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Response header value",
												},
												// "cover": {
												// 	Type:        schema.TypeString,
												// 	Computed:    true,
												// 	Description: "on or off",
												// },
												// "only_hit": {
												// 	Type:     schema.TypeString,
												// 	Computed: true,
												// 	Description: "Only hit mode: \n" +
												// 		"on: On \n" +
												// 		"off: Off",
												// },
											},
										},
									},
								},
							},
						},
						"add_back_source_head": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Add origin request headers",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"head_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Origin request header name, maximum 20 can be added",
									},
									"head_value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Origin request header value",
									},
									"write_when_exists": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "yes",
										Description: "Whether to overwrite when the same request header exists \n" +
											"yes: Overwrite \n" +
											"no: Do not overwrite \n" +
											"Default value is yes when not specified.",
									},
								},
							},
						},
						"https": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "HTTPS configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_id": {
										Type:     schema.TypeInt,
										Required: true,
										Description: "Specify the bound certificate ID, which can be obtained through the certificate query interface. \n" +
											"When cert_id=0, HTTPS service will be disabled for the domain.",
									},
									"http2": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "HTTP2 feature: \n" +
											"on: Enable \n" +
											"off: Disable",
									},
									"force_https": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Redirect HTTP requests to HTTPS protocol \n" +
											"0: No redirect \n" +
											"302: HTTP request 302 redirect to HTTPS request \n" +
											"301: HTTP request 301 redirect to HTTPS request \n" +
											"Default value is 0 when not specified.",
									},
									"ocsp": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "OCSP[on,off] \n" +
											"No change when not specified",
									},
									"ssl_protocol": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Value range: [TLSv1,TLSv1.1,TLSv1.2,TLSv1.3]",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"head_control": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "HTTP header control",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"regex": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Required field, matching URL",
												},
												"head_op": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Required field, operation content, value range: [ADD,DEL,ALT]; ADD: add; DEL: delete; ALT: modify;",
												},
												"head_direction": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Required field, direction, value range: [CLI_REQ,CLI_REP,SER_REQ,SER_REP]; CLI_REQ: client request header; CLI_REP: client response header; SER_REQ: origin request header; SER_REP: origin response header;",
												},
												"head": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Required field, HTTP header name",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Required field, header value",
												},
												"order": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Required field, priority",
												},
											},
										},
									},
								},
							},
						},
						"timeout": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Origin read timeout",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, timeout duration, unit s, value range: [5-300]",
									},
								},
							},
						},
						"connect_timeout": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Origin connection timeout",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_connect_timeout": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, origin connection timeout duration, unit s, value range: [5-60]",
									},
								},
							},
						},
						"cache_share": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Shared cache",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"share_way": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Required field, sharing method, value range: [inner_share,cross_single_share,cross_all_share] \n" +
											"inner_share: HTTP and HTTPS share cache within this domain; \n" +
											"cross_single_share: HTTP and HTTPS separately share cache between different domains \n" +
											"cross_all_share: HTTP and HTTPS all share cache between different domains",
									},
									"domain": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, domain to be shared, this item only takes effect when share_way is cross_single_share and cross_all_share",
									},
								},
							},
						},
						"deny_url": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Block illegal URLs",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"urls": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Required field, blocked URL list",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"speed_limit": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Speed limit",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "URL matching type,\n" +
											"ext: file extension,\n" +
											"dir: directory,\n" +
											"route: full path matching,\n" +
											"regex: regular expression,\n" +
											"all: all",
									},
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
										Description: "URL matching rules, multiple separated by commas, e.g.: \n" +
											"When type=all: only supports .* \n" +
											"When type=ext: jpg,png,gif \n" +
											"When type=dir: /product/index,/test/index,/user/index \n" +
											"When type=route: /index.html,/test/*.jpg,/user/get?index \n" +
											"When type=regex: set the corresponding regex, after setting regex, match the request URL against the regex, if matched then use this speed limit rule.",
									},
									"speed": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Speed limit value, unit Kbps, actual effect will be converted to KB",
									},
									"begin_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Speed limit effective start time, format HH:ii, e.g. (08:30) 24-hour format",
									},
									"end_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Speed limit effective end time, format HH:ii, e.g. (08:30) 24-hour format",
									},
								},
							},
						},
						"rate_limit": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Rate limit",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_rate_count": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Required field, rate limit value",
									},
									"leading_flow_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Required field, how many bytes at the beginning are not rate limited",
									},
									"leading_flow_unit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Required field, unit for how many bytes at the beginning are not rate limited",
									},
									"max_rate_unit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Required field, rate limit unit",
									},
								},
							},
						},
						"compress_response": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Compress response",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_type": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Corresponding headers, e.g.: ['text/plain','application/x-javascript']",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_size": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Minimum size, used with min_size_unit, indicates the minimum file size to start compression",
									},
									"min_size_unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Minimum size unit, optional values are (KB,MB) note: this parameter is required when min_size=0",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceDomainConfigCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	domain := d.Get("domain").(string)
	config := buildConfigFromResource(d)

	// 1. Create domain
	req := DomainCreateRequest{
		Domain: domain,
		Area:   d.Get("area").(string),
		Type:   d.Get("type").(string),
		Config: config,
	}

	log.Printf("[INFO] Creating CDN domain: %s", req.Domain)
	_, err := service.CreateDomain(req)
	if err != nil {
		return fmt.Errorf("failed to create CDN domain: %w", err)
	}

	// 2. Set configuration items
	log.Printf("[INFO] Creating domain configuration: %s, config: %v", domain, config)

	// Set all configuration items directly when creating
	_, err = service.SetDomainConfig(domain, config)
	if err != nil {
		// Delete domain if configuration creation fails
		_ = service.DeleteDomain(domain)
		return fmt.Errorf("failed to create domain configuration: %w", err)
	}

	// Set resource ID (using domain as ID)
	d.SetId(domain)

	log.Printf("[INFO] Domain configuration created successfully: %s", d.Id())
	return resourceDomainConfigRead(d, m)
}

func resourceDomainConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	domain := d.Id()

	log.Printf("[INFO] Updating domain configuration: %s", domain)

	// 2. Get desired configuration
	desiredConfig := buildConfigFromResource(d)

	// 3. Intelligent comparison and update
	err := compareAndUpdateConfig(service, d, domain, desiredConfig)
	if err != nil {
		return fmt.Errorf("failed to update domain configuration: %w", err)
	}

	log.Printf("[INFO] Domain configuration updated successfully: %s", d.Id())
	return resourceDomainConfigRead(d, m)
}

func resourceDomainConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	// 1. Get domain information
	domain := d.Id()
	log.Printf("[INFO] Reading CDN domain: %s", domain)

	response, err := service.GetDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to read CDN domain: %w", err)
	}

	if len(response.Data) == 0 {
		log.Printf("[WARN] Domain does not exist: %s", domain)
		d.SetId("")
		return nil
	}

	domainData := response.Data[0]

	// Set all computed fields
	if err := d.Set("domain", domainData.Domain); err != nil {
		return fmt.Errorf("error setting domain: %w", err)
	}
	if err := d.Set("type", domainData.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("status", domainData.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("icp_num", domainData.IcpNum); err != nil {
		return fmt.Errorf("error setting icp_num: %w", err)
	}
	if err := d.Set("icp_status", domainData.IcpStatus); err != nil {
		return fmt.Errorf("error setting icp_status: %w", err)
	}
	if err := d.Set("area", domainData.Area); err != nil {
		return fmt.Errorf("error setting area: %w", err)
	}
	if err := d.Set("cname", domainData.Cname); err != nil {
		return fmt.Errorf("error setting cname: %w", err)
	}
	if err := d.Set("https", domainData.Https); err != nil {
		return fmt.Errorf("error setting https: %w", err)
	}
	if err := d.Set("create_time", domainData.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time: %w", err)
	}
	if err := d.Set("update_time", domainData.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time: %w", err)
	}

	// 2. Get domain configuration
	log.Printf("[INFO] Reading domain configuration: %s", domain)
	configItem := getNonEmptyConfigItemsFromResource(d)
	response2, err := service.GetDomainConfig(domain, configItem)
	if err != nil {
		return fmt.Errorf("failed to read domain configuration: %w", err)
	}

	if len(response2.Data) == 0 {
		log.Printf("[WARN] Domain configuration does not exist: %s", domain)
		d.SetId("")
		return nil
	}

	// Set resource ID
	d.SetId(domain)

	apiConfig := response2.Data[0].Config
	// Need to convert API returned configuration to the format expected by Terraform schema
	terraformConfig := convertAPIConfigToTerraform(apiConfig)
	configList := []map[string]interface{}{terraformConfig}
	if err := d.Set("config", configList); err != nil {
		return fmt.Errorf("error setting config: %w", err)
	}

	log.Printf("[INFO] Domain and configuration read successfully: %s", domain)
	return nil
}

func resourceDomainConfigDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	domain := d.Id()
	log.Printf("[INFO] Deleting CDN domain and configuration: %s", domain)

	err := service.DeleteDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to delete CDN domain and configuration: %w", err)
	}
	d.SetId("")
	log.Printf("[INFO] CDN domain and configuration deleted successfully: %s", domain)
	return nil
}
