package elb

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceENELBL7Rules lists L7 rules for an L7 policy (POST .../l7rules/list).
func DataSourceENELBL7Rules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENELBL7RulesRead,
		Description: "Lists EdgeNext ELB L7 rules for an L7 policy using POST /elb/openapi/v2/l7rules/list " +
			"(console RegisterL7RuleHTTPServer).",
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "L7 policy ID to list rules for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1000,
				Description:  "Maximum rules to return (0-1000).",
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field key.",
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort direction: asc, desc, or empty.",
				ValidateFunc: validation.StringInSlice([]string{
					"", "asc", "desc",
				}, false),
			},
			"filter_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional filter: rule id.",
			},
			"filter_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional filter: HOSTNAME or PATH.",
			},
			"filter_provisioning_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional filter: provisioning_status.",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Rules returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule type.",
						},
						"compare_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compare type.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Match value.",
						},
						"invert": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Invert match.",
						},
						"admin_state_up": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Administrative state.",
						},
						"provisioning_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provisioning status.",
						},
						"operating_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating status.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project ID.",
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
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of L7 rules matching the query (from the API when present).",
			},
			"hasmore": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether more pages exist.",
			},
			"is_complete_task": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether async task is complete (from API is_complete_task).",
			},
		},
	}
}

func dataSourceENELBL7RulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"l7policy_id": strings.TrimSpace(d.Get("policy_id").(string)),
		"limit":       d.Get("limit").(int),
	}
	if sk := strings.TrimSpace(d.Get("sort_key").(string)); sk != "" {
		req["sort_key"] = sk
	}
	if sd := d.Get("sort_dir").(string); sd != "" {
		req["sort_dir"] = sd
	}
	if v := strings.TrimSpace(d.Get("filter_id").(string)); v != "" {
		req["id"] = v
	}
	if v := strings.TrimSpace(d.Get("filter_type").(string)); v != "" {
		req["type"] = v
	}
	if v := strings.TrimSpace(d.Get("filter_provisioning_status").(string)); v != "" {
		req["provisioning_status"] = v
	}

	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7RuleListPath, req, &resp); err != nil {
		return diag.Errorf("failed to list ELB L7 rules: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB L7 rules list response: %s", err)
	}

	out := make([]interface{}, 0)
	for _, raw := range helper.ListFromMap(payload, "rules") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, l7RuleListItemToState(row))
	}
	if err := d.Set("rules", out); err != nil {
		return diag.FromErr(err)
	}
	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(out) > 0 {
		total = len(out)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if v, ok := payload["hasmore"]; ok && v != nil {
		_ = d.Set("hasmore", v == true || v == "true" || fmt.Sprintf("%v", v) == "1")
	} else {
		_ = d.Set("hasmore", false)
	}
	if v, ok := payload["is_complete_task"]; ok && v != nil {
		_ = d.Set("is_complete_task", v == true || v == "true" || fmt.Sprintf("%v", v) == "1")
	} else {
		_ = d.Set("is_complete_task", false)
	}

	helper.SetDataSourceStableID(d, "policy_id", "limit", "sort_key", "sort_dir",
		"filter_id", "filter_type", "filter_provisioning_status")
	return nil
}

func l7RuleListItemToState(m map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{
		"id":                  helper.StringFromMap(m, "id"),
		"type":                helper.StringFromMap(m, "type"),
		"compare_type":        helper.StringFromMap(m, "compare_type"),
		"value":               helper.StringFromMap(m, "value"),
		"provisioning_status": helper.StringFromMap(m, "provisioning_status"),
		"operating_status":    helper.StringFromMap(m, "operating_status"),
		"project_id":          helper.StringFromMap(m, "project_id"),
		"created_at":          helper.IntFromMap(m, "created_at"),
		"updated_at":          helper.IntFromMap(m, "updated_at"),
	}
	if v, ok := m["invert"]; ok && v != nil {
		out["invert"] = v == true || v == "true" || fmt.Sprintf("%v", v) == "1"
	} else {
		out["invert"] = false
	}
	if v, ok := m["admin_state_up"]; ok && v != nil {
		out["admin_state_up"] = v == true || v == "true" || fmt.Sprintf("%v", v) == "1"
	} else {
		out["admin_state_up"] = false
	}
	return out
}
