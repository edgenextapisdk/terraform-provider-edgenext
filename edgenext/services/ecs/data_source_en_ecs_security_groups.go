package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSSecurityGroups returns the data source schema for ECS security_groups.
func DataSourceENECSSecurityGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSSecurityGroupsRead,
		Description: "Data source to query EdgeNext ECS security_groups.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name to filter security_groups.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Maximum number of security_groups to return.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS security_groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the security_group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the security_group.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tenant ID.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group description.",
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
						"security_group_rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Security group rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule ID.",
									},
									"tenant_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule tenant ID.",
									},
									"security_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Security group ID.",
									},
									"ethertype": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP version type.",
									},
									"direction": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Traffic direction.",
									},
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol name.",
									},
									"port_range_min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum port.",
									},
									"port_range_max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum port.",
									},
									"remote_ip_prefix": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Remote IP CIDR.",
									},
									"remote_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Remote security group ID.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule description.",
									},
									"tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Rule tags.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule creation time.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule update time.",
									},
									"revision_number": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Rule revision number.",
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Rule project ID.",
									},
								},
							},
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched security groups.",
			},
		},
	}
}

func dataSourceENECSSecurityGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"name":  d.Get("name").(string),
		"limit": d.Get("limit").(int),
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS security_groups: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS security_groups response: %s", err)
	}
	dataList := helper.ListFromMap(payload, "security_groups")
	items := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, securityGroupAttrsFromMap(row))
	}
	if err := d.Set("total", helper.IntFromMap(payload, "count")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "name", "limit")
	if err := d.Set("security_groups", items); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func securityGroupAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                   helper.StringFromMap(m, "id"),
		"name":                 helper.StringFromMap(m, "name"),
		"tenant_id":            helper.StringFromMap(m, "tenant_id"),
		"description":          helper.StringFromMap(m, "description"),
		"security_group_rules": normalizeSecurityGroupRules(helper.ListFromMap(m, "security_group_rules")),
		"tags":                 helper.InterfaceToStringSlice(m["tags"]),
		"created_at":           helper.StringFromMap(m, "created_at"),
		"updated_at":           helper.StringFromMap(m, "updated_at"),
		"revision_number":      helper.IntFromMap(m, "revision_number"),
		"project_id":           helper.StringFromMap(m, "project_id"),
	}
}

func normalizeSecurityGroupRules(rules []interface{}) []interface{} {
	out := make([]interface{}, 0, len(rules))
	for _, raw := range rules {
		rule, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":                helper.StringFromMap(rule, "id"),
			"tenant_id":         helper.StringFromMap(rule, "tenant_id"),
			"security_group_id": helper.StringFromMap(rule, "security_group_id"),
			"ethertype":         helper.StringFromMap(rule, "ethertype"),
			"direction":         helper.StringFromMap(rule, "direction"),
			"protocol":          helper.StringFromMap(rule, "protocol"),
			"port_range_min":    helper.IntFromMap(rule, "port_range_min"),
			"port_range_max":    helper.IntFromMap(rule, "port_range_max"),
			"remote_ip_prefix":  helper.StringFromMap(rule, "remote_ip_prefix"),
			"remote_group_id":   helper.StringFromMap(rule, "remote_group_id"),
			"description":       helper.StringFromMap(rule, "description"),
			"tags":              helper.InterfaceToStringSlice(rule["tags"]),
			"created_at":        helper.StringFromMap(rule, "created_at"),
			"updated_at":        helper.StringFromMap(rule, "updated_at"),
			"revision_number":   helper.IntFromMap(rule, "revision_number"),
			"project_id":        helper.StringFromMap(rule, "project_id"),
		})
	}
	return out
}
