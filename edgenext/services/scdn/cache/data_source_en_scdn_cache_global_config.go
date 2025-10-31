package cache

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCacheGlobalConfig returns the SCDN cache global config data source
func DataSourceEdgenextScdnCacheGlobalConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCacheGlobalConfigRead,

		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule ID (as string to match Terraform's resource ID format)",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule name",
			},
			"conf": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cache configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nocache": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Cache eligibility (true: bypass cache, false: cache)",
						},
						"cache_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Edge TTL cache configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cachetime": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cache time",
									},
									"ignore_cache_time": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Ignore source cache time",
									},
									"ignore_nocache_header": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Ignore no-cache header",
									},
									"no_cache_control_op": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "No cache control operation",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cache action: 'default', 'nocache', 'cachetime', or 'force'",
									},
								},
							},
						},
						"browser_cache_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Browser cache configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cachetime": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cache time",
									},
									"ignore_cache_time": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Ignore source cache time (cache-control)",
									},
									"nocache": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to cache",
									},
								},
							},
						},
						"cache_errstatus": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Status code cache configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cachetime": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Status code cache time",
									},
									"err_status": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Status code array",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"cache_url_rewrite": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Custom cache key configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sort_args": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Parameter sorting",
									},
									"ignore_case": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Ignore case",
									},
									"queries": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Query string processing",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"args_method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'",
												},
												"items": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Parameter keys",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"cookies": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cookie processing",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"args_method": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'",
												},
												"items": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Cookie keys",
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
						"cache_share": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cache sharing configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scheme": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP/HTTPS cache sharing method: 'http' or 'https'",
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

func dataSourceScdnCacheGlobalConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	log.Printf("[INFO] Reading SCDN cache global config")
	response, err := service.GetCacheGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to read cache global config: %w", err)
	}

	// Validate response data exists
	if response == nil {
		return fmt.Errorf("cache global config response is nil")
	}

	// Log the raw ID value for debugging
	log.Printf("[DEBUG] Cache global config ID from API: %v (type: %T)", response.Data.ID, response.Data.ID)

	// Get the ID value - ensure it's a valid integer
	// Explicitly convert to int to ensure type safety for Terraform SDK
	ruleID := int(response.Data.ID)
	if ruleID < 0 {
		return fmt.Errorf("invalid cache global config ID: %d (must be non-negative)", ruleID)
	}

	// Set resource ID (use string format for resource ID)
	resourceID := fmt.Sprintf("global-cache-config-%d", ruleID)
	d.SetId(resourceID)
	log.Printf("[DEBUG] Set resource ID to: %s", resourceID)

	// Set computed fields - ensure ID is set as string type
	// Terraform schema expects TypeString to match resource ID format
	idValue := strconv.Itoa(ruleID)
	if err := d.Set("id", idValue); err != nil {
		log.Printf("[ERROR] Failed to set id field. Value: %s, Type: %T, Error: %v", idValue, idValue, err)
		return fmt.Errorf("failed to set id field (value: %s, type: %T): %w", idValue, idValue, err)
	}
	log.Printf("[DEBUG] Successfully set id field to: %s (type: %T)", idValue, idValue)
	if err := d.Set("name", response.Data.Name); err != nil {
		log.Printf("[WARN] Failed to set name: %v", err)
	}

	if response.Data.Conf != nil {
		confMap := map[string]interface{}{
			"nocache": response.Data.Conf.NoCache,
		}

		if response.Data.Conf.CacheRule != nil {
			cacheRuleMap := map[string]interface{}{
				"cachetime":             response.Data.Conf.CacheRule.CacheTime,
				"ignore_cache_time":     response.Data.Conf.CacheRule.IgnoreCacheTime,
				"ignore_nocache_header": response.Data.Conf.CacheRule.IgnoreNoCacheHeader,
				"no_cache_control_op":   response.Data.Conf.CacheRule.NoCacheControlOp,
				"action":                response.Data.Conf.CacheRule.Action,
			}
			confMap["cache_rule"] = []interface{}{cacheRuleMap}
		}

		if response.Data.Conf.BrowserCacheRule != nil {
			browserCacheMap := map[string]interface{}{
				"cachetime":         response.Data.Conf.BrowserCacheRule.CacheTime,
				"ignore_cache_time": response.Data.Conf.BrowserCacheRule.IgnoreCacheTime,
				"nocache":           response.Data.Conf.BrowserCacheRule.NoCache,
			}
			confMap["browser_cache_rule"] = []interface{}{browserCacheMap}
		}

		if len(response.Data.Conf.CacheErrStatus) > 0 {
			errStatusList := make([]map[string]interface{}, 0, len(response.Data.Conf.CacheErrStatus))
			for _, errStatus := range response.Data.Conf.CacheErrStatus {
				errStatusMap := map[string]interface{}{
					"cachetime":  errStatus.CacheTime,
					"err_status": errStatus.ErrStatus,
				}
				errStatusList = append(errStatusList, errStatusMap)
			}
			confMap["cache_errstatus"] = errStatusList
		}

		if response.Data.Conf.CacheURLRewrite != nil {
			urlRewriteMap := map[string]interface{}{
				"sort_args":   response.Data.Conf.CacheURLRewrite.SortArgs,
				"ignore_case": response.Data.Conf.CacheURLRewrite.IgnoreCase,
			}

			if response.Data.Conf.CacheURLRewrite.Queries != nil {
				queriesMap := map[string]interface{}{
					"args_method": response.Data.Conf.CacheURLRewrite.Queries.ArgsMethod,
					"items":       response.Data.Conf.CacheURLRewrite.Queries.Items,
				}
				urlRewriteMap["queries"] = []interface{}{queriesMap}
			}

			if response.Data.Conf.CacheURLRewrite.Cookies != nil {
				cookiesMap := map[string]interface{}{
					"args_method": response.Data.Conf.CacheURLRewrite.Cookies.ArgsMethod,
					"items":       response.Data.Conf.CacheURLRewrite.Cookies.Items,
				}
				urlRewriteMap["cookies"] = []interface{}{cookiesMap}
			}

			confMap["cache_url_rewrite"] = []interface{}{urlRewriteMap}
		}

		if response.Data.Conf.CacheShare != nil {
			cacheShareMap := map[string]interface{}{
				"scheme": response.Data.Conf.CacheShare.Scheme,
			}
			confMap["cache_share"] = []interface{}{cacheShareMap}
		}

		if err := d.Set("conf", []interface{}{confMap}); err != nil {
			log.Printf("[WARN] Failed to set conf: %v", err)
		}
	}

	// Handle result_output_file if provided
	if _, ok := d.GetOk("result_output_file"); ok {
		outputData := map[string]interface{}{
			"id":   idValue, // Use string format to match schema
			"name": response.Data.Name,
			"conf": response.Data.Conf,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
	}

	log.Printf("[INFO] Cache global config read successfully")
	return nil
}
