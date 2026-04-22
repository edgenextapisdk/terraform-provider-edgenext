package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSVpcs returns the data source schema for ECS vpcs.
func DataSourceENECSVpcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSVpcsRead,
		Description: "Data source to query EdgeNext ECS vpcs.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionDataSchema("region description"),
			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The network ID to filter vpcs.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name to filter vpcs.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of vpcs to return.",
			},
			"vpcs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS vpcs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the vpc.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the vpc.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the vpc.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project ID.",
						},
						"ipv4_cidrs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of IPv4 CIDRs.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The vpc description.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched vpcs.",
			},
		},
	}
}

func dataSourceENECSVpcsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":     helper.NormalizeRegion(d.Get("region").(string)),
		"network_id": d.Get("network_id").(string),
		"name":       d.Get("name").(string),
		"limit":      d.Get("limit").(int),
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS vpcs: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS vpcs response: %s", err)
	}
	dataList := helper.ListFromMap(payload, "networks")
	items := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, map[string]interface{}{
			"id":          helper.StringFromMap(row, "id"),
			"name":        helper.StringFromMap(row, "name"),
			"status":      helper.StringFromMap(row, "status"),
			"project_id":  helper.StringFromMap(row, "project_id"),
			"ipv4_cidrs":  helper.InterfaceToStringSlice(row["ipv4_cidrs"]),
			"description": helper.StringFromMap(row, "description"),
			"created_at":  helper.StringFromMap(row, "created_at"),
			"updated_at":  helper.StringFromMap(row, "updated_at"),
		})
	}
	if err := d.Set("total", helper.IntFromMap(payload, "count")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "region", "network_id", "name", "limit")
	if err := d.Set("vpcs", items); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
