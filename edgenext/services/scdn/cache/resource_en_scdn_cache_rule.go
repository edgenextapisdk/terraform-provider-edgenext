package cache

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCacheRule returns the SCDN cache rule resource
func ResourceEdgenextScdnCacheRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCacheRuleCreate,
		Read:   resourceScdnCacheRuleRead,
		Update: resourceScdnCacheRuleUpdate,
		Delete: resourceScdnCacheRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Business ID (template ID for 'tpl' type, domain ID for 'domain' type)",
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business type: 'tpl' (template) or 'domain'",
			},
			"rule_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule ID for updating existing rule. If provided, this will update the rule instead of creating a new one.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name",
			},
			"expr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Wirefilter rule. Empty string means 'allow all'. If not set (null), keeps existing value.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule remark",
			},
			"conf": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Cache configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nocache": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Cache eligibility (true: bypass cache, false: cache)",
						},
						"cache_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Edge TTL cache configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cachetime": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache time",
									},
									"ignore_cache_time": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Ignore source cache time",
									},
									"ignore_nocache_header": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Ignore no-cache header",
									},
									"no_cache_control_op": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "No cache control operation",
									},
									"action": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cache action: 'default', 'nocache', 'cachetime', or 'force'",
									},
								},
							},
						},
						"browser_cache_rule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Browser cache configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cachetime": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Cache time",
									},
									"ignore_cache_time": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Ignore source cache time (cache-control)",
									},
									"nocache": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether to cache",
									},
								},
							},
						},
						"cache_errstatus": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Status code cache configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cachetime": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Status code cache time",
									},
									"err_status": {
										Type:        schema.TypeList,
										Required:    true,
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
							Optional:    true,
							MaxItems:    1,
							Description: "Custom cache key configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sort_args": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Parameter sorting",
									},
									"ignore_case": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Ignore case",
									},
									"queries": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Query string processing",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"args_method": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'",
												},
												"items": {
													Type:        schema.TypeList,
													Required:    true,
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
										Optional:    true,
										MaxItems:    1,
										Description: "Cookie processing",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"args_method": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'",
												},
												"items": {
													Type:        schema.TypeList,
													Required:    true,
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
							Optional:    true,
							MaxItems:    1,
							Description: "Cache sharing configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scheme": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "HTTP/HTTPS cache sharing method: '', 'http' or 'https'",
									},
								},
							},
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
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type: 'domain', 'tpl', or 'global'",
			},
		},
	}
}

func resourceScdnCacheRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	// Check if rule_id is provided - if so, read existing rule instead of creating
	if ruleIDVal, ok := d.GetOk("rule_id"); ok {
		ruleID := ruleIDVal.(int)
		log.Printf("[INFO] rule_id provided, reading existing rule instead of creating: rule_id=%d", ruleID)

		// Set the composite ID first
		d.SetId(fmt.Sprintf("%d-%s-%d", businessID, businessType, ruleID))

		// Check if name is provided - if not, just read the rule
		name := d.Get("name").(string)
		expr, exprSet := d.GetOk("expr")
		exprStr := ""
		if exprSet {
			exprStr = expr.(string)
		}

		// If name is empty, just read the rule (query scenario)
		// This is the typical case when importing an existing rule
		if name == "" {
			log.Printf("[INFO] name is empty, reading rule without update")
			return resourceScdnCacheRuleRead(d, m)
		}

		// If name is provided, this means user wants to update the rule
		// Call Update function to handle the update
		// Note: expr is optional - null means keep existing, empty string means "allow all"
		log.Printf("[INFO] name provided, updating rule (expr set: %v, expr=%q)", exprSet, exprStr)
		return resourceScdnCacheRuleUpdate(d, m)
	}

	// Build conf from schema
	conf, err := buildCacheRuleConfFromSchema(d)
	if err != nil {
		return fmt.Errorf("failed to build cache rule conf: %w", err)
	}

	// Get expr - if not set, use empty string as default (means "allow all")
	expr := ""
	if exprVal, ok := d.GetOk("expr"); ok {
		expr = exprVal.(string)
	}

	req := scdn.CacheRuleCreateRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		Name:         d.Get("name").(string),
		Expr:         expr,
		Conf:         conf,
	}

	if remark, ok := d.GetOk("remark"); ok {
		req.Remark = remark.(string)
	}

	log.Printf("[INFO] Creating SCDN cache rule: business_id=%d, business_type=%s, name=%s", businessID, businessType, req.Name)
	response, err := service.CreateCacheRule(req)
	if err != nil {
		return fmt.Errorf("failed to create cache rule: %w", err)
	}

	// Set composite ID format: business_id-business_type-rule_id
	ruleID := response.Data.ID
	d.SetId(fmt.Sprintf("%d-%s-%d", businessID, businessType, ruleID))

	// Set basic fields from creation response to avoid read issues immediately after creation
	if err := d.Set("business_id", businessID); err != nil {
		log.Printf("[WARN] Failed to set business_id: %v", err)
	}
	if err := d.Set("business_type", businessType); err != nil {
		log.Printf("[WARN] Failed to set business_type: %v", err)
	}
	if err := d.Set("rule_id", ruleID); err != nil {
		log.Printf("[WARN] Failed to set rule_id: %v", err)
	}
	if err := d.Set("name", req.Name); err != nil {
		log.Printf("[WARN] Failed to set name: %v", err)
	}
	if err := d.Set("expr", req.Expr); err != nil {
		log.Printf("[WARN] Failed to set expr: %v", err)
	}
	if req.Remark != "" {
		if err := d.Set("remark", req.Remark); err != nil {
			log.Printf("[WARN] Failed to set remark: %v", err)
		}
	}

	// Set conf from request
	if req.Conf != nil {
		confMap := buildCacheRuleConfToSchema(req.Conf)
		if err := d.Set("conf", []interface{}{confMap}); err != nil {
			log.Printf("[WARN] Failed to set conf: %v", err)
		}
	}

	log.Printf("[INFO] SCDN cache rule created successfully: rule_id=%d", ruleID)

	// Try to read full details, but don't fail if rule is not immediately available
	// This can happen due to API eventual consistency
	readErr := resourceScdnCacheRuleRead(d, m)
	if readErr != nil {
		log.Printf("[WARN] Failed to read cache rule immediately after creation (this is normal due to API delay): %v", readErr)
		// Don't return error, we've already set the basic fields from creation response
		// The next terraform apply will sync the state
	}

	return nil
}

func resourceScdnCacheRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	var businessID int
	var businessType string
	var ruleID int
	var err error

	// Parse ID if exists (format: business_id-business_type-rule_id)
	if d.Id() != "" {
		parts := strings.Split(d.Id(), "-")
		if len(parts) >= 3 {
			businessID, _ = strconv.Atoi(parts[0])
			businessType = parts[1]
			ruleID, _ = strconv.Atoi(parts[2])
		} else {
			// Fallback: try to get from state
			businessID = d.Get("business_id").(int)
			businessType = d.Get("business_type").(string)
			ruleID, _ = strconv.Atoi(d.Id())
		}
	} else {
		businessID = d.Get("business_id").(int)
		businessType = d.Get("business_type").(string)
		if ruleIDStr, ok := d.GetOk("rule_id"); ok {
			ruleID = ruleIDStr.(int)
		} else {
			return fmt.Errorf("rule_id is required for reading")
		}
	}

	// Get rule by id parameter (API supports id query parameter)
	req := scdn.CacheRuleGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ID:           ruleID, // Convert rule_id to id query parameter
	}

	log.Printf("[INFO] Reading SCDN cache rule: rule_id=%d (using id query parameter)", ruleID)
	response, err := service.GetCacheRules(req)
	if err != nil {
		return fmt.Errorf("failed to read cache rule: %w", err)
	}

	// Find the rule (API should return filtered list when id parameter is used)
	var foundRule *scdn.CacheRuleInfo
	if len(response.Data.List) > 0 {
		// When id parameter is provided, API should return only matching rule(s)
		for _, rule := range response.Data.List {
			if rule.ID == ruleID {
				foundRule = &rule
				break
			}
		}
	}

	if foundRule == nil {
		log.Printf("[WARN] Cache rule not found: %d", ruleID)
		// Don't clear the ID if it was just created - this might be due to API delay
		// Only clear ID if we're sure the rule doesn't exist (e.g., during normal read operations)
		// For now, keep the ID and let the next apply sync the state
		if d.Id() != "" {
			log.Printf("[INFO] Keeping resource ID %s - rule may not be immediately available due to API delay", d.Id())
			return nil
		}
		d.SetId("")
		return nil
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s-%d", businessID, businessType, ruleID))

	// Set basic fields
	if err := d.Set("business_id", businessID); err != nil {
		log.Printf("[WARN] Failed to set business_id: %v", err)
	}
	if err := d.Set("business_type", businessType); err != nil {
		log.Printf("[WARN] Failed to set business_type: %v", err)
	}
	if err := d.Set("rule_id", ruleID); err != nil {
		log.Printf("[WARN] Failed to set rule_id: %v", err)
	}
	if err := d.Set("name", foundRule.Name); err != nil {
		log.Printf("[WARN] Failed to set name: %v", err)
	}
	if err := d.Set("remark", foundRule.Remark); err != nil {
		log.Printf("[WARN] Failed to set remark: %v", err)
	}
	if err := d.Set("expr", foundRule.Expr); err != nil {
		log.Printf("[WARN] Failed to set expr: %v", err)
	}
	if err := d.Set("status", foundRule.Status); err != nil {
		log.Printf("[WARN] Failed to set status: %v", err)
	}
	if err := d.Set("weight", foundRule.Weight); err != nil {
		log.Printf("[WARN] Failed to set weight: %v", err)
	}
	if err := d.Set("type", foundRule.Type); err != nil {
		log.Printf("[WARN] Failed to set type: %v", err)
	}

	// Set conf
	if foundRule.Conf != nil {
		// Debug: log the API response to understand what we're receiving
		log.Printf("[DEBUG] API returned conf: %+v", foundRule.Conf)
		log.Printf("[DEBUG] API returned conf.CacheRule: %+v", foundRule.Conf.CacheRule)
		log.Printf("[DEBUG] API returned conf.CacheRule != nil: %v", foundRule.Conf.CacheRule != nil)
		if foundRule.Conf.CacheRule != nil {
			log.Printf("[DEBUG] API returned cache_rule details: cachetime=%d, action=%s",
				foundRule.Conf.CacheRule.CacheTime, foundRule.Conf.CacheRule.Action)
		}

		confMap := buildCacheRuleConfToSchema(foundRule.Conf)
		log.Printf("[DEBUG] Built confMap cache_rule: %+v", confMap["cache_rule"])
		log.Printf("[DEBUG] Built confMap cache_rule type: %T", confMap["cache_rule"])
		if cacheRuleList, ok := confMap["cache_rule"].([]interface{}); ok {
			log.Printf("[DEBUG] Built confMap cache_rule length: %d", len(cacheRuleList))
		}
		log.Printf("[DEBUG] API returned cache_errstatus: %+v (length: %d)", foundRule.Conf.CacheErrStatus, len(foundRule.Conf.CacheErrStatus))
		log.Printf("[DEBUG] Built confMap cache_errstatus: %+v", confMap["cache_errstatus"])
		if cacheErrStatusList, ok := confMap["cache_errstatus"].([]interface{}); ok {
			log.Printf("[DEBUG] Built confMap cache_errstatus length: %d", len(cacheErrStatusList))
		}

		// Preserve optional fields from state if API returned null/empty
		// This handles cases where API doesn't return these fields even though they exist in state
		if stateConf, ok := d.GetOk("conf.0"); ok {
			stateConfMap := stateConf.(map[string]interface{})

			// Preserve cache_rule if API returned nil but state has it
			if foundRule.Conf.CacheRule == nil {
				log.Printf("[DEBUG] API returned cache_rule as nil, checking state...")
				if stateCacheRule, ok := stateConfMap["cache_rule"]; ok {
					log.Printf("[DEBUG] State cache_rule found: %+v", stateCacheRule)
					if cacheRuleList, ok := stateCacheRule.([]interface{}); ok && len(cacheRuleList) > 0 {
						confMap["cache_rule"] = cacheRuleList
						log.Printf("[DEBUG] Preserved cache_rule from state: %+v", cacheRuleList)
					}
				}
			} else {
				log.Printf("[DEBUG] API returned cache_rule, using API value")
			}

			// Preserve browser_cache_rule if API returned nil but state has it
			if foundRule.Conf.BrowserCacheRule == nil {
				if stateBrowserCacheRule, ok := stateConfMap["browser_cache_rule"]; ok {
					if browserCacheRuleList, ok := stateBrowserCacheRule.([]interface{}); ok && len(browserCacheRuleList) > 0 {
						confMap["browser_cache_rule"] = browserCacheRuleList
					}
				}
			}

			// Note: cache_errstatus is an array that can contain multiple items
			// If API returns empty array, it means "no configuration", so we should use empty array
			// Do NOT preserve from state - always use API value to ensure consistency
			// The buildCacheRuleConfToSchema function already initializes it as empty array if API returns nil/empty

			// Preserve cache_url_rewrite if API returned nil but state has it
			if foundRule.Conf.CacheURLRewrite == nil {
				if stateCacheURLRewrite, ok := stateConfMap["cache_url_rewrite"]; ok {
					if cacheURLRewriteList, ok := stateCacheURLRewrite.([]interface{}); ok && len(cacheURLRewriteList) > 0 {
						confMap["cache_url_rewrite"] = cacheURLRewriteList
					}
				}
			}
		}

		if err := d.Set("conf", []interface{}{confMap}); err != nil {
			log.Printf("[WARN] Failed to set conf: %v", err)
		}
	}

	log.Printf("[INFO] Cache rule read successfully")
	return nil
}

func resourceScdnCacheRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	// Check if rule_id is provided in config - if not, create a new rule instead of updating
	if _, ok := d.GetOk("rule_id"); !ok {
		// No rule_id in config means user wants to create a new rule
		log.Printf("[INFO] rule_id not provided in config, creating new rule instead of updating")
		// Clear the old resource ID so Create function will create a new rule
		d.SetId("")
		return resourceScdnCacheRuleCreate(d, m)
	}

	var ruleID int
	var err error

	// Parse rule ID from resource ID or config
	if d.Id() != "" {
		parts := strings.Split(d.Id(), "-")
		if len(parts) >= 3 {
			ruleID, _ = strconv.Atoi(parts[2])
		} else {
			if ruleIDVal, ok := d.GetOk("rule_id"); ok {
				ruleID = ruleIDVal.(int)
			} else {
				// No rule_id and invalid resource ID, create new rule
				log.Printf("[INFO] Invalid resource ID and no rule_id, creating new rule")
				d.SetId("")
				return resourceScdnCacheRuleCreate(d, m)
			}
		}
	} else {
		if ruleIDVal, ok := d.GetOk("rule_id"); ok {
			ruleID = ruleIDVal.(int)
		} else {
			// No rule_id and no resource ID, create new rule
			log.Printf("[INFO] No resource ID and no rule_id, creating new rule")
			return resourceScdnCacheRuleCreate(d, m)
		}
	}

	// Verify the rule exists before updating
	// If rule doesn't exist, create a new one instead
	readReq := scdn.CacheRuleGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ID:           ruleID,
	}
	readResponse, err := service.GetCacheRules(readReq)
	if err != nil {
		log.Printf("[WARN] Failed to verify rule existence, creating new rule: %v", err)
		d.SetId("")
		return resourceScdnCacheRuleCreate(d, m)
	}

	// Check if rule exists
	ruleExists := false
	for _, rule := range readResponse.Data.List {
		if rule.ID == ruleID {
			ruleExists = true
			break
		}
	}

	if !ruleExists {
		log.Printf("[INFO] Rule %d does not exist, creating new rule instead of updating", ruleID)
		d.SetId("")
		return resourceScdnCacheRuleCreate(d, m)
	}

	// Check if conf or expr changed
	// If conf changed or expr changed, use UpdateCacheRuleConfig (supports expr and conf)
	// Otherwise use UpdateCacheRule (only supports name/remark)
	hasConfChange := d.HasChange("conf")
	hasExprChange := d.HasChange("expr")

	// Also check if conf is explicitly set in config (even if not changed)
	// If user provides conf block, we should use UpdateCacheRuleConfig to ensure conf is sent
	_, hasConfInConfig := d.GetOk("conf")

	if hasConfChange || hasExprChange || hasConfInConfig {
		// Update configuration
		conf, err := buildCacheRuleConfFromSchema(d)
		if err != nil {
			return fmt.Errorf("failed to build cache rule conf: %w", err)
		}

		req := scdn.CacheRuleUpdateConfigRequest{
			ID:           ruleID,
			BusinessID:   businessID,
			BusinessType: businessType,
			Conf:         conf,
		}

		// Get name and expr from configuration
		// expr is Optional: null means keep existing value, empty string means "allow all"
		// Note: d.GetOk returns false for empty string (zero value), so we need to check HasChange
		name := d.Get("name").(string)

		// Check if expr was explicitly set in the configuration (including empty string)
		// HasChange checks if the value changed from previous state OR if it's set in new config
		exprExplicitlySet := d.HasChange("expr")
		exprStr := ""
		if exprExplicitlySet {
			// Value changed or was set - get the new value (even if empty string)
			_, newVal := d.GetChange("expr")
			if newVal != nil {
				exprStr = newVal.(string)
			}
		} else {
			// No change - check if it exists in config
			// If it exists in config but hasn't changed, use the current value
			if exprVal, ok := d.GetOkExists("expr"); ok {
				exprStr = exprVal.(string)
				exprExplicitlySet = true
			}
		}

		// If name is empty or expr is not explicitly set, read current rule to get them
		if name == "" || !exprExplicitlySet {
			readReq := scdn.CacheRuleGetRulesRequest{
				BusinessID:   businessID,
				BusinessType: businessType,
				ID:           ruleID,
			}
			readResponse, err := service.GetCacheRules(readReq)
			if err != nil {
				return fmt.Errorf("failed to read cache rule to get name/expr: %w", err)
			}
			// Find the rule
			var foundRule *scdn.CacheRuleInfo
			for _, rule := range readResponse.Data.List {
				if rule.ID == ruleID {
					foundRule = &rule
					break
				}
			}
			if foundRule == nil {
				return fmt.Errorf("cache rule not found: rule_id=%d", ruleID)
			}
			if name == "" {
				name = foundRule.Name
			}
			if !exprExplicitlySet {
				exprStr = foundRule.Expr
			}
		}

		req.Name = name
		// Set expr: if explicitly set (even if empty string), use it
		// Empty string means "allow all", null means keep existing value
		req.Expr = exprStr
		log.Printf("[DEBUG] Setting expr=%q (explicitly set: %v, empty string means 'allow all')", exprStr, exprExplicitlySet)

		if remark, ok := d.GetOk("remark"); ok {
			req.Remark = remark.(string)
		}

		log.Printf("[INFO] Updating SCDN cache rule config: rule_id=%d", ruleID)
		_, err = service.UpdateCacheRuleConfig(req)
		if err != nil {
			return fmt.Errorf("failed to update cache rule config: %w", err)
		}
	} else {
		// Update name/remark only (no conf or expr changes)
		// Note: UpdateCacheRule API does not support expr field
		// If expr needs to be updated, it should go through UpdateCacheRuleConfig above
		req := scdn.CacheRuleUpdateRequest{
			ID: ruleID,
		}

		if name, ok := d.GetOk("name"); ok {
			req.Name = name.(string)
		}
		if remark, ok := d.GetOk("remark"); ok {
			req.Remark = remark.(string)
		}

		log.Printf("[INFO] Updating SCDN cache rule name/remark only (no conf/expr changes): rule_id=%d", ruleID)
		_, err = service.UpdateCacheRule(req)
		if err != nil {
			return fmt.Errorf("failed to update cache rule: %w", err)
		}
	}

	return resourceScdnCacheRuleRead(d, m)
}

func resourceScdnCacheRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	var ruleID int
	var err error

	// Parse rule ID from resource ID
	if d.Id() != "" {
		parts := strings.Split(d.Id(), "-")
		if len(parts) >= 3 {
			ruleID, _ = strconv.Atoi(parts[2])
		} else {
			if ruleIDVal, ok := d.GetOk("rule_id"); ok {
				ruleID = ruleIDVal.(int)
			} else {
				return fmt.Errorf("rule_id is required for delete")
			}
		}
	} else {
		if ruleIDVal, ok := d.GetOk("rule_id"); ok {
			ruleID = ruleIDVal.(int)
		} else {
			return fmt.Errorf("rule_id is required for delete")
		}
	}

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	req := scdn.CacheRuleDeleteRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		IDs:          []int{ruleID},
	}

	log.Printf("[INFO] Deleting SCDN cache rule: rule_id=%d", ruleID)
	_, err = service.DeleteCacheRule(req)
	if err != nil {
		return fmt.Errorf("failed to delete cache rule: %w", err)
	}

	d.SetId("")
	log.Printf("[INFO] Cache rule deleted successfully")
	return nil
}

// Helper functions

func buildCacheRuleConfFromSchema(d *schema.ResourceData) (*scdn.CacheRuleConf, error) {
	confList := d.Get("conf").([]interface{})
	if len(confList) == 0 {
		return nil, fmt.Errorf("conf is required")
	}

	confMap := confList[0].(map[string]interface{})
	conf := &scdn.CacheRuleConf{
		NoCache: confMap["nocache"].(bool),
	}

	// Build cache_rule
	if cacheRuleList, ok := confMap["cache_rule"].([]interface{}); ok && len(cacheRuleList) > 0 {
		cacheRuleMap := cacheRuleList[0].(map[string]interface{})
		conf.CacheRule = &scdn.CacheRule{
			CacheTime: cacheRuleMap["cachetime"].(int),
		}
		if val, ok := cacheRuleMap["ignore_cache_time"].(bool); ok {
			conf.CacheRule.IgnoreCacheTime = val
		}
		if val, ok := cacheRuleMap["ignore_nocache_header"].(bool); ok {
			conf.CacheRule.IgnoreNoCacheHeader = val
		}
		if val, ok := cacheRuleMap["no_cache_control_op"].(string); ok && val != "" {
			conf.CacheRule.NoCacheControlOp = val
		}
		if val, ok := cacheRuleMap["action"].(string); ok && val != "" {
			conf.CacheRule.Action = val
		}
	}

	// Build browser_cache_rule
	if browserCacheList, ok := confMap["browser_cache_rule"].([]interface{}); ok && len(browserCacheList) > 0 {
		browserCacheMap := browserCacheList[0].(map[string]interface{})
		conf.BrowserCacheRule = &scdn.BrowserCacheRule{
			CacheTime:       browserCacheMap["cachetime"].(int),
			IgnoreCacheTime: browserCacheMap["ignore_cache_time"].(bool),
			NoCache:         browserCacheMap["nocache"].(bool),
		}
	}

	// Build cache_errstatus
	// cache_errstatus is an array that can contain multiple items
	// Empty array means "no configuration", which is valid
	if errStatusList, ok := confMap["cache_errstatus"].([]interface{}); ok {
		log.Printf("[DEBUG] Building cache_errstatus from schema: length=%d", len(errStatusList))
		if len(errStatusList) > 0 {
			conf.CacheErrStatus = make([]scdn.CacheErrStatus, 0, len(errStatusList))
			for i, errStatusItem := range errStatusList {
				errStatusMap := errStatusItem.(map[string]interface{})
				errStatus := scdn.CacheErrStatus{
					CacheTime: errStatusMap["cachetime"].(int),
				}
				if errStatusList, ok := errStatusMap["err_status"].([]interface{}); ok {
					errStatus.ErrStatus = make([]int, 0, len(errStatusList))
					for _, v := range errStatusList {
						if v != nil {
							if intVal, ok := v.(int); ok {
								errStatus.ErrStatus = append(errStatus.ErrStatus, intVal)
							}
						}
					}
				}
				conf.CacheErrStatus = append(conf.CacheErrStatus, errStatus)
				log.Printf("[DEBUG] Added cache_errstatus[%d]: cachetime=%d, err_status=%v", i, errStatus.CacheTime, errStatus.ErrStatus)
			}
		} else {
			// Empty array - explicitly set to empty slice
			conf.CacheErrStatus = []scdn.CacheErrStatus{}
			log.Printf("[DEBUG] cache_errstatus is empty array, setting to empty slice")
		}
	} else {
		// Not present or not a list - set to empty slice
		conf.CacheErrStatus = []scdn.CacheErrStatus{}
		log.Printf("[DEBUG] cache_errstatus not found or invalid type, setting to empty slice")
	}

	// Build cache_url_rewrite
	if urlRewriteList, ok := confMap["cache_url_rewrite"].([]interface{}); ok && len(urlRewriteList) > 0 {
		urlRewriteMap := urlRewriteList[0].(map[string]interface{})
		conf.CacheURLRewrite = &scdn.CacheURLRewrite{
			SortArgs:   urlRewriteMap["sort_args"].(bool),
			IgnoreCase: urlRewriteMap["ignore_case"].(bool),
		}

		// Build queries
		if queriesList, ok := urlRewriteMap["queries"].([]interface{}); ok && len(queriesList) > 0 {
			queriesMap := queriesList[0].(map[string]interface{})
			conf.CacheURLRewrite.Queries = &scdn.CacheURLRewriteQueries{
				ArgsMethod: queriesMap["args_method"].(string),
			}
			if itemsList, ok := queriesMap["items"].([]interface{}); ok {
				conf.CacheURLRewrite.Queries.Items = make([]string, 0, len(itemsList))
				for _, v := range itemsList {
					if v != nil {
						if strVal, ok := v.(string); ok {
							conf.CacheURLRewrite.Queries.Items = append(conf.CacheURLRewrite.Queries.Items, strVal)
						}
					}
				}
			}
		}

		// Build cookies
		if cookiesList, ok := urlRewriteMap["cookies"].([]interface{}); ok && len(cookiesList) > 0 {
			cookiesMap := cookiesList[0].(map[string]interface{})
			conf.CacheURLRewrite.Cookies = &scdn.CacheURLRewriteCookies{
				ArgsMethod: cookiesMap["args_method"].(string),
			}
			if itemsList, ok := cookiesMap["items"].([]interface{}); ok {
				conf.CacheURLRewrite.Cookies.Items = make([]string, 0, len(itemsList))
				for _, v := range itemsList {
					if v != nil {
						if strVal, ok := v.(string); ok {
							conf.CacheURLRewrite.Cookies.Items = append(conf.CacheURLRewrite.Cookies.Items, strVal)
						}
					}
				}
			}
		}
	}

	// Build cache_share (optional)
	if cacheShareList, ok := confMap["cache_share"].([]interface{}); ok && len(cacheShareList) > 0 {
		cacheShareMap := cacheShareList[0].(map[string]interface{})
		conf.CacheShare = &scdn.CacheShare{
			Scheme: cacheShareMap["scheme"].(string),
		}
	} else {
		log.Printf("[DEBUG] cache_share not provided in configuration")
	}

	return conf, nil
}

