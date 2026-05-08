package rds

import (
	"context"
	"sort"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	rdsBackupPolicyStopTypeOptions      = []string{"end_time", "never", "count"}
	rdsBackupPolicyRetentionTypeOptions = []string{"permanent", "days", "quantity"}
)

// ResourceENRDSBackupPolicy returns the resource schema for an RDS backup policy.
func ResourceENRDSBackupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENRDSBackupPolicyCreate,
		ReadContext:   resourceENRDSBackupPolicyRead,
		UpdateContext: resourceENRDSBackupPolicyUpdate,
		DeleteContext: resourceENRDSBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENRDSBackupPolicyImport,
		},
		Description: "Manages an EdgeNext RDS backup policy.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Backup policy name.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"cycle_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Backup cycle type (for example weekly).",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"weekdays": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Weekday numbers for scheduling (for example 1..7).",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"schedule_times": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Schedule time strings (for example 00:00:00).",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
			"retention_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Retention type. Allowed values: permanent, days, quantity.",
				ValidateFunc: validation.StringInSlice(rdsBackupPolicyRetentionTypeOptions, false),
			},
			"retention_value": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Retention value associated with retention_type.",
			},
			"stop_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Stop condition type. Allowed values: end_time, never, count.",
				ValidateFunc: validation.StringInSlice(rdsBackupPolicyStopTypeOptions, false),
			},
			"stop_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Stop condition value, for example a datetime when stop_type is end_time.",
			},
			"schedule_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable the schedule.",
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
			"hours": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Hour numbers in the schedule.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
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
		},
	}
}

func resourceENRDSBackupPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	d.SetId(strings.TrimSpace(d.Id()))
	diags := resourceENRDSBackupPolicyRead(ctx, d, m)
	if diags.HasError() {
		return nil, diagToError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENRDSBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := rdsBackupPolicyRequestFromResourceData(d, "", true)
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/create", req, &resp); err != nil {
		return diag.Errorf("failed to create RDS backup policy: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse RDS backup policy create response: %s", err)
	}

	policy := helper.MapFromMap(payload, "policy")
	if policy == nil {
		return diag.Errorf("RDS backup policy create response missing policy")
	}
	policyID := helper.StringFromMap(policy, "id")
	if policyID == "" {
		return diag.Errorf("RDS backup policy create response missing policy id")
	}
	d.SetId(policyID)
	return resourceENRDSBackupPolicyRead(ctx, d, m)
}

func resourceENRDSBackupPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Id())
	if policyID == "" {
		return nil
	}

	req := map[string]interface{}{
		"id": policyID,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/detail", req, &resp); err != nil {
		return diag.Errorf("failed to query RDS backup policy detail: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse RDS backup policy detail response: %s", err)
	}

	policy := helper.MapFromMap(payload, "policy")
	if policy == nil {
		d.SetId("")
		return nil
	}

	attrs := rdsBackupPolicyAttrsFromMap(policy)
	for _, key := range []string{
		"name", "region_id", "status", "schedule_enabled", "cycle_type",
		"weekdays", "hours", "retention_type", "retention_value", "stop_type", "stop_value",
		"next_trigger_time_utc", "created", "schedule_description", "schedule_times",
	} {
		if err := d.Set(key, attrs[key]); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(policyID)
	return nil
}

func resourceENRDSBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	// API contract: when schedule_enabled changes, call /enabled endpoint.
	if d.HasChanges(
		"name",
		"cycle_type",
		"weekdays",
		"schedule_times",
		"retention_type",
		"retention_value",
		"stop_type",
		"stop_value",
	) {
		req := rdsBackupPolicyRequestFromResourceData(d, d.Id(), false)
		var resp map[string]interface{}
		if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/update", req, &resp); err != nil {
			return diag.Errorf("failed to update RDS backup policy: %s", err)
		}
		if _, err := helper.ParseAPIResponseMap(resp); err != nil {
			return diag.Errorf("failed to parse RDS backup policy update response: %s", err)
		}
	}

	if d.HasChange("schedule_enabled") {
		if diags := rdsBackupPolicySetEnabled(ctx, rdsClient, d.Id(), d.Get("schedule_enabled").(bool)); diags.HasError() {
			return diags
		}
	}

	return resourceENRDSBackupPolicyRead(ctx, d, m)
}

func resourceENRDSBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"ids": []string{d.Id()},
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/delete", req, &resp); err != nil {
		return diag.Errorf("failed to delete RDS backup policy: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS backup policy delete response: %s", err)
	}

	d.SetId("")
	return nil
}

func rdsBackupPolicyRequestFromResourceData(d *schema.ResourceData, policyID string, includeScheduleEnabled bool) map[string]interface{} {
	req := map[string]interface{}{
		"name":            d.Get("name").(string),
		"cycle_type":      d.Get("cycle_type").(string),
		"weekdays":        intSetToSlice(d.Get("weekdays").(*schema.Set)),
		"schedule_times":  stringSetToSlice(d.Get("schedule_times").(*schema.Set)),
		"retention_type":  d.Get("retention_type").(string),
		"retention_value": d.Get("retention_value").(int),
		"stop_type":       d.Get("stop_type").(string),
		"stop_value":      d.Get("stop_value").(string),
	}
	if includeScheduleEnabled {
		req["schedule_enabled"] = d.Get("schedule_enabled").(bool)
	}
	if policyID != "" {
		req["id"] = policyID
	}
	return req
}

func rdsBackupPolicySetEnabled(ctx context.Context, rdsClient *connectivity.RDSClient, policyID string, enabled bool) diag.Diagnostics {
	scheduleEnabled := 0
	if enabled {
		scheduleEnabled = 1
	}
	req := map[string]interface{}{
		"id":               policyID,
		"schedule_enabled": scheduleEnabled,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/enabled", req, &resp); err != nil {
		return diag.Errorf("failed to update RDS backup policy schedule_enabled: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS backup policy enabled response: %s", err)
	}
	return nil
}

func intSetToSlice(set *schema.Set) []int {
	raw := set.List()
	out := make([]int, 0, len(raw))
	for _, v := range raw {
		if n, ok := v.(int); ok {
			out = append(out, n)
			continue
		}
		if n, ok := policyToInt(v); ok {
			out = append(out, n)
		}
	}
	sort.Ints(out)
	return out
}

func stringSetToSlice(set *schema.Set) []string {
	raw := set.List()
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
