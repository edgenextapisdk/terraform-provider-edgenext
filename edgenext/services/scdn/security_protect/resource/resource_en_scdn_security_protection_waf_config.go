package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnSecurityProtectionWafConfig returns the SCDN security protection WAF config resource
func ResourceEdgenextScdnSecurityProtectionWafConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnSecurityProtectionWafConfigCreate,
		Read:   resourceScdnSecurityProtectionWafConfigRead,
		Update: resourceScdnSecurityProtectionWafConfigUpdate,
		Delete: resourceScdnSecurityProtectionWafConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Business ID",
			},
			"waf_rule_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "WAF rule configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: on, off, keep",
						},
						"ai_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AI status: on, off",
						},
						"waf_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protection level: general, strict, keep",
						},
						"waf_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protection mode: off, active, block, ban, keep",
						},
					},
				},
			},
			"waf_intercept_page": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "WAF intercept page configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: on, off",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Page type: custom, default, keep",
						},
						"content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom content",
						},
					},
				},
			},
		},
	}
}

func resourceScdnSecurityProtectionWafConfigCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.WafRuleConfigUpdateRequest{
		BusinessID: businessID,
	}

	// Build waf_rule_config
	if v, ok := d.GetOk("waf_rule_config"); ok {
		wafRuleList := v.([]interface{})
		if len(wafRuleList) > 0 {
			wafRuleMap := wafRuleList[0].(map[string]interface{})
			wafRule := &scdn.WafRuleConfig{}
			if val, ok := wafRuleMap["status"].(string); ok && val != "" {
				wafRule.Status = val
			}
			if val, ok := wafRuleMap["ai_status"].(string); ok && val != "" {
				wafRule.AIStatus = val
			}
			if val, ok := wafRuleMap["waf_level"].(string); ok && val != "" {
				wafRule.WafLevel = val
			}
			if val, ok := wafRuleMap["waf_mode"].(string); ok && val != "" {
				wafRule.WafMode = val
			}
			req.WafRuleConfig = wafRule
		}
	}

	// Build waf_intercept_page
	if v, ok := d.GetOk("waf_intercept_page"); ok {
		interceptPageList := v.([]interface{})
		if len(interceptPageList) > 0 {
			interceptPageMap := interceptPageList[0].(map[string]interface{})
			interceptPage := &scdn.WafInterceptPage{}
			if val, ok := interceptPageMap["status"].(string); ok && val != "" {
				interceptPage.Status = val
			}
			if val, ok := interceptPageMap["type"].(string); ok && val != "" {
				interceptPage.Type = val
			}
			if val, ok := interceptPageMap["content"].(string); ok && val != "" {
				interceptPage.Content = val
			}
			req.WafInterceptPage = interceptPage
		}
	}

	log.Printf("[INFO] Creating/Updating SCDN security protection WAF config: business_id=%d", businessID)
	_, err := service.UpdateWafRuleConfig(req)
	if err != nil {
		return fmt.Errorf("failed to create/update WAF rule config: %w", err)
	}

	d.SetId(fmt.Sprintf("waf-config-%d", businessID))
	return resourceScdnSecurityProtectionWafConfigRead(d, m)
}

func resourceScdnSecurityProtectionWafConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.WafRuleConfigGetRequest{
		BusinessID: businessID,
	}

	log.Printf("[INFO] Reading SCDN security protection WAF config: business_id=%d", businessID)
	response, err := service.GetWafRuleConfig(req)
	if err != nil {
		return fmt.Errorf("failed to read WAF rule config: %w", err)
	}

	// Set waf_rule_config
	if response.Data.WafRuleConfig != nil {
		wafRule := []map[string]interface{}{
			{
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
				"status":  response.Data.WafInterceptPage.Status,
				"type":    response.Data.WafInterceptPage.Type,
				"content": response.Data.WafInterceptPage.Content,
			},
		}
		if err := d.Set("waf_intercept_page", interceptPage); err != nil {
			return fmt.Errorf("error setting waf_intercept_page: %w", err)
		}
	}

	return nil
}

func resourceScdnSecurityProtectionWafConfigUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnSecurityProtectionWafConfigCreate(d, m)
}

func resourceScdnSecurityProtectionWafConfigDelete(d *schema.ResourceData, m interface{}) error {
	// WAF config cannot be deleted, only reset
	// Set all values to default/off
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.WafRuleConfigUpdateRequest{
		BusinessID: businessID,
		WafRuleConfig: &scdn.WafRuleConfig{
			Status:   "off",
			AIStatus: "off",
			WafLevel: "general",
			WafMode:  "off",
		},
		WafInterceptPage: &scdn.WafInterceptPage{
			Status:  "off",
			Type:    "default",
			Content: "",
		},
	}

	log.Printf("[INFO] Resetting SCDN security protection WAF config: business_id=%d", businessID)
	_, err := service.UpdateWafRuleConfig(req)
	if err != nil {
		return fmt.Errorf("failed to reset WAF rule config: %w", err)
	}

	d.SetId("")
	return nil
}
