package networkspeed

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnNetworkSpeedRules returns the SCDN network speed rules data source
func DataSourceEdgenextScdnNetworkSpeedRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnNetworkSpeedRulesRead,

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Business ID",
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Business type: 'tpl' or 'global'",
			},
			"config_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'",
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
				Description: "List of rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule ID",
						},
						"business_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Business type",
						},
						"business_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Business ID",
						},
						"config_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config group",
						},
						"custom_page": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Custom page rule",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Status code",
									},
									"page_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Page type",
									},
									"page_content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Page content",
									},
								},
							},
						},
						"upstream_uri_change_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Upstream URI change rule",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"typ": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type",
									},
									"action": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action",
									},
									"match": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Match value",
									},
									"target": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target value",
									},
								},
							},
						},
						"resp_headers_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Response headers rule",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type",
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content",
									},
									"remark": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Remark",
									},
								},
							},
						},
						"customized_req_headers_rule": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Customized request headers rule",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type",
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Content",
									},
									"remark": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Remark",
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

func dataSourceScdnNetworkSpeedRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)
	configGroup := d.Get("config_group").(string)

	req := scdn.NetworkSpeedGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ConfigGroup:  configGroup,
	}

	log.Printf("[INFO] Reading SCDN network speed rules: business_id=%d, business_type=%s, config_group=%s", businessID, businessType, configGroup)
	response, err := service.GetNetworkSpeedRules(req)
	if err != nil {
		return fmt.Errorf("failed to read network speed rules: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s-%s", businessID, businessType, configGroup))

	// Set computed fields
	if err := d.Set("total", response.Data.Total); err != nil {
		log.Printf("[WARN] Failed to set total: %v", err)
	}

	// Convert list to schema format
	ruleList := make([]map[string]interface{}, 0, len(response.Data.List))
	for _, rule := range response.Data.List {
		ruleMap := map[string]interface{}{
			"id":            rule.ID,
			"business_type": rule.BusinessType,
			"business_id":   rule.BusinessID,
			"config_group":  rule.ConfigGroup,
		}

		if rule.CustomPage != nil {
			customPageMap := map[string]interface{}{
				"status_code":  rule.CustomPage.StatusCode,
				"page_type":    rule.CustomPage.PageType,
				"page_content": rule.CustomPage.PageContent,
			}
			ruleMap["custom_page"] = []interface{}{customPageMap}
		}

		if rule.UpstreamURIChangeRule != nil {
			uriChangeMap := map[string]interface{}{
				"typ":    rule.UpstreamURIChangeRule.Type,
				"action": rule.UpstreamURIChangeRule.Action,
				"match":  rule.UpstreamURIChangeRule.Match,
				"target": rule.UpstreamURIChangeRule.Target,
			}
			ruleMap["upstream_uri_change_rule"] = []interface{}{uriChangeMap}
		}

		if rule.RespHeadersRule != nil {
			respHeadersMap := map[string]interface{}{
				"type":    rule.RespHeadersRule.Type,
				"content": rule.RespHeadersRule.Content,
				"remark":  rule.RespHeadersRule.Remark,
			}
			ruleMap["resp_headers_rule"] = []interface{}{respHeadersMap}
		}

		if rule.CustomizedReqHeadersRule != nil {
			reqHeadersMap := map[string]interface{}{
				"type":    rule.CustomizedReqHeadersRule.Type,
				"content": rule.CustomizedReqHeadersRule.Content,
				"remark":  rule.CustomizedReqHeadersRule.Remark,
			}
			ruleMap["customized_req_headers_rule"] = []interface{}{reqHeadersMap}
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
			"config_group":  configGroup,
			"total":         response.Data.Total,
			"list":          ruleList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
	}

	log.Printf("[INFO] Network speed rules read successfully: total=%d", response.Data.Total)
	return nil
}
