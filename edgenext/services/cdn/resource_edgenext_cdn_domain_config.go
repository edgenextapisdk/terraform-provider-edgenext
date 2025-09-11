package cdn

import (
	"fmt"
	"log"
	"reflect"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

// getConfigItemsFromResource extracts configuration item names from Terraform resource configuration
func getConfigItemsFromResource(d *schema.ResourceData) []string {
	var configItems []string

	if config, ok := d.GetOk("config.0"); ok {
		configMap := config.(map[string]interface{})

		// Check if each configuration item is defined in the resource
		for key, value := range configMap {
			if isNonEmpty(value) {
				configItems = append(configItems, key)
			}
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
			// if isNonEmpty(value) {
			convertedValue := convertListTypeConfig(key, value)
			if convertedValue != nil {
				config[key] = convertedValue
			}
			// }
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

	// Determine if it is a List type
	list, ok := value.([]interface{})
	if !ok || len(list) == 0 {
		return value // Non-List type returns directly
	}

	// Define configuration items that need to be converted to a single Map (MaxItems: 1)
	singleMapConfigs := map[string]bool{
		"origin":                 true,
		"origin_host":            true,
		"referer":                true,
		"ip_black_list":          true,
		"ip_white_list":          true,
		"add_response_head":      true,
		"https":                  true,
		"visit_timestamp":        true,
		"forbid_http_x":          true,
		"cache_error_code":       true,
		"video_drag":             true,
		"compress_response":      true,
		"extend":                 true,
		"rate_limit":             true,
		"cache_share":            true,
		"head_control":           true,
		"timeout":                true,
		"connect_timeout":        true,
		"qiniu_origin_auth":      true,
		"forward_status":         true,
		"error_page_rewrite":     true,
		"post_upload_size_limit": true,
		"deny_url":               true,
		"tos_origin":             true,
	}

	// Define configuration items that need to be converted to Map array
	arrayMapConfigs := map[string]bool{
		"cache_rule":           true,
		"cache_rule_list":      true,
		"add_back_source_head": true,
		"speed_limit":          true,
		"visit_deny_whitelist": true,
		"new_origin":           true,
		"source_url_rewrite":   true,
		"combined_ban":         true,
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

// convertMapToOriginItem converts map to OriginItem struct
func convertMapToOriginItem(originMap interface{}) OriginItem {
	if originMap == nil {
		return OriginItem{}
	}

	m, ok := originMap.(map[string]interface{})
	if !ok {
		return OriginItem{}
	}

	origin := OriginItem{}
	if v, ok := m["default_master"].(string); ok {
		origin.DefaultMaster = v
	}
	if v, ok := m["default_slave"].(string); ok {
		origin.DefaultSlave = v
	}
	if v, ok := m["origin_mode"].(string); ok {
		origin.OriginMode = v
	}
	if v, ok := m["ori_https"].(string); ok {
		origin.OriHttps = v
	}
	if v, ok := m["port"].(string); ok {
		origin.Port = v
	}

	return origin
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

	// Define single Map configuration items that need to be converted to List (MaxItems: 1)
	singleMapConfigs := map[string]bool{
		"origin":                 true,
		"origin_host":            true,
		"referer":                true,
		"ip_black_list":          true,
		"ip_white_list":          true,
		"add_response_head":      true,
		"https":                  true,
		"visit_timestamp":        true,
		"forbid_http_x":          true,
		"cache_error_code":       true,
		"video_drag":             true,
		"compress_response":      true,
		"extend":                 true,
		"rate_limit":             true,
		"cache_share":            true,
		"head_control":           true,
		"timeout":                true,
		"connect_timeout":        true,
		"qiniu_origin_auth":      true,
		"forward_status":         true,
		"error_page_rewrite":     true,
		"post_upload_size_limit": true,
		"deny_url":               true,
		"tos_origin":             true,
	}

	// Define Map array configuration items (maintain array format)
	arrayMapConfigs := map[string]bool{
		"cache_rule":           true,
		"cache_rule_list":      true,
		"add_back_source_head": true,
		"speed_limit":          true,
		"visit_deny_whitelist": true,
		"new_origin":           true,
		"source_url_rewrite":   true,
		"combined_ban":         true,
	}

	for key, value := range apiConfig {
		if value == nil {
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
func compareAndUpdateConfig(service *CdnService, domain string, currentConfig, desiredConfig map[string]interface{}) error {
	// 1. Find configuration items that need to be deleted (in current config but not in desired config)
	var toDelete []string
	for key := range currentConfig {
		if _, exists := desiredConfig[key]; !exists {
			toDelete = append(toDelete, key)
		}
	}

	// 2. Find configuration items that need to be set/updated
	toSet := make(map[string]interface{})
	for key, desiredValue := range desiredConfig {
		currentValue, exists := currentConfig[key]

		// If configuration item doesn't exist or value has changed, it needs to be set
		if !exists || !reflect.DeepEqual(currentValue, desiredValue) {
			toSet[key] = desiredValue
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
		log.Printf("[INFO] Setting configuration items: %v", toSet)
		_, err := service.SetDomainConfig(domain, toSet)
		if err != nil {
			return fmt.Errorf("failed to set configuration items: %w", err)
		}
	}

	return nil
}

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
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain ID",
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
										Description: "This value needs to be set when origin_mode=custom. \n" +
											"HTTPS protocol origin: \n" +
											"yes: Yes \n" +
											"no: No",
									},
									"port": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "This value needs to be set when origin_mode=custom. \n" +
											"Origin port, valid value range (0-65535).",
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
										Description: "Query parameter operation mode, default value is no, optional parameters: no,yes,customer",
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Sort value, lower priority value means higher priority, duplicates are not allowed",
									},
									"cache_or_not": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to cache, optional values: ['yes','no'], if not passed, defaults to determining cache based on expire",
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
									"mode": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "IP list mode: \n" +
											"append: Append mode \n" +
											"cover: Cover mode, default cover",
									},
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
									"mode": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "IP list mode: \n" +
											"append: Append mode \n" +
											"cover: Cover mode, default cover",
									},
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
													Optional:    true,
													Description: "Response header value",
												},
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
						"visit_timestamp": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Timestamp anti-hotlinking",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pattern": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Regular expression for matching URLs, URL format is $scheme://$domain/$uri?$args, regex should be set according to the URL to be matched, e.g. .*, matches all accessed URL addresses",
									},
									"time_format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Time format, optional values are timestamp_hex(hexadecimal timestamp), date_minute(date, format YYYYmmddHHii, e.g. 201805211010), timestamp(decimal timestamp)",
									},
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Secret key for generating signature, multiple keys separated by spaces",
									},
									"deadtime": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "URL lifetime, e.g. 3600, parameter must be greater than 0, unit is seconds(s), URLs are considered invalid after this time",
									},
									"req_uri_type": {
										Type:     schema.TypeInt,
										Required: true,
										Description: "URL matching pattern, currently 4 URL patterns: \n" +
											"1: $scheme://$domain/$uri?$args&{keyname}=$key&{timename}=$time&$args \n" +
											"2: $scheme://$domain/$uri?$args&{timename}=$time&{keyname}=$key&$args \n" +
											"3: $scheme://$domain/$time/$key/$uri?$args \n" +
											"4: $scheme://$domain/$key/$time/$uri?$args",
									},
									"origin_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Whether origin request carries parameters, optional values are 1,2, default is 1, 1 means origin request without encryption string, 2: origin request with encryption string",
									},
									"style": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Parameter encryption sorting combination, $ourkey$uri$time, changing the order can rearrange the ($ourkey,$uri,$time) three variables, additional fields and information are not currently supported",
									},
									"timename": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Default value is time, parameter name used for passing parameters in URL, only needed when passing parameters through ?",
									},
									"keyname": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Default value is key, generate signature using md5 with (path,key,time), passed through parameter name in URL",
									},
								},
							},
						},
						"forbid_http_x": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Forbid HTTP or HTTPS access",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Protocol to forbid access, optional values are http,https",
									},
								},
							},
						},
						"cache_error_code": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Cache error codes",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Status code list, multiple values can be set separated by commas, optional values are 400,403,404,414,500,501,502,503,504,506,5xx, where 5xx includes 500,501,502,503,504,506.",
									},
									"bcache": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to enable caching, optional values are on,off",
									},
									"cache_time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache time, used with cache_unit, maximum time not exceeding 2 years, note: this parameter is required when bcache=on",
									},
									"cache_unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Cache time unit, optional values are (year:year,month:month,day:day,hour:hour,minute:minute,second:second) note: this parameter is required when bcache=on",
									},
								},
							},
						},
						"video_drag": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Video seeking",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Matching URL rules, e.g.: (mp4 flv f4v m4v)",
									},
									"mp4": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to enable MP4 seeking, optional values are on,off",
									},
									"flv": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to enable FLV seeking, optional values are on,off",
									},
									"start": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Start parameter, e.g.: start, parameter name allowed characters: 'letters, numbers, underscore, and cannot start with a number'",
									},
									"end": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "End parameter, e.g.: end, parameter name allowed characters: 'letters, numbers, underscore, and cannot start with a number'",
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
						"visit_deny_whitelist": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Anti-hotlinking whitelist",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, set URL path: /user/get?index",
									},
									"turn_on": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Optional values: true, false, true=enable current URL whitelist, false=remove current URL anti-hotlinking whitelist setting",
									},
								},
							},
						},
						"range_back_source": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Range back-to-source. Default enables range back-to-source, configure off to disable range back-to-source. Note: configuring on will not display in domain details, configuring off can be displayed in domain details",
						},
						"extend": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Squid extended configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"squid": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, Squid extended configuration",
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
										Required:    true,
										Description: "Required field, how many bytes at the beginning are not rate limited",
									},
									"leading_flow_unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, unit for how many bytes at the beginning are not rate limited",
									},
									"max_rate_unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, rate limit unit",
									},
								},
							},
						},
						"new_origin": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Domain origin (new)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, origin server address",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Optional, origin port, determined by protocol parameter when not specified, protocol=http uses port=80; protocol=https uses port=443; protocol=default ignores port value and fixes to 0, parameter range (0, 65535)",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, origin protocol; optional values [default,http,https]; default value default; default: follow request port and protocol for origin (request xx port uses origin xx port, https request uses https protocol for origin); http: origin fixed to use http protocol; https: origin fixed to use https protocol",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, origin host; default value \"\", must be in domain format; when origin host is not set, the origin host header will be consistent with the accelerated domain",
									},
									"level": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Optional, primary-backup level, higher values have higher priority; optional values [10, 20]; default value 20; Note: when multiple origin servers have the same primary-backup level, they are all primary sources; if multiple origin servers have different primary-backup levels, the one with the highest value is the primary source, others are backup sources, and among backup sources, higher primary-backup level values have higher priority; For example, if an accelerated domain has 4 origin IPs configured, A and B both have level 20, C has level 15, D has level 10, then when A and B sources are normal, priority goes to AB sources; when both A and B are abnormal, priority goes to C source; when A, B, C are all abnormal, goes to D source",
									},
									"weight_level": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Optional, weight level; optional values [1, 10000]; default value 1; Note: higher weight level values mean greater weight, when primary-backup levels are the same, origin requests are distributed based on weight ratio; For example: if an accelerated domain has 2 origin servers configured, both A and B sources have primary-backup level 20, A source has weight level 60, B source has weight level 40, then A source's origin request count is approximately 60% of the domain's total origin requests, B source's origin request count is approximately 40% of the domain's total origin requests",
									},
									"isp": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, ISP, specify ISP for origin; optional values [default.dx.lt.yd]; default value default; Note: when this configuration is set, the default source default value must be set; For example, if an accelerated domain sets Telecom source A, then default source B must be set, Telecom requests go to A origin server, other requests go to B origin server;",
									},
									"connect_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Optional, TCP connection time, TCP connection time refers to the TCP connection timeout when going back to origin, default 15 seconds, can be set to positive integers between 3~60 seconds",
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
						"qiniu_origin_auth": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Qiniu origin authentication",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, only supports URLs starting with http protocol, e.g. http://test.com/",
									},
									"match_method": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, default [default], [default] all, [ext] only matches corresponding suffix, [regex] only matches regex content",
									},
									"pattern": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, default empty matches all, when match_method=ext, suffix separated by commas is required, when match_method=regex, regex matching is required",
									},
								},
							},
						},
						"forward_status": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Continue to fetch content after redirect",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codes": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Required field, parameter type for continued redirection; 301: 301 redirect; 302: 302 redirect, 301 and 302 are int type",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"error_page_rewrite": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Custom error page",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_status_code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Required field, error status code, value range 400<=x<=599",
									},
									"redirect_status_code": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Required field, redirect status code, value range 301,302",
									},
									"redirect_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, redirect URL, must start with http:// or https://, and conform to URL format",
									},
								},
							},
						},
						"post_upload_size_limit": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "POST upload size limit",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_value": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Required field, limit size value, value range 1<=x<=1024, unit M",
									},
								},
							},
						},
						"source_url_rewrite": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Origin URL rewrite",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, original Path regex before rewriting. Note: starts with /, matches original client access path + request parameters, and requires using regular expressions (e.g.: /test/a.jpg?a=1), regex special characters can be escaped. You can also specify groups in the regex, groups are wrapped by (). These groups can be referenced using $n in the target origin path. For example, /aaa/bbb/(.*) represents all directories and files under path /aaa/bbb/. This example contains one group.",
									},
									"target_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, target Path regex after rewriting. Note: starts with /, contains origin path + request parameters, (e.g.: /newtest/b.jpg?a=1). Example: if the origin path to be rewritten is set to /test/(.*)/(.*).jpg, and the target origin path is set to /newtest/$1/$2.apk, then when user accesses with origin path /test/a/b.jpg, $1 will capture the content of the first regex parentheses, which is a; $2 will capture the content of the second regex parentheses, which is b, so the actual origin path will be rewritten to /newtest/a/b.apk.",
									},
								},
							},
						},
						"combined_ban": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Combined-Ban combined blocking",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, default: ban, note currently only supports ban method",
									},
									"configs": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: "Required field, rule group configuration info key",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"method": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to match, if not matched then treated as non-range, default: yes, options: yes, no",
															},
															"case_insensitive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to ignore case, default: yes, options: yes, no",
															},
															"method_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Required field, blocking type key, currently only supports method, ip, referer, url, ua",
															},
															"list": {
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
												"ip": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to match, if not matched then treated as non-range, default: yes, options: yes, no",
															},
															"case_insensitive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to ignore case, default: yes, options: yes, no",
															},
															"method_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Required field, blocking type key, currently only supports method, ip, referer, url, ua",
															},
															"list": {
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
												"referer": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to match, if not matched then treated as non-range, default: yes, options: yes, no",
															},
															"case_insensitive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to ignore case, default: yes, options: yes, no",
															},
															"method_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Required field, blocking type key, currently only supports method, ip, referer, url, ua",
															},
															"list": {
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
												"url": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to match, if not matched then treated as non-range, default: yes, options: yes, no",
															},
															"case_insensitive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to ignore case, default: yes, options: yes, no",
															},
															"method_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Required field, blocking type key, currently only supports method, ip, referer, url, ua",
															},
															"list": {
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
												"ua": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to match, if not matched then treated as non-range, default: yes, options: yes, no",
															},
															"case_insensitive": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Optional, whether to ignore case, default: yes, options: yes, no",
															},
															"method_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Required field, blocking type key, currently only supports method, ip, referer, url, ua",
															},
															"list": {
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
											},
										},
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
						"tos_origin": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "TOS domain origin",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"isp": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, ISP, must be default ISP default, optional range: [default, dx, lt, yd]",
									},
									"ips": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Required field, origin server, can configure IP or domain name, note: same configuration [] can only have one origin type",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"group_sort": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Required field, primary-backup order, value range 1~10",
									},
									"weight": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Optional, weight, value range 1~10000",
									},
									"origin_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, origin mode default is default, optional range: [default, http, https, custom], note: selecting custom allows customizing origin protocol and port",
									},
									"protocol": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, origin protocol",
									},
									"port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Optional, origin port",
									},
									"host_mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, origin host mode default is default, optional range: [default, custom], note: selecting custom allows customizing origin host",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, origin host",
									},
									"auth_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, authentication type default is default (no authentication), optional range: [default, oss, tos], note: selecting oss allows configuring [auth_bucket_name], selecting tos allows configuring [auth_expire, auth_cdn_tag], common authentication parameters auth_secret_key, auth_access_key",
									},
									"auth_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, secret_key credential",
									},
									"auth_access_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, access_key credential",
									},
									"auth_bucket_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, bucket_name storage bucket name",
									},
									"auth_expire": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Optional, expire expiration time, e.g.: 300 expires in five minutes",
									},
									"auth_cdn_tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, cdn_tag identifier",
									},
									"parse_priority": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Optional, origin priority resolution setting default is default, optional range: [default, v4, v4v6, v6v4, v6], note: v4v6 means v4 first then v6",
									},
								},
							},
						},
						"client_real_ip": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Client real IP",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"head": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, client real IP header",
									},
								},
							},
						},
						"user_agent": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "User-Agent blacklist/whitelist",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ua_list": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, User-Agent list",
									},
									"url_pattern": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, URL rule",
									},
									"url_case_insensitive": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, URL case insensitive, value range [yes, no]",
									},
									"allow_empty_ua": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, allow empty User-Agent, value range [yes, no]",
									},
									"ua_case_insensitive": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, User-Agent case insensitive, value range [yes, no]",
									},
									"url_match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, URL match type, value range [all(all), ext(suffix), dir(directory), route(path), regex(regular expression)]",
									},
									"ua_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, User-Agent blacklist/whitelist type, value range [black(UA blacklist mode), white(UA whitelist mode)]",
									},
									"ua_match_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, User-Agent match mode, value range [pattern(fuzzy match), exact(exact match)]",
									},
								},
							},
						},
						"visit_areas_limit": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Access area blacklist/whitelist",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, limit type, value range [white(whitelist), black(blacklist)]",
									},
									"country_list": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Required field, country list, multiple countries separated by English comma (,)",
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

	// 1. Create domain
	var originConfig OriginItem
	if originData := d.Get("config.0.origin.0"); originData != nil {
		originConfig = convertMapToOriginItem(originData)
	}

	req := DomainCreateRequest{
		Domain: d.Get("domain").(string),
		Area:   d.Get("area").(string),
		Type:   d.Get("type").(string),
		Config: DomainConfig{
			Origin: originConfig,
		},
	}

	log.Printf("[INFO] Creating CDN domain: %s", req.Domain)
	_, err := service.CreateDomain(req)
	if err != nil {
		return fmt.Errorf("failed to create CDN domain: %w", err)
	}

	// 2. Set configuration items
	domain := d.Get("domain").(string)
	config := buildConfigFromResource(d)

	log.Printf("[INFO] Creating domain configuration: %s, config: %v", domain, config)

	if len(config) == 0 {
		return fmt.Errorf("at least one configuration item must be specified")
	}

	// Set all configuration items directly when creating
	_, err = service.SetDomainConfig(domain, config)
	if err != nil {
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

	// 1. Get current configuration state
	// managedConfigItems := getConfigItemsFromResource(d)
	// currentResponse, err := service.GetDomainConfig(domain, managedConfigItems)
	currentResponse, err := service.GetDomainConfig(domain, nil)
	if err != nil {
		return fmt.Errorf("failed to get current configuration: %w", err)
	}

	var currentConfig map[string]interface{}
	if len(currentResponse.Data) > 0 {
		currentConfig = currentResponse.Data[0].Config
	} else {
		currentConfig = make(map[string]interface{})
	}
	currentConfig = convertAPIConfigToTerraform(currentConfig)
	log.Printf("[DEBUG] Current configuration: %+v", currentConfig)

	// 2. Get desired configuration
	desiredConfig := buildConfigFromResource(d)

	// 3. Intelligent comparison and update
	err = compareAndUpdateConfig(service, domain, currentConfig, desiredConfig)
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
	d.Set("id", domainData.ID)
	d.Set("domain", domainData.Domain)
	d.Set("type", domainData.Type)
	d.Set("status", domainData.Status)
	d.Set("icp_num", domainData.IcpNum)
	d.Set("icp_status", domainData.IcpStatus)
	d.Set("area", domainData.Area)
	d.Set("cname", domainData.Cname)
	d.Set("https", domainData.Https)
	d.Set("create_time", domainData.CreateTime)
	d.Set("update_time", domainData.UpdateTime)

	// 2. Get domain configuration
	log.Printf("[INFO] Reading domain configuration: %s", domain)

	// Get the list of configuration items managed by the resource
	// managedConfigItems := getConfigItemsFromResource(d)

	// Query current configuration (only query the configuration items managed by the resource)
	// response2, err := service.GetDomainConfig(domain, managedConfigItems)
	response2, err := service.GetDomainConfig(domain, nil)
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
	// Build configuration list, only including configuration items managed by the resource
	// Need to convert API returned configuration to the format expected by Terraform schema
	apiConfig := response2.Data[0].Config
	terraformConfig := convertAPIConfigToTerraform(apiConfig)
	configList := []map[string]interface{}{terraformConfig}
	d.Set("config", configList)

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
