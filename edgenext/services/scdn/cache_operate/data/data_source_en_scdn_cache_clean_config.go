package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCacheCleanConfig returns the SCDN cache clean config data source
func DataSourceEdgenextScdnCacheCleanConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCacheCleanConfigRead,

		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Config ID",
			},
			"wholesite": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Whole site config",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"specialurl": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Special URL config",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"specialdir": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Special directory config",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceScdnCacheCleanConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.CacheCleanGetConfigRequest{}

	log.Printf("[INFO] Querying SCDN cache clean config")
	response, err := service.GetCacheCleanConfig(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN cache clean config: %w", err)
	}

	// Set all fields
	if err := d.Set("id", response.Data.ID); err != nil {
		return fmt.Errorf("error setting id: %w", err)
	}
	if err := d.Set("wholesite", response.Data.Wholesite); err != nil {
		return fmt.Errorf("error setting wholesite: %w", err)
	}
	if err := d.Set("specialurl", response.Data.Specialurl); err != nil {
		return fmt.Errorf("error setting specialurl: %w", err)
	}
	if err := d.Set("specialdir", response.Data.Specialdir); err != nil {
		return fmt.Errorf("error setting specialdir: %w", err)
	}

	// Set the config ID as the resource ID
	d.SetId(response.Data.ID)

	log.Printf("[INFO] SCDN cache clean config queried successfully: %s", response.Data.ID)
	return nil
}
