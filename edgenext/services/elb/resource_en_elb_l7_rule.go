package elb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	elbL7RuleCreatePath = "/elb/openapi/v2/l7rules/create"
	elbL7RuleUpdatePath = "/elb/openapi/v2/l7rules/update"
	elbL7RuleDeletePath = "/elb/openapi/v2/l7rules/delete"
	elbL7RuleListPath   = "/elb/openapi/v2/l7rules/list"
)

// ResourceENELBL7Rule manages a standalone L7 rule on an L7 policy (console RegisterL7RuleHTTPServer).
// There is no dedicated get API; read uses l7rules/list with l7policy_id and id filter.
func ResourceENELBL7Rule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENELBL7RuleCreate,
		ReadContext:   resourceENELBL7RuleRead,
		UpdateContext: resourceENELBL7RuleUpdate,
		DeleteContext: resourceENELBL7RuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENELBL7RuleImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Description: "Manages an EdgeNext ELB L7 rule under an L7 policy (POST .../l7rules/create, list/update/delete). " +
			"Read uses POST .../l7rules/list with id filter. Import id is l7policy_id/l7rule_id. " +
			"Deleting a rule is asynchronous on the load balancer; ForceNew replacement waits until the old rule is gone before create returns. " +
			"If replace fails with load balancer immutable, increase timeouts.delete and re-apply.",
		Schema: map[string]*schema.Schema{
			"l7policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "L7 policy ID this rule belongs to.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Rule type: HOSTNAME or PATH.",
				ValidateFunc: validation.StringInSlice([]string{"HOSTNAME", "PATH"}, false),
			},
			"compare_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Compare type (for example EQUAL_TO, REGEX, STARTS_WITH, ENDS_WITH, CONTAINS).",
				ValidateFunc: validation.StringInSlice([]string{
					"EQUAL_TO", "REGEX", "STARTS_WITH", "ENDS_WITH", "CONTAINS",
				}, false),
			},
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Match value (1-255 characters).",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"invert": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the rule match is inverted (computed).",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Administrative up/down state (computed).",
			},
			"operating_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating status from the API.",
			},
			"provisioning_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provisioning status from the API.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID from the API.",
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
}

func resourceENELBL7RuleImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	policyID, ruleID, err := parseL7RuleImportID(strings.TrimSpace(d.Id()))
	if err != nil {
		return nil, err
	}
	if err := d.Set("l7policy_id", policyID); err != nil {
		return nil, err
	}
	d.SetId(ruleID)
	diags := resourceENELBL7RuleRead(ctx, d, m)
	if diags.HasError() {
		return nil, elbListenerDiagError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func parseL7RuleImportID(id string) (policyID, ruleID string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return "", "", fmt.Errorf("import id must be l7policy_id/l7rule_id (two segments separated by '/')")
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func resourceENELBL7RuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"l7policy_id":  strings.TrimSpace(d.Get("l7policy_id").(string)),
		"type":         strings.TrimSpace(d.Get("type").(string)),
		"compare_type": strings.TrimSpace(d.Get("compare_type").(string)),
		"value":        strings.TrimSpace(d.Get("value").(string)),
	}

	policyID := strings.TrimSpace(d.Get("l7policy_id").(string))
	if err := elbWaitL7PolicyReadyForRule(ctx, elbClient, policyID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	var resp map[string]interface{}
	if err := elbPOSTWithOctaviaImmutableRetry(ctx, elbClient, elbL7RuleCreatePath, req, &resp); err != nil {
		return diag.Errorf("failed to create ELB L7 rule: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB L7 rule create response: %s", err)
	}
	rule := helper.MapFromMap(payload, "rule")
	if rule == nil {
		rule = payload
	}
	ruleID := helper.StringFromMap(rule, "id")
	if ruleID == "" {
		return diag.Errorf("create response missing L7 rule id")
	}
	d.SetId(ruleID)
	return resourceENELBL7RuleRead(ctx, d, m)
}

func resourceENELBL7RuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Get("l7policy_id").(string))
	ruleID := strings.TrimSpace(d.Id())
	if policyID == "" || ruleID == "" {
		d.SetId("")
		return nil
	}

	rule, diags := elbL7RuleFetchByPolicyAndID(ctx, elbClient, policyID, ruleID)
	if diags.HasError() {
		return diags
	}
	if rule == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("type", helper.StringFromMap(rule, "type"))
	_ = d.Set("compare_type", helper.StringFromMap(rule, "compare_type"))
	_ = d.Set("value", helper.StringFromMap(rule, "value"))
	if v, ok := rule["invert"]; ok && v != nil {
		_ = d.Set("invert", v == true || v == "true" || fmt.Sprintf("%v", v) == "1")
	} else {
		_ = d.Set("invert", false)
	}
	if v, ok := rule["admin_state_up"]; ok && v != nil {
		_ = d.Set("admin_state_up", v == true || v == "true" || fmt.Sprintf("%v", v) == "1")
	}
	_ = d.Set("operating_status", helper.StringFromMap(rule, "operating_status"))
	_ = d.Set("provisioning_status", helper.StringFromMap(rule, "provisioning_status"))
	_ = d.Set("project_id", helper.StringFromMap(rule, "project_id"))
	_ = d.Set("created_at", helper.IntFromMap(rule, "created_at"))
	_ = d.Set("updated_at", helper.IntFromMap(rule, "updated_at"))

	d.SetId(helper.StringFromMap(rule, "id"))
	return nil
}

