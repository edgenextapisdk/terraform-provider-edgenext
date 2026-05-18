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

const elbL7PolicyListPath = "/elb/openapi/v2/l7policies/list"

// DataSourceENELBL7Policies lists L7 policies for a listener (POST .../l7policies/list).
func DataSourceENELBL7Policies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENELBL7PoliciesRead,
		Description: "Lists EdgeNext ELB L7 policies for a listener using POST /elb/openapi/v2/l7policies/list " +
			"(console RegisterL7PolicyHTTPServer).",
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Listener ID to list L7 policies for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1000,
				Description:  "Maximum policies to return (0-1000).",
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
			"l7_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Policy list items returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "L7 policy ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy description.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy action.",
						},
						"position": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy position (priority).",
						},
						"operating_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating status.",
						},
						"provisioning_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provisioning status.",
						},
						"rules_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of rules.",
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
				Description: "Total number of L7 policies matching the query (from the API when present).",
			},
			"hasmore": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether more pages exist (from API hasmore).",
			},
		},
	}
}

func dataSourceENELBL7PoliciesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"listener_id": strings.TrimSpace(d.Get("listener_id").(string)),
		"limit":       d.Get("limit").(int),
	}
	if sk := strings.TrimSpace(d.Get("sort_key").(string)); sk != "" {
		req["sort_key"] = sk
	}
	if sd := d.Get("sort_dir").(string); sd != "" {
		req["sort_dir"] = sd
	}

	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7PolicyListPath, req, &resp); err != nil {
		return diag.Errorf("failed to list ELB L7 policies: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB L7 policies list response: %s", err)
	}

	out := make([]interface{}, 0)
	for _, raw := range helper.ListFromMap(payload, "l7policies") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, l7PolicyListItemFromMap(row))
	}

	if err := d.Set("l7_policies", out); err != nil {
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

	helper.SetDataSourceStableID(d, "listener_id", "limit", "sort_key", "sort_dir")
	return nil
}

func l7PolicyListItemFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                  helper.StringFromMap(m, "id"),
		"name":                helper.StringFromMap(m, "name"),
		"description":         helper.StringFromMap(m, "description"),
		"action":              helper.StringFromMap(m, "action"),
		"position":            helper.IntFromMap(m, "position"),
		"operating_status":    helper.StringFromMap(m, "operating_status"),
		"provisioning_status": helper.StringFromMap(m, "provisioning_status"),
		"rules_count":         helper.IntFromMap(m, "rules_count"),
		"created_at":          helper.IntFromMap(m, "created_at"),
		"updated_at":          helper.IntFromMap(m, "updated_at"),
	}
}
