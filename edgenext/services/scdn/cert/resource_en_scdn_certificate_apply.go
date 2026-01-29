package cert

import (
	"fmt"
	"log"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCertificateApply returns the SCDN certificate apply resource
func ResourceEdgenextScdnCertificateApply() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCertificateApplyCreate,
		Read:   resourceScdnCertificateApplyRead,
		Delete: resourceScdnCertificateApplyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "The list of domains to apply for certificate",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the certificate application",
			},
			"ca_id_domains": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The mapping of domain_id to domain",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ca_id_names": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The mapping of ca_id to ca_name",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceScdnCertificateApplyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Get domains from schema
	domainsInterface := d.Get("domain").([]interface{})
	domains := make([]string, len(domainsInterface))
	for i, domain := range domainsInterface {
		domains[i] = domain.(string)
	}

	// Build create request
	req := scdn.CAApplyAddRequest{
		Domain: domains,
	}

	log.Printf("[INFO] Applying for SCDN certificate for domains: %v", domains)
	response, err := service.ApplyCertificate(req)
	if err != nil {
		return fmt.Errorf("failed to apply for SCDN certificate: %w", err)
	}

	log.Printf("[DEBUG] Certificate application response: %+v", response)

	// Convert maps for Terraform
	caIDDomains := make(map[string]interface{})
	caIds := make([]string, 0)
	for k, v := range response.Data.CAIDDomains {
		caIDDomains[k] = v
		caIds = append(caIds, k)
	}

	// Generate resource ID from domain list
	resourceID := strings.Join(caIds, ",")
	d.SetId(resourceID)

	caIDNames := make(map[string]interface{})
	for k, v := range response.Data.CAIDNames {
		caIDNames[k] = v
	}

	// Set computed fields
	if err := d.Set("ca_id_domains", caIDDomains); err != nil {
		log.Printf("[WARN] Failed to set ca_id_domains: %v", err)
	}
	if err := d.Set("ca_id_names", caIDNames); err != nil {
		log.Printf("[WARN] Failed to set ca_id_names: %v", err)
	}

	log.Printf("[INFO] SCDN certificate application created successfully: %s", d.Id())
	return nil
}

func resourceScdnCertificateApplyRead(d *schema.ResourceData, m interface{}) error {
	// Certificate application is a one-time operation
	// The read operation just returns the current state
	// In practice, you might want to query the certificate status
	log.Printf("[DEBUG] Reading SCDN certificate application: %s", d.Id())
	return nil
}

func resourceScdnCertificateApplyDelete(d *schema.ResourceData, m interface{}) error {
	// Certificate application cannot be deleted via API
	// This is a no-op, the resource will just be removed from state
	log.Printf("[INFO] Deleting SCDN certificate application from state: %s", d.Id())
	return nil
}
