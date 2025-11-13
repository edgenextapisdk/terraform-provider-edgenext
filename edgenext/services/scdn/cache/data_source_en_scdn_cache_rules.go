package cache

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCacheRules returns the SCDN cache rules data source
func DataSourceEdgenextScdnCacheRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCacheRulesRead,

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Business ID (template ID for 'tpl' type, domain ID for 'domain' type)",
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Business type: 'tpl' (template) or 'domain'",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page number",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Page size",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of rules",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of cache rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule name",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status (1: enabled, 2: disabled)",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Weight",
						},
						"expr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Wirefilter rule",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type: 'domain', 'tpl', or 'global'",
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
				},
			},
		},
	}
}

func dataSourceScdnCacheRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	req := scdn.CacheRuleGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
	}

	if page, ok := d.GetOk("page"); ok {
		req.Page = page.(int)
	}
	if pageSize, ok := d.GetOk("page_size"); ok {
		req.PageSize = pageSize.(int)
	}

	log.Printf("[INFO] Reading SCDN cache rules: business_id=%d, business_type=%s", businessID, businessType)
	response, err := service.GetCacheRules(req)
	if err != nil {
		return fmt.Errorf("failed to read cache rules: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s", businessID, businessType))

	// Set computed fields
	if err := d.Set("total", response.Data.Total); err != nil {
		log.Printf("[WARN] Failed to set total: %v", err)
	}

	// Convert list to schema format
	ruleList := make([]map[string]interface{}, 0, len(response.Data.List))
	for _, rule := range response.Data.List {
		ruleMap := map[string]interface{}{
			"id":     rule.ID,
			"name":   rule.Name,
			"remark": rule.Remark,
			"status": rule.Status,
			"weight": rule.Weight,
			"expr":   rule.Expr,
			"type":   rule.Type,
		}

		if rule.Conf != nil {
			confMap := map[string]interface{}{
				"nocache": rule.Conf.NoCache,
			}

			if rule.Conf.CacheRule != nil {
				cacheRuleMap := map[string]interface{}{
					"cachetime":             rule.Conf.CacheRule.CacheTime,
					"ignore_cache_time":     rule.Conf.CacheRule.IgnoreCacheTime,
					"ignore_nocache_header": rule.Conf.CacheRule.IgnoreNoCacheHeader,
					"no_cache_control_op":   rule.Conf.CacheRule.NoCacheControlOp,
					"action":                rule.Conf.CacheRule.Action,
				}
				confMap["cache_rule"] = []interface{}{cacheRuleMap}
			}

			if rule.Conf.BrowserCacheRule != nil {
				browserCacheMap := map[string]interface{}{
					"cachetime":         rule.Conf.BrowserCacheRule.CacheTime,
					"ignore_cache_time": rule.Conf.BrowserCacheRule.IgnoreCacheTime,
					"nocache":           rule.Conf.BrowserCacheRule.NoCache,
				}
				confMap["browser_cache_rule"] = []interface{}{browserCacheMap}
			}

			if len(rule.Conf.CacheErrStatus) > 0 {
				errStatusList := make([]map[string]interface{}, 0, len(rule.Conf.CacheErrStatus))
				for _, errStatus := range rule.Conf.CacheErrStatus {
					errStatusMap := map[string]interface{}{
						"cachetime":  errStatus.CacheTime,
						"err_status": errStatus.ErrStatus,
					}
					errStatusList = append(errStatusList, errStatusMap)
				}
				confMap["cache_errstatus"] = errStatusList
			}

			if rule.Conf.CacheURLRewrite != nil {
				urlRewriteMap := map[string]interface{}{
					"sort_args":   rule.Conf.CacheURLRewrite.SortArgs,
					"ignore_case": rule.Conf.CacheURLRewrite.IgnoreCase,
				}

				if rule.Conf.CacheURLRewrite.Queries != nil {
					queriesMap := map[string]interface{}{
						"args_method": rule.Conf.CacheURLRewrite.Queries.ArgsMethod,
						"items":       rule.Conf.CacheURLRewrite.Queries.Items,
					}
					urlRewriteMap["queries"] = []interface{}{queriesMap}
				}

				if rule.Conf.CacheURLRewrite.Cookies != nil {
					cookiesMap := map[string]interface{}{
						"args_method": rule.Conf.CacheURLRewrite.Cookies.ArgsMethod,
						"items":       rule.Conf.CacheURLRewrite.Cookies.Items,
					}
					urlRewriteMap["cookies"] = []interface{}{cookiesMap}
				}

				confMap["cache_url_rewrite"] = []interface{}{urlRewriteMap}
			}

			if rule.Conf.CacheShare != nil {
				cacheShareMap := map[string]interface{}{
					"scheme": rule.Conf.CacheShare.Scheme,
				}
				confMap["cache_share"] = []interface{}{cacheShareMap}
			}

			ruleMap["conf"] = []interface{}{confMap}
		}

		ruleList = append(ruleList, ruleMap)
	}

	if err := d.Set("list", ruleList); err != nil {
		log.Printf("[WARN] Failed to set list: %v", err)
	}

	// Handle result_output_file if provided
	if _, ok := d.GetOk("result_output_file"); ok {
		outputData := map[string]interface{}{
			"business_id":   businessID,
			"business_type": businessType,
			"total":         response.Data.Total,
			"list":          ruleList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
	}

	log.Printf("[INFO] Cache rules read successfully: total=%d", response.Data.Total)
	return nil
}
