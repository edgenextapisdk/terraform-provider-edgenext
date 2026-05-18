package elb

import (
	"context"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const elbTargetGroupsListPath = "/elb/openapi/v2/pools/list"

// DataSourceENELBTargetGroups lists target groups for a load balancer.
// Target identifiers come from each row in the list response; use edgenext_elb_target_group_attachment to manage targets in Terraform.
// When healthmonitor_id is set and a target group has no embedded health monitor, a separate health monitor read may be used to resolve details.
func DataSourceENELBTargetGroups() *schema.Resource {
	healthMonitorElem := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Health monitor name.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Health monitor type.",
			},
			"max_retries": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum retries.",
			},
			"delay": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Seconds between health checks.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Health check timeout in seconds.",
			},
			"http_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "HTTP method for HTTP(S) checks.",
			},
			"url_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL path for HTTP(S) checks.",
			},
			"expected_codes": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expected HTTP status codes.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Health monitor ID.",
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
	}

	return &schema.Resource{
		ReadContext: dataSourceENELBTargetGroupsRead,
		Description: "Lists EdgeNext ELB target groups for a load balancer.",
		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Load balancer ID to list target groups for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"target_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Target groups for the load balancer. Includes health_monitor when resolvable and per-target identifiers when present (targets are managed with edgenext_elb_target_group_attachment).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group description.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target group protocol.",
						},
						"lb_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Load balancing algorithm.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "First listener id from the target group.",
						},
						"loadbalancer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "First load balancer id from the target group.",
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
						"healthmonitor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health monitor id (same as health_monitor.0.id when present).",
						},
						"health_monitor": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resolved health monitor for this target group when available; at most one element.",
							Elem:        healthMonitorElem,
						},
						"member_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Target identifiers when listed (manage targets in Terraform with edgenext_elb_target_group_attachment).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
		},
	}
}

func dataSourceENELBTargetGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	lbID := strings.TrimSpace(d.Get("loadbalancer_id").(string))
	req := map[string]interface{}{
		"loadbalancer_id": lbID,
	}

	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbTargetGroupsListPath, req, &resp); err != nil {
		return diag.Errorf("failed to list target groups for load balancer %q: %s", lbID, err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse target groups list response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "pools")
	targetGroups := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		sum, sumDiags := elbTargetGroupSummaryFromListRow(ctx, elbClient, row)
		if sumDiags.HasError() {
			return sumDiags
		}
		targetGroups = append(targetGroups, sum)
	}

	if err := d.Set("target_groups", targetGroups); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "loadbalancer_id")
	return nil
}

func elbTargetGroupSummaryFromListRow(ctx context.Context, elbClient *connectivity.ELBClient, tg map[string]interface{}) (map[string]interface{}, diag.Diagnostics) {
	summary := map[string]interface{}{
		"id":                  helper.StringFromMap(tg, "id"),
		"name":                helper.StringFromMap(tg, "name"),
		"description":         helper.StringFromMap(tg, "description"),
		"protocol":            helper.StringFromMap(tg, "protocol"),
		"lb_algorithm":        helper.StringFromMap(tg, "lb_algorithm"),
		"listener_id":         firstRefIDFromList(tg, "listeners"),
		"loadbalancer_id":     firstRefIDFromList(tg, "loadbalancers"),
		"provisioning_status": helper.StringFromMap(tg, "provisioning_status"),
		"operating_status":    helper.StringFromMap(tg, "operating_status"),
		"healthmonitor_id":    helper.StringFromMap(tg, "healthmonitor_id"),
		"member_ids":          targetIDsFromTargetGroupPayload(tg),
		"created_at":          helper.IntFromMap(tg, "created_at"),
		"updated_at":          helper.IntFromMap(tg, "updated_at"),
	}
	hmList, hmDiags := healthMonitorListForTargetGroupState(ctx, elbClient, tg, true)
	if hmDiags.HasError() {
		return nil, hmDiags
	}
	summary["health_monitor"] = hmList
	return summary, nil
}
