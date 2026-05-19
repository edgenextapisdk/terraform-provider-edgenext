package elb

import (
	"context"
	"encoding/json"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const elbLoadBalancersListPath = "/elb/openapi/v2/load_balancers/list"

// DataSourceENELBLoadBalancers returns the data source schema for ELB load balancers.
func DataSourceENELBLoadBalancers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENELBLoadBalancersRead,
		Description: "Data source to query EdgeNext ELB load balancers.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Filter by load balancer name. Use empty string to omit the filter.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of load balancers to return in one request.",
			},
			"balancers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Load balancers returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancer ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancer name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"provisioning_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provisioning status (for example ACTIVE).",
						},
						"operating_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating status (for example ONLINE, ERROR).",
						},
						"listeners": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Listener references.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener ID.",
									},
								},
							},
						},
						"target_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target group references.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Target group ID.",
									},
								},
							},
						},
						"flavor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flavor ID.",
						},
						"vip_network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VIP network ID.",
						},
						"vip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VIP address.",
						},
						"vip_port_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VIP port ID.",
						},
						"vip_subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VIP subnet ID.",
						},
						"elastic_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Elastic IP if bound.",
						},
						"created_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation time as Unix timestamp (seconds).",
						},
						"updated_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last update time as Unix timestamp (seconds).",
						},
						"iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM resource ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancer type (for example application, network).",
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mode.",
						},
						"spec": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specification details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Spec ID.",
									},
									"spec_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Spec name.",
									},
									"flavor_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flavor ID in spec.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Spec description.",
									},
									"architecture": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Architecture (for example single).",
									},
									"max_connections": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum connections.",
									},
									"new_connections_per_sec": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "New connections per second.",
									},
									"qps": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Queries per second.",
									},
									"sort_order": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Sort order.",
									},
									"created_at": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Spec creation time as Unix timestamp.",
									},
									"updated_at": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Spec update time as Unix timestamp.",
									},
								},
							},
						},
						"floating_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Floating IP when present; JSON-encoded if the API returns an object, otherwise empty.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of load balancers matching the query (from the API when present).",
			},
			"hasmore": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether more results exist beyond this page.",
			},
		},
	}
}

func dataSourceENELBLoadBalancersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"name":  d.Get("name").(string),
		"limit": d.Get("limit").(int),
	}

	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbLoadBalancersListPath, req, &resp); err != nil {
		return diag.Errorf("failed to list load balancers: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse load balancers response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "loadbalancers")
	balancers := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		balancers = append(balancers, elbLoadBalancerAttrsFromMap(row))
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(balancers) > 0 {
		total = len(balancers)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hasmore", boolFromInterface(payload["hasmore"])); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("balancers", balancers); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "name", "limit")
	return nil
}

func elbLoadBalancerAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	listeners := idRefListFromMapList(m, "listeners")
	targetGroups := idRefListFromMapList(m, "pools")

	spec := make([]interface{}, 0)
	if s := helper.MapFromMap(m, "spec"); s != nil {
		spec = append(spec, map[string]interface{}{
			"id":                      helper.IntFromMap(s, "id"),
			"spec_name":               helper.StringFromMap(s, "spec_name"),
			"flavor_id":               helper.StringFromMap(s, "flavor_id"),
			"description":             helper.StringFromMap(s, "description"),
			"architecture":            helper.StringFromMap(s, "architecture"),
			"max_connections":         helper.IntFromMap(s, "max_connections"),
			"new_connections_per_sec": helper.IntFromMap(s, "new_connections_per_sec"),
			"qps":                     helper.IntFromMap(s, "qps"),
			"sort_order":              helper.IntFromMap(s, "sort_order"),
			"created_at":              helper.IntFromMap(s, "created_at"),
			"updated_at":              helper.IntFromMap(s, "updated_at"),
		})
	}

	return map[string]interface{}{
		"id":                  helper.StringFromMap(m, "id"),
		"name":                helper.StringFromMap(m, "name"),
		"description":         helper.StringFromMap(m, "description"),
		"provisioning_status": helper.StringFromMap(m, "provisioning_status"),
		"operating_status":    helper.StringFromMap(m, "operating_status"),
		"listeners":           listeners,
		"target_groups":       targetGroups,
		"flavor_id":           helper.StringFromMap(m, "flavor_id"),
		"vip_network_id":      helper.StringFromMap(m, "vip_network_id"),
		"vip_address":         helper.StringFromMap(m, "vip_address"),
		"vip_port_id":         helper.StringFromMap(m, "vip_port_id"),
		"vip_subnet_id":       helper.StringFromMap(m, "vip_subnet_id"),
		"elastic_ip":          helper.StringFromMap(m, "elastic_ip"),
		"created_at":          helper.IntFromMap(m, "created_at"),
		"updated_at":          helper.IntFromMap(m, "updated_at"),
		"iam_id":              helper.StringFromMap(m, "iam_id"),
		"type":                helper.StringFromMap(m, "type"),
		"mode":                helper.StringFromMap(m, "mode"),
		"spec":                spec,
		"floating_ip":         floatingIPString(m["floating_ip"]),
	}
}

func idRefListFromMapList(m map[string]interface{}, key string) []interface{} {
	out := make([]interface{}, 0)
	for _, raw := range helper.ListFromMap(m, key) {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"id": helper.StringFromMap(row, "id"),
		})
	}
	return out
}

func floatingIPString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func boolFromInterface(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
