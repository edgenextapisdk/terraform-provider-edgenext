package domain

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnBriefDomains returns the SCDN brief domains data source
func DataSourceEdgenextScdnBriefDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnBriefDomainsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of domain IDs to query (optional, queries all if not specified)",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of brief domain information",
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
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of domains",
			},
		},
	}
}

func dataSourceScdnBriefDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.BriefDomainListRequest{}

	// Get domain IDs if specified
	if v, ok := d.GetOk("ids"); ok {
		idsRaw := v.([]interface{})
		ids := make([]int, len(idsRaw))
		for i, idRaw := range idsRaw {
			ids[i] = idRaw.(int)
		}
		req.IDs = ids
	}

	log.Printf("[INFO] Querying SCDN brief domains")
	response, err := service.ListBriefDomains(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN brief domains: %w", err)
	}

	// Convert domains to the format expected by Terraform
	domainsList := make([]map[string]interface{}, len(response.Data.List))
	ids := make([]string, len(response.Data.List))
	for i, domain := range response.Data.List {
		domainMap := map[string]interface{}{
			"id":     domain.ID,
			"domain": domain.Domain,
		}
		// Note: member_id is not available in BriefDomainInfo
		domainsList[i] = domainMap
		ids[i] = fmt.Sprintf("%d", domain.ID)
	}

	// Set the resource ID
	d.SetId(helper.DataResourceIdsHash(ids))

	// Set the domains list
	if err := d.Set("list", domainsList); err != nil {
		return fmt.Errorf("error setting list: %w", err)
	}

	// Set the total count
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total": response.Data.Total,
			"list":  domainsList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN brief domains query successful, found %d domains", len(domainsList))
	return nil
}
