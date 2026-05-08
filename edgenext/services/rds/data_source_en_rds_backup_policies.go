package rds

import (
	"context"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENRDSBackupPolicies returns the data source schema for RDS backup policies.
func DataSourceENRDSBackupPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENRDSBackupPoliciesRead,
		Description: "Data source to query EdgeNext RDS backup policies.",
		Schema: map[string]*schema.Schema{
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for the list request. The API request body uses the field name page_number.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Page size for the list request.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by policy ID. Maps to the API field id. Use empty string to omit the filter.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by policy name. Use empty string to omit the filter.",
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Backup policies returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy name.",
						},
						"region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy status.",
						},
						"schedule_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the schedule is enabled.",
						},
						"cycle_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup cycle type (for example weekly).",
						},
						"weekdays": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Weekday numbers in the schedule.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"hours": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Hour numbers in the schedule.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"retention_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Retention type (for example days or quantity).",
						},
						"retention_value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retention value from the API.",
						},
						"stop_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stop condition type.",
						},
						"stop_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stop condition value (format depends on stop_type).",
						},
						"next_trigger_time_utc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Next trigger time in UTC when returned by the API.",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"schedule_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Human-readable schedule description.",
						},
						"schedule_times": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Scheduled time strings (for example HH:MM:SS).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of policies matching the query.",
			},
		},
	}
}

func dataSourceENRDSBackupPoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"page_number": d.Get("page_num").(int),
		"page_size":   d.Get("page_size").(int),
		"id":          d.Get("policy_id").(string),
		"name":        d.Get("name").(string),
	}

	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/list", req, &resp); err != nil {
		return diag.Errorf("failed to list RDS backup policies: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse RDS backup policies response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "policies")
	policies := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		policies = append(policies, rdsBackupPolicyAttrsFromMap(row))
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(policies) > 0 {
		total = len(policies)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policies", policies); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "page_num", "page_size", "policy_id", "name")
	return nil
}

func rdsBackupPolicyAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                    helper.StringFromMap(m, "id"),
		"name":                  helper.StringFromMap(m, "name"),
		"region_id":             helper.StringFromMap(m, "region_id"),
		"status":                helper.StringFromMap(m, "status"),
		"schedule_enabled":      policyBoolFromInterface(m["schedule_enabled"]),
		"cycle_type":            helper.StringFromMap(m, "cycle_type"),
		"weekdays":              policyIntListFromMap(m, "weekdays"),
		"hours":                 policyIntListFromMap(m, "hours"),
		"retention_type":        helper.StringFromMap(m, "retention_type"),
		"retention_value":       helper.IntFromMap(m, "retention_value"),
		"stop_type":             helper.StringFromMap(m, "stop_type"),
		"stop_value":            policyStopValueString(m["stop_value"]),
		"next_trigger_time_utc": helper.StringFromMap(m, "next_trigger_time_utc"),
		"created":               helper.StringFromMap(m, "created"),
		"schedule_description":  helper.StringFromMap(m, "schedule_description"),
		"schedule_times":        policyStringListFromMap(m, "schedule_times"),
	}
}

func policyBoolFromInterface(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func policyIntListFromMap(m map[string]interface{}, key string) []interface{} {
	raw := helper.ListFromMap(m, key)
	out := make([]interface{}, 0, len(raw))
	for _, v := range raw {
		if n, ok := policyToInt(v); ok {
			out = append(out, n)
		}
	}
	return out
}

func policyToInt(v interface{}) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int32:
		return int(t), true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}

func policyStopValueString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func policyStringListFromMap(m map[string]interface{}, key string) []interface{} {
	raw := helper.ListFromMap(m, key)
	out := make([]interface{}, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
