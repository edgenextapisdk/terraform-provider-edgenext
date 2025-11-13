package domain

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnDomainTemplates returns the SCDN domain templates data source
func DataSourceEdgenextScdnDomainTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnDomainTemplatesRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the domain to query templates",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"binded_templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of binded templates",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"business_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Business ID",
						},
						"business_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Business type",
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application type",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template name",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnDomainTemplatesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	req := scdn.DomainTemplatesRequest{
		DomainID: domainID,
	}

	log.Printf("[INFO] Querying SCDN domain templates for domain: %d", domainID)
	response, err := service.GetDomainTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN domain templates: %w", err)
	}

	// Set domain ID as resource ID
	d.SetId(fmt.Sprintf("%d", domainID))

	// Convert binded templates list
	templatesList := make([]map[string]interface{}, len(response.Data.BindedTemplates))
	for i, template := range response.Data.BindedTemplates {
		templatesList[i] = map[string]interface{}{
			"business_id":   template.BusinessID,
			"business_type": template.BusinessType,
			"app_type":      template.AppType,
			"name":          template.Name,
		}
	}

	if err := d.Set("binded_templates", templatesList); err != nil {
		return fmt.Errorf("error setting binded_templates: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"domain_id":        domainID,
			"binded_templates": templatesList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN domain templates query successful for domain: %d, found %d templates", domainID, len(templatesList))
	return nil
}
