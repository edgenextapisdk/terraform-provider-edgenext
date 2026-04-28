package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSFloatingIps returns the data source schema for ECS floating_ips.
func DataSourceENECSFloatingIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSFloatingIpsRead,
		Description: "Data source to query EdgeNext ECS floating_ips.",
		Schema: map[string]*schema.Schema{
			"floating_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The floating IP ID to filter.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The floating IP address to filter.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of floating IPs to return.",
			},
			"floating_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS floating_ips.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the floating_ip.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tenant ID.",
						},
						"floating_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The floating IP address.",
						},
						"floating_network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The floating network ID.",
						},
						"router_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The router ID.",
						},
						"network_interface_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network interface ID.",
						},
						"fixed_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fixed IP address.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description.",
						},
						"qos_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The QoS policy ID.",
						},
						"port_forwardings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Port forwarding entries.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of tag strings.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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
						"revision_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Revision number.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth in Mbps.",
						},
						"charge_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge mode.",
						},
						"floating_network_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Floating network name.",
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time.",
						},
						"network_interface_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network interface name.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"billing_model": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing model.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched floating IPs.",
			},
		},
	}
}

func dataSourceENECSFloatingIpsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id":                  d.Get("floating_ip_id").(string),
		"limit":               d.Get("limit").(int),
		"floating_ip_address": d.Get("floating_ip_address").(string),
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/floatingips/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS floating_ips: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS floating_ips response: %s", err)
	}
	dataList := helper.ListFromMap(payload, "floating_ip")
	items := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, map[string]interface{}{
			"id":                     helper.StringFromMap(row, "id"),
			"tenant_id":              helper.StringFromMap(row, "tenant_id"),
			"floating_ip_address":    helper.StringFromMap(row, "floating_ip_address"),
			"floating_network_id":    helper.StringFromMap(row, "floating_network_id"),
			"router_id":              helper.StringFromMap(row, "router_id"),
			"network_interface_id":   helper.StringFromMap(row, "port_id"),
			"fixed_ip_address":       helper.StringFromMap(row, "fixed_ip_address"),
			"status":                 helper.StringFromMap(row, "status"),
			"description":            helper.StringFromMap(row, "description"),
			"qos_policy_id":          helper.StringFromMap(row, "qos_policy_id"),
			"port_forwardings":       helper.InterfaceToStringSlice(row["port_forwardings"]),
			"tags":                   helper.InterfaceToStringSlice(row["tags"]),
			"created_at":             helper.StringFromMap(row, "created_at"),
			"updated_at":             helper.StringFromMap(row, "updated_at"),
			"revision_number":        helper.IntFromMap(row, "revision_number"),
			"project_id":             helper.StringFromMap(row, "project_id"),
			"bandwidth":              helper.IntFromMap(row, "bandwidth"),
			"charge_mode":            helper.StringFromMap(row, "charge_mode"),
			"floating_network_name":  helper.StringFromMap(row, "floating_network_name"),
			"expiration_time":        helper.StringFromMap(row, "expiration_time"),
			"network_interface_name": helper.StringFromMap(row, "port_name"),
			"instance_name":          helper.StringFromMap(row, "instance_name"),
			"billing_model":          helper.IntFromMap(row, "billing_model"),
		})
	}
	if err := d.Set("total", helper.IntFromMap(payload, "count")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "floating_ip_id", "floating_ip_address", "limit")
	if err := d.Set("floating_ips", items); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
