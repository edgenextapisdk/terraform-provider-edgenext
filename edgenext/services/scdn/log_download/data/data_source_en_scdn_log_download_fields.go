package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnLogDownloadFields returns the SCDN log download fields data source
func DataSourceEdgenextScdnLogDownloadFields() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnLogDownloadFieldsRead,

		Schema: map[string]*schema.Schema{
			"data_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by data source: ng, cc, waf. If not specified, returns all data sources",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Field configurations by data source",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source key: ng, cc, waf",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source name",
						},
						"download_fields": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available download fields",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"search_terms": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available search terms",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnLogDownloadFieldsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	log.Printf("[DEBUG] Reading SCDN log download fields")
	response, err := service.GetLogDownloadFields()
	if err != nil {
		return fmt.Errorf("failed to get log download fields: %w", err)
	}

	// Filter by data source if specified
	dataSourceFilter := ""
	if v, ok := d.GetOk("data_source"); ok {
		dataSourceFilter = v.(string)
	}

	// Convert to list format for Terraform
	configsList := make([]map[string]interface{}, 0)
	for dsKey, config := range response.Data {
		// Apply filter if specified
		if dataSourceFilter != "" && dsKey != dataSourceFilter {
			continue
		}

		configMap := map[string]interface{}{
			"data_source":     dsKey,
			"name":            config.Name,
			"download_fields": config.DownloadFields,
			"search_terms":    config.SearchTerms,
		}
		configsList = append(configsList, configMap)
	}

	d.Set("configs", configsList)
	d.SetId("log-download-fields")

	// Save to file if specified
	if v, ok := d.GetOk("result_output_file"); ok {
		outputFile := v.(string)
		outputData := map[string]interface{}{
			"configs": configsList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
		_ = outputFile // Suppress unused variable warning
	}

	return nil
}
