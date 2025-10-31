package cert

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCertificates returns the SCDN certificates data source
func DataSourceEdgenextScdnCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCertificatesRead,

		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The page number for pagination",
			},
			"per_page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "The page size for pagination",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by domain name",
			},
			"product_flag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by product flag",
			},
			"ca_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by certificate name",
			},
			"binded": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by binding status: true-bound, false-unbound",
			},
			"apply_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by application status: 1-applying, 2-issued, 3-review failed, 4-uploaded",
			},
			"issuer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by issuer",
			},
			"expiry_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by expiry status: true-expired, false-not expired, inno-about to expire (within 30 days)",
			},
			"is_exact_search": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "off",
				Description: "Whether to use exact search: on-yes, off-no",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of certificates",
			},
			"issuer_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of issuers",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func dataSourceScdnCertificatesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request
	req := scdn.CASelfListRequest{
		Page:          d.Get("page").(int),
		PerPage:       d.Get("per_page").(int),
		Domain:        d.Get("domain").(string),
		ProductFlag:   d.Get("product_flag").(string),
		CAName:        d.Get("ca_name").(string),
		Binded:        d.Get("binded").(string),
		ApplyStatus:   d.Get("apply_status").(string),
		Issuer:        d.Get("issuer").(string),
		ExpiryTime:    d.Get("expiry_time").(string),
		IsExactSearch: d.Get("is_exact_search").(string),
	}

	log.Printf("[INFO] Querying SCDN certificates with filters: %+v", req)
	response, err := service.ListCertificates(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN certificates: %w", err)
	}

	// Convert certificates to the format expected by Terraform
	certificates := make([]map[string]interface{}, len(response.Data.List))
	ids := make([]string, len(response.Data.List))
	for i, cert := range response.Data.List {
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

	// Set the total count
	total, err := strconv.Atoi(response.Data.Total)
	if err == nil {
		if err := d.Set("total", total); err != nil {
			return fmt.Errorf("error setting total: %w", err)
		}
	}

	// Set the issuer list
	if err := d.Set("issuer_list", response.Data.IssuerList); err != nil {
		return fmt.Errorf("error setting issuer_list: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":        response.Data.Total,
			"issuer_list":  response.Data.IssuerList,
			"certificates": certificates,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN certificates queried successfully, %d certificates found", len(response.Data.List))
	return nil
}
