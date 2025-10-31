package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionTemplateUnboundDomains returns the SCDN security protection template unbound domains data source
func DataSourceEdgenextScdnSecurityProtectionTemplateUnboundDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionTemplateUnboundDomainsRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain filter",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "Page size",
			},
			"member_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Member ID",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of unbound domains",
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Unbound domain list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Domain ID",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain type",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnSecurityProtectionTemplateUnboundDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.SecurityProtectionTemplateUnboundDomainSearchRequest{
		Page:     d.Get("page").(int),
		PageSize: d.Get("page_size").(int),
	}

	if domain, ok := d.GetOk("domain"); ok {
		req.Domain = domain.(string)
	}
	if memberID, ok := d.GetOk("member_id"); ok {
		req.MemberID = memberID.(int)
	}

	log.Printf("[INFO] Querying SCDN security protection template unbound domains")
	response, err := service.GetSecurityProtectionTemplateUnboundDomains(req)
	if err != nil {
		return fmt.Errorf("failed to query unbound template domains: %w", err)
	}

	// Set resource ID
	d.SetId("template-unbound-domains")

	// Set total
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Convert domains to schema format
	domainList := make([]map[string]interface{}, 0, len(response.Data.Domains))
	for _, domain := range response.Data.Domains {
		domainMap := map[string]interface{}{
			"id":         domain.ID,
			"domain":     domain.Domain,
			"type":       domain.Type,
			"created_at": domain.CreatedAt,
			"remark":     domain.Remark,
		}
		domainList = append(domainList, domainMap)
	}

	if err := d.Set("domains", domainList); err != nil {
		return fmt.Errorf("error setting domains: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":   response.Data.Total,
			"domains": domainList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN security protection template unbound domains queried successfully: total=%d", response.Data.Total)
	return nil
}
