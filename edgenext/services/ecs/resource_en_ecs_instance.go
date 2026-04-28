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

// ResourceENECSInstance returns the resource schema for ECS instance.
func ResourceENECSInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSInstanceCreate,
		ReadContext:   resourceENECSInstanceRead,
		UpdateContext: resourceENECSInstanceUpdate,
		DeleteContext: resourceENECSInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSInstanceImport,
		},
		Description: "Provides an EdgeNext ECS instance resource.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name description",
			},
			"flavor_ref": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "flavor_ref description",
			},
			"image_ref": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "image_ref description",
			},
			"admin_pass": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "admin_pass description",
			},
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "key_name description",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "project_id description",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "bandwidth description",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "status description",
			},
			"networks": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "networks description",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "security_groups description",
			},
		},
	}
}

func resourceENECSInstanceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	instanceID := strings.TrimSpace(d.Id())
	if instanceID == "" {
		return nil, fmt.Errorf("expected import id as instance_id, got %q", d.Id())
	}
	d.SetId(instanceID)
	if diags := resourceENECSInstanceRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("instance %q not found", instanceID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"flavor_ref": d.Get("flavor_ref").(string),
		"image_ref":  d.Get("image_ref").(string),
		"admin_pass": d.Get("admin_pass").(string),
		"key_name":   d.Get("key_name").(string),
		"project_id": d.Get("project_id").(string),
		"bandwidth":  d.Get("bandwidth").(int),
		"networks": func() []string {
			lst := d.Get("networks").([]interface{})
			res := make([]string, len(lst))
			for i, v := range lst {
				res[i] = v.(string)
			}
			return res
		}(),
		"security_groups": func() []string {
			lst := d.Get("security_groups").([]interface{})
			res := make([]string, len(lst))
			for i, v := range lst {
				res[i] = v.(string)
			}
			return res
		}(),
	}
	if n, ok := d.GetOk("name"); ok {
		req["name"] = n.(string)
	}
	var resp map[string]interface{}

	path := "/ecs/openapi/v2/instance/create"

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS instance: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS instance create response: %s", err)
	}
	fallbackID := "created-instance"
	if n, ok := d.GetOk("name"); ok {
		fallbackID = n.(string)
	}
	d.SetId(helper.ExtractIDFromPayload(payload, fallbackID))

	return resourceENECSInstanceRead(ctx, d, m)
}

func resourceENECSInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/instance/detail", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS instance detail response: %s", err)
	}

	if name, ok := payload["name"].(string); ok {
		d.Set("name", name)
	}
	if val, ok := payload["flavor_ref"]; ok {
		d.Set("flavor_ref", val)
	}
	if val, ok := payload["image_ref"]; ok {
		d.Set("image_ref", val)
	}
	if val, ok := payload["admin_pass"]; ok {
		d.Set("admin_pass", val)
	}
	if val, ok := payload["key_name"]; ok {
		d.Set("key_name", val)
	}
	if val, ok := payload["project_id"]; ok {
		d.Set("project_id", val)
	}
	if val, ok := payload["bandwidth"]; ok {
		d.Set("bandwidth", val)
	}
	if val, ok := payload["status"]; ok {
		d.Set("status", val)
	}
	if val, ok := payload["networks"]; ok {
		d.Set("networks", val)
	}
	if val, ok := payload["security_groups"]; ok {
		d.Set("security_groups", val)
	}

	return nil
}

func resourceENECSInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

		err = ecsClient.Post(ctx, "/ecs/openapi/v2/instance/update", req, &resp)
		if err != nil {
			return diag.Errorf("failed to update ECS instance: %s", err)
		}
		if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
			return diag.Errorf("failed to parse ECS instance update response: %s", err)
		}
	}

	return resourceENECSInstanceRead(ctx, d, m)
}

func resourceENECSInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/instance/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS instance: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS instance delete response: %s", err)
	}

	return nil
}
