package template

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnRuleTemplate returns the SCDN rule template resource
func ResourceEdgenextScdnRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnRuleTemplateCreate,
		Read:   resourceScdnRuleTemplateRead,
		Update: resourceScdnRuleTemplateUpdate,
		Delete: resourceScdnRuleTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The template ID for updating an existing template. If provided, this will update the template instead of creating a new one.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule template name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The rule template description",
			},
			"app_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application type (e.g., 'network_speed')",
			},
			"tpl_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The template type",
			},
			"from_tpl_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Existing template ID to copy from",
			},
			"bind_domain": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Domain binding information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_domain": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "If true, bind to all domains",
						},
						"domain_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of domain IDs to bind",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"domain_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of domain group IDs to bind",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"domains": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of domain names to bind",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_bind": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to bind domains",
						},
					},
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the rule template",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template creation timestamp",
			},
		},
	}
}

func resourceScdnRuleTemplateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Check if template_id is provided (for updating existing template)
	if templateIDStr, hasTemplateID := d.GetOk("template_id"); hasTemplateID {
		// Update existing template
		templateID, err := strconv.Atoi(templateIDStr.(string))
		if err != nil {
			return fmt.Errorf("invalid template_id: %w", err)
		}

		req := scdn.RuleTemplateUpdateRequest{
			ID:          templateID,
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}

		log.Printf("[INFO] Updating SCDN rule template: %+v", req)
		response, err := service.UpdateRuleTemplate(req)
		if err != nil {
			return fmt.Errorf("failed to update rule template: %w", err)
		}

		// Verify the updated template ID matches what we expected
		if response.Data.ID != templateID {
			return fmt.Errorf("template ID mismatch after update: expected %d, but got %d", templateID, response.Data.ID)
		}

		// Set the ID to match the updated template ID
		d.SetId(templateIDStr.(string))

		return resourceScdnRuleTemplateRead(d, m)
	}

	// Build create request for new template
	req := scdn.RuleTemplateCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		AppType:     d.Get("app_type").(string),
		TplType:     d.Get("tpl_type").(string),
	}

	if fromTplID, ok := d.GetOk("from_tpl_id"); ok {
		req.FromTplID = fromTplID.(int)
	}

	// Handle bind_domain
	if bindDomainList, ok := d.GetOk("bind_domain"); ok && len(bindDomainList.([]interface{})) > 0 {
		bindDomainMap := bindDomainList.([]interface{})[0].(map[string]interface{})
		bindDomain := &scdn.RuleTemplateBindDomain{
			AllDomain: bindDomainMap["all_domain"].(bool),
			IsBind:    bindDomainMap["is_bind"].(bool),
		}

		if domainIDs, ok := bindDomainMap["domain_ids"].([]interface{}); ok && len(domainIDs) > 0 {
			bindDomain.DomainIDs = make([]int, len(domainIDs))
			for i, v := range domainIDs {
				bindDomain.DomainIDs[i] = v.(int)
			}
		}

		if domainGroupIDs, ok := bindDomainMap["domain_group_ids"].([]interface{}); ok && len(domainGroupIDs) > 0 {
			bindDomain.DomainGroupIDs = make([]int, len(domainGroupIDs))
			for i, v := range domainGroupIDs {
				bindDomain.DomainGroupIDs[i] = v.(int)
			}
		}

		if domains, ok := bindDomainMap["domains"].([]interface{}); ok && len(domains) > 0 {
			bindDomain.Domains = make([]string, len(domains))
			for i, v := range domains {
				bindDomain.Domains[i] = v.(string)
			}
		}

		req.BindDomain = bindDomain
	}

	log.Printf("[INFO] Creating SCDN rule template: %s", req.Name)
	response, err := service.CreateRuleTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN rule template: %w", err)
	}

	log.Printf("[DEBUG] Rule template creation response: %+v", response)

	// Set the template ID as the resource ID
	templateIDStr := strconv.Itoa(response.Data.ID)
	d.SetId(templateIDStr)

	log.Printf("[INFO] SCDN rule template created successfully: %s", d.Id())

	// Call read to get full details
	return resourceScdnRuleTemplateRead(d, m)
}

func resourceScdnRuleTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	// Get app_type from state or schema
	appType := d.Get("app_type").(string)
	if appType == "" {
		// Try to get from state
		appType = "network_speed" // default, but should be set
	}

	// List templates with filter to find this specific one
	req := scdn.RuleTemplateListRequest{
		Page:     1,
		PageSize: 1000,
		AppType:  appType,
	}

	log.Printf("[INFO] Reading SCDN rule template: %d", templateID)
	response, err := service.ListRuleTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to list rule templates: %w", err)
	}

	// Find the template with matching ID
	var template *scdn.RuleTemplateInfo
	for _, tpl := range response.Data.List {
		if tpl.ID == templateID {
			template = &tpl
			break
		}
	}

	if template == nil {
		log.Printf("[WARN] Rule template not found: %d", templateID)
		d.SetId("")
		return nil
	}

	// Verify the found template ID matches what we're looking for
	if template.ID != templateID {
		return fmt.Errorf("template ID mismatch: expected %d, but found %d", templateID, template.ID)
	}

	// Set computed fields - use the ID from d.Id() to ensure consistency
	// The id field should match the resource ID (d.Id())
	if err := d.Set("id", d.Id()); err != nil {
		log.Printf("[WARN] Failed to set template id: %v", err)
	}
	if err := d.Set("name", template.Name); err != nil {
		log.Printf("[WARN] Failed to set template name: %v", err)
	}
	if err := d.Set("description", template.Description); err != nil {
		log.Printf("[WARN] Failed to set template description: %v", err)
	}
	if err := d.Set("app_type", template.AppType); err != nil {
		log.Printf("[WARN] Failed to set template app_type: %v", err)
	}
	if err := d.Set("created_at", template.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set template created_at: %v", err)
	}

	// Set bind_domain if available
	if len(template.BindDomains) > 0 {
		bindDomainMap := map[string]interface{}{
			"is_bind": true,
		}

		domainIDs := make([]int, 0, len(template.BindDomains))
		for _, bd := range template.BindDomains {
			domainIDs = append(domainIDs, bd.DomainID)
		}
		if len(domainIDs) > 0 {
			bindDomainMap["domain_ids"] = domainIDs
		}

		bindDomainList := []interface{}{bindDomainMap}
		if err := d.Set("bind_domain", bindDomainList); err != nil {
			log.Printf("[WARN] Failed to set template bind_domain: %v", err)
		}
	}

	log.Printf("[INFO] SCDN rule template read successfully: %s", d.Id())
	return nil
}

func resourceScdnRuleTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	// Build update request
	req := scdn.RuleTemplateUpdateRequest{
		ID: templateID,
	}

	if d.HasChange("name") {
		req.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		req.Description = d.Get("description").(string)
	}

	log.Printf("[INFO] Updating SCDN rule template: %d", templateID)
	_, err = service.UpdateRuleTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to update SCDN rule template: %w", err)
	}

	log.Printf("[INFO] SCDN rule template updated successfully: %s", d.Id())

	// Call read to refresh state
	return resourceScdnRuleTemplateRead(d, m)
}

func resourceScdnRuleTemplateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	req := scdn.RuleTemplateDeleteRequest{
		ID: templateID,
	}

	log.Printf("[INFO] Deleting SCDN rule template: %d", templateID)
	_, err = service.DeleteRuleTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to delete SCDN rule template: %w", err)
	}

	log.Printf("[INFO] SCDN rule template deleted successfully: %s", d.Id())
	d.SetId("")
	return nil
}
