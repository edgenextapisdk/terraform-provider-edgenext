package resource

import (
	"fmt"
	"log"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnSecurityProtectionDomainTemplate returns the SCDN security protection domain template resource
// This resource creates a domain-level security protection template.
// Unlike the multi-domain template, this resource:
// - Only supports creation, not update
// - Deletion rebinds the domain to the global template
// - Requires template_source_id (obtained from global template data source)
func ResourceEdgenextScdnSecurityProtectionDomainTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnSecurityProtectionDomainTemplateCreate,
		Read:   resourceScdnSecurityProtectionDomainTemplateRead,
		Delete: resourceScdnSecurityProtectionDomainTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Domain ID to create template for.",
			},
			"template_source_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Source template ID. Use data source edgenext_scdn_security_protection_member_global_template to get the global template ID.",
			},
			// Computed fields
			"business_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Created business ID (template ID) for the domain.",
			},
		},
	}
}

func resourceScdnSecurityProtectionDomainTemplateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request - single domain only
	domainID := d.Get("domain_id").(int)
	req := scdn.SecurityProtectionTemplateCreateDomainRequest{
		DomainIDs:        []int{domainID},
		TemplateSourceID: d.Get("template_source_id").(int),
	}

	log.Printf("[INFO] Creating SCDN security protection domain template: domain_id=%d, template_source_id=%d", domainID, req.TemplateSourceID)
	response, err := service.CreateSecurityProtectionTemplateDomain(req)
	if err != nil {
		return fmt.Errorf("failed to create domain template: %w", err)
	}

	// Set resource ID based on domain ID
	d.SetId(fmt.Sprintf("domain-template-%d", domainID))

	// Report fail_domains if any
	if len(response.Data.FailDomains) > 0 {
		var failMsg []string
		for k, v := range response.Data.FailDomains {
			failMsg = append(failMsg, fmt.Sprintf("%s: %s", k, v))
		}
		return fmt.Errorf("domain template creation failed: %s", strings.Join(failMsg, "; "))
	}

	// Set business_id (first one since we only create for single domain)
	var businessID int
	if len(response.Data.BusinessIDs) > 0 {
		businessID = response.Data.BusinessIDs[0]
	} else {
		// If business_ids is empty, the domain might already have a template
		// Search for the template by domain_id
		log.Printf("[INFO] business_ids is empty, searching for existing template for domain %d", domainID)
		searchReq := scdn.SecurityProtectionTemplateSearchRequest{
			TplType:  "only_domain",
			Page:     1,
			PageSize: 100,
		}
		searchResp, err := service.SearchSecurityProtectionTemplates(searchReq)
		if err != nil {
			log.Printf("[WARN] Failed to search templates: %v", err)
		} else {
			// Find template with matching domain_id
			for _, tpl := range searchResp.Data.Templates {
				if tpl.DomainID == domainID {
					businessID = tpl.ID
					log.Printf("[INFO] Found existing template %d for domain %d", businessID, domainID)
					break
				}
			}
		}
	}

	if businessID > 0 {
		if err := d.Set("business_id", businessID); err != nil {
			log.Printf("[WARN] Failed to set business_id: %v", err)
		}
	} else {
		log.Printf("[WARN] Could not determine business_id for domain %d", domainID)
	}

	return resourceScdnSecurityProtectionDomainTemplateRead(d, m)
}

func resourceScdnSecurityProtectionDomainTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	log.Printf("[INFO] Reading SCDN security protection domain template: %s", d.Id())

	// Parse domain_id from resource ID (format: domain-template-{domain_id})
	resourceID := d.Id()
	if !strings.HasPrefix(resourceID, "domain-template-") {
		log.Printf("[WARN] Invalid resource ID format: %s, removing from state", resourceID)
		d.SetId("")
		return nil
	}

	domainIDStr := strings.TrimPrefix(resourceID, "domain-template-")
	var domainID int
	if _, err := fmt.Sscanf(domainIDStr, "%d", &domainID); err != nil {
		log.Printf("[WARN] Failed to parse domain_id from resource ID: %s, error: %v", resourceID, err)
		d.SetId("")
		return nil
	}

	// Search for domain templates (only_domain type)
	searchReq := scdn.SecurityProtectionTemplateSearchRequest{
		TplType:  "only_domain",
		Page:     1,
		PageSize: 100,
	}
	searchResp, err := service.SearchSecurityProtectionTemplates(searchReq)
	if err != nil {
		return fmt.Errorf("failed to search domain templates: %w", err)
	}

	// Find template with matching domain_id
	var foundTemplate *scdn.SecurityProtectionTemplateInfo
	for i := range searchResp.Data.Templates {
		if searchResp.Data.Templates[i].DomainID == domainID {
			foundTemplate = &searchResp.Data.Templates[i]
			break
		}
	}

	if foundTemplate == nil {
		log.Printf("[WARN] Domain template not found for domain_id %d, removing from state", domainID)
		d.SetId("")
		return nil
	}

	// Set attributes
	if err := d.Set("domain_id", domainID); err != nil {
		return fmt.Errorf("failed to set domain_id: %w", err)
	}

	if err := d.Set("business_id", foundTemplate.ID); err != nil {
		return fmt.Errorf("failed to set business_id: %w", err)
	}

	// Note: template_source_id cannot be retrieved from API
	// It's only used during creation. For imported resources,
	// this field will be empty in state but that's acceptable
	// since the resource doesn't support updates.
	log.Printf("[INFO] Found domain template: domain_id=%d, business_id=%d", domainID, foundTemplate.ID)

	return nil
}

func resourceScdnSecurityProtectionDomainTemplateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Get business ID from state
	businessID := d.Get("business_id").(int)
	if businessID == 0 {
		log.Printf("[WARN] No business ID found, removing from state only")
		d.SetId("")
		return nil
	}

	// Get global template ID for rebinding
	globalTemplateResp, err := service.GetMemberGlobalTemplate()
	if err != nil {
		return fmt.Errorf("failed to get global template for rebinding: %w", err)
	}

	if globalTemplateResp.Data.Template == nil {
		return fmt.Errorf("global template not found, cannot rebind domains")
	}

	globalTemplateID := globalTemplateResp.Data.Template.ID

	// Rebind domain to global template by calling bind API
	// This effectively "deletes" the domain-specific template
	req := scdn.SecurityProtectionTemplateBindDomainRequest{
		BusinessID:      globalTemplateID,
		BindBusinessIDs: []int{businessID},
	}

	log.Printf("[INFO] Rebinding domain template %d to global template %d", businessID, globalTemplateID)
	_, err = service.BindSecurityProtectionTemplateDomain(req)
	if err != nil {
		log.Printf("[WARN] Failed to rebind domain template %d to global template: %v", businessID, err)
	}

	d.SetId("")
	return nil
}
