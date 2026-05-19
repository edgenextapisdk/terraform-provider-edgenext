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

// DataSourceENELBTargetGroupAttachments lists backend targets for a target group.
func DataSourceENELBTargetGroupAttachments() *schema.Resource {
	targetElem := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target name.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target IP address.",
			},
			"protocol_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Protocol port.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Target weight.",
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
		ReadContext: dataSourceENELBTargetGroupAttachmentsRead,
		Description: "Lists backend targets for an ELB target group.",
		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Target group ID to list targets for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Backend targets for the target group.",
				Elem:        targetElem,
			},
		},
	}
}

func dataSourceENELBTargetGroupAttachmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	targetGroupID := strings.TrimSpace(d.Get("target_group_id").(string))
	req := map[string]interface{}{
		"pool_id": targetGroupID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbTargetGroupTargetListPath, req, &resp); err != nil {
		return diag.Errorf("failed to list targets for target group %q: %s", targetGroupID, err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB target group targets list response: %s", err)
	}

	out := make([]interface{}, 0)
	for _, raw := range helper.ListFromMap(payload, "members") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, targetSummaryFromListRow(row))
	}

	if err := d.Set("targets", out); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(targetGroupID)
	return nil
}

func targetSummaryFromListRow(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                  helper.StringFromMap(m, "id"),
		"name":                helper.StringFromMap(m, "name"),
		"address":             helper.StringFromMap(m, "address"),
		"protocol_port":       helper.IntFromMap(m, "protocol_port"),
		"weight":              helper.IntFromMap(m, "weight"),
		"provisioning_status": helper.StringFromMap(m, "provisioning_status"),
		"operating_status":    helper.StringFromMap(m, "operating_status"),
		"created_at":          helper.IntFromMap(m, "created_at"),
		"updated_at":          helper.IntFromMap(m, "updated_at"),
	}
}
