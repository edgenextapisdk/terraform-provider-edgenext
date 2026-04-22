package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSNetworkInterfaces returns the data source schema for ECS network_interfaces.
func DataSourceENECSNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSNetworkInterfacesRead,
		Description: "Data source to query EdgeNext ECS network_interfaces (Neutron ports via extension list API).",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionDataSchema("region description"),
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by port name (partial match per API behavior).",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of ports to return.",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count from the API response (data.count).",
			},
			"network_interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS network interfaces (Neutron ports).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port name.",
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network (VPC) ID.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tenant ID (may be filled from origin_data when top-level is empty).",
						},
						"admin_state_up": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Administrative state of the port.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port status.",
						},
						"device_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Attached device (e.g. instance) ID.",
						},
						"device_owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Device owner string (e.g. compute:nova).",
						},
						"fixed_ips": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Fixed IP assignments.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet ID.",
									},
									"ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Fixed IP address.",
									},
									"floating_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Associated floating IP if present.",
									},
								},
							},
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
						"qos_policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "QoS policy ID.",
						},
						"port_security_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether port security is enabled.",
						},
						"security_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Security group IDs (merged from origin_data when top-level is null).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Tags.",
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
							Description: "Last update time (may be filled from origin_data when top-level is empty).",
						},
						"revision_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Revision number (may be filled from origin_data when top-level is zero).",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC address.",
						},
						"binding_vnic_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VNIC binding type (from binding:vnic_type; may be filled from origin_data).",
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resolved server (instance) name.",
						},
						"network_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resolved network name.",
						},
						"ipv4": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IPv4 addresses.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ipv6": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IPv6 addresses.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceENECSNetworkInterfacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"name":   d.Get("name").(string),
		"limit":  d.Get("limit").(int),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/ports/extension/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS network_interfaces: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS network_interfaces response: %s", err)
	}
	dataList := helper.ListFromMap(payload, "ports")
	items := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, flattenNetworkInterfacePort(row))
	}

	if err := d.Set("total", helper.IntFromMap(payload, "count")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "region", "name", "limit")
	if err := d.Set("network_interfaces", items); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenNetworkInterfacePort(row map[string]interface{}) map[string]interface{} {
	od := helper.MapFromMap(row, "origin_data")

	tenantID := helper.StringFromMap(row, "tenant_id")
	if tenantID == "" && od != nil {
		tenantID = helper.StringFromMap(od, "tenant_id")
	}

	sgs := helper.InterfaceToStringSlice(row["security_groups"])
	if len(sgs) == 0 && od != nil {
		sgs = helper.InterfaceToStringSlice(od["security_groups"])
	}

	bindingVnic := helper.StringFromMap(row, "binding:vnic_type")
	if bindingVnic == "" && od != nil {
		bindingVnic = helper.StringFromMap(od, "binding:vnic_type")
	}

	updatedAt := helper.StringFromMap(row, "updated_at")
	if updatedAt == "" && od != nil {
		updatedAt = helper.StringFromMap(od, "updated_at")
	}

	rev := helper.IntFromMap(row, "revision_number")
	if rev == 0 && od != nil {
		rev = helper.IntFromMap(od, "revision_number")
	}

	return map[string]interface{}{
		"id":                    helper.StringFromMap(row, "id"),
		"name":                  helper.StringFromMap(row, "name"),
		"network_id":            helper.StringFromMap(row, "network_id"),
		"tenant_id":             tenantID,
		"admin_state_up":        networkInterfaceBoolFromMap(row, "admin_state_up"),
		"status":                helper.StringFromMap(row, "status"),
		"device_id":             helper.StringFromMap(row, "device_id"),
		"device_owner":          helper.StringFromMap(row, "device_owner"),
		"fixed_ips":             flattenNetworkInterfaceFixedIPs(row["fixed_ips"]),
		"project_id":            helper.StringFromMap(row, "project_id"),
		"qos_policy_id":         helper.StringFromMap(row, "qos_policy_id"),
		"port_security_enabled": networkInterfaceBoolFromMap(row, "port_security_enabled"),
		"security_groups":       sgs,
		"description":           helper.StringFromMap(row, "description"),
		"tags":                  helper.InterfaceToStringSlice(row["tags"]),
		"created_at":            helper.StringFromMap(row, "created_at"),
		"updated_at":            updatedAt,
		"revision_number":       rev,
		"mac_address":           helper.StringFromMap(row, "mac_address"),
		"binding_vnic_type":     bindingVnic,
		"server_name":           helper.StringFromMap(row, "server_name"),
		"network_name":          helper.StringFromMap(row, "network_name"),
		"ipv4":                  helper.InterfaceToStringSlice(row["ipv4"]),
		"ipv6":                  helper.InterfaceToStringSlice(row["ipv6"]),
	}
}

func flattenNetworkInterfaceFixedIPs(v interface{}) []interface{} {
	raw := helper.InterfaceToList(v)
	out := make([]interface{}, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"subnet_id":   helper.StringFromMap(m, "subnet_id"),
			"ip_address":  helper.StringFromMap(m, "ip_address"),
			"floating_ip": helper.StringFromMap(m, "floating_ip"),
		})
	}
	return out
}

func networkInterfaceBoolFromMap(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
