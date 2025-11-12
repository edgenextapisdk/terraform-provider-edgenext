package ssl

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// suppressCertificateFormatDiff suppresses certificate format differences
// When comparing certificate content, ignore newline format differences (\n vs real newlines)
func suppressCertificateFormatDiff(k, old, new string, d *schema.ResourceData) bool {
	// Normalize newlines: convert \n to real newlines
	normalizedOld := strings.ReplaceAll(old, "\\n", "\n")
	normalizedNew := strings.ReplaceAll(new, "\\n", "\n")

	// Remove all whitespace characters for comparison
	cleanOld := strings.ReplaceAll(strings.ReplaceAll(normalizedOld, "\n", ""), " ", "")
	cleanNew := strings.ReplaceAll(strings.ReplaceAll(normalizedNew, "\n", ""), " ", "")

	// If content is the same (ignoring format), suppress the difference
	return cleanOld == cleanNew
}

func ResourceEdgenextSslCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSslCertificateCreate,
		Read:   resourceSslCertificateRead,
		Update: resourceSslCertificateUpdate,
		Delete: resourceSslCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SSL certificate name",
			},
			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "SSL certificate content",
				DiffSuppressFunc: suppressCertificateFormatDiff,
			},
			"key": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				Description:      "SSL certificate private key content",
				DiffSuppressFunc: suppressCertificateFormatDiff,
			},
			"cert_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate ID",
			},
			"bind_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of bound domains",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cert_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate start time",
			},
			"cert_expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate end time",
			},
		},
	}
}

func resourceSslCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewSslCertificateService(client)

	certificate := d.Get("certificate").(string)
	privateKey := d.Get("key").(string)
	name := d.Get("name").(string)

	log.Printf("[INFO] Creating SSL certificate: name=%s", name)

	// Build certificate request
	request := SslCertificateRequest{
		Certificate: certificate,
		Key:         privateKey,
		Name:        name,
	}

	// Call create or update certificate API
	response, err := service.CreateOrUpdateSslCertificate(request)
	if err != nil {
		return fmt.Errorf("failed to create SSL certificate: %w", err)
	}

	// Set resource ID
	d.SetId(response.Data.CertID)

	log.Printf("[INFO] SSL certificate created successfully: %s", response.Data.CertID)
	return resourceSslCertificateRead(d, m)
}

func resourceSslCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewSslCertificateService(client)

	certID := d.Id()
	if certID == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}

	// Convert string ID to integer
	certIDInt, err := strconv.Atoi(certID)
	if err != nil {
		return fmt.Errorf("invalid certificate ID format: %s", certID)
	}

	log.Printf("[INFO] Reading SSL certificate: %s", certID)

	// Query certificate information
	response, err := service.GetSslCertificate(certIDInt)
	if err != nil {
		return fmt.Errorf("failed to read SSL certificate: %w", err)
	}

	// Set resource ID
	d.SetId(response.Data.CertID)
	// Set response data
	d.Set("name", response.Data.Name)
	d.Set("cert_id", response.Data.CertID)
	d.Set("cert_start_time", response.Data.CertStartTime)
	d.Set("cert_expire_time", response.Data.CertExpireTime)
	d.Set("bind_domains", response.Data.BindDomains)

	// Set certificate and key fields (use DiffSuppressFunc to handle format differences)
	d.Set("certificate", response.Data.Certificate)
	d.Set("key", response.Data.Key)

	return nil
}

func resourceSslCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewSslCertificateService(client)

	certificate := d.Get("certificate").(string)
	privateKey := d.Get("key").(string)
	name := d.Get("name").(string)

	// Get certificate ID for update
	certID := d.Id()
	if certID == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}

	certIDInt, err := strconv.Atoi(certID)
	if err != nil {
		return fmt.Errorf("invalid certificate ID format: %s", certID)
	}

	log.Printf("[INFO] Updating SSL certificate: %s", name)

	// Build certificate request
	request := SslCertificateRequest{
		Certificate: certificate,
		Key:         privateKey,
		Name:        name,
		CertID:      &certIDInt, // Specify the certificate ID to update
	}

	// Call create or update certificate API
	response, err := service.CreateOrUpdateSslCertificate(request)
	if err != nil {
		return fmt.Errorf("failed to update SSL certificate: %w", err)
	}

	log.Printf("[INFO] SSL certificate updated successfully: %s", response.Data.CertID)
	return resourceSslCertificateRead(d, m)
}

func resourceSslCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewSslCertificateService(client)

	certID := d.Id()
	if certID == "" {
		return fmt.Errorf("certificate ID cannot be empty")
	}

	certIDInt, err := strconv.Atoi(certID)
	if err != nil {
		return fmt.Errorf("invalid certificate ID format: %s", certID)
	}

	log.Printf("[INFO] Deleting SSL certificate: %s", certID)

	// Build delete request
	deleteReq := DeleteSslCertificateRequest{
		CertID: certIDInt,
	}

	// Call delete certificate API
	err = service.DeleteSslCertificate(deleteReq)
	if err != nil {
		return fmt.Errorf("failed to delete SSL certificate: %w", err)
	}
	d.SetId("")
	log.Printf("[INFO] SSL certificate deleted successfully: %s", certID)
	return nil
}
