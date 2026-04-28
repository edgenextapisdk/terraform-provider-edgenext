package ecs

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

// ResourceENECSInstancePower manages ECS instance power state.
func ResourceENECSInstancePower() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSInstancePowerCreate,
		ReadContext:   resourceENECSInstancePowerRead,
		UpdateContext: resourceENECSInstancePowerUpdate,
		DeleteContext: resourceENECSInstancePowerDelete,
		Description:   "Provides an EdgeNext ECS instance power control resource.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance ID.",
			},
			"desired_state": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Desired instance power state. Valid values: ACTIVE, SHUTOFF.",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "SHUTOFF"}, false),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current instance status from detail API.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance name.",
			},
		},
	}
}

func resourceENECSInstancePowerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId(strings.TrimSpace(d.Get("instance_id").(string)))
	return resourceENECSInstancePowerEnsureState(ctx, d, m)
}

func resourceENECSInstancePowerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	server, err := resourceENECSInstancePowerDetail(ctx, ecsClient, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("status", helper.StringFromMap(server, "status"))
	_ = d.Set("instance_name", helper.StringFromMap(server, "name"))
	return nil
}

func resourceENECSInstancePowerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("desired_state") {
		return resourceENECSInstancePowerRead(ctx, d, m)
	}
	return resourceENECSInstancePowerEnsureState(ctx, d, m)
}

func resourceENECSInstancePowerDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Deleting this resource only removes Terraform state.
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "State-only delete for instance power resource",
			Detail:   fmt.Sprintf("Deleting edgenext_ecs_instance_power for instance %q only removes Terraform state and does not send a stop action.", d.Id()),
		},
	}
}

func resourceENECSInstancePowerEnsureState(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Id()
	desiredState := strings.ToUpper(strings.TrimSpace(d.Get("desired_state").(string)))

	server, err := resourceENECSInstancePowerDetail(ctx, ecsClient, instanceID)
	if err != nil {
		return diag.Errorf("failed to query ECS instance detail before power action: %s", err)
	}
	currentStatus := strings.ToUpper(strings.TrimSpace(helper.StringFromMap(server, "status")))
	if resourceENECSInstancePowerMatchesDesired(currentStatus, desiredState) {
		return resourceENECSInstancePowerRead(ctx, d, m)
	}

	action := "start"
	if desiredState == "SHUTOFF" {
		action = "stop"
	}
	if err := resourceENECSInstancePowerAction(ctx, ecsClient, instanceID, action); err != nil {
		return diag.FromErr(err)
	}

	const retries = 24
	for i := 0; i < retries; i++ {
		time.Sleep(5 * time.Second)
		server, err = resourceENECSInstancePowerDetail(ctx, ecsClient, instanceID)
		if err != nil {
			continue
		}
		currentStatus = strings.ToUpper(strings.TrimSpace(helper.StringFromMap(server, "status")))
		if resourceENECSInstancePowerMatchesDesired(currentStatus, desiredState) {
			return resourceENECSInstancePowerRead(ctx, d, m)
		}
	}
	return diag.Errorf("instance %q did not reach desired_state %q within timeout", instanceID, desiredState)
}

func resourceENECSInstancePowerAction(ctx context.Context, ecsClient *connectivity.ECSClient, instanceID, action string) error {
	req := map[string]interface{}{
		"server_ids": []string{instanceID},
		"action":     action,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/instance/action", req, &resp); err != nil {
		return err
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return err
	}
	if result := helper.StringFromMap(payload, instanceID); !strings.EqualFold(result, "ok") {
		return fmt.Errorf("instance action result for %s: %s", instanceID, result)
	}
	return nil
}

func resourceENECSInstancePowerDetail(ctx context.Context, ecsClient *connectivity.ECSClient, instanceID string) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"server_id": instanceID,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/instance/detail", req, &resp); err != nil {
		return nil, err
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, err
	}
	for _, raw := range helper.ListFromMap(payload, "servers") {
		server, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if strings.TrimSpace(helper.StringFromMap(server, "id")) == instanceID {
			return server, nil
		}
	}
	return nil, fmt.Errorf("instance %q not found", instanceID)
}

func resourceENECSInstancePowerMatchesDesired(currentStatus, desiredState string) bool {
	current := strings.ToUpper(strings.TrimSpace(currentStatus))
	desired := strings.ToUpper(strings.TrimSpace(desiredState))
	if desired == "ACTIVE" {
		return current == "ACTIVE" || current == "RUNNING"
	}
	return current == desired
}
