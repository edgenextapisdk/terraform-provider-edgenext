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
	elbTargetGroupCreatePath   = "/elb/openapi/v2/pools"
	elbTargetGroupUpdatePath   = "/elb/openapi/v2/pools/update"
	elbTargetGroupGetPath      = "/elb/openapi/v2/pools/get"
	elbTargetGroupDeletePath   = "/elb/openapi/v2/pools/delete"
	elbHealthMonitorCreatePath = "/elb/openapi/v2/healthmonitors"
	elbHealthMonitorGetPath    = "/elb/openapi/v2/healthmonitors/get"
)

// ResourceENELBTargetGroup manages an EdgeNext ELB target group with health checks.
// Backend targets are managed with edgenext_elb_target_group_attachment.
func ResourceENELBTargetGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENELBTargetGroupCreate,
		ReadContext:   resourceENELBTargetGroupRead,
		UpdateContext: resourceENELBTargetGroupUpdate,
		DeleteContext: resourceENELBTargetGroupDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Description: "Manages an EdgeNext ELB target group attached to a listener, including its health check configuration. " +
			"Add or remove targets with edgenext_elb_target_group_attachment. " +
			"After create, only name and description can be updated; protocol, algorithm, listener, and health_monitor require replacement. " +
			"Replacement deletes the old target group asynchronously; if recreate fails, increase timeouts.delete and re-apply.",
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Listener this target group is attached to.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Target group name.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Target group description (updatable).",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Traffic protocol between the load balancer and backends (for example HTTP).",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"lb_algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Load balancing algorithm (for example ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP).",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"health_monitor": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    true,
				Description: "Health check for this target group. Single block; changing it forces replacement.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Health monitor name.",
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Health monitor type (for example HTTP, HTTPS, TCP, PING, TLS-HELLO).",
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"max_retries": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Maximum retries before marking the target down.",
							ValidateFunc: validation.IntAtLeast(0),
						},
						"delay": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Seconds between health checks.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"timeout": {
							Type:         schema.TypeInt,
							Required:     true,
							Description:  "Health check timeout in seconds.",
							ValidateFunc: validation.IntAtLeast(1),
						},
						"http_method": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "HTTP method for HTTP(S) checks.",
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"url_path": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "URL path for HTTP(S) checks.",
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"expected_codes": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Expected HTTP status codes (for example 200).",
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health monitor identifier (computed).",
						},
						"provisioning_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health monitor provisioning status.",
						},
						"operating_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health monitor operating status.",
						},
						"created_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health monitor creation time as Unix timestamp (seconds).",
						},
						"updated_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health monitor last update time as Unix timestamp (seconds).",
						},
					},
				},
			},
			"loadbalancer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer this target group belongs to (computed).",
			},
			"provisioning_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target group provisioning status.",
			},
			"operating_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target group operating status.",
			},
			"healthmonitor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Health monitor identifier for this target group (computed).",
			},
			"target_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Identifiers of targets (backends) attached to this target group (computed). Use edgenext_elb_target_group_attachment to add or remove targets.",
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
	}
}

func resourceENELBTargetGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tgCreateBody := map[string]interface{}{
		"listener_id":  strings.TrimSpace(d.Get("listener_id").(string)),
		"name":         strings.TrimSpace(d.Get("name").(string)),
		"protocol":     strings.TrimSpace(d.Get("protocol").(string)),
		"lb_algorithm": strings.TrimSpace(d.Get("lb_algorithm").(string)),
		"description":  d.Get("description").(string),
	}
	reqCreate := map[string]interface{}{
		"pool": tgCreateBody,
	}
	var respCreate map[string]interface{}
	if err := elbPOSTWithOctaviaImmutableRetry(ctx, elbClient, elbTargetGroupCreatePath, reqCreate, &respCreate); err != nil {
		return diag.Errorf("failed to create ELB target group: %s", err)
	}
	payloadCreate, err := helper.ParseAPIResponseMap(respCreate)
	if err != nil {
		return diag.Errorf("failed to parse ELB target group create response: %s", err)
	}
	tgObj := helper.MapFromMap(payloadCreate, "pool")
	if tgObj == nil {
		return diag.Errorf("create response missing target group in server payload")
	}
	tgID := helper.StringFromMap(tgObj, "id")
	if tgID == "" {
		return diag.Errorf("create response missing target group id")
	}
	if err := elbWaitTargetGroupReadyForHealthMonitor(ctx, elbClient, tgID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}
	hmBody, diags := expandHealthMonitorForCreate(d)
	if diags.HasError() {
		return diags
	}
	hmBody["pool_id"] = tgID
	reqHM := map[string]interface{}{
		"healthmonitor": hmBody,
	}
	var respHM map[string]interface{}
	if err := elbPOSTWithOctaviaImmutableRetry(ctx, elbClient, elbHealthMonitorCreatePath, reqHM, &respHM); err != nil {
		return diag.Errorf("failed to create ELB health monitor: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(respHM); err != nil {
		return diag.Errorf("failed to parse ELB health monitor create response: %s", err)
	}

	d.SetId(tgID)
	return resourceENELBTargetGroupRead(ctx, d, m)
}