func buildCacheRuleConfToSchema(conf *scdn.CacheRuleConf) map[string]interface{} {
	// Initialize confMap with all optional fields as empty arrays
	// This ensures schema consistency: even if API returns null, fields are present
	confMap := map[string]interface{}{
		"nocache":            conf.NoCache,
		"cache_rule":         []interface{}{},
		"browser_cache_rule": []interface{}{},
		"cache_errstatus":    []interface{}{},
		"cache_url_rewrite":  []interface{}{},
		"cache_share":        []interface{}{},
	}

	// Override with actual values if API returned them
	if conf.CacheRule != nil {
		cacheRuleMap := map[string]interface{}{
			"cachetime":             conf.CacheRule.CacheTime,
			"ignore_cache_time":     conf.CacheRule.IgnoreCacheTime,
			"ignore_nocache_header": conf.CacheRule.IgnoreNoCacheHeader,
			"no_cache_control_op":   conf.CacheRule.NoCacheControlOp,
			"action":                conf.CacheRule.Action,
		}
		confMap["cache_rule"] = []interface{}{cacheRuleMap}
	}

	if conf.BrowserCacheRule != nil {
		browserCacheMap := map[string]interface{}{
			"cachetime":         conf.BrowserCacheRule.CacheTime,
			"ignore_cache_time": conf.BrowserCacheRule.IgnoreCacheTime,
			"nocache":           conf.BrowserCacheRule.NoCache,
		}
		confMap["browser_cache_rule"] = []interface{}{browserCacheMap}
	}

	if len(conf.CacheErrStatus) > 0 {
		errStatusList := make([]map[string]interface{}, 0, len(conf.CacheErrStatus))
		for _, errStatus := range conf.CacheErrStatus {
			errStatusMap := map[string]interface{}{
				"cachetime":  errStatus.CacheTime,
				"err_status": errStatus.ErrStatus,
			}
			errStatusList = append(errStatusList, errStatusMap)
		}
		confMap["cache_errstatus"] = errStatusList
	}

	if conf.CacheURLRewrite != nil {
		urlRewriteMap := map[string]interface{}{
			"sort_args":   conf.CacheURLRewrite.SortArgs,
			"ignore_case": conf.CacheURLRewrite.IgnoreCase,
		}

		if conf.CacheURLRewrite.Queries != nil {
			queriesMap := map[string]interface{}{
				"args_method": conf.CacheURLRewrite.Queries.ArgsMethod,
				"items":       conf.CacheURLRewrite.Queries.Items,
			}
			urlRewriteMap["queries"] = []interface{}{queriesMap}
		}

		if conf.CacheURLRewrite.Cookies != nil {
			cookiesMap := map[string]interface{}{
				"args_method": conf.CacheURLRewrite.Cookies.ArgsMethod,
				"items":       conf.CacheURLRewrite.Cookies.Items,
			}
			urlRewriteMap["cookies"] = []interface{}{cookiesMap}
		}

		confMap["cache_url_rewrite"] = []interface{}{urlRewriteMap}
	}

	// cache_share is optional, so it might be nil
	if conf.CacheShare != nil {
		cacheShareMap := map[string]interface{}{
			"scheme": conf.CacheShare.Scheme,
		}
		confMap["cache_share"] = []interface{}{cacheShareMap}
	}

	return confMap
}
