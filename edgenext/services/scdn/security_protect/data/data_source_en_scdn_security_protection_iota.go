package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionIota returns the SCDN security protection iota data source
func DataSourceEdgenextScdnSecurityProtectionIota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionIotaRead,

		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"iota": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Enum key-value pairs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceScdnSecurityProtectionIotaRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	log.Printf("[INFO] Querying SCDN security protection iota")
	response, err := service.GetSecurityProtectionIota()
	if err != nil {
		return fmt.Errorf("failed to query security protection iota: %w", err)
	}

	// Set resource ID
	d.SetId("security-protection-iota")

	// Convert iota map to schema format
	iotaMap := make(map[string]interface{})
	for k, v := range response.Data.Iota {
		iotaMap[k] = v
	}

	if err := d.Set("iota", iotaMap); err != nil {
		return fmt.Errorf("error setting iota: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"iota": iotaMap,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN security protection iota queried successfully")
	return nil
}
