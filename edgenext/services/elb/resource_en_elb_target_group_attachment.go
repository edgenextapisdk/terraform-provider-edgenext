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

// ELB OpenAPI paths for listing, creating, reading, updating, and deleting targets on a target group.
const (
	elbTargetGroupTargetListPath   = "/elb/openapi/v2/pools/members/list"
	elbTargetGroupTargetCreatePath = "/elb/openapi/v2/pools/members"
	elbTargetGroupTargetGetPath    = "/elb/openapi/v2/pools/members/get"
	elbTargetGroupTargetUpdatePath = "/elb/openapi/v2/pools/members/update"
	elbTargetGroupTargetDeletePath = "/elb/openapi/v2/pools/members/delete"
)

// ResourceENELBTargetGroupAttachment manages a single target (backend) on an ELB target group.
func ResourceENELBTargetGroupAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENELBTargetGroupAttachmentCreate,
		ReadContext:   resourceENELBTargetGroupAttachmentRead,
		UpdateContext: resourceENELBTargetGroupAttachmentUpdate,
		DeleteContext: resourceENELBTargetGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENELBTargetGroupAttachmentImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Description: "Manages one target (backend) on an EdgeNext ELB target group (create, read, update name, weight, protocol_port, delete). " +
			"Deleting a target is asynchronous on the load balancer; ForceNew replacement waits until the old target is gone before create returns. " +
			"If replace fails with pool immutable, increase timeouts.delete and re-apply.",
		Schema: map[string]*schema.Schema{
			"target_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Target group this target belongs to (same id as edgenext_elb_target_group).",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "IP address of the target (backend).",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"protocol_port": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Backend protocol port (updatable).",
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "Display name for the target (optional on create; updatable).",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				Description:  "Target weight.",
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"provisioning_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Provisioning status from the API.",
			},
			"operating_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating status from the API.",
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

func resourceENELBTargetGroupAttachmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(d.Id())
	if id == "" {
		return nil, fmt.Errorf("import id must be <target_group_id>/<target_id>")
	}
	tgID, targetID, err := parseTargetGroupTargetImportID(id)
	if err != nil {
		return nil, err
	}
	if err := d.Set("target_group_id", tgID); err != nil {
		return nil, err
	}
	d.SetId(targetID)
	diags := resourceENELBTargetGroupAttachmentRead(ctx, d, m)
	if diags.HasError() {
		return nil, fmt.Errorf("%s", diagFlat(diags))
	}
	return []*schema.ResourceData{d}, nil
}

func parseTargetGroupTargetImportID(id string) (tgID, targetID string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return "", "", fmt.Errorf("import id must be <target_group_id>/<target_id> (two segments separated by '/')")
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func diagFlat(diags diag.Diagnostics) string {
	if len(diags) == 0 {
		return "operation failed"
	}
	e := diags[0]
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Summary, e.Detail)
	}
	return e.Summary
}

func resourceENELBTargetGroupAttachmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tgID := strings.TrimSpace(d.Get("target_group_id").(string))
	targetCreateBody := map[string]interface{}{
		"address":       strings.TrimSpace(d.Get("address").(string)),
		"protocol_port": d.Get("protocol_port").(int),
		"weight":        d.Get("weight").(int),
	}
	if n := strings.TrimSpace(d.Get("name").(string)); n != "" {
		targetCreateBody["name"] = n
	}

	req := map[string]interface{}{
		"pool_id": tgID,
		"member":  targetCreateBody,
	}

	var resp map[string]interface{}
	if err := elbPOSTWithOctaviaImmutableRetry(ctx, elbClient, elbTargetGroupTargetCreatePath, req, &resp); err != nil {
		return diag.Errorf("failed to create ELB target: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB target create response: %s", err)
	}
	targetPayload := helper.MapFromMap(payload, "member")
	if targetPayload == nil {
		targetPayload = payload
	}
	targetID := helper.StringFromMap(targetPayload, "id")
	if targetID == "" {
		return diag.Errorf("create response missing target id")
	}
	d.SetId(targetID)
	return resourceENELBTargetGroupAttachmentRead(ctx, d, m)
}

func resourceENELBTargetGroupAttachmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tgID := strings.TrimSpace(d.Get("target_group_id").(string))
	targetID := strings.TrimSpace(d.Id())
	if tgID == "" || targetID == "" {
		d.SetId("")
		return nil
	}

	req := map[string]interface{}{
		"pool_id":   tgID,
		"member_id": targetID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbTargetGroupTargetGetPath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read ELB target: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB target get response: %s", err)
	}
	targetPayload := helper.MapFromMap(payload, "member")
	if targetPayload == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("address", helper.StringFromMap(targetPayload, "address"))
	_ = d.Set("protocol_port", helper.IntFromMap(targetPayload, "protocol_port"))
	_ = d.Set("name", helper.StringFromMap(targetPayload, "name"))
	_ = d.Set("weight", helper.IntFromMap(targetPayload, "weight"))
	_ = d.Set("provisioning_status", helper.StringFromMap(targetPayload, "provisioning_status"))
	_ = d.Set("operating_status", helper.StringFromMap(targetPayload, "operating_status"))
	_ = d.Set("created_at", helper.IntFromMap(targetPayload, "created_at"))
	_ = d.Set("updated_at", helper.IntFromMap(targetPayload, "updated_at"))
	_ = d.Set("target_group_id", tgID)

	d.SetId(helper.StringFromMap(targetPayload, "id"))
	return nil
}

func resourceENELBTargetGroupAttachmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("name") && !d.HasChange("weight") && !d.HasChange("protocol_port") {
		return nil
	}

	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	targetUpdateBody := map[string]interface{}{
		"weight":        d.Get("weight").(int),
		"protocol_port": d.Get("protocol_port").(int),
	}
	if n := strings.TrimSpace(d.Get("name").(string)); n != "" {
		targetUpdateBody["name"] = n
	} else {
		targetUpdateBody["name"] = ""
	}

	req := map[string]interface{}{
		"pool_id":   strings.TrimSpace(d.Get("target_group_id").(string)),
		"member_id": strings.TrimSpace(d.Id()),
		"member":    targetUpdateBody,
	}

	var resp map[string]interface{}
	if err := elbPOSTWithOctaviaImmutableRetry(ctx, elbClient, elbTargetGroupTargetUpdatePath, req, &resp); err != nil {
		return diag.Errorf("failed to update ELB target: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse ELB target update response: %s", err)
	}
	return resourceENELBTargetGroupAttachmentRead(ctx, d, m)
}

func resourceENELBTargetGroupAttachmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tgID := strings.TrimSpace(d.Get("target_group_id").(string))
	targetID := strings.TrimSpace(d.Id())
	if tgID == "" || targetID == "" {
		return nil
	}

	req := map[string]interface{}{
		"pool_id":   tgID,
		"member_id": targetID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbTargetGroupTargetDeletePath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to delete ELB target: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB target delete response: %s", err)
	}
	if err := elbWaitUntilTargetAttachmentDeleted(ctx, elbClient, tgID, targetID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

// elbWaitUntilTargetAttachmentDeleted polls pools/members/get until the member is gone so a following create does not hit Octavia pool immutable.
func elbWaitUntilTargetAttachmentDeleted(ctx context.Context, elbClient *connectivity.ELBClient, tgID, targetID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	const poll = 2 * time.Second
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout after %v waiting for target %q to finish deleting on target group %q (pool may still be updating); increase timeouts { delete = \"30m\" } on edgenext_elb_target_group_attachment or re-run terraform apply", timeout, targetID, tgID)
		}
		req := map[string]interface{}{
			"pool_id":   tgID,
			"member_id": targetID,
		}
		var resp map[string]interface{}
		err := elbClient.Post(ctx, elbTargetGroupTargetGetPath, req, &resp)
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
		member := helper.MapFromMap(payload, "member")
		if member == nil {
			return nil
		}
		id := strings.TrimSpace(helper.StringFromMap(member, "id"))
		if id == "" || id != targetID {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(poll):
		}
	}
}
