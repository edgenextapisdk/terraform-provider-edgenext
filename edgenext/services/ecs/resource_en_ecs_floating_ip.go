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

// ResourceENECSFloatingIp returns the resource schema for ECS floating_ip.
func ResourceENECSFloatingIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSFloatingIpCreate,
		ReadContext:   resourceENECSFloatingIpRead,
		UpdateContext: resourceENECSFloatingIpUpdate,
		DeleteContext: resourceENECSFloatingIpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSFloatingIpImport,
		},
		Description: "Provides an EdgeNext ECS floating_ip resource.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionResourceSchema("region description"),
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "bandwidth description",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ip_address description",
			},
		},
	}
}

func resourceENECSFloatingIpImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import id as region/floating_ip_id, got %q", d.Id())
	}
	region := helper.NormalizeRegion(parts[0])
	floatingIPID := strings.TrimSpace(parts[1])
	if region == "" || floatingIPID == "" {
		return nil, fmt.Errorf("expected import id as region/floating_ip_id, got %q", d.Id())
	}
	if err := d.Set("region", region); err != nil {
		return nil, err
	}
	d.SetId(floatingIPID)
	if diags := resourceENECSFloatingIpRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("floating IP %q not found in region %q", floatingIPID, region)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSFloatingIpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":    helper.NormalizeRegion(d.Get("region").(string)),
		"bandwidth": d.Get("bandwidth").(int),
	}
	if n, ok := d.GetOk("name"); ok {
		req["name"] = n.(string)
	}
	var resp map[string]interface{}

	path := "/ecs/openapi/v2/floatingIp/create"
	if "floating_ip" == "instance" {
		path = "/ecs/openapi/v2/floatingIp/create_order"
	}

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS floating_ip: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS floating_ip create response: %s", err)
	}
	fallbackID := "created-floating_ip"
	if n, ok := d.GetOk("name"); ok {
		fallbackID = n.(string)
	}
	d.SetId(helper.ExtractIDFromPayload(payload, fallbackID))

	return resourceENECSFloatingIpRead(ctx, d, m)
}

func resourceENECSFloatingIpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/floatingIp/detail", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS floating_ip detail response: %s", err)
	}

	if name, ok := payload["name"].(string); ok {
		d.Set("name", name)
	}
	if val, ok := payload["bandwidth"]; ok {
		d.Set("bandwidth", val)
	}
	if val, ok := payload["ip_address"]; ok {
		d.Set("ip_address", val)
	}

	return nil
}

func resourceENECSFloatingIpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

		err = ecsClient.Post(ctx, "/ecs/openapi/v2/floatingIp/update", req, &resp)
		if err != nil {
			return diag.Errorf("failed to update ECS floating_ip: %s", err)
		}
		if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
			return diag.Errorf("failed to parse ECS floating_ip update response: %s", err)
		}
	}

	return resourceENECSFloatingIpRead(ctx, d, m)
}

func resourceENECSFloatingIpDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/floatingIp/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS floating_ip: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS floating_ip delete response: %s", err)
	}

	return nil
}
