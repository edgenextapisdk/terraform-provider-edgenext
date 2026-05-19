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

const (
	elbL7PolicyCreatePath = "/elb/openapi/v2/l7policies/create"
	elbL7PolicyGetPath    = "/elb/openapi/v2/l7policies/get"
	elbL7PolicyUpdatePath = "/elb/openapi/v2/l7policies/update"
	elbL7PolicyDeletePath = "/elb/openapi/v2/l7policies/delete"
)

// ResourceENELBL7Policy manages an EdgeNext ELB L7 policy (console RegisterL7PolicyHTTPServer).
func ResourceENELBL7Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENELBL7PolicyCreate,
		ReadContext:   resourceENELBL7PolicyRead,
		UpdateContext: resourceENELBL7PolicyUpdate,
		DeleteContext: resourceENELBL7PolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENELBL7PolicyImport,
		},
		Description: "Manages an EdgeNext ELB L7 policy on a listener. Use edgenext_elb_l7_rule for match rules.",
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Listener ID this L7 policy is attached to.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Policy name.",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "Policy description.",
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy action: REDIRECT_TO_URL or REDIRECT_TO_TARGET_GROUP.",
				ValidateFunc: validation.StringInSlice([]string{
					"REDIRECT_TO_URL", "REDIRECT_TO_TARGET_GROUP",
				}, false),
			},
			"redirect_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Redirect URL when action is REDIRECT_TO_URL.",
			},
			"redirect_target_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target group ID when action is REDIRECT_TO_TARGET_GROUP.",
			},
			"position": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Evaluation priority among L7 policies on the listener (computed).",
			},
			"redirect_target_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the redirect target group when present in read results.",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Administrative up/down state (computed).",
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

func resourceENELBL7PolicyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(d.Id())
	if id == "" {
		return nil, fmt.Errorf("import id must be the L7 policy id (non-empty)")
	}
	d.SetId(id)
	diags := resourceENELBL7PolicyRead(ctx, d, m)
	if diags.HasError() {
		return nil, elbListenerDiagError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENELBL7PolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	action := strings.TrimSpace(d.Get("action").(string))
	redirectURL := strings.TrimSpace(d.Get("redirect_url").(string))
	redirectTGID := strings.TrimSpace(d.Get("redirect_target_group_id").(string))
	if action == "REDIRECT_TO_URL" && redirectURL == "" {
		return diag.Errorf("redirect_url is required when action is REDIRECT_TO_URL")
	}
	if action == "REDIRECT_TO_TARGET_GROUP" && redirectTGID == "" {
		return diag.Errorf("redirect_target_group_id is required when action is REDIRECT_TO_TARGET_GROUP")
	}

	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"listener_id": strings.TrimSpace(d.Get("listener_id").(string)),
		"name":        strings.TrimSpace(d.Get("name").(string)),
		"action":      l7PolicyActionForAPI(action),
	}
	if desc := strings.TrimSpace(d.Get("description").(string)); desc != "" {
		req["description"] = desc
	}
	if redirectURL != "" {
		req["redirect_url"] = redirectURL
	}
	if redirectTGID != "" {
		req["redirect_pool_id"] = redirectTGID
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7PolicyCreatePath, req, &resp); err != nil {
		return diag.Errorf("failed to create ELB L7 policy: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB L7 policy create response: %s", err)
	}
	pol := helper.MapFromMap(payload, "l7policy")
	if pol == nil {
		pol = payload
	}
	policyID := helper.StringFromMap(pol, "id")
	if policyID == "" {
		return diag.Errorf("create response missing L7 policy id")
	}
	d.SetId(policyID)
	return resourceENELBL7PolicyRead(ctx, d, m)
}

