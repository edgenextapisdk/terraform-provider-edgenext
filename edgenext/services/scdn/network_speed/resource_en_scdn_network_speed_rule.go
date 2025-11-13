package networkspeed

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnNetworkSpeedRule returns the SCDN network speed rule resource
func ResourceEdgenextScdnNetworkSpeedRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnNetworkSpeedRuleCreate,
		Read:   resourceScdnNetworkSpeedRuleRead,
		Update: resourceScdnNetworkSpeedRuleUpdate,
		Delete: resourceScdnNetworkSpeedRuleDelete,

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
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business type: 'tpl' or 'global'",
			},
			"config_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule ID for updating existing rule. If provided, this will update the rule instead of creating a new one.",
			},
			// Rule types - only one should be set based on config_group
			"custom_page": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Custom page rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Status code",
						},
						"page_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Page type",
						},
						"page_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Page content",
						},
					},
				},
			},
			"upstream_uri_change_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Upstream URI change rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"typ": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type",
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action",
						},
						"match": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Match value",
						},
						"target": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Target value",
						},
					},
				},
			},
			"resp_headers_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Response headers rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Content",
						},
						"remark": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remark",
						},
					},
				},
			},
			"customized_req_headers_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Customized request headers rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Content",
						},
						"remark": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remark",
						},
					},
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID",
			},
		},
	}
}

func resourceScdnNetworkSpeedRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)
	configGroup := d.Get("config_group").(string)

	// Check if rule_id is provided - if so, update existing rule instead of creating
	if ruleIDVal, ok := d.GetOk("rule_id"); ok {
		ruleID := ruleIDVal.(int)
		log.Printf("[INFO] rule_id provided, updating existing rule instead of creating: rule_id=%d", ruleID)

		// Set the composite ID first so Update knows which rule to update
		d.SetId(fmt.Sprintf("%d-%s-%s-%d", businessID, businessType, configGroup, ruleID))

		// Call Update function
		return resourceScdnNetworkSpeedRuleUpdate(d, m)
	}

	req := scdn.NetworkSpeedCreateRuleRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ConfigGroup:  configGroup,
	}

	// Build rule based on config_group
	if customPageList, ok := d.GetOk("custom_page"); ok && len(customPageList.([]interface{})) > 0 {
		customPageMap := customPageList.([]interface{})[0].(map[string]interface{})
		req.CustomPage = &scdn.CustomPageRule{
			StatusCode:  customPageMap["status_code"].(int),
			PageType:    customPageMap["page_type"].(string),
			PageContent: customPageMap["page_content"].(string),
		}
	}

	if uriChangeList, ok := d.GetOk("upstream_uri_change_rule"); ok && len(uriChangeList.([]interface{})) > 0 {
		uriChangeMap := uriChangeList.([]interface{})[0].(map[string]interface{})
		req.UpstreamURIChangeRule = &scdn.UpstreamURIChangeRule{
			Type:   uriChangeMap["typ"].(string),
			Action: uriChangeMap["action"].(string),
			Match:  uriChangeMap["match"].(string),
			Target: uriChangeMap["target"].(string),
		}
	}

	if respHeadersList, ok := d.GetOk("resp_headers_rule"); ok && len(respHeadersList.([]interface{})) > 0 {
		respHeadersMap := respHeadersList.([]interface{})[0].(map[string]interface{})
		req.RespHeadersRule = &scdn.RespHeadersRule{
			Type:    respHeadersMap["type"].(string),
			Content: respHeadersMap["content"].(string),
		}
		if remark, ok := respHeadersMap["remark"]; ok {
			req.RespHeadersRule.Remark = remark.(string)
		}
	}

	if reqHeadersList, ok := d.GetOk("customized_req_headers_rule"); ok && len(reqHeadersList.([]interface{})) > 0 {
		reqHeadersMap := reqHeadersList.([]interface{})[0].(map[string]interface{})
		req.CustomizedReqHeadersRule = &scdn.CustomizedReqHeadersRule{
			Type:    reqHeadersMap["type"].(string),
			Content: reqHeadersMap["content"].(string),
		}
		if remark, ok := reqHeadersMap["remark"]; ok {
			req.CustomizedReqHeadersRule.Remark = remark.(string)
		}
	}

	log.Printf("[INFO] Creating SCDN network speed rule: business_id=%d, business_type=%s, config_group=%s", businessID, businessType, configGroup)
	response, err := service.CreateNetworkSpeedRule(req)
	if err != nil {
		return fmt.Errorf("failed to create network speed rule: %w", err)
	}

	// Set composite ID format: business_id-business_type-config_group-rule_id
	d.SetId(fmt.Sprintf("%d-%s-%s-%d", businessID, businessType, configGroup, response.Data.ID))
	return resourceScdnNetworkSpeedRuleRead(d, m)
}

func resourceScdnNetworkSpeedRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	var businessID int
	var businessType string
	var configGroup string
	var ruleID int
	var err error

	// Parse ID if exists (format: business_id-business_type-config_group-rule_id)
	if d.Id() != "" {
		parts := strings.Split(d.Id(), "-")
		if len(parts) >= 4 {
			businessID, _ = strconv.Atoi(parts[0])
			businessType = parts[1]
			configGroup = parts[2]
			ruleID, _ = strconv.Atoi(parts[3])
		} else {
			// Fallback: try to get from state
			businessID = d.Get("business_id").(int)
			businessType = d.Get("business_type").(string)
			configGroup = d.Get("config_group").(string)
			ruleID, _ = strconv.Atoi(d.Id())
		}
	} else {
		businessID = d.Get("business_id").(int)
		businessType = d.Get("business_type").(string)
		configGroup = d.Get("config_group").(string)
		if ruleIDStr, ok := d.GetOk("rule_id"); ok {
			ruleID = ruleIDStr.(int)
		} else {
			return fmt.Errorf("rule_id is required for reading")
		}
	}

	// Get rules list to find the specific rule
	req := scdn.NetworkSpeedGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ConfigGroup:  configGroup,
	}

	log.Printf("[INFO] Reading SCDN network speed rule: rule_id=%d", ruleID)
	response, err := service.GetNetworkSpeedRules(req)
	if err != nil {
		return fmt.Errorf("failed to read network speed rule: %w", err)
	}

	// Find the rule
	var foundRule *scdn.NetworkSpeedRuleInfo
	for _, rule := range response.Data.List {
		if rule.ID == ruleID {
			foundRule = &rule
			break
		}
	}

	if foundRule == nil {
		log.Printf("[WARN] Network speed rule not found: %d", ruleID)
		d.SetId("")
		return nil
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s-%s-%d", businessID, businessType, configGroup, ruleID))

	// Set fields
	if err := d.Set("business_id", businessID); err != nil {
		log.Printf("[WARN] Failed to set business_id: %v", err)
	}
	if err := d.Set("business_type", businessType); err != nil {
		log.Printf("[WARN] Failed to set business_type: %v", err)
	}
	if err := d.Set("config_group", configGroup); err != nil {
		log.Printf("[WARN] Failed to set config_group: %v", err)
	}
	if err := d.Set("rule_id", ruleID); err != nil {
		log.Printf("[WARN] Failed to set rule_id: %v", err)
	}

	// Set rule-specific fields
	if foundRule.CustomPage != nil {
		customPageMap := map[string]interface{}{
			"status_code":  foundRule.CustomPage.StatusCode,
			"page_type":    foundRule.CustomPage.PageType,
			"page_content": foundRule.CustomPage.PageContent,
		}
		if err := d.Set("custom_page", []interface{}{customPageMap}); err != nil {
			log.Printf("[WARN] Failed to set custom_page: %v", err)
		}
	}

	if foundRule.UpstreamURIChangeRule != nil {
		uriChangeMap := map[string]interface{}{
			"typ":    foundRule.UpstreamURIChangeRule.Type,
			"action": foundRule.UpstreamURIChangeRule.Action,
			"match":  foundRule.UpstreamURIChangeRule.Match,
			"target": foundRule.UpstreamURIChangeRule.Target,
		}
		if err := d.Set("upstream_uri_change_rule", []interface{}{uriChangeMap}); err != nil {
			log.Printf("[WARN] Failed to set upstream_uri_change_rule: %v", err)
		}
	}

	if foundRule.RespHeadersRule != nil {
		respHeadersMap := map[string]interface{}{
			"type":    foundRule.RespHeadersRule.Type,
			"content": foundRule.RespHeadersRule.Content,
			"remark":  foundRule.RespHeadersRule.Remark,
		}
		if err := d.Set("resp_headers_rule", []interface{}{respHeadersMap}); err != nil {
			log.Printf("[WARN] Failed to set resp_headers_rule: %v", err)
		}
	}

	if foundRule.CustomizedReqHeadersRule != nil {
		reqHeadersMap := map[string]interface{}{
			"type":    foundRule.CustomizedReqHeadersRule.Type,
			"content": foundRule.CustomizedReqHeadersRule.Content,
			"remark":  foundRule.CustomizedReqHeadersRule.Remark,
		}
		if err := d.Set("customized_req_headers_rule", []interface{}{reqHeadersMap}); err != nil {
			log.Printf("[WARN] Failed to set customized_req_headers_rule: %v", err)
		}
	}

	log.Printf("[INFO] Network speed rule read successfully")
	return nil
}

func resourceScdnNetworkSpeedRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	ruleIDStr := d.Id()
	if ruleIDStr == "" {
		if ruleID, ok := d.GetOk("rule_id"); ok {
			ruleIDStr = strconv.Itoa(ruleID.(int))
		} else {
			return fmt.Errorf("rule_id is required for update")
		}
	}

	ruleID, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		// Try to parse from composite ID
		parts := strings.Split(ruleIDStr, "-")
		if len(parts) >= 4 {
			ruleID, _ = strconv.Atoi(parts[3])
		} else {
			return fmt.Errorf("invalid rule ID: %s", ruleIDStr)
		}
	}

	configGroup := d.Get("config_group").(string)

	req := scdn.NetworkSpeedUpdateRuleRequest{
		ID:          ruleID,
		ConfigGroup: configGroup,
	}

	// Build rule based on config_group (same as create)
	if customPageList, ok := d.GetOk("custom_page"); ok && len(customPageList.([]interface{})) > 0 {
		customPageMap := customPageList.([]interface{})[0].(map[string]interface{})
		req.CustomPage = &scdn.CustomPageRule{
			StatusCode:  customPageMap["status_code"].(int),
			PageType:    customPageMap["page_type"].(string),
			PageContent: customPageMap["page_content"].(string),
		}
	}

	if uriChangeList, ok := d.GetOk("upstream_uri_change_rule"); ok && len(uriChangeList.([]interface{})) > 0 {
		uriChangeMap := uriChangeList.([]interface{})[0].(map[string]interface{})
		req.UpstreamURIChangeRule = &scdn.UpstreamURIChangeRule{
			Type:   uriChangeMap["typ"].(string),
			Action: uriChangeMap["action"].(string),
			Match:  uriChangeMap["match"].(string),
			Target: uriChangeMap["target"].(string),
		}
	}

	if respHeadersList, ok := d.GetOk("resp_headers_rule"); ok && len(respHeadersList.([]interface{})) > 0 {
		respHeadersMap := respHeadersList.([]interface{})[0].(map[string]interface{})
		req.RespHeadersRule = &scdn.RespHeadersRule{
			Type:    respHeadersMap["type"].(string),
			Content: respHeadersMap["content"].(string),
		}
		if remark, ok := respHeadersMap["remark"]; ok {
			req.RespHeadersRule.Remark = remark.(string)
		}
	}

	if reqHeadersList, ok := d.GetOk("customized_req_headers_rule"); ok && len(reqHeadersList.([]interface{})) > 0 {
		reqHeadersMap := reqHeadersList.([]interface{})[0].(map[string]interface{})
		req.CustomizedReqHeadersRule = &scdn.CustomizedReqHeadersRule{
			Type:    reqHeadersMap["type"].(string),
			Content: reqHeadersMap["content"].(string),
		}
		if remark, ok := reqHeadersMap["remark"]; ok {
			req.CustomizedReqHeadersRule.Remark = remark.(string)
		}
	}

	log.Printf("[INFO] Updating SCDN network speed rule: rule_id=%d", ruleID)
	_, err = service.UpdateNetworkSpeedRule(req)
	if err != nil {
		return fmt.Errorf("failed to update network speed rule: %w", err)
	}

	return resourceScdnNetworkSpeedRuleRead(d, m)
}

func resourceScdnNetworkSpeedRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)
	configGroup := d.Get("config_group").(string)

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		// Try to parse from composite ID
		parts := strings.Split(d.Id(), "-")
		if len(parts) >= 4 {
			ruleID, _ = strconv.Atoi(parts[3])
		} else {
			return fmt.Errorf("invalid rule ID: %s", d.Id())
		}
	}

	req := scdn.NetworkSpeedDeleteRuleRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ConfigGroup:  configGroup,
		IDs:          []int{ruleID},
	}

	log.Printf("[INFO] Deleting SCDN network speed rule: rule_id=%d", ruleID)
	_, err = service.DeleteNetworkSpeedRule(req)
	if err != nil {
		return fmt.Errorf("failed to delete network speed rule: %w", err)
	}

	log.Printf("[INFO] Network speed rule deleted successfully")
	return nil
}
