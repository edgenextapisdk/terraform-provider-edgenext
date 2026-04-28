package ecs

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceENECSNetworkInterfaceInstanceBinding binds an instance to a network interface.
// No UpdateContext: all arguments are ForceNew; SDK rejects a superfluous Update in that case.
func ResourceENECSNetworkInterfaceInstanceBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSNetworkInterfaceInstanceBindingCreate,
		ReadContext:   resourceENECSNetworkInterfaceInstanceBindingRead,
		DeleteContext: resourceENECSNetworkInterfaceInstanceBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSNetworkInterfaceInstanceBindingImport,
		},
		Description: "Provides an EdgeNext ECS network interface instance binding resource.",
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network interface ID. Changing this forces a new resource.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance ID to bind. Changing this forces a new resource.",
			},
		},
	}
}

func resourceENECSNetworkInterfaceInstanceBindingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import id as network_interface_id/instance_id, got %q", d.Id())
	}
	portID := strings.TrimSpace(parts[0])
	instanceID := strings.TrimSpace(parts[1])
	if portID == "" || instanceID == "" {
		return nil, fmt.Errorf("expected import id as network_interface_id/instance_id, got %q", d.Id())
	}
	_ = d.Set("network_interface_id", portID)
	_ = d.Set("instance_id", instanceID)
	d.SetId(fmt.Sprintf("%s/%s", portID, instanceID))
	if diags := resourceENECSNetworkInterfaceInstanceBindingRead(ctx, d, meta); diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("network interface instance binding not found for import id %q", fmt.Sprintf("%s/%s", portID, instanceID))
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSNetworkInterfaceInstanceBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}
	portID := strings.TrimSpace(d.Get("network_interface_id").(string))
	instanceID := strings.TrimSpace(d.Get("instance_id").(string))
	if err := resourceENECSNetworkInterfaceRelationServer(ctx, ecsClient, portID, "add", instanceID); err != nil {
		return diag.Errorf("failed to bind instance to ECS network_interface: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", portID, instanceID))
	return resourceENECSNetworkInterfaceInstanceBindingRead(ctx, d, m)
}

func resourceENECSNetworkInterfaceInstanceBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}
	portID := strings.TrimSpace(d.Get("network_interface_id").(string))
	expectedInstanceID := strings.TrimSpace(d.Get("instance_id").(string))

	port, err := resourceENECSNetworkInterfacePortDetail(ctx, ecsClient, portID)
	if err != nil {
		d.SetId("")
		return nil
	}
	actualInstanceID := strings.TrimSpace(helper.StringFromMap(port, "device_id"))
	if expectedInstanceID == "" || actualInstanceID != expectedInstanceID {
		d.SetId("")
		return nil
	}
	return nil
}

func resourceENECSNetworkInterfaceInstanceBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}
	portID := strings.TrimSpace(d.Get("network_interface_id").(string))
	instanceID := strings.TrimSpace(d.Get("instance_id").(string))
	if instanceID != "" {
		if err := resourceENECSNetworkInterfaceRelationServer(ctx, ecsClient, portID, "remove", instanceID); err != nil {
			return diag.Errorf("failed to unbind instance from ECS network_interface: %s", err)
		}
	}
	return nil
}
