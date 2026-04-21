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

// ResourceENECSRouter returns the resource schema for ECS router.
func ResourceENECSRouter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSRouterCreate,
		ReadContext:   resourceENECSRouterRead,
		UpdateContext: resourceENECSRouterUpdate,
		DeleteContext: resourceENECSRouterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSRouterImport,
		},
		Description: "Provides an EdgeNext ECS router resource.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionResourceSchema("region description"),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name description",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description description",
			},
			"external_network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "external_network_id description",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "tenant_id description",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "admin_state_up description",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "status description",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "project_id description",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "created_at description",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "updated_at description",
			},
			"revision_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "revision_number description",
			},
		},
	}
}

func resourceENECSRouterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import id as region/router_id, got %q", d.Id())
	}

	region := helper.NormalizeRegion(parts[0])
	if err := d.Set("region", region); err != nil {
		return nil, err
	}
	d.SetId(parts[1])

	if diags := resourceENECSRouterRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("router %q not found in region %q", parts[1], region)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceENECSRouterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"router": map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		},
	}
	if v, ok := d.GetOk("external_network_id"); ok && v.(string) != "" {
		req["router"].(map[string]interface{})["external_network_id"] = v.(string)
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/routers/add", req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS router: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS router create response: %s", err)
	}
	router := helper.MapFromMap(payload, "router")
	if router == nil {
		return diag.Errorf("failed to parse ECS router create response: missing router")
	}
	routerID := helper.StringFromMap(router, "id")
	if routerID == "" {
		return diag.Errorf("failed to parse ECS router create response: missing router.id")
	}
	d.SetId(routerID)

	return resourceENECSRouterRead(ctx, d, m)
}

func resourceENECSRouterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"id":     d.Id(),
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/routers/detail", req, &resp); err != nil {
		d.SetId("")
		return nil
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS router detail response: %s", err)
	}
	router := helper.MapFromMap(payload, "router")
	if router == nil {
		d.SetId("")
		return nil
	}

	if name, ok := router["name"].(string); ok {
		_ = d.Set("name", name)
	}
	if val, ok := router["description"]; ok {
		_ = d.Set("description", val)
	}
	_ = d.Set("tenant_id", helper.StringFromMap(router, "tenant_id"))
	_ = d.Set("admin_state_up", routerBoolFromMap(router, "admin_state_up"))
	_ = d.Set("status", helper.StringFromMap(router, "status"))
	_ = d.Set("project_id", helper.StringFromMap(router, "project_id"))
	_ = d.Set("created_at", helper.StringFromMap(router, "created_at"))
	_ = d.Set("updated_at", helper.StringFromMap(router, "updated_at"))
	_ = d.Set("revision_number", helper.IntFromMap(router, "revision_number"))

	externalGatewayInfo := helper.MapFromMap(router, "external_gateway_info")
	_ = d.Set("external_network_id", helper.StringFromMap(externalGatewayInfo, "network_id"))

	return nil
}

func resourceENECSRouterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.HasChanges("name", "description", "external_network_id") {
		return resourceENECSRouterRead(ctx, d, m)
	}

	if d.HasChanges("name", "description") {
		req := map[string]interface{}{
			"region": helper.NormalizeRegion(d.Get("region").(string)),
			"router": map[string]interface{}{
				"id":          d.Id(),
				"name":        d.Get("name").(string),
				"description": d.Get("description").(string),
			},
		}
		var resp map[string]interface{}
		err = ecsClient.Post(ctx, "/ecs/openapi/v2/routers/add", req, &resp)
		if err != nil {
			return diag.Errorf("failed to update ECS router: %s", err)
		}
		if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
			return diag.Errorf("failed to parse ECS router update response: %s", err)
		}
	}

	if d.HasChange("external_network_id") {
		req := map[string]interface{}{
			"region": helper.NormalizeRegion(d.Get("region").(string)),
			"router": map[string]interface{}{
				"id":                  d.Id(),
				"external_network_id": d.Get("external_network_id").(string),
			},
		}
		var resp map[string]interface{}
		err = ecsClient.Post(ctx, "/ecs/openapi/v2/routers/gateway", req, &resp)
		if err != nil {
			return diag.Errorf("failed to update ECS router gateway: %s", err)
		}
		if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
			return diag.Errorf("failed to parse ECS router gateway update response: %s", err)
		}
	}

	return resourceENECSRouterRead(ctx, d, m)
}

func resourceENECSRouterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"ids":    []string{d.Id()},
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/routers/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS router: %s", err)
	}
	payload, err := helper.ParseAPIResponsePayload(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS router delete response: %s", err)
	}
	if m, ok := payload.(map[string]interface{}); ok {
		if status, ok := m[d.Id()].(string); !ok || status != "ok" {
			return diag.Errorf("ECS router delete: unexpected status for id %q: %v", d.Id(), m[d.Id()])
		}
	}

	return nil
}

func routerBoolFromMap(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}
