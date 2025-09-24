package cdn

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"

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
		for key := range configMap {
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
			convertedValue := convertTerraformValueToAPI(value, key)
			if convertedValue != nil {
				config[key] = convertedValue
			}
		}
	}
	// set the type to reset
	if config["add_response_head"] != nil {
		config["add_response_head"].(map[string]interface{})["type"] = "reset"
	}

	log.Printf("[DEBUG] Final built configuration: %+v", config)
	return config
}

// convertTerraformValueToAPI recursively converts Terraform values to API format using reflection
func convertTerraformValueToAPI(value interface{}, key string) interface{} {
	return convertTerraformValueToAPIWithContext(value, key, "")
}

// convertTerraformValueToAPIWithContext recursively converts Terraform values with context awareness
func convertTerraformValueToAPIWithContext(value interface{}, key string, parentKey string) interface{} {
	if isEmpty(value) {
		return nil
	}

	// Use reflection to determine the actual type
	v := reflect.ValueOf(value)
	kind := v.Kind()

	// log.Printf("[DEBUG] convertTerraformValueToAPIWithContext key:%s, parentKey:%s, value:%+v, kind:%v", key, parentKey, value, kind)

	switch kind {
	case reflect.Slice, reflect.Array:
		// Handle []interface{} - common in Terraform
		if sliceVal, ok := value.([]interface{}); ok {
			if len(sliceVal) == 0 {
				return nil
			}

			// Check the type of first element to determine conversion strategy
			firstElem := sliceVal[0]
			firstKind := reflect.ValueOf(firstElem).Kind()

			if firstKind == reflect.Map {
				// Use Schema metadata to determine if this is a single map config (MaxItems: 1)
				if isSingleMapConfigFromSchema(key, parentKey) {
					// Single map config (MaxItems: 1) - take first element and unwrap from array
					if mapVal, ok := firstElem.(map[string]interface{}); ok {
						return convertTerraformMapToAPIWithContext(mapVal, key)
					}
				} else {
					// Array of maps - convert all elements and maintain array structure
					var result []map[string]interface{}
					for _, item := range sliceVal {
						if mapVal, ok := item.(map[string]interface{}); ok {
							convertedMap := convertTerraformMapToAPIWithContext(mapVal, key)
							if convertedMap != nil {
								result = append(result, convertedMap)
							}
						}
					}
					if len(result) > 0 {
						return result
					}
				}
			} else {
				// Array of primitives - check if all are strings
				var stringArray []string
				allStrings := true

				for _, item := range sliceVal {
					if str, ok := item.(string); ok && str != "" {
						stringArray = append(stringArray, str)
					} else if item != nil && item != "" {
						allStrings = false
						break
					}
				}

				if allStrings && len(stringArray) > 0 {
					return stringArray
				} else if !allStrings {
					// Mixed types or non-strings, return as-is
					return sliceVal
				}
			}
		}

	case reflect.Map:
		// Handle map[string]interface{}
		if mapVal, ok := value.(map[string]interface{}); ok {
			return convertTerraformMapToAPIWithContext(mapVal, key)
		}

	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		// Primitive types - return as-is if not empty
		if str, ok := value.(string); ok && str == "" {
			return nil
		}
		return value

	default:
		// Unknown type - return as-is with warning
		log.Printf("[WARN] convertTerraformValueToAPI: unknown type %v for key %s, returning as-is", kind, key)
		return value
	}

	return nil
}

// convertTerraformMapToAPIWithContext recursively converts a Terraform map to API format with context
func convertTerraformMapToAPIWithContext(mapVal map[string]interface{}, parentKey string) map[string]interface{} {
	if len(mapVal) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range mapVal {
		convertedValue := convertTerraformValueToAPIWithContext(v, k, parentKey)
		if convertedValue != nil {
			result[k] = convertedValue
		}
	}

	if len(result) > 0 {
		return result
	}
	return nil
}

// SchemaFieldInfo contains metadata about schema fields
type SchemaFieldInfo struct {
	MaxItems int
	Type     string
}

