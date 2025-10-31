package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionWafConfig returns the SCDN security protection WAF config data source
func DataSourceEdgenextScdnSecurityProtectionWafConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionWafConfigRead,

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Business ID",
			},
			"keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify config keys",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"waf_rule_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "WAF rule configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: on, off, keep",
						},
						"ai_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AI status: on, off",
						},
						"waf_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protection level: general, strict, keep",
						},
						"waf_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protection mode: off, active, block, ban, keep",
						},
					},
				},
			},
			"waf_intercept_page": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "WAF intercept page configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: on, off",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Page type: custom, default, keep",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom content",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnSecurityProtectionWafConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.WafRuleConfigGetRequest{
		BusinessID: businessID,
	}

	// Get keys if provided
	if v, ok := d.GetOk("keys"); ok {
		keysList := v.([]interface{})
		keys := make([]string, len(keysList))
		for i, v := range keysList {
			keys[i] = v.(string)
		}
		req.Keys = keys
	}

	log.Printf("[INFO] Querying SCDN security protection WAF config: business_id=%d", businessID)
	response, err := service.GetWafRuleConfig(req)
	if err != nil {
		return fmt.Errorf("failed to query WAF rule config: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("waf-config-%d", businessID))

	// Set waf_rule_config
	if response.Data.WafRuleConfig != nil {
		wafRule := []map[string]interface{}{
			{
				"id":        response.Data.WafRuleConfig.ID,
				"status":    response.Data.WafRuleConfig.Status,
				"ai_status": response.Data.WafRuleConfig.AIStatus,
				"waf_level": response.Data.WafRuleConfig.WafLevel,
				"waf_mode":  response.Data.WafRuleConfig.WafMode,
			},
		}
		if err := d.Set("waf_rule_config", wafRule); err != nil {
			return fmt.Errorf("error setting waf_rule_config: %w", err)
		}
	}

	// Set waf_intercept_page
	if response.Data.WafInterceptPage != nil {
		interceptPage := []map[string]interface{}{
			{
				"id":      response.Data.WafInterceptPage.ID,
				"status":  response.Data.WafInterceptPage.Status,
				"type":    response.Data.WafInterceptPage.Type,
				"content": response.Data.WafInterceptPage.Content,
			},
		}
		if err := d.Set("waf_intercept_page", interceptPage); err != nil {
			return fmt.Errorf("error setting waf_intercept_page: %w", err)
		}
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"waf_rule_config":    response.Data.WafRuleConfig,
			"waf_intercept_page": response.Data.WafInterceptPage,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN security protection WAF config queried successfully: business_id=%d", businessID)
	return nil
}
