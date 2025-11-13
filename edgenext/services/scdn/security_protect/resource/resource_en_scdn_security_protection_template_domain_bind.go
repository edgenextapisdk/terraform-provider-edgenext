package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnSecurityProtectionTemplateDomainBind returns the SCDN security protection template domain bind resource
func ResourceEdgenextScdnSecurityProtectionTemplateDomainBind() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnSecurityProtectionTemplateDomainBindCreate,
		Read:   resourceScdnSecurityProtectionTemplateDomainBindRead,
		Delete: resourceScdnSecurityProtectionTemplateDomainBindDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Business ID (template ID)",
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Domain ID list",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"bind_business_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Bind business ID list",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Group ID list",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
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

func resourceScdnSecurityProtectionTemplateDomainBindCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.SecurityProtectionTemplateBindDomainRequest{
		BusinessID: businessID,
	}

	if domainIDs, ok := d.GetOk("domain_ids"); ok {
		domainIDsList := domainIDs.([]interface{})
		req.DomainIDs = make([]int, len(domainIDsList))
		for i, v := range domainIDsList {
			req.DomainIDs[i] = v.(int)
		}
	}

	if bindBusinessIDs, ok := d.GetOk("bind_business_ids"); ok {
		bindBusinessIDsList := bindBusinessIDs.([]interface{})
		req.BindBusinessIDs = make([]int, len(bindBusinessIDsList))
		for i, v := range bindBusinessIDsList {
			req.BindBusinessIDs[i] = v.(int)
		}
	}

	if groupIDs, ok := d.GetOk("group_ids"); ok {
		groupIDsList := groupIDs.([]interface{})
		req.GroupIDs = make([]int, len(groupIDsList))
		for i, v := range groupIDsList {
			req.GroupIDs[i] = v.(int)
		}
	}

	log.Printf("[INFO] Binding SCDN security protection template domain: business_id=%d", businessID)
	response, err := service.BindSecurityProtectionTemplateDomain(req)
	if err != nil {
		return fmt.Errorf("failed to bind template domain: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("template-domain-bind-%d", businessID))

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

	return resourceScdnSecurityProtectionTemplateDomainBindRead(d, m)
}

func resourceScdnSecurityProtectionTemplateDomainBindRead(d *schema.ResourceData, m interface{}) error {
	// Parse business_id from resource ID if not set in state
	businessID := d.Get("business_id").(int)
	if businessID == 0 {
		// Try to parse from resource ID (format: template-domain-bind-{business_id})
		id := d.Id()
		if id != "" {
			var parsedID int
			_, err := fmt.Sscanf(id, "template-domain-bind-%d", &parsedID)
			if err == nil {
				businessID = parsedID
				if err := d.Set("business_id", businessID); err != nil {
					log.Printf("[WARN] Failed to set business_id from ID: %v", err)
				}
			} else {
				// If we can't parse the ID and business_id is not set, mark as removed
				d.SetId("")
				return nil
			}
		} else {
			d.SetId("")
			return nil
		}
	}

	// Binding is a one-time operation, we can't really "read" the binding state from API
	// Just verify the resource exists by ensuring business_id is set
	return nil
}

func resourceScdnSecurityProtectionTemplateDomainBindDelete(d *schema.ResourceData, m interface{}) error {
	// Unbinding is not directly supported by the API
	// The binding can be removed by updating the template or domain configuration
	// For now, we just remove the resource from state
	log.Printf("[INFO] Unbinding SCDN security protection template domain: business_id=%d", d.Get("business_id").(int))
	d.SetId("")
	return nil
}