// getSchemaMetadata returns metadata extracted from the Terraform Schema
func getSchemaMetadata() map[string]SchemaFieldInfo {
	// This metadata is extracted from the actual Schema definition in ResourceEdgenextCdnDomainConfig
	// Fields with MaxItems: 1 should be wrapped as single-element arrays
	return map[string]SchemaFieldInfo{
		// Top-level single map configs (MaxItems: 1)
		"origin":            {MaxItems: 1, Type: "List"},
		"origin_host":       {MaxItems: 1, Type: "List"},
		"referer":           {MaxItems: 1, Type: "List"},
		"ip_black_list":     {MaxItems: 1, Type: "List"},
		"ip_white_list":     {MaxItems: 1, Type: "List"},
		"add_response_head": {MaxItems: 1, Type: "List"},
		"https":             {MaxItems: 1, Type: "List"},
		"compress_response": {MaxItems: 1, Type: "List"},
		"rate_limit":        {MaxItems: 1, Type: "List"},
		"cache_share":       {MaxItems: 1, Type: "List"},
		"head_control":      {MaxItems: 1, Type: "List"},
		"timeout":           {MaxItems: 1, Type: "List"},
		"connect_timeout":   {MaxItems: 1, Type: "List"},
		"deny_url":          {MaxItems: 1, Type: "List"},

		// Array configs (no MaxItems limit)
		"cache_rule":           {MaxItems: 0, Type: "List"},
		"cache_rule_list":      {MaxItems: 0, Type: "List"},
		"add_back_source_head": {MaxItems: 0, Type: "List"},
		"speed_limit":          {MaxItems: 0, Type: "List"},

		// Nested fields in array contexts should not be wrapped (example)
		// "cache_rule.origin":   {MaxItems: 0, Type: "String"},
		// "cache_rule.timeout":  {MaxItems: 0, Type: "String"},
		// "speed_limit.origin":  {MaxItems: 0, Type: "String"},
		// "speed_limit.timeout": {MaxItems: 0, Type: "String"},
	}
}

