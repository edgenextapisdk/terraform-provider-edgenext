package domain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCertBinding returns the SCDN certificate binding resource
func ResourceEdgenextScdnCertBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCertBindingCreate,
		Read:   resourceScdnCertBindingRead,
		Delete: resourceScdnCertBindingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain to bind the certificate to",
			},
			"ca_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the certificate to bind",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this certificate binding",
			},
		},
	}
}

func resourceScdnCertBindingCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)
	caID := d.Get("ca_id").(int)

	req := scdn.DomainCertBindRequest{
		DomainID: domainID,
		CAID:     caID,
	}

	log.Printf("[INFO] Creating SCDN certificate binding: %+v", req)
	response, err := service.BindDomainCert(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN certificate binding: %w", err)
	}

	// Create a unique ID for this binding
	bindingID := fmt.Sprintf("%d-%d", domainID, response.Data.CAID)
	d.SetId(bindingID)

	log.Printf("[INFO] SCDN certificate binding created successfully: %s", d.Id())
	return resourceScdnCertBindingRead(d, m)
}

func resourceScdnCertBindingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Parse the binding ID to get domain ID and certificate ID
	bindingID := d.Id()
	domainID, caID, err := parseCertBindingID(bindingID)
	if err != nil {
		return fmt.Errorf("invalid binding ID: %w", err)
	}

	// Verify the binding exists by checking the domain's certificate status
	req := scdn.DomainListRequest{
		Page:     1,
		PageSize: 100,
	}

	response, err := service.ListDomains(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN certificate binding: %w", err)
	}

	var domainInfo *scdn.DomainInfo
	for _, domain := range response.Data.List {
		if domain.ID == domainID {
			domainInfo = &domain
			break
		}
	}

	if domainInfo == nil {
		log.Printf("[WARN] SCDN domain not found: %d", domainID)
		d.SetId("")
		return nil
	}

	// Check if the certificate is bound to this domain
	if domainInfo.CAID != caID {
		log.Printf("[WARN] Certificate %d is not bound to domain %d", caID, domainID)
		d.SetId("")
		return nil
	}

	// Set basic fields
	if err := d.Set("domain_id", domainID); err != nil {
		return fmt.Errorf("error setting domain_id: %w", err)
	}
	if err := d.Set("ca_id", caID); err != nil {
		return fmt.Errorf("error setting ca_id: %w", err)
	}

	log.Printf("[INFO] SCDN certificate binding read successfully: %s", d.Id())
	return nil
}

func resourceScdnCertBindingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Parse the binding ID to get domain ID and certificate ID
	bindingID := d.Id()
	domainID, caID, err := parseCertBindingID(bindingID)
	if err != nil {
		return fmt.Errorf("invalid binding ID: %w", err)
	}

	req := scdn.DomainCertUnbindRequest{
		DomainID: domainID,
		CAID:     caID,
	}

	log.Printf("[INFO] Deleting SCDN certificate binding: %+v", req)
	_, err = service.UnbindDomainCert(req)
	if err != nil {
		return fmt.Errorf("failed to delete SCDN certificate binding: %w", err)
	}

	d.SetId("")
	log.Printf("[INFO] SCDN certificate binding deleted successfully: %s", bindingID)
	return nil
}

// parseCertBindingID parses a certificate binding ID in the format "domainID-certID"
func parseCertBindingID(bindingID string) (domainID, caID int, err error) {
	// The binding ID format is "domainID-certID"
	parts := make([]string, 0, 2)
	for i := 0; i < len(bindingID); i++ {
		if bindingID[i] == '-' {
			parts = append(parts, bindingID[:i], bindingID[i+1:])
			break
		}
	}

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid binding ID format: %s", bindingID)
	}

	domainID, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid domain ID: %w", err)
	}

	caID, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid certificate ID: %w", err)
	}

	return domainID, caID, nil
}
