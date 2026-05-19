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
	elbListenerCreatePath = "/elb/openapi/v2/listeners/create"
	elbListenerUpdatePath = "/elb/openapi/v2/listeners/update"
	elbListenerDeletePath = "/elb/openapi/v2/listeners/delete"
	elbListenerGetPath    = "/elb/openapi/v2/listeners/get"
)

// ResourceENELBListener manages an EdgeNext ELB listener.
func ResourceENELBListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENELBListenerCreate,
		ReadContext:   resourceENELBListenerRead,
		UpdateContext: resourceENELBListenerUpdate,
		DeleteContext: resourceENELBListenerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENELBListenerImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		Description: "Manages an EdgeNext ELB listener (create, read, update name/description/default_tls_certificate_ref, delete).",
		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Load balancer ID this listener belongs to.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Listener name.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Listener description (updatable).",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Listener protocol (for example TERMINATED_HTTPS, HTTP).",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"protocol_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "Protocol port.",
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     -1,
				ForceNew:    true,
				Description: "Connection limit (-1 for unlimited). Create only.",
			},
			"insert_headers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Description: "Insert headers for the listener (string keys and values, for example X-Forwarded-For, X-Forwarded-Port, and X-Forwarded-Proto set to true). " +
					"Create only.",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"default_tls_certificate_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default TLS certificate (Barbican container) reference URL (for example TERMINATED_HTTPS). Updatable via listeners/update.",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Administrative up/down state from the API.",
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
			"default_target_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default target group ID when set.",
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

func resourceENELBListenerImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(d.Id())
	if id == "" {
		return nil, fmt.Errorf("import id must be the listener id (non-empty)")
	}
	d.SetId(id)
	diags := resourceENELBListenerRead(ctx, d, m)
	if diags.HasError() {
		return nil, elbListenerDiagError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func elbListenerDiagError(diags diag.Diagnostics) error {
	if len(diags) == 0 {
		return fmt.Errorf("operation failed")
	}
	e := diags[0]
	if e.Detail != "" {
		return fmt.Errorf("%s: %s", e.Summary, e.Detail)
	}
	return fmt.Errorf("%s", e.Summary)
}

func resourceENELBListenerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"loadbalancer_id":  strings.TrimSpace(d.Get("loadbalancer_id").(string)),
		"name":             strings.TrimSpace(d.Get("name").(string)),
		"protocol":         strings.TrimSpace(d.Get("protocol").(string)),
		"protocol_port":    d.Get("protocol_port").(int),
		"connection_limit": d.Get("connection_limit").(int),
	}
	if desc := strings.TrimSpace(d.Get("description").(string)); desc != "" {
		req["description"] = desc
	}
	if ih := listenerInsertHeadersForCreate(d); len(ih) > 0 {
		req["insert_headers"] = ih
	}
	if ref := strings.TrimSpace(d.Get("default_tls_certificate_ref").(string)); ref != "" {
		req["default_tls_container_ref"] = ref
	}

	var resp map[string]interface{}
	if err := elbPOSTWithOctaviaImmutableRetry(ctx, elbClient, elbListenerCreatePath, req, &resp); err != nil {
		return diag.Errorf("failed to create ELB listener: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB listener create response: %s", err)
	}
	listenerID := helper.StringFromMap(payload, "id")
	if listenerID == "" {
		return diag.Errorf("create response missing listener id")
	}
	d.SetId(listenerID)
	return resourceENELBListenerRead(ctx, d, m)
}

