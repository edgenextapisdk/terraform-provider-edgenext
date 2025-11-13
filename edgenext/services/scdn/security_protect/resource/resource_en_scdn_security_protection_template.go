package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnSecurityProtectionTemplate returns the SCDN security protection template resource
func ResourceEdgenextScdnSecurityProtectionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnSecurityProtectionTemplateCreate,
		Read:   resourceScdnSecurityProtectionTemplateRead,
		Update: resourceScdnSecurityProtectionTemplateUpdate,
		Delete: resourceScdnSecurityProtectionTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Business ID (template ID). Required for update/delete, computed for create.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template name",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template remark",
			},
			"template_source_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Source template ID",
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domain ID list",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Group ID list",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domain list",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"bind_all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Bind all domains",
			},
			// Computed fields
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template type: global, only_domain, more_domain",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time",
			},
			"bind_domain_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Bind domain count",
			},
			"fail_domains": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Failed domains",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceScdnSecurityProtectionTemplateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.SecurityProtectionTemplateCreateRequest{
		Name: d.Get("name").(string),
	}

	if remark, ok := d.GetOk("remark"); ok {
		req.Remark = remark.(string)
	}
	if templateSourceID, ok := d.GetOk("template_source_id"); ok {
		req.TemplateSourceID = templateSourceID.(int)
	}
	if domainIDs, ok := d.GetOk("domain_ids"); ok {
		domainIDsList := domainIDs.([]interface{})
		req.DomainIDs = make([]int, len(domainIDsList))
		for i, v := range domainIDsList {
			req.DomainIDs[i] = v.(int)
		}
	}
	if groupIDs, ok := d.GetOk("group_ids"); ok {
		groupIDsList := groupIDs.([]interface{})
		req.GroupIDs = make([]int, len(groupIDsList))
		for i, v := range groupIDsList {
			req.GroupIDs[i] = v.(int)
		}
	}
	if domains, ok := d.GetOk("domains"); ok {
		domainsList := domains.([]interface{})
		req.Domains = make([]string, len(domainsList))
		for i, v := range domainsList {
			req.Domains[i] = v.(string)
		}
	}
	if bindAll, ok := d.GetOk("bind_all"); ok {
		req.BindAll = bindAll.(bool)
	}

	log.Printf("[INFO] Creating SCDN security protection template: name=%s", req.Name)
	response, err := service.CreateSecurityProtectionTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to create security protection template: %w", err)
	}

	businessID := response.Data.BusinessID
	d.SetId(fmt.Sprintf("template-%d", businessID))
	if err := d.Set("business_id", businessID); err != nil {
		return fmt.Errorf("error setting business_id: %w", err)
	}

	// Set fail_domains if any
	if len(response.Data.FailDomains) > 0 {
		failDomainsMap := make(map[string]interface{})
		for k, v := range response.Data.FailDomains {
			failDomainsMap[k] = v
		}
		if err := d.Set("fail_domains", failDomainsMap); err != nil {
			log.Printf("[WARN] Failed to set fail_domains: %v", err)
		}
	}

	return resourceScdnSecurityProtectionTemplateRead(d, m)
}

func resourceScdnSecurityProtectionTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	if businessID == 0 {
		// Try to parse from resource ID
		return fmt.Errorf("business_id is required for reading template")
	}

	// Search for template by ID
	req := scdn.SecurityProtectionTemplateSearchRequest{
		TplType:  "global", // Default, can be adjusted
		Page:     1,
		PageSize: 100,
	}

	response, err := service.SearchSecurityProtectionTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to search templates: %w", err)
	}

	// Find the template by business_id
	var foundTemplate *scdn.SecurityProtectionTemplateInfo
	for _, template := range response.Data.Templates {
		if template.ID == businessID {
			foundTemplate = &template
			break
		}
	}

	if foundTemplate == nil {
		// Try to get global template
		globalResp, err := service.GetMemberGlobalTemplate()
		if err != nil {
			return fmt.Errorf("failed to get global template: %w", err)
		}
		if globalResp.Data.Template != nil && globalResp.Data.Template.ID == businessID {
			foundTemplate = globalResp.Data.Template
			if err := d.Set("bind_domain_count", globalResp.Data.BindDomainCount); err != nil {
				log.Printf("[WARN] Failed to set bind_domain_count: %v", err)
			}
		}
	}

	if foundTemplate == nil {
		d.SetId("")
		return fmt.Errorf("template not found: business_id=%d", businessID)
	}

	// Set fields
	if err := d.Set("name", foundTemplate.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("remark", foundTemplate.Remark); err != nil {
		log.Printf("[WARN] Failed to set remark: %v", err)
	}
	if err := d.Set("type", foundTemplate.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("created_at", foundTemplate.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set created_at: %v", err)
	}

	return nil
}

func resourceScdnSecurityProtectionTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	if businessID == 0 {
		return fmt.Errorf("business_id is required for updating template")
	}

	req := scdn.SecurityProtectionTemplateEditRequest{
		BusinessID: businessID,
		Name:       d.Get("name").(string),
	}

	if remark, ok := d.GetOk("remark"); ok {
		req.Remark = remark.(string)
	}

	log.Printf("[INFO] Updating SCDN security protection template: business_id=%d", businessID)
	_, err := service.EditSecurityProtectionTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to update security protection template: %w", err)
	}

	return resourceScdnSecurityProtectionTemplateRead(d, m)
}

func resourceScdnSecurityProtectionTemplateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	if businessID == 0 {
		return fmt.Errorf("business_id is required for deleting template")
	}

	req := scdn.SecurityProtectionTemplateDeleteRequest{
		BusinessID: businessID,
	}

	log.Printf("[INFO] Deleting SCDN security protection template: business_id=%d", businessID)
	_, err := service.DeleteSecurityProtectionTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to delete security protection template: %w", err)
	}

	d.SetId("")
	return nil
}
