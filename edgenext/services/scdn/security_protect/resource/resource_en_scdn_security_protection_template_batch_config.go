package resource

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnSecurityProtectionTemplateBatchConfig returns the SCDN security protection template batch config resource
func ResourceEdgenextScdnSecurityProtectionTemplateBatchConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnSecurityProtectionTemplateBatchConfigCreate,
		Read:   resourceScdnSecurityProtectionTemplateBatchConfigRead,
		Update: resourceScdnSecurityProtectionTemplateBatchConfigUpdate,
		Delete: resourceScdnSecurityProtectionTemplateBatchConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"template_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Template ID list",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"all": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "All flag (0 or 1)",
			},
			"domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domain list",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domain ID list",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"ddos_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "DDoS protection configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_ddos_protection": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Application layer DDoS protection configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status: on, off, keep",
									},
									"ai_cc_status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AI protection status: on, off",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Protection type: default, normal, strict, captcha, keep",
									},
									"need_attack_detection": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Attack detection switch: 0 or 1",
									},
									"ai_status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "AI status: on, off",
									},
								},
							},
						},
						"visitor_authentication": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Visitor authentication configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status: on, off",
									},
									"auth_token": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Authentication token",
									},
									"pass_still_check": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Pass still check: 0 or 1",
									},
								},
							},
						},
					},
				},
			},
			"waf_rule_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "WAF rule configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"waf_rule_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "WAF rule config",
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
									"waf_strategy_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "WAF strategy ID",
									},
								},
							},
						},
						"waf_intercept_page": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "WAF intercept page config",
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
										Description: "Type: custom, default, keep",
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
				},
			},
			"bot_management_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Bot management configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"business_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Business ID",
						},
						"ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "ID list",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"precise_access_control_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Precise access control configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action: add, cover",
						},
						"policies": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Policy list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Policy type",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Policy action",
									},
									"action_data": {
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Action data",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"rules": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Rules list",
										Elem: &schema.Schema{
											Type: schema.TypeMap,
										},
									},
									"from": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "From source",
									},
									"status": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Status",
									},
								},
							},
						},
					},
				},
			},
			"fail_templates": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Failed templates",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceScdnSecurityProtectionTemplateBatchConfigCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnSecurityProtectionTemplateBatchConfigUpdate(d, m)
}

func resourceScdnSecurityProtectionTemplateBatchConfigRead(d *schema.ResourceData, m interface{}) error {
	// Batch config is a one-time operation, we can't really "read" the config state from API
	// Just verify the resource exists by checking if ID is set
	if d.Id() == "" {
		d.SetId("")
		return nil
	}

	// Try to parse template_ids from ID if not set in state
	if _, ok := d.GetOk("template_ids"); !ok {
		// ID format: template-batch-config-[template_ids]
		// For simplicity, we'll just ensure the resource exists
		// The actual template_ids should be preserved in state
	}

	return nil
}