func listenerInsertHeadersForCreate(d *schema.ResourceData) map[string]interface{} {
	v := d.Get("insert_headers")
	if v == nil {
		return nil
	}
	raw, ok := v.(map[string]interface{})
	if !ok || len(raw) == 0 {
		return nil
	}
	out := make(map[string]interface{}, len(raw))
	for k, v := range raw {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

func resourceENELBListenerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	listenerID := strings.TrimSpace(d.Id())
	if listenerID == "" {
		d.SetId("")
		return nil
	}

	req := map[string]interface{}{
		"listener_id": listenerID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbListenerGetPath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read ELB listener: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB listener get response: %s", err)
	}

	_ = d.Set("loadbalancer_id", helper.StringFromMap(payload, "loadbalancer_id"))
	_ = d.Set("name", helper.StringFromMap(payload, "name"))
	_ = d.Set("description", helper.StringFromMap(payload, "description"))
	_ = d.Set("protocol", helper.StringFromMap(payload, "protocol"))
	_ = d.Set("protocol_port", helper.IntFromMap(payload, "protocol_port"))
	_ = d.Set("connection_limit", helper.IntFromMap(payload, "connection_limit"))
	_ = d.Set("admin_state_up", boolFromInterface(payload["admin_state_up"]))
	_ = d.Set("provisioning_status", helper.StringFromMap(payload, "provisioning_status"))
	_ = d.Set("operating_status", helper.StringFromMap(payload, "operating_status"))
	_ = d.Set("default_target_group_id", helper.StringFromMap(payload, "default_pool_id"))
	_ = d.Set("created_at", helper.IntFromMap(payload, "created_at"))
	_ = d.Set("updated_at", helper.IntFromMap(payload, "updated_at"))
	_ = d.Set("insert_headers", elbInsertHeadersStringMap(payload["insert_headers"]))
	if ref := helper.StringFromMap(payload, "default_tls_container_ref"); ref != "" {
		_ = d.Set("default_tls_certificate_ref", ref)
	}
	d.SetId(helper.StringFromMap(payload, "id"))
	return nil
}

func resourceENELBListenerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("name") && !d.HasChange("description") && !d.HasChange("default_tls_certificate_ref") {
		return nil
	}

	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	ref := strings.TrimSpace(d.Get("default_tls_certificate_ref").(string))
	req := map[string]interface{}{
		"listener_id": strings.TrimSpace(d.Id()),
		"name":        strings.TrimSpace(d.Get("name").(string)),
		"description": d.Get("description").(string),
	}
	if d.HasChange("default_tls_certificate_ref") {
		req["default_tls_container_ref"] = ref
	} else if ref != "" {
		req["default_tls_container_ref"] = ref
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbListenerUpdatePath, req, &resp); err != nil {
		return diag.Errorf("failed to update ELB listener: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse ELB listener update response: %s", err)
	}
	return resourceENELBListenerRead(ctx, d, m)
}

func resourceENELBListenerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	listenerID := strings.TrimSpace(d.Id())
	if listenerID == "" {
		return nil
	}

	req := map[string]interface{}{
		"listener_id": listenerID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbListenerDeletePath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to delete ELB listener: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB listener delete response: %s", err)
	}
	if err := elbWaitUntilListenerDeleted(ctx, elbClient, listenerID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func elbWaitUntilListenerDeleted(ctx context.Context, elbClient *connectivity.ELBClient, listenerID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	const poll = 2 * time.Second
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout after %v waiting for listener %q to finish deleting (load balancer may still be updating); increase timeouts { delete = \"30m\" } on edgenext_elb_listener or re-run terraform apply", timeout, listenerID)
		}
		req := map[string]interface{}{
			"listener_id": listenerID,
		}
		var resp map[string]interface{}
		err := elbClient.Post(ctx, elbListenerGetPath, req, &resp)
		if err != nil {
			if elbListenerAPIGone(err) {
				return nil
			}
		} else {
			payload, perr := helper.ParseAPIResponseMap(resp)
			if perr != nil {
				if elbListenerAPIGone(perr) {
					return nil
				}
			} else {
				id := strings.TrimSpace(helper.StringFromMap(payload, "id"))
				if id == "" {
					return nil
				}
				if id != listenerID {
					return nil
				}
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(poll):
		}
	}
}

func elbListenerAPIGone(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "404") ||
		strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "not exist")
}
