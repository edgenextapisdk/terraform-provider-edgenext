package domain

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnAccessProgress returns the SCDN access progress data source
func DataSourceEdgenextScdnAccessProgress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnAccessProgressRead,

		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"progress": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of access progress status options",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Progress key/identifier",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Progress status name",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnAccessProgressRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	log.Printf("[INFO] Querying SCDN access progress status list")
	response, err := service.GetAccessProgress()
	if err != nil {
		return fmt.Errorf("failed to query SCDN access progress: %w", err)
	}

	// Set a fixed ID for this data source
	d.SetId("scdn-access-progress")

	// Convert progress list
	progressList := make([]map[string]interface{}, len(response.Data.Progress))
	for i, progress := range response.Data.Progress {
		progressList[i] = map[string]interface{}{
			"key":  progress.Key,
			"name": progress.Name,
		}
	}

	if err := d.Set("progress", progressList); err != nil {
		return fmt.Errorf("error setting progress: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"progress": progressList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN access progress query successful, found %d status options", len(progressList))
	return nil
}
