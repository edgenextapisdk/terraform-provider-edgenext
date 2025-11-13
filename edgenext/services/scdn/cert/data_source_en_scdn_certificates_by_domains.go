package cert

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCertificatesByDomains returns the SCDN certificates by domains data source
func DataSourceEdgenextScdnCertificatesByDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCertificatesByDomainsRead,

		Schema: map[string]*schema.Schema{
			"domains": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The list of domain names to query certificates for",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of certificates",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate ID",
						},
						"member_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The member ID",
						},
						"ca_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate name",
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
						"issuer_expiry_time_auto_renew_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The certificate auto-renew status",
						},
						"renew_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The renewal status",
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
							Description: "The application status",
						},
						"ca_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate type",
						},
						"ca_type_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate domain type",
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The application error code",
						},
						"msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The application error message",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnCertificatesByDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Get domains from schema
	domainsInterface := d.Get("domains").([]interface{})
	domains := make([]string, len(domainsInterface))
	for i, domain := range domainsInterface {
		domains[i] = domain.(string)
	}

	// Build request
	req := scdn.CABatchListRequest{
		Domains: domains,
	}

	log.Printf("[INFO] Querying SCDN certificates by domains: %v", domains)
	response, err := service.ListCertificatesByDomains(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN certificates by domains: %w", err)
	}

	// Convert certificates to the format expected by Terraform
	certificates := make([]map[string]interface{}, len(response.Data))
	ids := make([]string, len(response.Data))
	for i, cert := range response.Data {
		certMap := map[string]interface{}{
			"id":                                   cert.ID,
			"member_id":                            cert.MemberID,
			"ca_name":                              cert.CAName,
			"issuer":                               cert.Issuer,
			"issuer_start_time":                    cert.IssuerStartTime,
			"issuer_expiry_time":                   cert.IssuerExpiryTime,
			"issuer_expiry_time_desc":              cert.IssuerExpiryTimeDesc,
			"issuer_expiry_time_auto_renew_status": cert.IssuerExpiryTimeAutoRenewStatus,
			"renew_status":                         cert.RenewStatus,
			"binded":                               cert.Binded,
			"ca_domain":                            cert.CADomain,
			"apply_status":                         cert.ApplyStatus,
			"ca_type":                              cert.CAType,
			"ca_type_domain":                       cert.CATypeDomain,
			"code":                                 cert.Code,
			"msg":                                  cert.Msg,
		}

		certificates[i] = certMap
		ids[i] = cert.ID
	}

	// Set the resource ID
	d.SetId(helper.DataResourceIdsHash(ids))

	// Set the certificates list
	if err := d.Set("certificates", certificates); err != nil {
		return fmt.Errorf("error setting certificates: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"certificates": certificates,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN certificates queried successfully, %d certificates found", len(response.Data))
	return nil
}