func resourceENELBL7PolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Id())
	if policyID == "" {
		d.SetId("")
		return nil
	}

	req := map[string]interface{}{
		"l7policy_id": policyID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7PolicyGetPath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read ELB L7 policy: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB L7 policy get response: %s", err)
	}
	pol := helper.MapFromMap(payload, "l7policy")
	if pol == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("listener_id", helper.StringFromMap(pol, "listener_id"))
	_ = d.Set("name", helper.StringFromMap(pol, "name"))
	_ = d.Set("description", helper.StringFromMap(pol, "description"))
	_ = d.Set("action", l7PolicyActionFromAPI(helper.StringFromMap(pol, "action")))
	_ = d.Set("redirect_url", helper.StringFromMap(pol, "redirect_url"))
	_ = d.Set("redirect_target_group_id", l7PolicyRedirectTargetGroupIDFromPayload(pol))
	_ = d.Set("redirect_target_group_name", l7PolicyRedirectTargetGroupNameFromPayload(pol))
	_ = d.Set("position", helper.IntFromMap(pol, "position"))
	if v, ok := pol["admin_state_up"]; ok && v != nil {
		_ = d.Set("admin_state_up", v == true || v == "true" || fmt.Sprintf("%v", v) == "1")
	} else {
		_ = d.Set("admin_state_up", false)
	}
	_ = d.Set("operating_status", helper.StringFromMap(pol, "operating_status"))
	_ = d.Set("provisioning_status", helper.StringFromMap(pol, "provisioning_status"))
	_ = d.Set("project_id", helper.StringFromMap(pol, "project_id"))
	_ = d.Set("created_at", helper.IntFromMap(pol, "created_at"))
	_ = d.Set("updated_at", helper.IntFromMap(pol, "updated_at"))

	d.SetId(helper.StringFromMap(pol, "id"))
	return nil
}

func resourceENELBL7PolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("name") && !d.HasChange("description") && !d.HasChange("action") &&
		!d.HasChange("redirect_url") && !d.HasChange("redirect_target_group_id") {
		return nil
	}

	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"l7policy_id": strings.TrimSpace(d.Id()),
	}
	if d.HasChange("name") {
		req["name"] = strings.TrimSpace(d.Get("name").(string))
	}
	if d.HasChange("description") {
		req["description"] = d.Get("description").(string)
	}
	if d.HasChange("action") {
		req["action"] = l7PolicyActionForAPI(strings.TrimSpace(d.Get("action").(string)))
	}
	if d.HasChange("redirect_url") {
		if u := strings.TrimSpace(d.Get("redirect_url").(string)); u != "" {
			req["redirect_url"] = u
		}
	}
	if d.HasChange("redirect_target_group_id") {
		if p := strings.TrimSpace(d.Get("redirect_target_group_id").(string)); p != "" {
			req["redirect_pool_id"] = p
		}
	}

	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7PolicyUpdatePath, req, &resp); err != nil {
		return diag.Errorf("failed to update ELB L7 policy: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse ELB L7 policy update response: %s", err)
	}

	return resourceENELBL7PolicyRead(ctx, d, m)
}

func resourceENELBL7PolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Id())
	if policyID == "" {
		return nil
	}

	req := map[string]interface{}{
		"l7policy_id": policyID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbL7PolicyDeletePath, req, &resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to delete ELB L7 policy: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		if elbListenerAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB L7 policy delete response: %s", err)
	}
	d.SetId("")
	return nil
}

// l7PolicyActionForAPI maps Terraform action values to the EdgeNext/Octavia API action string.
func l7PolicyActionForAPI(action string) string {
	if action == "REDIRECT_TO_TARGET_GROUP" {
		return "REDIRECT_TO_POOL"
	}
	return action
}

// l7PolicyActionFromAPI maps API action strings to Terraform schema values.
func l7PolicyActionFromAPI(action string) string {
	switch strings.ToUpper(strings.TrimSpace(action)) {
	case "REDIRECT_TO_POOL":
		return "REDIRECT_TO_TARGET_GROUP"
	default:
		return strings.TrimSpace(action)
	}
}

func l7PolicyRedirectTargetGroupIDFromPayload(pol map[string]interface{}) string {
	if id := helper.StringFromMap(pol, "redirect_target_group_id"); id != "" {
		return id
	}
	return helper.StringFromMap(pol, "redirect_pool_id")
}

func l7PolicyRedirectTargetGroupNameFromPayload(pol map[string]interface{}) string {
	if name := helper.StringFromMap(pol, "redirect_target_group_name"); name != "" {
		return name
	}
	return helper.StringFromMap(pol, "redirect_pool_name")
}
