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

// DataSourceEdgenextScdnCertificate returns the SCDN certificate data source
func DataSourceEdgenextScdnCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCertificateRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The certificate ID",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"ca_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate name",
			},
			"member_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The member ID",
			},
			"ca_sn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate serial number",
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
			"issuer_organization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The issuer organization",
			},
			"issuer_organization_element": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The issuer organization element",
			},
			"serial_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate serial number",
			},
			"issuer_object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The issuer object",
			},
			"use_organization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The use organization",
			},
			"use_organization_element": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The use organization element",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The city",
			},
			"province": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The province",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The country",
			},
			"authentication_usable_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authentication usable domain",
			},
		},
	}
}

func dataSourceScdnCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	certIDStr := d.Get("id").(string)
	certID, err := strconv.Atoi(certIDStr)
	if err != nil {
		return fmt.Errorf("invalid certificate ID: %w", err)
	}

	req := scdn.CASelfDetailRequest{
		ID: certID,
	}

	log.Printf("[INFO] Querying SCDN certificate: %d", certID)
	response, err := service.GetCertificateDetail(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN certificate: %w", err)
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
	if err := d.Set("issuer_expiry_time_auto_renew_status", response.Data.IssuerExpiryTimeAutoRenewStatus); err != nil {
		return fmt.Errorf("error setting issuer_expiry_time_auto_renew_status: %w", err)
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
	if err := d.Set("code", response.Data.Code); err != nil {
		return fmt.Errorf("error setting code: %w", err)
	}
	if err := d.Set("msg", response.Data.Msg); err != nil {
		return fmt.Errorf("error setting msg: %w", err)
	}
	if err := d.Set("created_at", response.Data.CreatedAt); err != nil {
		return fmt.Errorf("error setting created_at: %w", err)
	}
	if err := d.Set("updated_at", response.Data.UpdatedAt); err != nil {
		return fmt.Errorf("error setting updated_at: %w", err)
	}
	if err := d.Set("issuer_organization", response.Data.IssuerOrganization); err != nil {
		return fmt.Errorf("error setting issuer_organization: %w", err)
	}
	if err := d.Set("issuer_organization_element", response.Data.IssuerOrganizationElement); err != nil {
		return fmt.Errorf("error setting issuer_organization_element: %w", err)
	}
	if err := d.Set("serial_number", response.Data.SerialNumber); err != nil {
		return fmt.Errorf("error setting serial_number: %w", err)
	}
	if err := d.Set("issuer_object", response.Data.IssuerObject); err != nil {
		return fmt.Errorf("error setting issuer_object: %w", err)
	}
	if err := d.Set("use_organization", response.Data.UseOrganization); err != nil {
		return fmt.Errorf("error setting use_organization: %w", err)
	}
	if err := d.Set("use_organization_element", response.Data.UseOrganizationElement); err != nil {
		return fmt.Errorf("error setting use_organization_element: %w", err)
	}
	if err := d.Set("city", response.Data.City); err != nil {
		return fmt.Errorf("error setting city: %w", err)
	}
	if err := d.Set("province", response.Data.Province); err != nil {
		return fmt.Errorf("error setting province: %w", err)
	}
	if err := d.Set("country", response.Data.Country); err != nil {
		return fmt.Errorf("error setting country: %w", err)
	}
	if err := d.Set("authentication_usable_domain", response.Data.AuthenticationUsableDomain); err != nil {
		return fmt.Errorf("error setting authentication_usable_domain: %w", err)
	}

	// Set the certificate ID as the resource ID
	d.SetId(certIDStr)

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"id":                                   response.Data.ID,
			"ca_name":                              response.Data.CAName,
			"member_id":                            response.Data.MemberID,
			"issuer":                               response.Data.Issuer,
			"issuer_start_time":                    response.Data.IssuerStartTime,
			"issuer_expiry_time":                   response.Data.IssuerExpiryTime,
			"issuer_expiry_time_desc":              response.Data.IssuerExpiryTimeDesc,
			"issuer_expiry_time_auto_renew_status": response.Data.IssuerExpiryTimeAutoRenewStatus,
			"renew_status":                         response.Data.RenewStatus,
			"binded":                               response.Data.Binded,
			"ca_domain":                            response.Data.CADomain,
			"apply_status":                         response.Data.ApplyStatus,
			"ca_type":                              response.Data.CAType,
			"ca_type_domain":                       response.Data.CATypeDomain,
			"code":                                 response.Data.Code,
			"msg":                                  response.Data.Msg,
			"created_at":                           response.Data.CreatedAt,
			"updated_at":                           response.Data.UpdatedAt,
			"issuer_organization":                  response.Data.IssuerOrganization,
			"issuer_organization_element":          response.Data.IssuerOrganizationElement,
			"serial_number":                        response.Data.SerialNumber,
			"issuer_object":                        response.Data.IssuerObject,
			"use_organization":                     response.Data.UseOrganization,
			"use_organization_element":             response.Data.UseOrganizationElement,
			"city":                                 response.Data.City,
			"province":                             response.Data.Province,
			"country":                              response.Data.Country,
			"authentication_usable_domain":         response.Data.AuthenticationUsableDomain,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN certificate queried successfully: %s", certIDStr)
	return nil
}