// isSingleMapConfigFromSchema checks if a field should be treated as single map config (MaxItems: 1)
// This is used for Terraform -> API conversion to determine if array should be unwrapped
func isSingleMapConfigFromSchema(key string, parentKey string) bool {
	schemaInfo := getSchemaMetadata()

	// Build the field path for lookup
	fieldPath := key
	if parentKey != "" {
		fieldPath = parentKey + "." + key
	}

	// Check if this field is defined with MaxItems: 1
	if info, exists := schemaInfo[fieldPath]; exists {
		// log.Printf("[DEBUG] isSingleMapConfigFromSchema: field %s has MaxItems=%d", fieldPath, info.MaxItems)
		return info.MaxItems == 1
	}

	// Check direct key lookup (for top-level fields)
	if info, exists := schemaInfo[key]; exists {
		// log.Printf("[DEBUG] isSingleMapConfigFromSchema: field %s has MaxItems=%d", key, info.MaxItems)
		return info.MaxItems == 1
	}

	return false
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Acceleration area.\n" +
					"	- mainland_china\n" +
					"	- outside_mainland_china\n" +
					"	- global",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Domain type.\n" +
					"	- page\n" +
					"	- download\n" +
					"	- video_demand\n" +
					"	- dynamic",
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
							Required:    true,
							MaxItems:    1,
							Description: "Origin configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_master": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Primary origin address, can fill multiple IPs or domains.\n" +
											"Multiple IPs or domains separated by comma(,); primary and backup origins cannot have same IP or domain.",
									},
									"default_slave": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Backup origin address, can fill multiple IPs or domains.\n" +
											"Multiple IPs or domains separated by comma(,); primary and backup origins cannot have same IP or domain.",
									},
									"origin_mode": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Origin mode: \n" +
											"	- default: Origin with user request protocol and port\n" +
											"	- http: Origin with http protocol on port 80\n" +
											"	- https: Origin with https protocol on port 443\n" +
											"	- custom: Origin with custom protocol(ori_https) and port(port)",
									},
									"ori_https": {
										Type:     schema.TypeString,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// get origin_mode path from k
											// example: "config.0.origin.0.ori_https" -> "config.0.origin.0.origin_mode"
											originModePath := strings.Replace(k, ".ori_https", ".origin_mode", 1)
											originModeValue := d.Get(originModePath)
											var originMode string
											if originModeValue != nil {
												originMode = originModeValue.(string)
											}
											// when origin_mode is not custom, ignore ori_https change
											if originMode == "http" || originMode == "https" || originMode == "default" {
												return true
											}
											return old == new
										},
										Description: "Whether to enable HTTPS protocol origin, this value needs to be set when origin_mode=custom. \n" +
											"	- yes \n" +
											"	- no",
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
											// get origin_mode path from k
											// example: "config.0.origin.0.port" -> "config.0.origin.0.origin_mode"
											originModePath := strings.Replace(k, ".port", ".origin_mode", 1)
											originModeValue := d.Get(originModePath)

											var originMode string
											if originModeValue != nil {
												originMode = originModeValue.(string)
											}
											// when origin_mode is not custom, ignore port change
											if originMode == "http" || originMode == "https" || originMode == "default" {
												return true
											}
											return old == new
										},
										Description: "This value needs to be set when origin_mode=custom. \n" +
											"Origin port, valid value range [1-65535].",
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
										Description: "sharing method. \n" +
											"	- inner_share: HTTP and HTTPS share cache within this domain; \n" +
											"	- cross_single_share: HTTP and HTTPS separately share cache between different domains \n" +
											"	- cross_all_share: HTTP and HTTPS all share cache between different domains",
									},
									"domain": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "domain to be shared, this item only takes effect when share_way is cross_single_share and cross_all_share",
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
										Type:     schema.TypeInt,
										Required: true,
										Description: "Cache type. \n" +
											"	- 1: file extension \n" +
											"	- 2: directory \n" +
											"	- 3: full path matching \n" +
											"	- 4: regex",
									},
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Cache rules, multiple separated by commas. Default is cached with parameters, ignoring expiration time, not ignoring no-cache headers. For example: \n" +
											"	- When type=1: jpg,png,gif \n" +
											"	- When type=2: /product/index,/test/index,/user/index \n" +
											"	- When type=3: /index.html,/test/*.jpg,/user/get?index \n" +
											"	- When type=4: set the corresponding regex. After setting regex, the request URL is matched against the regex, and if matched, this cache rule is used.",
									},
									"time": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache time, used with timeunit, maximum time not exceeding 2 years, when time=0 no caching for specified pattern",
									},
									"timeunit": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "s",
										Description: "Cache time unit, default value is s. \n" +
											"	- Y: year \n" +
											"	- M: month \n" +
											"	- D: day \n" +
											"	- h: hour \n" +
											"	- i: minute \n" +
											"	- s: second",
									},
									"ignore_no_cache": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "off",
										Description: "Valid when cache time is greater than 0, ignore origin server no-cache headers, default value is off. \n" +
											"	- on \n" +
											"	- off",
									},
									"ignore_expired": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "on",
										Description: "Valid when cache time is greater than 0, ignore origin server expiration time, default value is on. \n" +
											"	- on \n" +
											"	- off",
									},
									"ignore_query": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "off",
										Description: "Valid when cache time is greater than 0, ignore parameters for caching and ignore parameters for origin requests, default value is off. \n" +
											"	- on \n" +
											"	- off",
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
										Type:     schema.TypeString,
										Required: true,
										Description: "Cache type. \n" +
											"	- ext: file extension \n" +
											"	- dir: directory \n" +
											"	- route: full path matching \n" +
											"	- regex: regular expression",
									},
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Cache rules, multiple separated by commas. Default is cached with parameters, ignoring expiration time, not ignoring no-cache headers. \n" +
											"	- when type=ext: jpg,png,gif \n" +
											"	- when type=dir: /product/index,/test/index,/user/index \n" +
											"	- when type=route: /index.html,/test/*.jpg,/user/get?index \n" +
											"	- when type=regex: set the corresponding regex as described below. After setting regex, the request URL is matched against the regex, and if matched, this cache rule is used.",
									},
									"case_ignore": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Whether to ignore case, Default is yes. \n" +
											"	- yes: Ignore \n" +
											"	- no: Do not ignore",
									},
									"expire": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache time, used with expire_unit, maximum time not exceeding 2 years. When time=0, no caching for the specified pattern (i.e., caching disabled)",
									},
									"expire_unit": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Cache time unit, default value is s. \n" +
											"	- Y: year \n" +
											"	- M: month \n" +
											"	- D: day \n" +
											"	- h: hour \n" +
											"	- i: minute \n" +
											"	- s: second",
									},
									"ignore_no_cache_headers": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Whether to ignore no-cache information in origin server response headers, such as Cache-Control:no-cache, default value is no. \n" +
											"	- yes\n" +
											"	- no",
									},
									"follow_expired": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Whether to follow origin server cache time, default value is no. \n" +
											"	- yes\n" +
											"	- no",
									},
									"priority": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Sort value, lower priority value means higher priority, duplicates are not allowed",
									},
									"cache_or_not": {
										Type:     schema.TypeString,
										Optional: true,
										// 不传默认按照expire来判断是否开启缓存
										Description: "Whether to cache, not set means use expire to determine whether to cache. \n" +
											"	- yes\n" +
											"	- no",
									},
									"query_params_op": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "do_nothing",
										Description: "Query parameter operation mode, default value is no. \n" +
											"	- no\n" +
											"	- yes\n" +
											"	- customer",
									},
									"query_params_op_way": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Only takes effect when query_params_op=customer. \n" +
											"	- keep \n" +
											"	- remove",
									},
									"query_params_op_when": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Only takes effect when query_params_op=customer. \n" +
											"	- cache: only process during caching \n" +
											"	- cache_back_source: process both caching and origin requests",
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
											"	- 1: referer blacklist \n" +
											"	- 2: referer whitelist",
									},
									"list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Referer list, maximum 200 entries, multiple separated by commas; regex not supported; for wildcard domains, start with *, e.g., *.example2.com, including any matching host headers and empty host headers.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"allow_empty": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
										Description: "Whether to allow empty referer, default value is true. \n" +
											"	- true\n" +
											"	- false",
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
									"list": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Response header list",
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
												"cover": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "on or off",
												},
												"only_hit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "on or off",
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
										Default:  "yes",
										Description: "Whether to overwrite when the same request header exists, default value is yes. \n" +
											"	- yes: Overwrite \n" +
											"	- no: Do not overwrite",
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
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											if old == "" && new == "off" {
												return true
											}
											return old == new
										},
										Description: "HTTP2 feature: \n" +
											"	- on: Enable \n" +
											"	- off: Disable",
									},
									"force_https": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "Redirect HTTP requests to HTTPS protocol, Default value is 0. \n" +
											"0: No redirect \n" +
											"302: HTTP request 302 redirect to HTTPS request \n" +
											"301: HTTP request 301 redirect to HTTPS request",
									},
									"ocsp": {
										Type:     schema.TypeString,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											if old == "" && new == "off" {
												return true
											}
											return old == new
										},
										Description: "No change when not specified. \n" +
											"	- on: Enable \n" +
											"	- off: Disable",
									},
									"ssl_protocol": {
										Type:     schema.TypeList,
										Optional: true,
										Description: "ssl protocol. \n" +
											"	- TLSv1 \n" +
											"	- TLSv1.1 \n" +
											"	- TLSv1.2 \n" +
											"	- TLSv1.3",
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
										Type:        schema.TypeList,
										Required:    true,
										Description: "HTTP header control rules list",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"regex": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "matching URL",
												},
												"head_op": {
													Type:     schema.TypeString,
													Required: true,
													Description: "operation content. \n" +
														"	- ADD: add \n" +
														"	- DEL: delete \n" +
														"	- ALT: modify",
												},
												"head_direction": {
													Type:     schema.TypeString,
													Required: true,
													Description: "direction. \n" +
														"	- CLI_REQ: client request header \n" +
														"	- CLI_REP: client response header \n" +
														"	- SER_REQ: origin request header \n" +
														"	- SER_REP: origin response header",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "http header value",
												},
												"head": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "http header name",
												},
												"order": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "priority",
												},
												"fun_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Computed: true,
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
										Description: "timeout duration, unit is s, value range: [5-300]",
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
										Description: "origin connection timeout duration, unit is s, value range: [5-60]",
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
										Description: "blocked URL list",
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
										Description: "URL matching type. \n" +
											"	- ext: file extension\n" +
											"	- dir: directory\n" +
											"	- route: full path matching\n" +
											"	- regex: regular expression\n" +
											"	- all: all",
									},
									"pattern": {
										Type:     schema.TypeString,
										Required: true,
										Description: "URL matching rules, multiple separated by commas. \n" +
											"	- When type=all: only supports .* \n" +
											"	- When type=ext: jpg,png,gif \n" +
											"	- When type=dir: /product/index,/test/index,/user/index \n" +
											"	- When type=route: /index.html,/test/*.jpg,/user/get?index \n" +
											"	- When type=regex: set the corresponding regex, after setting regex, match the request URL against the regex, if matched then use this speed limit rule.",
									},
									"speed": {
										Type:     schema.TypeInt,
										Required: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											newInt, _ := strconv.Atoi(new)
											newNewInt := int(math.Max(1, math.Ceil(float64(newInt)/8.0))) * 8
											if oldInt, _ := strconv.Atoi(old); oldInt == newNewInt {
												return true
											}
											return old == new
										},
										Description: "Speed limit value, unit is Kbps, actual effect will be converted to KB",
									},
									"begin_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Speed limit effective start time, format is HH:ii, e.g., (08:30) 24-hour format",
									},
									"end_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Speed limit effective end time, format is HH:ii, e.g., (08:30) 24-hour format",
									},
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Speed limit priority",
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
										Description: "rate limit value",
									},
									"leading_flow_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "how many bytes at the beginning are not rate limited",
									},
									"leading_flow_unit": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "MB",
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// ignore case
											old = strings.ToLower(old)
											new = strings.ToLower(new)
											return old == new
										},
										Description: "unit for how many bytes at the beginning are not rate limited, default value is MB. \n" +
											"	- KB\n" +
											"	- MB",
									},
									"max_rate_unit": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "MB",
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// ignore case
											old = strings.ToLower(old)
											new = strings.ToLower(new)
											return old == new
										},
										Description: "rate limit unit, default value is MB. \n" +
											"	- KB\n" +
											"	- MB",
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
										Description: "Corresponding headers, e.g., ['text/plain','application/x-javascript']",
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
										Type:     schema.TypeString,
										Required: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											// ignore case
											old = strings.ToLower(old)
											new = strings.ToLower(new)
											return old == new
										},
										Description: "Minimum size unit. \n" +
											"	- KB\n" +
											"	- MB",
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
	domainConfig := make(map[string]interface{})
	for key, value := range config {
		if key == "origin" {
			domainConfig[key] = value
			break
		}
	}

	// 1. Create domain
	req := DomainCreateRequest{
		Domain: domain,
		Area:   d.Get("area").(string),
		Type:   d.Get("type").(string),
		Config: domainConfig,
	}

	log.Printf("[INFO] Creating CDN domain: %+v", req)
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
	// return nil
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
	// return nil
	return resourceDomainConfigRead(d, m)
}