// elbL7RuleFetchByPolicyAndID loads one rule via l7rules/list with id filter.
func elbL7RuleFetchByPolicyAndID(ctx context.Context, elbClient *connectivity.ELBClient, policyID, ruleID string) (map[string]interface{}, diag.Diagnostics) {
	req := map[string]interface{}{
		"l7policy_id": policyID,
		"id":          ruleID,
		"limit":       50,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7RuleListPath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			return nil, nil
		}
		return nil, diag.Errorf("failed to list ELB L7 rule %q: %s", ruleID, err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		if elbListenerAPIGone(err) {
			return nil, nil
		}
		return nil, diag.Errorf("failed to parse ELB L7 rules list response: %s", err)
	}
	for _, raw := range helper.ListFromMap(payload, "rules") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(row, "id") == ruleID {
			return row, nil
		}
	}
	return nil, nil
}

func resourceENELBL7RuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("type") && !d.HasChange("compare_type") && !d.HasChange("value") {
		return nil
	}

	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"l7policy_id": strings.TrimSpace(d.Get("l7policy_id").(string)),
		"l7rule_id":   strings.TrimSpace(d.Id()),
	}
	if d.HasChange("type") {
		req["type"] = strings.TrimSpace(d.Get("type").(string))
	}
	if d.HasChange("compare_type") {
		req["compare_type"] = strings.TrimSpace(d.Get("compare_type").(string))
	}
	if d.HasChange("value") {
		req["value"] = strings.TrimSpace(d.Get("value").(string))
	}

	var resp map[string]interface{}
	if err := elbPOSTWithOctaviaImmutableRetry(ctx, elbClient, elbL7RuleUpdatePath, req, &resp); err != nil {
		return diag.Errorf("failed to update ELB L7 rule: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse ELB L7 rule update response: %s", err)
	}
	return resourceENELBL7RuleRead(ctx, d, m)
}

func resourceENELBL7RuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Get("l7policy_id").(string))
	ruleID := strings.TrimSpace(d.Id())
	if policyID == "" || ruleID == "" {
		return nil
	}

	req := map[string]interface{}{
		"l7policy_id": policyID,
		"l7rule_id":   ruleID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7RuleDeletePath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to delete ELB L7 rule: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB L7 rule delete response: %s", err)
	}
	if err := elbWaitUntilL7RuleDeleted(ctx, elbClient, policyID, ruleID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// elbWaitL7PolicyReadyForRule polls until the L7 policy reaches provisioning_status ACTIVE before creating a rule.
func elbWaitL7PolicyReadyForRule(ctx context.Context, elbClient *connectivity.ELBClient, policyID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	const poll = 2 * time.Second
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout after %v waiting for L7 policy %q to reach provisioning_status ACTIVE (required before creating a rule); increase timeouts { create = \"30m\" } on edgenext_elb_l7_rule", timeout, policyID)
		}
		req := map[string]interface{}{
			"l7policy_id": policyID,
		}
		var resp map[string]interface{}
		err := elbClient.Post(ctx, elbL7PolicyGetPath, req, &resp)
		if err != nil {
			if elbListenerAPIGone(err) {
				return fmt.Errorf("L7 policy %q not found while waiting to create rule", policyID)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(poll):
			}
			continue
		}
		payload, perr := helper.ParseAPIResponseMap(resp)
		if perr != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(poll):
			}
			continue
		}
		pol := helper.MapFromMap(payload, "l7policy")
		if pol == nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(poll):
			}
			continue
		}
		status := strings.ToUpper(strings.TrimSpace(helper.StringFromMap(pol, "provisioning_status")))
		switch status {
		case "ACTIVE":
			return nil
		case "ERROR", "DELETED", "PENDING_DELETE":
			return fmt.Errorf("L7 policy %q provisioning_status is %q, cannot create rule", policyID, status)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(poll):
		}
	}
}

// elbWaitUntilL7RuleDeleted polls l7rules/list until the rule is gone so a following create does not hit Octavia load balancer immutable.
func elbWaitUntilL7RuleDeleted(ctx context.Context, elbClient *connectivity.ELBClient, policyID, ruleID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	const poll = 2 * time.Second
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout after %v waiting for L7 rule %q to finish deleting on policy %q (load balancer may still be updating); increase timeouts { delete = \"30m\" } on edgenext_elb_l7_rule or re-run terraform apply", timeout, ruleID, policyID)
		}
		rule, diags := elbL7RuleFetchByPolicyAndID(ctx, elbClient, policyID, ruleID)
		if diags.HasError() {
			return fmt.Errorf("%s", diagFlat(diags))
		}
		if rule == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(poll):
		}
	}
}
