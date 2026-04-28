package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSVpcSubnets returns the data source schema for ECS vpc subnets.
func DataSourceENECSVpcSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSVpcSubnetsRead,
		Description: "Data source to query EdgeNext ECS vpc subnets.",
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC ID to filter subnets.",
			},
			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of VPC subnets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant ID.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"ip_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP version.",
						},
						"subnetpool_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet pool ID.",
						},
						"enable_dhcp": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether DHCP is enabled.",
						},
						"ipv6_ra_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 RA mode.",
						},
						"ipv6_address_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IPv6 address mode.",
						},
						"gateway_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Gateway IP.",
						},
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDR block.",
						},
						"allocation_pools": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Allocation pools.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start IP.",
									},
									"end": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End IP.",
									},
								},
							},
						},
						"host_routes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Host routes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route destination.",
									},
									"nexthop": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route next hop.",
									},
								},
							},
						},
						"dns_nameservers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "DNS nameservers.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"service_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Service types.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
						"used_ips": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used IP count.",
						},
						"total_ips": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total IP count.",
						},
						"port_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port count.",
						},
						"not_bind_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reason if subnet is not bindable.",
						},
						"router_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bound router ID.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched subnets.",
			},
		},
	}
}

func dataSourceENECSVpcSubnetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"network_id": d.Get("vpc_id").(string),
	}
	var resp map[string]interface{}

	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/subnets_list", req, &resp); err != nil {
		return diag.Errorf("failed to read ECS vpc subnets: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS vpc subnets response: %s", err)
	}
	subnetsRaw := helper.ListFromMap(payload, "subnets")
	subnets := make([]interface{}, 0, len(subnetsRaw))
	for _, raw := range subnetsRaw {
		subnet, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		subnets = append(subnets, map[string]interface{}{
			"id":                helper.StringFromMap(subnet, "id"),
			"name":              helper.StringFromMap(subnet, "name"),
			"tenant_id":         helper.StringFromMap(subnet, "tenant_id"),
			"vpc_id":            helper.StringFromMap(subnet, "network_id"),
			"ip_version":        helper.IntFromMap(subnet, "ip_version"),
			"subnetpool_id":     helper.StringFromMap(subnet, "subnetpool_id"),
			"enable_dhcp":       subnetBool(subnet, "enable_dhcp"),
			"ipv6_ra_mode":      helper.StringFromMap(subnet, "ipv6_ra_mode"),
			"ipv6_address_mode": helper.StringFromMap(subnet, "ipv6_address_mode"),
			"gateway_ip":        helper.StringFromMap(subnet, "gateway_ip"),
			"cidr":              helper.StringFromMap(subnet, "cidr"),
			"allocation_pools":  subnetAllocationPools(subnet),
			"host_routes":       subnetHostRoutes(subnet),
			"dns_nameservers":   helper.InterfaceToStringSlice(subnet["dns_nameservers"]),
			"description":       helper.StringFromMap(subnet, "description"),
			"service_types":     helper.InterfaceToStringSlice(subnet["service_types"]),
			"tags":              helper.InterfaceToStringSlice(subnet["tags"]),
			"created_at":        helper.StringFromMap(subnet, "created_at"),
			"updated_at":        helper.StringFromMap(subnet, "updated_at"),
			"revision_number":   helper.IntFromMap(subnet, "revision_number"),
			"project_id":        helper.StringFromMap(subnet, "project_id"),
			"used_ips":          helper.IntFromMap(subnet, "used_ips"),
			"total_ips":         helper.IntFromMap(subnet, "total_ips"),
			"port_num":          helper.IntFromMap(subnet, "port_num"),
			"not_bind_reason":   helper.StringFromMap(subnet, "not_bind_reason"),
			"router_id":         helper.StringFromMap(subnet, "router_id"),
		})
	}
	d.Set("total", len(subnets))
	helper.SetDataSourceStableID(d, "vpc_id")
	if err := d.Set("subnets", subnets); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func subnetBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

func subnetAllocationPools(subnet map[string]interface{}) []interface{} {
	raw := helper.ListFromMap(subnet, "allocation_pools")
	out := make([]interface{}, 0, len(raw))
	for _, poolRaw := range raw {
		pool, ok := poolRaw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"start": helper.StringFromMap(pool, "start"),
			"end":   helper.StringFromMap(pool, "end"),
		})
	}
	return out
}

func subnetHostRoutes(subnet map[string]interface{}) []interface{} {
	raw := helper.ListFromMap(subnet, "host_routes")
	out := make([]interface{}, 0, len(raw))
	for _, routeRaw := range raw {
		route, ok := routeRaw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"destination": helper.StringFromMap(route, "destination"),
			"nexthop":     helper.StringFromMap(route, "nexthop"),
		})
	}
	return out
}