// elbWaitTargetGroupReadyForHealthMonitor polls until the target group reaches provisioning_status ACTIVE,
// so the health monitor can be created without Octavia rejecting the operation while provisioning is still in progress.
func elbWaitTargetGroupReadyForHealthMonitor(ctx context.Context, elbClient *connectivity.ELBClient, tgID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	const poll = 2 * time.Second
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout after %v waiting for target group %q to reach provisioning_status ACTIVE (required before creating health monitor); increase timeouts { create = \"30m\" } on edgenext_elb_target_group", timeout, tgID)
		}
		req := map[string]interface{}{
			"pool_id": tgID,
		}
		var resp map[string]interface{}
		err := elbClient.Post(ctx, elbTargetGroupGetPath, req, &resp)
		if err != nil {
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
		tgPayload := helper.MapFromMap(payload, "pool")
		if tgPayload == nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(poll):
			}
			continue
		}
		status := strings.ToUpper(strings.TrimSpace(helper.StringFromMap(tgPayload, "provisioning_status")))
		switch status {
		case "ACTIVE":
			return nil
		case "ERROR", "DELETED", "PENDING_DELETE":
			return fmt.Errorf("target group %q provisioning_status is %q, cannot create health monitor", tgID, status)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(poll):
		}
	}
}

func expandHealthMonitorForCreate(d *schema.ResourceData) (map[string]interface{}, diag.Diagnostics) {
	raw := d.Get("health_monitor").([]interface{})
	if len(raw) == 0 {
		return nil, diag.Errorf("health_monitor block is required")
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return nil, diag.Errorf("invalid health_monitor block")
	}
	return map[string]interface{}{
		"name":           strings.TrimSpace(m["name"].(string)),
		"type":           strings.TrimSpace(m["type"].(string)),
		"max_retries":    m["max_retries"].(int),
		"delay":          m["delay"].(int),
		"timeout":        m["timeout"].(int),
		"http_method":    strings.TrimSpace(m["http_method"].(string)),
		"url_path":       strings.TrimSpace(m["url_path"].(string)),
		"expected_codes": strings.TrimSpace(m["expected_codes"].(string)),
	}, nil
}

func firstRefIDFromList(row map[string]interface{}, key string) string {
	for _, raw := range helper.ListFromMap(row, key) {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if id := helper.StringFromMap(row, "id"); id != "" {
			return id
		}
	}
	return ""
}

func targetIDsFromTargetGroupPayload(tg map[string]interface{}) []interface{} {
	out := make([]interface{}, 0)
	for _, raw := range helper.ListFromMap(tg, "members") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if id := helper.StringFromMap(row, "id"); id != "" {
			out = append(out, id)
		}
	}
	return out
}

func resourceENELBTargetGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tgID := strings.TrimSpace(d.Id())
	if tgID == "" {
		d.SetId("")
		return nil
	}

	req := map[string]interface{}{
		"pool_id": tgID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbTargetGroupGetPath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read ELB target group: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB target group get response: %s", err)
	}
	tgPayload := helper.MapFromMap(payload, "pool")
	if tgPayload == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", helper.StringFromMap(tgPayload, "name"))
	_ = d.Set("description", helper.StringFromMap(tgPayload, "description"))
	_ = d.Set("protocol", helper.StringFromMap(tgPayload, "protocol"))
	_ = d.Set("lb_algorithm", helper.StringFromMap(tgPayload, "lb_algorithm"))
	_ = d.Set("provisioning_status", helper.StringFromMap(tgPayload, "provisioning_status"))
	_ = d.Set("operating_status", helper.StringFromMap(tgPayload, "operating_status"))
	_ = d.Set("healthmonitor_id", helper.StringFromMap(tgPayload, "healthmonitor_id"))
	_ = d.Set("created_at", helper.IntFromMap(tgPayload, "created_at"))
	_ = d.Set("updated_at", helper.IntFromMap(tgPayload, "updated_at"))

	if lid := firstRefIDFromList(tgPayload, "listeners"); lid != "" {
		_ = d.Set("listener_id", lid)
	}
	if lbid := firstRefIDFromList(tgPayload, "loadbalancers"); lbid != "" {
		_ = d.Set("loadbalancer_id", lbid)
	}

	targetIDs := targetIDsFromTargetGroupPayload(tgPayload)
	if err := d.Set("target_ids", targetIDs); err != nil {
		return diag.FromErr(err)
	}

	hmList, hmDiags := healthMonitorListForTargetGroupState(ctx, elbClient, tgPayload, false)
	if hmDiags.HasError() {
		return hmDiags
	}
	if err := d.Set("health_monitor", hmList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(helper.StringFromMap(tgPayload, "id"))
	return nil
}

