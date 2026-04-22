package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSDisks returns the data source schema for ECS disks.
func DataSourceENECSDisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSDisksRead,
		Description: "Data source to query EdgeNext ECS disks.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionDataSchema("region description"),
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name to filter disks.",
			},
			"ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of disk IDs.",
			},
			"disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the disk.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the disk.",
						},
					},
				},
			},
		},
	}
}

func dataSourceENECSDisksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
	}
	if ids, ok := d.GetOk("ids"); ok {
		req["ids"] = ids.([]interface{})
	}
	if name, ok := d.GetOk("name"); ok {
		req["name"] = name.(string)
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/volume/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS disks: %s", err)
	}

	dataList, err := helper.ParseAPIResponseList(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS disks response: %s", err)
	}
	helper.SetDataSourceStableID(d, "region", "name", "ids")
	if err := d.Set("disks", dataList); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