func resourceScdnSecurityProtectionTemplateBatchConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.SecurityProtectionTemplateBatchConfigRequest{}

	// Set template_ids
	templateIDs := d.Get("template_ids").([]interface{})
	req.TemplateIDs = make([]int, len(templateIDs))
	for i, v := range templateIDs {
		req.TemplateIDs[i] = v.(int)
	}

	// Set all flag
	if all, ok := d.GetOk("all"); ok {
		req.All = all.(int)
	}

	// Set domains
	if domains, ok := d.GetOk("domains"); ok {
		domainsList := domains.([]interface{})
		req.Domains = make([]string, len(domainsList))
		for i, v := range domainsList {
			req.Domains[i] = v.(string)
		}
	}

	// Set domain_ids
	if domainIDs, ok := d.GetOk("domain_ids"); ok {
		domainIDsList := domainIDs.([]interface{})
		req.DomainIDs = make([]int, len(domainIDsList))
		for i, v := range domainIDsList {
			req.DomainIDs[i] = v.(int)
		}
	}

	// Build ddos_config
	if v, ok := d.GetOk("ddos_config"); ok {
		ddosConfigList := v.([]interface{})
		if len(ddosConfigList) > 0 {
			ddosConfigMap := ddosConfigList[0].(map[string]interface{})
			ddosConfig := &scdn.DdosProtectionGetConfigData{}

			if appDdos, ok := ddosConfigMap["application_ddos_protection"].([]interface{}); ok && len(appDdos) > 0 {
				appDdosMap := appDdos[0].(map[string]interface{})
				appDdosProtection := &scdn.ApplicationDdosProtection{}
				if val, ok := appDdosMap["status"].(string); ok && val != "" {
					appDdosProtection.Status = val
				}
				if val, ok := appDdosMap["ai_cc_status"].(string); ok && val != "" {
					appDdosProtection.AICCStatus = val
				}
				if val, ok := appDdosMap["type"].(string); ok && val != "" {
					appDdosProtection.Type = val
				}
				if val, ok := appDdosMap["need_attack_detection"].(int); ok {
					appDdosProtection.NeedAttackDetection = val
				}
				if val, ok := appDdosMap["ai_status"].(string); ok && val != "" {
					appDdosProtection.AIStatus = val
				}
				ddosConfig.ApplicationDdosProtection = appDdosProtection
			}

			if visitorAuth, ok := ddosConfigMap["visitor_authentication"].([]interface{}); ok && len(visitorAuth) > 0 {
				visitorAuthMap := visitorAuth[0].(map[string]interface{})
				visitorAuthentication := &scdn.VisitorAuthentication{}
				if val, ok := visitorAuthMap["status"].(string); ok && val != "" {
					visitorAuthentication.Status = val
				}
				if val, ok := visitorAuthMap["auth_token"].(string); ok && val != "" {
					visitorAuthentication.AuthToken = val
				}
				if val, ok := visitorAuthMap["pass_still_check"].(int); ok {
					visitorAuthentication.PassStillCheck = val
				}
				ddosConfig.VisitorAuthentication = visitorAuthentication
			}

			req.DdosConfig = ddosConfig
		}
	}

	// Build waf_rule_config
	if v, ok := d.GetOk("waf_rule_config"); ok {
		wafConfigList := v.([]interface{})
		if len(wafConfigList) > 0 {
			wafConfigMap := wafConfigList[0].(map[string]interface{})
			wafRuleConfig := &scdn.BatchUpdateWafRuleConfigRequest{}

			if wafRule, ok := wafConfigMap["waf_rule_config"].([]interface{}); ok && len(wafRule) > 0 {
				wafRuleMap := wafRule[0].(map[string]interface{})
				wafRuleCfg := &scdn.WafRuleConfig{}
				if val, ok := wafRuleMap["status"].(string); ok && val != "" {
					wafRuleCfg.Status = val
				}
				if val, ok := wafRuleMap["ai_status"].(string); ok && val != "" {
					wafRuleCfg.AIStatus = val
				}
				if val, ok := wafRuleMap["waf_level"].(string); ok && val != "" {
					wafRuleCfg.WafLevel = val
				}
				if val, ok := wafRuleMap["waf_mode"].(string); ok && val != "" {
					wafRuleCfg.WafMode = val
				}
				if val, ok := wafRuleMap["waf_strategy_id"].(int); ok {
					wafRuleCfg.WafStrategyID = val
				}
				wafRuleConfig.WafRuleConfig = wafRuleCfg
			}

			if wafInterceptPage, ok := wafConfigMap["waf_intercept_page"].([]interface{}); ok && len(wafInterceptPage) > 0 {
				wafInterceptPageMap := wafInterceptPage[0].(map[string]interface{})
				wafInterceptPageCfg := &scdn.WafInterceptPage{}
				if val, ok := wafInterceptPageMap["status"].(string); ok && val != "" {
					wafInterceptPageCfg.Status = val
				}
				if val, ok := wafInterceptPageMap["type"].(string); ok && val != "" {
					wafInterceptPageCfg.Type = val
				}
				if val, ok := wafInterceptPageMap["content"].(string); ok && val != "" {
					wafInterceptPageCfg.Content = val
				}
				wafRuleConfig.WafInterceptPage = wafInterceptPageCfg
			}

			req.WafRuleConfig = wafRuleConfig
		}
	}

	// Build bot_management_config
	if v, ok := d.GetOk("bot_management_config"); ok {
		botConfigList := v.([]interface{})
		if len(botConfigList) > 0 {
			botConfigMap := botConfigList[0].(map[string]interface{})
			botManagementConfig := &scdn.UpdateBotManagementConfigRequest{}
			if val, ok := botConfigMap["business_id"].(int); ok {
				botManagementConfig.BusinessID = val
			}
			if ids, ok := botConfigMap["ids"].([]interface{}); ok {
				botManagementConfig.IDs = make([]int, len(ids))
				for i, v := range ids {
					botManagementConfig.IDs[i] = v.(int)
				}
			}
			req.BotManagementConfig = botManagementConfig
		}
	}

	// Build precise_access_control_config
	if v, ok := d.GetOk("precise_access_control_config"); ok {
		preciseConfigList := v.([]interface{})
		if len(preciseConfigList) > 0 {
			preciseConfigMap := preciseConfigList[0].(map[string]interface{})
			preciseAccessControlConfig := &scdn.UpdatePreciseAccessControlConfigRequest{
				Action: preciseConfigMap["action"].(string),
			}
			if policies, ok := preciseConfigMap["policies"].([]interface{}); ok {
				preciseAccessControlConfig.Policies = make([]scdn.PreciseAccessControlPolicy, len(policies))
				for i, policy := range policies {
					policyMap := policy.(map[string]interface{})
					policyCfg := scdn.PreciseAccessControlPolicy{}
					if val, ok := policyMap["type"].(string); ok {
						policyCfg.Type = val
					}
					if val, ok := policyMap["action"].(string); ok {
						policyCfg.Action = val
					}
					if val, ok := policyMap["action_data"].(map[string]interface{}); ok {
						policyCfg.ActionData = make(map[string]interface{})
						for k, v := range val {
							policyCfg.ActionData[k] = v
						}
					}
					if val, ok := policyMap["rules"].([]interface{}); ok {
						policyCfg.Rules = make([]map[string]interface{}, len(val))
						for j, rule := range val {
							if ruleMap, ok := rule.(map[string]interface{}); ok {
								policyCfg.Rules[j] = ruleMap
							}
						}
					}
					if val, ok := policyMap["from"].(string); ok {
						policyCfg.From = val
					}
					if val, ok := policyMap["status"].(int); ok {
						policyCfg.Status = val
					}
					preciseAccessControlConfig.Policies[i] = policyCfg
				}
			}
			req.PreciseAccessControlConfig = preciseAccessControlConfig
		}
	}

	log.Printf("[INFO] Batch configuring SCDN security protection templates: template_ids=%v", req.TemplateIDs)
	response, err := service.BatchConfigSecurityProtectionTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to batch config security protection template: %w", err)
	}

	// Set resource ID - use sorted template IDs for stable ID generation
	templateIDsCopy := make([]int, len(req.TemplateIDs))
	copy(templateIDsCopy, req.TemplateIDs)
	sort.Ints(templateIDsCopy)
	idParts := make([]string, len(templateIDsCopy))
	for i, id := range templateIDsCopy {
		idParts[i] = fmt.Sprintf("%d", id)
	}
	d.SetId(fmt.Sprintf("template-batch-config-%s", strings.Join(idParts, "-")))

	// Set fail_templates if any
	if len(response.Data.FailTemplates) > 0 {
		failTemplatesMap := make(map[string]interface{})
		for k, v := range response.Data.FailTemplates {
			failTemplatesMap[k] = v
		}
		if err := d.Set("fail_templates", failTemplatesMap); err != nil {
			log.Printf("[WARN] Failed to set fail_templates: %v", err)
		}
	}

	return nil
}

func resourceScdnSecurityProtectionTemplateBatchConfigDelete(d *schema.ResourceData, m interface{}) error {
	// Batch config cannot be deleted, only reset
	// For now, we just remove the resource from state
	log.Printf("[INFO] Removing SCDN security protection template batch config from state")
	d.SetId("")
	return nil
}
