package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSRouterPorts returns the data source schema for ECS router ports.
func DataSourceENECSRouterPorts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSRouterPortsRead,
		Description: "Data source to query EdgeNext ECS router ports.",
		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The router ID.",
			},
			"ports": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of router ports.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port name.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port IP address.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port MAC address.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port status.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of router ports.",
			},
		},
	}
}

func dataSourceENECSRouterPortsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Get("router_id").(string),
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/routers/port_list", req, &resp); err != nil {
		return diag.Errorf("failed to read ECS router ports: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS router ports response: %s", err)
	}
	rawPorts := helper.ListFromMap(payload, "ports")
	ports := make([]interface{}, 0, len(rawPorts))
	for _, raw := range rawPorts {
		port, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		ports = append(ports, map[string]interface{}{
			"id":          helper.StringFromMap(port, "id"),
			"name":        helper.StringFromMap(port, "name"),
			"ip_address":  helper.StringFromMap(port, "ip_address"),
			"mac_address": helper.StringFromMap(port, "mac_address"),
			"vpc_name":    helper.StringFromMap(port, "network_name"),
			"status":      helper.StringFromMap(port, "status"),
			"created_at":  helper.StringFromMap(port, "created_at"),
		})
	}

	if err := d.Set("total", helper.IntFromMap(payload, "count")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "router_id")
	if err := d.Set("ports", ports); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