func readDomain(d *schema.ResourceData, service *CdnService, domain string) error {
	response, err := service.GetDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to read CDN domain: %w", err)
	}

	if len(response.Data) == 0 {
		log.Printf("[WARN] Domain does not exist: %s", domain)
		d.SetId("")
		return fmt.Errorf("domain does not exist: %s", domain)
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
	return nil
}

func readDomainConfig(d *schema.ResourceData, service *CdnService, domain string, configItem []string) error {
	response, err := service.GetDomainConfig(domain, configItem)
	if err != nil {
		return fmt.Errorf("failed to read domain configuration: %w", err)
	}

	if len(response.Data) == 0 {
		log.Printf("[WARN] Domain configuration does not exist: %s", domain)
		d.SetId("")
		return fmt.Errorf("domain configuration does not exist: %s", domain)
	}

	// Need to convert API returned configuration to the format expected by Terraform schema
	apiConfig := response.Data[0].Config

	config := make(map[string]interface{})
	if apiConfig.Origin != nil {
		origin := make(map[string]interface{})
		origin["default_master"] = apiConfig.Origin.DefaultMaster
		origin["default_slave"] = apiConfig.Origin.DefaultSlave
		origin["origin_mode"] = apiConfig.Origin.OriginMode
		origin["ori_https"] = apiConfig.Origin.OriHttps
		origin["port"] = apiConfig.Origin.Port.Int()
		config["origin"] = []map[string]interface{}{origin}
	}

	if apiConfig.OriginHost != nil {
		originHost := make(map[string]interface{})
		originHost["host"] = apiConfig.OriginHost.Host
		config["origin_host"] = []map[string]interface{}{originHost}
	}

	if len(apiConfig.CacheRule) > 0 {
		cacheRule := make([]map[string]interface{}, len(apiConfig.CacheRule))
		for i, rule := range apiConfig.CacheRule {
			if rule != nil {
				cacheRule[i] = make(map[string]interface{})
				cacheRule[i]["type"] = rule.Type
				cacheRule[i]["pattern"] = rule.Pattern
				cacheRule[i]["time"] = rule.Time
				cacheRule[i]["timeunit"] = rule.TimeUnit
				cacheRule[i]["ignore_no_cache"] = rule.IgnoreNoCache
				cacheRule[i]["ignore_expired"] = rule.IgnoreExpired
				cacheRule[i]["ignore_query"] = rule.IgnoreQuery
			}
		}
		config["cache_rule"] = cacheRule
	}

	if len(apiConfig.CacheRuleList) > 0 {
		var existingCacheRuleList []interface{}
		if existingConfig, ok := d.GetOk("config.0.cache_rule_list"); ok {
			existingCacheRuleList = existingConfig.([]interface{})
		}
		cacheRuleList := make([]map[string]interface{}, len(apiConfig.CacheRuleList))
		for i, rule := range apiConfig.CacheRuleList {
			if rule != nil {
				cacheRuleList[i] = make(map[string]interface{})
				cacheRuleList[i]["match_method"] = rule.MatchMethod
				cacheRuleList[i]["pattern"] = rule.Pattern
				cacheRuleList[i]["expire"] = rule.Expire
				cacheRuleList[i]["expire_unit"] = rule.ExpireUnit
				cacheRuleList[i]["ignore_no_cache_headers"] = rule.IgnoreNoCacheHeaders
				cacheRuleList[i]["follow_expired"] = rule.FollowExpired
				cacheRuleList[i]["query_params_op"] = rule.QueryParamsOp
				cacheRuleList[i]["priority"] = rule.Priority
				cacheRuleList[i]["query_params_op_way"] = rule.QueryParamsOpWay
				cacheRuleList[i]["query_params_op_when"] = rule.QueryParamsOpWhen
				cacheRuleList[i]["params"] = rule.Params

				if len(existingCacheRuleList) > i && existingCacheRuleList[i] != nil {
					if existingItem, ok := existingCacheRuleList[i].(map[string]interface{}); ok {
						cacheRuleList[i]["case_ignore"] = existingItem["case_ignore"]
						cacheRuleList[i]["cache_or_not"] = existingItem["cache_or_not"]
					}
				}
			}
		}
		config["cache_rule_list"] = cacheRuleList
	}

	if apiConfig.Referer != nil {
		referer := make(map[string]interface{})
		referer["type"] = apiConfig.Referer.Type
		referer["list"] = apiConfig.Referer.List
		referer["allow_empty"] = apiConfig.Referer.AllowEmpty
		config["referer"] = []map[string]interface{}{referer}
	}

	if apiConfig.IPBlackList != nil {
		ipBlackList := make(map[string]interface{})
		ipBlackList["list"] = apiConfig.IPBlackList.List
		config["ip_black_list"] = []map[string]interface{}{ipBlackList}
	}

	if apiConfig.IPWhiteList != nil {
		ipWhiteList := make(map[string]interface{})
		ipWhiteList["list"] = apiConfig.IPWhiteList.List
		config["ip_white_list"] = []map[string]interface{}{ipWhiteList}
	}

	if apiConfig.AddResponseHead != nil {
		addResponseHead := make(map[string]interface{})
		if len(apiConfig.AddResponseHead.List) > 0 {
			list := make([]map[string]interface{}, len(apiConfig.AddResponseHead.List))
			for i, head := range apiConfig.AddResponseHead.List {
				if head != nil {
					list[i] = make(map[string]interface{})
					list[i]["name"] = head.Name
					list[i]["value"] = head.Value
					list[i]["cover"] = head.Cover
					list[i]["only_hit"] = head.OnlyHit
				}
			}
			addResponseHead["list"] = list
		}
		config["add_response_head"] = []map[string]interface{}{addResponseHead}
	}

	if len(apiConfig.AddBackSourceHead) > 0 {
		addBackSourceHead := make([]map[string]interface{}, len(apiConfig.AddBackSourceHead))
		for i, head := range apiConfig.AddBackSourceHead {
			if head != nil {
				addBackSourceHead[i] = make(map[string]interface{})
				addBackSourceHead[i]["head_name"] = head.Name
				addBackSourceHead[i]["head_value"] = head.Value
				addBackSourceHead[i]["write_when_exists"] = head.WriteWhenExists
			}
		}
		config["add_back_source_head"] = addBackSourceHead
	}

	if apiConfig.HTTPS != nil {
		var existingHttps []interface{}
		if existingConfig, ok := d.GetOk("config.0.https"); ok {
			existingHttps = existingConfig.([]interface{})
		}
		https := make(map[string]interface{})
		https["http2"] = apiConfig.HTTPS.HTTP2
		https["force_https"] = apiConfig.HTTPS.ForceHTTPS
		https["cert_id"] = apiConfig.HTTPS.CertID

		if len(existingHttps) > 0 && existingHttps[0] != nil {
			if existingItem, ok := existingHttps[0].(map[string]interface{}); ok {
				https["ocsp"] = existingItem["ocsp"]
				https["ssl_protocol"] = existingItem["ssl_protocol"]
			}
		}
		config["https"] = []map[string]interface{}{https}
	}

	if apiConfig.CompressResponse != nil {
		compressResponse := make(map[string]interface{})
		compressResponse["content_type"] = apiConfig.CompressResponse.ContentType
		compressResponse["min_size"] = apiConfig.CompressResponse.MinSize
		compressResponse["min_size_unit"] = apiConfig.CompressResponse.MinSizeUnit
		config["compress_response"] = []map[string]interface{}{compressResponse}
	}

	if apiConfig.RateLimit != nil {
		rateLimit := make(map[string]interface{})
		rateLimit["max_rate_count"] = apiConfig.RateLimit.MaxRateCount
		rateLimit["leading_flow_count"] = apiConfig.RateLimit.LeadingFlowCount
		rateLimit["leading_flow_unit"] = apiConfig.RateLimit.LeadingFlowUnit
		rateLimit["max_rate_unit"] = apiConfig.RateLimit.MaxRateUnit
		config["rate_limit"] = []map[string]interface{}{rateLimit}
	}

	if apiConfig.CacheShare != nil {
		cacheShare := make(map[string]interface{})
		cacheShare["share_way"] = apiConfig.CacheShare.ShareWay
		if apiConfig.CacheShare.ShareWay == "inner_share" {
			cacheShare["domain"] = domain
		} else {
			cacheShare["domain"] = apiConfig.CacheShare.Domain
		}
		config["cache_share"] = []map[string]interface{}{cacheShare}
	}

	if apiConfig.HeadControl != nil {
		headControl := make(map[string]interface{})
		if len(apiConfig.HeadControl.List) > 0 {
			list := make([]map[string]interface{}, len(apiConfig.HeadControl.List))

			// Get existing head_control configuration to preserve order and head fields
			var existingHeadControl []interface{}
			if existingConfig, ok := d.GetOk("config.0.head_control.0.list"); ok {
				existingHeadControl = existingConfig.([]interface{})
			}

			for i, head := range apiConfig.HeadControl.List {
				if head != nil {
					list[i] = make(map[string]interface{})
					list[i]["regex"] = head.Regex
					list[i]["head_op"] = head.HeadOp
					list[i]["head_direction"] = head.HeadDirection
					list[i]["value"] = head.Value
					list[i]["fun_name"] = head.FunName
					list[i]["key"] = head.Key

					if len(existingHeadControl) > i && existingHeadControl[i] != nil {
						if existingItem, ok := existingHeadControl[i].(map[string]interface{}); ok {
							list[i]["order"] = existingItem["order"]
							list[i]["head"] = existingItem["head"]
						}
					}
				}
			}
			headControl["list"] = list
		}
		config["head_control"] = []map[string]interface{}{headControl}
	}

	if apiConfig.Timeout != nil {
		timeout := make(map[string]interface{})
		timeout["time"] = apiConfig.Timeout.Time
		config["timeout"] = []map[string]interface{}{timeout}
	}

	if apiConfig.ConnectTimeout != nil {
		connectTimeout := make(map[string]interface{})
		connectTimeout["origin_connect_timeout"] = apiConfig.ConnectTimeout.OriginConnectTimeout
		config["connect_timeout"] = []map[string]interface{}{connectTimeout}
	}

	if apiConfig.DenyURL != nil {
		denyURL := make(map[string]interface{})
		denyURL["urls"] = apiConfig.DenyURL.URLs
		config["deny_url"] = []map[string]interface{}{denyURL}
	}

	if len(apiConfig.SpeedLimit) > 0 {
		speedLimit := make([]map[string]interface{}, len(apiConfig.SpeedLimit))
		for i, limit := range apiConfig.SpeedLimit {
			if limit != nil {
				speedLimit[i] = make(map[string]interface{})
				speedLimit[i]["type"] = limit.Type
				speedLimit[i]["pattern"] = limit.Pattern
				speedLimit[i]["speed"] = limit.Speed
				speedLimit[i]["begin_time"] = limit.BeginTime
				speedLimit[i]["end_time"] = limit.EndTime
				speedLimit[i]["priority"] = limit.Priority
			}
		}
		config["speed_limit"] = speedLimit
	}

	configList := []map[string]interface{}{config}
	if err := d.Set("config", configList); err != nil {
		return fmt.Errorf("error setting config: %w", err)
	}
	return nil
}

func resourceDomainConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	// 1. Get domain information
	domain := d.Id()
	log.Printf("[INFO] Resource reading CDN domain: %s", domain)
	err := readDomain(d, service, domain)
	if err != nil {
		return fmt.Errorf("resource failed to read CDN domain: %w", err)
	}

	// 2. Get domain configuration
	log.Printf("[INFO] Resource reading domain configuration: %s", domain)
	err = readDomainConfig(d, service, domain, nil)
	if err != nil {
		return fmt.Errorf("resource failed to read domain configuration: %w", err)
	}
	// Set resource ID
	d.SetId(domain)
	log.Printf("[INFO] Resource read domain and configuration successfully: %s", domain)
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