// healthMonitorListForTargetGroupState returns a single-element list for schema health_monitor when details are known.
// If allowEmpty is false and no monitor can be resolved, returns a diagnostic (for managed resources).
func healthMonitorListForTargetGroupState(ctx context.Context, elbClient *connectivity.ELBClient, tg map[string]interface{}, allowEmpty bool) ([]interface{}, diag.Diagnostics) {
	hm, diags := resolveHealthMonitorMap(ctx, elbClient, tg)
	if diags.HasError() {
		return nil, diags
	}
	if hm == nil {
		if allowEmpty {
			return []interface{}{}, nil
		}
		return nil, diag.Errorf("no health monitor details for target group %q: missing healthmonitor_id and embedded health monitor on target group read response", helper.StringFromMap(tg, "id"))
	}
	return []interface{}{healthMonitorStateFromPayload(hm)}, nil
}

func resolveHealthMonitorMap(ctx context.Context, elbClient *connectivity.ELBClient, tg map[string]interface{}) (map[string]interface{}, diag.Diagnostics) {
	if nested := helper.MapFromMap(tg, "healthmonitor"); nested != nil {
		return nested, nil
	}
	if nested := helper.MapFromMap(tg, "health_monitor"); nested != nil {
		return nested, nil
	}
	hmID := strings.TrimSpace(helper.StringFromMap(tg, "healthmonitor_id"))
	if hmID == "" {
		return nil, nil
	}
	req := map[string]interface{}{
		"healthmonitor_id": hmID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbHealthMonitorGetPath, req, &resp); err != nil {
		return nil, diag.Errorf("failed to get health monitor %q: %s", hmID, err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, diag.Errorf("failed to parse health monitor get response: %s", err)
	}
	if hm := helper.MapFromMap(payload, "healthmonitor"); hm != nil {
		return hm, nil
	}
	return payload, nil
}

func healthMonitorStateFromPayload(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":                helper.StringFromMap(m, "name"),
		"type":                helper.StringFromMap(m, "type"),
		"max_retries":         helper.IntFromMap(m, "max_retries"),
		"delay":               helper.IntFromMap(m, "delay"),
		"timeout":             helper.IntFromMap(m, "timeout"),
		"http_method":         helper.StringFromMap(m, "http_method"),
		"url_path":            helper.StringFromMap(m, "url_path"),
		"expected_codes":      helper.StringFromMap(m, "expected_codes"),
		"id":                  helper.StringFromMap(m, "id"),
		"provisioning_status": helper.StringFromMap(m, "provisioning_status"),
		"operating_status":    helper.StringFromMap(m, "operating_status"),
		"created_at":          helper.IntFromMap(m, "created_at"),
		"updated_at":          helper.IntFromMap(m, "updated_at"),
	}
}

func resourceENELBTargetGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("name") && !d.HasChange("description") {
		return nil
	}

	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"pool_id": strings.TrimSpace(d.Id()),
		"pool": map[string]interface{}{
			"name":        strings.TrimSpace(d.Get("name").(string)),
			"description": d.Get("description").(string),
		},
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbTargetGroupUpdatePath, req, &resp); err != nil {
		return diag.Errorf("failed to update ELB target group: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse ELB target group update response: %s", err)
	}
	return resourceENELBTargetGroupRead(ctx, d, m)
}

func resourceENELBTargetGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tgID := strings.TrimSpace(d.Id())
	if tgID == "" {
		return nil
	}

	req := map[string]interface{}{
		"pool_id": tgID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbTargetGroupDeletePath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to delete ELB target group: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB target group delete response: %s", err)
	}
	if err := elbWaitUntilTargetGroupDeleted(ctx, elbClient, tgID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// elbWaitUntilTargetGroupDeleted polls pools/get until the target group is gone (same pattern as listener delete).
// Octavia finishes pool delete asynchronously; without this wait, a ForceNew recreate can hit immutable/conflict on the listener or load balancer.
func elbWaitUntilTargetGroupDeleted(ctx context.Context, elbClient *connectivity.ELBClient, tgID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	const poll = 2 * time.Second
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout after %v waiting for target group %q to finish deleting (listener or load balancer may still be updating); increase timeouts { delete = \"30m\" } on edgenext_elb_target_group or re-run terraform apply", timeout, tgID)
		}
		req := map[string]interface{}{
			"pool_id": tgID,
		}
		var resp map[string]interface{}
		err := elbClient.Post(ctx, elbTargetGroupGetPath, req, &resp)
		if err != nil {
			if elbListenerAPIGone(err) {
				return nil
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
			if elbListenerAPIGone(perr) {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(poll):
			}
			continue
		}
		tgPayload := helper.MapFromMap(payload, "pool")
		if tgPayload == nil {
			return nil
		}
		id := strings.TrimSpace(helper.StringFromMap(tgPayload, "id"))
		if id == "" || id != tgID {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(poll):
		}
	}
}
