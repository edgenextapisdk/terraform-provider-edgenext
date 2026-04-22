package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSRouters returns the data source schema for ECS routers.
func DataSourceENECSRouters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSRoutersRead,
		Description: "Data source to query EdgeNext ECS routers.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionDataSchema("region description"),
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name to filter routers.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of routers to return.",
			},
			"routers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS routers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the router.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the router.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tenant ID.",
						},
						"admin_state_up": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether router admin state is up.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Router status.",
						},
						"external_gateway_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "External gateway info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "External network ID.",
									},
									"network_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "External network name.",
									},
									"enable_snat": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether SNAT is enabled.",
									},
									"external_fixed_ips": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "External fixed IPs.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "External fixed IP subnet ID.",
												},
												"ip_address": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "External fixed IP address.",
												},
											},
										},
									},
								},
							},
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability zones.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"availability_zone_hints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability zone hints.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"routes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Static routes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route destination CIDR.",
									},
									"nexthop": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Route next hop.",
									},
								},
							},
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flavor ID.",
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
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched routers.",
			},
		},
	}
}

func dataSourceENECSRoutersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/routers/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS routers: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS routers response: %s", err)
	}
	dataList := helper.ListFromMap(payload, "routers")
	items := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		router, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, map[string]interface{}{
			"id":                      helper.StringFromMap(router, "id"),
			"name":                    helper.StringFromMap(router, "name"),
			"tenant_id":               helper.StringFromMap(router, "tenant_id"),
			"admin_state_up":          routerBool(router, "admin_state_up"),
			"status":                  helper.StringFromMap(router, "status"),
			"external_gateway_info":   routerExternalGatewayInfo(router),
			"description":             helper.StringFromMap(router, "description"),
			"availability_zones":      helper.InterfaceToStringSlice(router["availability_zones"]),
			"availability_zone_hints": helper.InterfaceToStringSlice(router["availability_zone_hints"]),
			"routes":                  routerRoutes(router),
			"flavor_id":               helper.StringFromMap(router, "flavor_id"),
			"tags":                    helper.InterfaceToStringSlice(router["tags"]),
			"created_at":              helper.StringFromMap(router, "created_at"),
			"updated_at":              helper.StringFromMap(router, "updated_at"),
			"revision_number":         helper.IntFromMap(router, "revision_number"),
			"project_id":              helper.StringFromMap(router, "project_id"),
		})
	}
	if err := d.Set("total", helper.IntFromMap(payload, "count")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "region", "name", "limit")
	if err := d.Set("routers", items); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func routerBool(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

func routerExternalGatewayInfo(router map[string]interface{}) []interface{} {
	rawGateway := helper.MapFromMap(router, "external_gateway_info")
	if rawGateway == nil {
		return []interface{}{}
	}

	rawFixedIPs := helper.ListFromMap(rawGateway, "external_fixed_ips")
	fixedIPs := make([]interface{}, 0, len(rawFixedIPs))
	for _, raw := range rawFixedIPs {
		ip, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		fixedIPs = append(fixedIPs, map[string]interface{}{
			"subnet_id":  helper.StringFromMap(ip, "subnet_id"),
			"ip_address": helper.StringFromMap(ip, "ip_address"),
		})
	}

	return []interface{}{
		map[string]interface{}{
			"network_id":         helper.StringFromMap(rawGateway, "network_id"),
			"network_name":       helper.StringFromMap(rawGateway, "network_name"),
			"enable_snat":        routerBool(rawGateway, "enable_snat"),
			"external_fixed_ips": fixedIPs,
		},
	}
}

func routerRoutes(router map[string]interface{}) []interface{} {
	rawRoutes := helper.ListFromMap(router, "routes")
	routes := make([]interface{}, 0, len(rawRoutes))
	for _, raw := range rawRoutes {
		route, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		routes = append(routes, map[string]interface{}{
			"destination": helper.StringFromMap(route, "destination"),
			"nexthop":     helper.StringFromMap(route, "nexthop"),
		})
	}
	return routes
}
