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

// ResourceENECSDisk returns the resource schema for ECS disk.
func ResourceENECSDisk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSDiskCreate,
		ReadContext:   resourceENECSDiskRead,
		UpdateContext: resourceENECSDiskUpdate,
		DeleteContext: resourceENECSDiskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSDiskImport,
		},
		Description: "Provides an EdgeNext ECS disk resource.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name description",
			},
			"volume_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "volume_type description",
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "size description",
			},
		},
	}
}

func resourceENECSDiskImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	diskID := strings.TrimSpace(d.Id())
	if diskID == "" {
		return nil, fmt.Errorf("expected import id as disk_id, got %q", d.Id())
	}
	d.SetId(diskID)
	if diags := resourceENECSDiskRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("disk %q not found", diskID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSDiskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"volume_type": d.Get("volume_type").(string),
		"size":        d.Get("size").(int),
	}
	if n, ok := d.GetOk("name"); ok {
		req["name"] = n.(string)
	}
	var resp map[string]interface{}

	path := "/ecs/openapi/v2/volume/create"
	if "disk" == "instance" {
		path = "/ecs/openapi/v2/volume/create_order"
	}

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS disk: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS disk create response: %s", err)
	}
	fallbackID := "created-disk"
	if n, ok := d.GetOk("name"); ok {
		fallbackID = n.(string)
	}
	d.SetId(helper.ExtractIDFromPayload(payload, fallbackID))

	return resourceENECSDiskRead(ctx, d, m)
}

func resourceENECSDiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/volume/detail", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS disk detail response: %s", err)
	}

	if name, ok := payload["name"].(string); ok {
		d.Set("name", name)
	}
	if val, ok := payload["volume_type"]; ok {
		d.Set("volume_type", val)
	}
	if val, ok := payload["size"]; ok {
		d.Set("size", val)
	}

	return nil
}

func resourceENECSDiskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

		err = ecsClient.Post(ctx, "/ecs/openapi/v2/volume/update", req, &resp)
		if err != nil {
			return diag.Errorf("failed to update ECS disk: %s", err)
		}
		if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
			return diag.Errorf("failed to parse ECS disk update response: %s", err)
		}
	}

	return resourceENECSDiskRead(ctx, d, m)
}

func resourceENECSDiskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/volume/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS disk: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS disk delete response: %s", err)
	}

	return nil
}
