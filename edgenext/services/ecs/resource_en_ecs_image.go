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

// ResourceENECSImage returns the resource schema for ECS image.
func ResourceENECSImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSImageCreate,
		ReadContext:   resourceENECSImageRead,
		UpdateContext: resourceENECSImageUpdate,
		DeleteContext: resourceENECSImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSImageImport,
		},
		Description: "Provides an EdgeNext ECS image resource.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name description",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "instance_id description",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description description",
			},
			"os_distro": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "os_distro description",
			},
		},
	}
}

func resourceENECSImageImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	imageID := strings.TrimSpace(d.Id())
	if imageID == "" {
		return nil, fmt.Errorf("expected import id as image_id, got %q", d.Id())
	}
	d.SetId(imageID)
	if diags := resourceENECSImageRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("image %q not found", imageID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"instance_id": d.Get("instance_id").(string),
		"description": d.Get("description").(string),
	}
	if n, ok := d.GetOk("name"); ok {
		req["name"] = n.(string)
	}
	var resp map[string]interface{}

	path := "/ecs/openapi/v2/image/create"
	if "image" == "instance" {
		path = "/ecs/openapi/v2/image/create_order"
	}

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS image: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS image create response: %s", err)
	}
	fallbackID := "created-image"
	if n, ok := d.GetOk("name"); ok {
		fallbackID = n.(string)
	}
	d.SetId(helper.ExtractIDFromPayload(payload, fallbackID))

	return resourceENECSImageRead(ctx, d, m)
}

func resourceENECSImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/image/detail", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS image detail response: %s", err)
	}

	if name, ok := payload["name"].(string); ok {
		d.Set("name", name)
	}
	if val, ok := payload["instance_id"]; ok {
		d.Set("instance_id", val)
	}
	if val, ok := payload["description"]; ok {
		d.Set("description", val)
	}
	if val, ok := payload["os_distro"]; ok {
		d.Set("os_distro", val)
	}

	return nil
}

func resourceENECSImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

		err = ecsClient.Post(ctx, "/ecs/openapi/v2/image/update", req, &resp)
		if err != nil {
			return diag.Errorf("failed to update ECS image: %s", err)
		}
		if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
			return diag.Errorf("failed to parse ECS image update response: %s", err)
		}
	}

	return resourceENECSImageRead(ctx, d, m)
}

func resourceENECSImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/image/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS image: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS image delete response: %s", err)
	}

	return nil
}
