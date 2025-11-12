package ssl

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextSslCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSslCertificateRead,

		Schema: map[string]*schema.Schema{
			"cert_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate ID",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate name",
			},
			"certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Certificate content",
			},
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Private key content",
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

func dataSourceSslCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewSslCertificateService(client)

	certID := d.Get("cert_id").(string)

	log.Printf("[INFO] Querying SSL certificate ID: %s", certID)

	// Query certificate information - need to find certificate ID by domain
	certIDInt, err := strconv.Atoi(certID)
	if err != nil {
		return fmt.Errorf("invalid certificate ID: %s", certID)
	}
	response, err := service.GetSslCertificate(certIDInt)
	if err != nil {
		return fmt.Errorf("failed to query SSL certificate: %w", err)
	}
	foundCert := response.Data
	// Set resource ID
	d.SetId(response.Data.CertID)

	// Set response data
	d.Set("cert_id", foundCert.CertID)
	d.Set("name", foundCert.Name)
	d.Set("certificate", foundCert.Certificate)
	d.Set("key", foundCert.Key)
	d.Set("associated_domains", foundCert.BindDomains)
	d.Set("cert_start_time", foundCert.CertStartTime)
	d.Set("cert_expire_time", foundCert.CertExpireTime)

	// Write result to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"cert_id":          foundCert.CertID,
			"name":             foundCert.Name,
			"certificate":      foundCert.Certificate,
			"key":              foundCert.Key,
			"bind_domains":     foundCert.BindDomains,
			"cert_start_time":  foundCert.CertStartTime,
			"cert_expire_time": foundCert.CertExpireTime,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SSL certificate query successful: %s", foundCert.CertID)
	return nil
}

// DataSourceEdgenextSslCertificates data source for querying multiple SSL certificates
func DataSourceEdgenextSslCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSslCertificatesRead,

		Schema: map[string]*schema.Schema{
			"page_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number, must be greater than 0 if specified",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					val := i.(int)
					if val <= 0 {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid page_number value",
							Detail:   fmt.Sprintf("page_number must be greater than 0, current value: %d", val),
						})
					}
					return diags
				},
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Number of items per page, range 1-500 if specified",
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					val := i.(int)
					if val < 1 || val > 500 {
						diags = append(diags, diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Invalid page_size value",
							Detail:   fmt.Sprintf("page_size range is 1-500, current value: %d", val),
						})
					}
					return diags
				},
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SSL certificate list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate name",
						},
						"associated_domains": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Associated domain list",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"include_domains": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Included domain list",
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
				},
			},
		},
	}
}

func dataSourceSslCertificatesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewSslCertificateService(client)

	log.Printf("[INFO] Querying all SSL certificates")

	pageNumber := d.Get("page_number").(int)
	pageSize := d.Get("page_size").(int)

	// Query certificate list
	response, err := service.ListSslCertificates(pageNumber, pageSize) // Get all certificates
	if err != nil {
		return fmt.Errorf("failed to query SSL certificate list: %w", err)
	}

	// Set certificate list
	var certificates []map[string]interface{}
	ids := make([]string, 0)
	for _, cert := range response.Data.List {
		certMap := map[string]interface{}{
			"cert_id":            cert.CertID,
			"name":               cert.Name,
			"associated_domains": cert.AssociatedDomains,
			"include_domains":    cert.IncludeDomains,
			"cert_start_time":    cert.CertStartTime,
			"cert_expire_time":   cert.CertExpireTime,
		}
		certificates = append(certificates, certMap)
		ids = append(ids, cert.CertID)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("list", certificates)
	if err != nil {
		log.Printf("[ERROR] Failed to set certificate list: %v", err)
		return err
	}

	// Write result to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"page_number": pageNumber,
			"page_size":   pageSize,
			"list":        certificates,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SSL certificate list query successful, total %d certificates", len(certificates))
	return nil
}
