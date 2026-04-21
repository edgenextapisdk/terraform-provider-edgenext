package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceENECSNetworkInterface returns the resource schema for ECS network_interface.
func ResourceENECSNetworkInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSNetworkInterfaceCreate,
		ReadContext:   resourceENECSNetworkInterfaceRead,
		UpdateContext: resourceENECSNetworkInterfaceUpdate,
		DeleteContext: resourceENECSNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Provides an EdgeNext ECS network_interface resource.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionResourceSchema("region description"),
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "network_id description",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet_id description",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "mac_address description",
			},
		},
	}
}

func resourceENECSNetworkInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":     d.Get("region").(string),
		"network_id": d.Get("network_id").(string),
		"subnet_id":  d.Get("subnet_id").(string),
	}
	if n, ok := d.GetOk("name"); ok {
		req["name"] = n.(string)
	}
	var resp map[string]interface{}

	path := "/ecs/openapi/v2/port/create"
	if "network_interface" == "instance" {
		path = "/ecs/openapi/v2/port/create_order"
	}

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS network_interface: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS network_interface create response: %s", err)
	}
	fallbackID := "created-network_interface"
	if n, ok := d.GetOk("name"); ok {
		fallbackID = n.(string)
	}
	d.SetId(helper.ExtractIDFromPayload(payload, fallbackID))

	return resourceENECSNetworkInterfaceRead(ctx, d, m)
}

func resourceENECSNetworkInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/port/detail", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS network_interface detail response: %s", err)
	}

	if name, ok := payload["name"].(string); ok {
		d.Set("name", name)
	}
	if val, ok := payload["network_id"]; ok {
		d.Set("network_id", val)
	}
	if val, ok := payload["subnet_id"]; ok {
		d.Set("subnet_id", val)
	}
	if val, ok := payload["mac_address"]; ok {
		d.Set("mac_address", val)
	}

	return nil
}

func resourceENECSNetworkInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	// Try name update
	if d.HasChange("name") {
		req := map[string]interface{}{
			"id":   d.Id(),
			"name": d.Get("name"),
		}
		var resp map[string]interface{}

		err = ecsClient.Post(ctx, "/ecs/openapi/v2/port/update", req, &resp)
		if err != nil {
			return diag.Errorf("failed to update ECS network_interface: %s", err)
		}
		if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
			return diag.Errorf("failed to parse ECS network_interface update response: %s", err)
		}
	}

	return resourceENECSNetworkInterfaceRead(ctx, d, m)
}

func resourceENECSNetworkInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/port/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS network_interface: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS network_interface delete response: %s", err)
	}

	return nil
}
