package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSExternalGateways returns the data source schema for ECS external gateways.
func DataSourceENECSExternalGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSExternalGatewaysRead,
		Description: "Data source to query EdgeNext ECS external gateways.",
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of external gateways to return.",
			},
			"external_gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of external gateway networks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network name.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant ID.",
						},
						"admin_state_up": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether admin state is up.",
						},
						"mtu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network MTU.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network status.",
						},
						"subnets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subnet IDs.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"shared": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether network is shared.",
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
						"router_external": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is an external gateway network.",
						},
						"provider_network_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provider network type.",
						},
						"provider_physical_network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provider physical network.",
						},
						"provider_segmentation_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Provider segmentation ID.",
						},
						"availability_zone_hints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability zone hints.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"is_default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is default network.",
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability zones.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ipv4_address_scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv4 address scope.",
						},
						"ipv6_address_scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 address scope.",
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
							Description: "Last update time.",
						},
						"revision_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Revision number.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched external gateways.",
			},
		},
	}
}

func dataSourceENECSExternalGatewaysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"is_all":          true,
		"limit":           d.Get("limit").(int),
		"router_external": "true",
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/dict", req, &resp); err != nil {
		return diag.Errorf("failed to read ECS external gateways: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS external gateways response: %s", err)
	}
	networksRaw := helper.ListFromMap(payload, "networks")
	networks := make([]interface{}, 0, len(networksRaw))
	for _, raw := range networksRaw {
		network, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		networks = append(networks, map[string]interface{}{
			"id":                        helper.StringFromMap(network, "id"),
			"name":                      helper.StringFromMap(network, "name"),
			"tenant_id":                 helper.StringFromMap(network, "tenant_id"),
			"admin_state_up":            extGatewayBool(network, "admin_state_up"),
			"mtu":                       helper.IntFromMap(network, "mtu"),
			"status":                    helper.StringFromMap(network, "status"),
			"subnets":                   helper.InterfaceToStringSlice(network["subnets"]),
			"shared":                    extGatewayBool(network, "shared"),
			"project_id":                helper.StringFromMap(network, "project_id"),
			"qos_policy_id":             helper.StringFromMap(network, "qos_policy_id"),
			"port_security_enabled":     extGatewayBool(network, "port_security_enabled"),
			"router_external":           extGatewayBool(network, "router:external"),
			"provider_network_type":     helper.StringFromMap(network, "provider:network_type"),
			"provider_physical_network": helper.StringFromMap(network, "provider:physical_network"),
			"provider_segmentation_id":  helper.IntFromMap(network, "provider:segmentation_id"),
			"availability_zone_hints":   helper.InterfaceToStringSlice(network["availability_zone_hints"]),
			"is_default":                extGatewayBool(network, "is_default"),
			"availability_zones":        helper.InterfaceToStringSlice(network["availability_zones"]),
			"ipv4_address_scope":        helper.StringFromMap(network, "ipv4_address_scope"),
			"ipv6_address_scope":        helper.StringFromMap(network, "ipv6_address_scope"),
			"description":               helper.StringFromMap(network, "description"),
			"tags":                      helper.InterfaceToStringSlice(network["tags"]),
			"created_at":                helper.StringFromMap(network, "created_at"),
			"updated_at":                helper.StringFromMap(network, "updated_at"),
			"revision_number":           helper.IntFromMap(network, "revision_number"),
		})
	}
	if err := d.Set("total", helper.IntFromMap(payload, "count")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "limit")
	if err := d.Set("external_gateways", networks); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func extGatewayBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}
