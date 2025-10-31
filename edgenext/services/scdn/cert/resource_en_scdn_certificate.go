package cert

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCertificate returns the SCDN certificate resource
func ResourceEdgenextScdnCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCertificateCreate,
		Read:   resourceScdnCertificateRead,
		Update: resourceScdnCertificateUpdate,
		Delete: resourceScdnCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The certificate ID for updating an existing certificate. If provided, this will update the certificate instead of creating a new one.",
			},
			"ca_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The certificate name",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The certificate public key (PEM format). Required for creation, optional for updates.",
			},
			"ca_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The certificate private key (PEM format). Required for creation, optional for updates.",
			},
			"product_flag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The product flag",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the certificate",
			},
			"ca_sn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate serial number",
			},
			"member_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The member ID",
			},
			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate issuer",
			},
			"issuer_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate start time",
			},
			"issuer_expiry_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate expiry time",
			},
			"issuer_expiry_time_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate expiry time description",
			},
			"renew_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The renewal status: 1-default, 2-renewing, 3-renewal failed, 4-renewal successful",
			},
			"binded": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the certificate is bound",
			},
			"ca_domain": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The domains in the certificate",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"apply_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The application status: 1-applying, 2-issued, 3-review failed, 4-uploaded",
			},
			"ca_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate type: 1-upload, 2-lets apply",
			},
			"ca_type_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate domain type: 1-single domain, 2-multiple domains, 3-wildcard domain",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update timestamp",
			},
		},
	}
}

func resourceScdnCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Check if certificate_id is provided (for updating existing certificate)
	certIDStr, hasCertID := d.GetOk("certificate_id")
	caCert, caCertOk := d.GetOk("ca_cert")
	caKey, caKeyOk := d.GetOk("ca_key")

	// If certificate_id is provided and no cert/key provided, this is a name-only update
	if hasCertID && !caCertOk && !caKeyOk {
		// Update existing certificate name only
		certID, err := strconv.Atoi(certIDStr.(string))
		if err != nil {
			return fmt.Errorf("invalid certificate_id: %w", err)
		}

		// Set the ID first so subsequent operations know this is an existing resource
		d.SetId(certIDStr.(string))

		req := scdn.CAEditNameRequest{
			ID:          certID,
			CAName:      d.Get("ca_name").(string),
			ProductFlag: d.Get("product_flag").(string),
		}

		log.Printf("[INFO] Updating SCDN certificate name: %+v", req)
		_, err = service.EditCertificateName(req)
		if err != nil {
			return fmt.Errorf("failed to update certificate name: %w", err)
		}

		return resourceScdnCertificateRead(d, m)
	}

	// Check if this is an imported resource (has ID in state) and only name update is needed
	if d.Id() != "" && !caCertOk && !caKeyOk {
		// This is an imported resource, only updating the name
		log.Printf("[INFO] Imported certificate detected, updating name only")
		certID, err := strconv.Atoi(d.Id())
		if err != nil {
			return fmt.Errorf("invalid certificate ID: %w", err)
		}

		req := scdn.CAEditNameRequest{
			ID:          certID,
			CAName:      d.Get("ca_name").(string),
			ProductFlag: d.Get("product_flag").(string),
		}

		log.Printf("[INFO] Updating SCDN certificate name: %+v", req)
		_, err = service.EditCertificateName(req)
		if err != nil {
			return fmt.Errorf("failed to update certificate name: %w", err)
		}

		return resourceScdnCertificateRead(d, m)
	}

	// Validate required fields for new certificate creation
	if !caCertOk || !caKeyOk {
		return fmt.Errorf("ca_cert and ca_key are required for certificate creation. To update an existing certificate, provide certificate_id instead of ca_cert and ca_key")
	}

	// Check for placeholder values
	caCertStr := caCert.(string)
	caKeyStr := caKey.(string)
	if caCertStr == "UPDATE_WITH_ACTUAL_CERTIFICATE_CONTENT" ||
		caKeyStr == "UPDATE_WITH_ACTUAL_PRIVATE_KEY_CONTENT" ||
		caCertStr == "" || caKeyStr == "" {
		return fmt.Errorf("ca_cert and ca_key must contain valid certificate content, not placeholder values")
	}

	// Build create request
	req := scdn.CATextSaveRequest{
		CAName:      d.Get("ca_name").(string),
		CACert:      caCertStr,
		CAKey:       caKeyStr,
		ProductFlag: d.Get("product_flag").(string),
	}

	log.Printf("[INFO] Creating SCDN certificate: %s", req.CAName)
	response, err := service.SaveCertificate(req)
	if err != nil {
		// Improve error message for certificate/key mismatch
		return fmt.Errorf("failed to create SCDN certificate: %w. Please ensure the certificate and private key match", err)
	}

	log.Printf("[DEBUG] Certificate creation response: %+v", response)

	// Set the certificate ID as the resource ID
	d.SetId(response.Data.ID)

	// Set basic fields directly from creation response
	if err := d.Set("id", response.Data.ID); err != nil {
		log.Printf("[WARN] Failed to set certificate id: %v", err)
	}
	if err := d.Set("ca_sn", response.Data.CASN); err != nil {
		log.Printf("[WARN] Failed to set certificate serial number: %v", err)
	}
	if err := d.Set("ca_name", req.CAName); err != nil {
		log.Printf("[WARN] Failed to set certificate name: %v", err)
	}

	log.Printf("[INFO] SCDN certificate created successfully: %s", d.Id())

	// Call read to get full details
	return resourceScdnCertificateRead(d, m)
}

func resourceScdnCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	certID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid certificate ID: %w", err)
	}

	req := scdn.CASelfDetailRequest{
		ID: certID,
	}

	log.Printf("[DEBUG] Reading SCDN certificate: %d", certID)
	response, err := service.GetCertificateDetail(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN certificate: %w", err)
	}

	if response.Data.ID == "" {
		log.Printf("[WARN] SCDN certificate not found: %d", certID)
		d.SetId("")
		return nil
	}

	// Set all fields
	if err := d.Set("id", response.Data.ID); err != nil {
		return fmt.Errorf("error setting id: %w", err)
	}
	if err := d.Set("ca_name", response.Data.CAName); err != nil {
		return fmt.Errorf("error setting ca_name: %w", err)
	}
	if err := d.Set("member_id", response.Data.MemberID); err != nil {
		return fmt.Errorf("error setting member_id: %w", err)
	}
	if err := d.Set("issuer", response.Data.Issuer); err != nil {
		return fmt.Errorf("error setting issuer: %w", err)
	}
	if err := d.Set("issuer_start_time", response.Data.IssuerStartTime); err != nil {
		return fmt.Errorf("error setting issuer_start_time: %w", err)
	}
	if err := d.Set("issuer_expiry_time", response.Data.IssuerExpiryTime); err != nil {
		return fmt.Errorf("error setting issuer_expiry_time: %w", err)
	}
	if err := d.Set("issuer_expiry_time_desc", response.Data.IssuerExpiryTimeDesc); err != nil {
		return fmt.Errorf("error setting issuer_expiry_time_desc: %w", err)
	}
	if err := d.Set("renew_status", response.Data.RenewStatus); err != nil {
		return fmt.Errorf("error setting renew_status: %w", err)
	}
	if err := d.Set("binded", response.Data.Binded); err != nil {
		return fmt.Errorf("error setting binded: %w", err)
	}
	if err := d.Set("ca_domain", response.Data.CADomain); err != nil {
		return fmt.Errorf("error setting ca_domain: %w", err)
	}
	if err := d.Set("apply_status", response.Data.ApplyStatus); err != nil {
		return fmt.Errorf("error setting apply_status: %w", err)
	}
	if err := d.Set("ca_type", response.Data.CAType); err != nil {
		return fmt.Errorf("error setting ca_type: %w", err)
	}
	if err := d.Set("ca_type_domain", response.Data.CATypeDomain); err != nil {
		return fmt.Errorf("error setting ca_type_domain: %w", err)
	}
	if err := d.Set("created_at", response.Data.CreatedAt); err != nil {
		return fmt.Errorf("error setting created_at: %w", err)
	}
	if err := d.Set("updated_at", response.Data.UpdatedAt); err != nil {
		return fmt.Errorf("error setting updated_at: %w", err)
	}

	// Note: ca_cert and ca_key are not returned by the API for security reasons
	// They are only set during creation/update

	log.Printf("[INFO] SCDN certificate read successfully: %s", d.Id())
	return nil
}

func resourceScdnCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	certID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid certificate ID: %w", err)
	}

	// Update certificate name if changed
	if d.HasChange("ca_name") {
		req := scdn.CAEditNameRequest{
			ID:          certID,
			CAName:      d.Get("ca_name").(string),
			ProductFlag: d.Get("product_flag").(string),
		}

		log.Printf("[INFO] Updating SCDN certificate name: %+v", req)
		_, err := service.EditCertificateName(req)
		if err != nil {
			return fmt.Errorf("failed to update certificate name: %w", err)
		}
	}

	// Update certificate content if changed
	if d.HasChange("ca_cert") || d.HasChange("ca_key") {
		// Both ca_cert and ca_key must be provided together when updating certificate content
		caCert, caCertOk := d.GetOk("ca_cert")
		caKey, caKeyOk := d.GetOk("ca_key")

		if !caCertOk || !caKeyOk {
			return fmt.Errorf("both ca_cert and ca_key must be provided when updating certificate content")
		}

		caCertStr := caCert.(string)
		caKeyStr := caKey.(string)

		// Check for placeholder values
		if caCertStr == "UPDATE_WITH_ACTUAL_CERTIFICATE_CONTENT" ||
			caKeyStr == "UPDATE_WITH_ACTUAL_PRIVATE_KEY_CONTENT" ||
			caCertStr == "" || caKeyStr == "" {
			return fmt.Errorf("ca_cert and ca_key must contain valid certificate content, not placeholder values")
		}

		req := scdn.CATextSaveRequest{
			ID:          certID,
			CAName:      d.Get("ca_name").(string),
			CACert:      caCertStr,
			CAKey:       caKeyStr,
			ProductFlag: d.Get("product_flag").(string),
		}

		log.Printf("[INFO] Updating SCDN certificate content: %+v", req)
		response, err := service.SaveCertificate(req)
		if err != nil {
			// Improve error message for certificate/key mismatch
			return fmt.Errorf("failed to update certificate content: %w. Please ensure the certificate and private key match", err)
		}

		// Update ca_sn if changed
		if response.Data.CASN != "" {
			if err := d.Set("ca_sn", response.Data.CASN); err != nil {
				log.Printf("[WARN] Failed to set certificate serial number: %v", err)
			}
		}
	}

	log.Printf("[INFO] SCDN certificate updated successfully: %s", d.Id())
	return resourceScdnCertificateRead(d, m)
}

func resourceScdnCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	certID := d.Id()

	req := scdn.CASelfDeleteRequest{
		IDs:         certID,
		ProductFlag: d.Get("product_flag").(string),
	}

	log.Printf("[INFO] Deleting SCDN certificate: %+v", req)
	_, err := service.DeleteCertificate(req)
	if err != nil {
		return fmt.Errorf("failed to delete SCDN certificate: %w", err)
	}

	d.SetId("")
	log.Printf("[INFO] SCDN certificate deleted successfully: %s", certID)
	return nil
}
