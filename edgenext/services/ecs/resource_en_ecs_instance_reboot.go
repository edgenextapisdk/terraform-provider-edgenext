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

// ResourceENECSInstanceReboot executes reboot action for an existing ECS instance.
func ResourceENECSInstanceReboot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSInstanceRebootCreate,
		ReadContext:   resourceENECSInstanceRebootRead,
		UpdateContext: resourceENECSInstanceRebootUpdate,
		DeleteContext: resourceENECSInstanceRebootDelete,
		Description:   "Provides an EdgeNext ECS instance reboot action resource.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance ID.",
			},
			"reboot_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "reboot_soft",
				Description:  "Reboot action type. Currently only reboot_soft is supported.",
				ValidateFunc: validation.StringInSlice([]string{"reboot_soft"}, false),
			},
			"trigger": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Update this field to trigger reboot again.",
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

func resourceENECSInstanceRebootCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId(strings.TrimSpace(d.Get("instance_id").(string)))
	return resourceENECSInstanceRebootActionAndWait(ctx, d, m)
}

func resourceENECSInstanceRebootRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceENECSInstanceRebootUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("trigger") && !d.HasChange("reboot_type") {
		return resourceENECSInstanceRebootRead(ctx, d, m)
	}
	return resourceENECSInstanceRebootActionAndWait(ctx, d, m)
}

func resourceENECSInstanceRebootDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Deleting this resource only removes Terraform state.
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "State-only delete for instance reboot resource",
			Detail:   fmt.Sprintf("Deleting edgenext_ecs_instance_reboot for instance %q only removes Terraform state and does not send any reboot action.", d.Id()),
		},
	}
}

func resourceENECSInstanceRebootActionAndWait(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Id()
	rebootType := strings.TrimSpace(d.Get("reboot_type").(string))
	if err := resourceENECSInstancePowerAction(ctx, ecsClient, instanceID, rebootType); err != nil {
		return diag.FromErr(err)
	}

	// During reboot, detail status may become REBOOT first, then ACTIVE.
	// Wait until the instance returns to a running state.
	const retries = 24
	for i := 0; i < retries; i++ {
		time.Sleep(5 * time.Second)
		server, err := resourceENECSInstancePowerDetail(ctx, ecsClient, instanceID)
		if err != nil {
			continue
		}
		status := strings.ToUpper(strings.TrimSpace(helper.StringFromMap(server, "status")))
		if status == "ACTIVE" || status == "RUNNING" {
			return resourceENECSInstanceRebootRead(ctx, d, m)
		}
	}
	return diag.Errorf("instance %q did not reach ACTIVE status after reboot within timeout", instanceID)
}
