package cert

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCertificateExport returns the SCDN certificate export data source
func DataSourceEdgenextScdnCertificateExport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCertificateExportRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The certificate ID (can be a single ID or comma-separated IDs)",
			},
			"product_flag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The product flag",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"exports": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of exported certificate data",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hash": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The export hash",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The export key",
						},
						"real_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The real URL for downloading the exported certificate",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnCertificateExportRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request
	req := scdn.CASelfExportRequest{
		ID:          d.Get("id").(string),
		ProductFlag: d.Get("product_flag").(string),
	}

	log.Printf("[INFO] Exporting SCDN certificate: %s", req.ID)
	response, err := service.ExportCertificate(req)
	if err != nil {
		return fmt.Errorf("failed to export SCDN certificate: %w", err)
	}

	// Convert exports to the format expected by Terraform
	exports := make([]map[string]interface{}, len(response.Data))
	for i, export := range response.Data {
		exportMap := map[string]interface{}{
			"hash":     export.Hash,
			"key":      export.Key,
			"real_url": export.RealURL,
		}
		exports[i] = exportMap
	}

	// Set the resource ID
	d.SetId(fmt.Sprintf("%s-%s", req.ID, req.ProductFlag))

	// Set the exports list
	if err := d.Set("exports", exports); err != nil {
		return fmt.Errorf("error setting exports: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"exports": exports,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN certificate exported successfully, %d export(s) found", len(response.Data))
	return nil
}
