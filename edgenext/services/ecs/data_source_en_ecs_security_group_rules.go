package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSSecurityGroupRules returns the data source schema for ECS security group rules.
func DataSourceENECSSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSSecurityGroupRulesRead,
		Description: "Data source to query EdgeNext ECS security group rules.",
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The security group ID to filter rules.",
			},
			"security_group_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of security group rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule ID.",
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule tenant ID.",
						},
						"security_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group ID this rule belongs to.",
						},
						"ethertype": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version.",
						},
						"direction": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Traffic direction (ingress/egress).",
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
							Description: "Remote IP prefix.",
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
							Description: "Rule revision number.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
						},
					},
				},
			},
		},
	}
}

func dataSourceENECSSecurityGroupRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Get("security_group_id").(string),
	}

	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/detail", req, &resp); err != nil {
		return diag.Errorf("failed to read ECS security group rules: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS security group detail response: %s", err)
	}

	sg := helper.MapFromMap(payload, "security_group")
	if sg == nil {
		return diag.Errorf("failed to parse ECS security group detail response: missing security_group")
	}

	rulesRaw := helper.ListFromMap(sg, "security_group_rules")
	rules := make([]interface{}, 0, len(rulesRaw))
	for _, raw := range rulesRaw {
		rule, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		rules = append(rules, map[string]interface{}{
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

	helper.SetDataSourceStableID(d, "security_group_id")
	if err := d.Set("security_group_rules", rules); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
