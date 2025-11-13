package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionTemplateDomains returns the SCDN security protection template domains data source
func DataSourceEdgenextScdnSecurityProtectionTemplateDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionTemplateDomainsRead,

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Business ID (template ID)",
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
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain filter",
			},
			"tpl_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template type: global, only_domain, more_domain",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of domains",
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Domain list",
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

func dataSourceScdnSecurityProtectionTemplateDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.SecurityProtectionTemplateBindDomainSearchRequest{
		BusinessID: d.Get("business_id").(int),
		Page:       d.Get("page").(int),
		PageSize:   d.Get("page_size").(int),
	}

	if domain, ok := d.GetOk("domain"); ok {
		req.Domain = domain.(string)
	}
	if tplType, ok := d.GetOk("tpl_type"); ok {
		req.TplType = tplType.(string)
	}

	log.Printf("[INFO] Querying SCDN security protection template domains: business_id=%d", req.BusinessID)
	response, err := service.GetSecurityProtectionTemplateBindDomains(req)
	if err != nil {
		return fmt.Errorf("failed to query template domains: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("template-domains-%d", req.BusinessID))

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

	log.Printf("[INFO] SCDN security protection template domains queried successfully: total=%d", response.Data.Total)
	return nil
}
