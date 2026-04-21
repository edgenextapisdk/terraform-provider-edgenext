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

// ResourceENECSRouterPort returns the resource schema for ECS router subnet attachment.
func ResourceENECSRouterPort() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSRouterPortCreate,
		ReadContext:   resourceENECSRouterPortRead,
		DeleteContext: resourceENECSRouterPortDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSRouterPortImport,
		},
		Description: "Provides an EdgeNext ECS router port attachment resource. There is no update API; changing any argument replaces the resource.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionResourceSchema("The region of the router."),
			"router_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The router ID.",
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network ID to attach.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet ID to attach.",
			},
			"port_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The created router port ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port name.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port IP address.",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port MAC address.",
			},
			"network_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port network name.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port status.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port creation time.",
			},
		},
	}
}

func resourceENECSRouterPortImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected import id as region/router_id/port_id, got %q", d.Id())
	}

	region := helper.NormalizeRegion(parts[0])
	if err := d.Set("region", region); err != nil {
		return nil, err
	}
	if err := d.Set("router_id", parts[1]); err != nil {
		return nil, err
	}
	d.SetId(parts[2])

	if diags := resourceENECSRouterPortRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("router port %q not found under router %q in region %q", d.Id(), parts[1], region)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceENECSRouterPortCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":     helper.NormalizeRegion(d.Get("region").(string)),
		"id":         d.Get("router_id").(string),
		"network_id": d.Get("network_id").(string),
		"subnet_id":  d.Get("subnet_id").(string),
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/routers/add_sub", req, &resp); err != nil {
		return diag.Errorf("failed to create ECS router port: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS router port create response: %s", err)
	}
	portID := helper.StringFromMap(payload, "port_id")
	if portID == "" {
		return diag.Errorf("failed to parse ECS router port create response: missing port_id")
	}
	d.SetId(portID)
	_ = d.Set("port_id", portID)

	return resourceENECSRouterPortRead(ctx, d, m)
}

func resourceENECSRouterPortRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"id":     d.Get("router_id").(string),
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/routers/port_list", req, &resp); err != nil {
		d.SetId("")
		return nil
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS router port read response: %s", err)
	}
	rawPorts := helper.ListFromMap(payload, "ports")
	var matched map[string]interface{}
	for _, raw := range rawPorts {
		port, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(port, "id") == d.Id() {
			matched = port
			break
		}
	}
	if matched == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("port_id", helper.StringFromMap(matched, "id"))
	_ = d.Set("name", helper.StringFromMap(matched, "name"))
	_ = d.Set("ip_address", helper.StringFromMap(matched, "ip_address"))
	_ = d.Set("mac_address", helper.StringFromMap(matched, "mac_address"))
	_ = d.Set("network_name", helper.StringFromMap(matched, "network_name"))
	_ = d.Set("status", helper.StringFromMap(matched, "status"))
	_ = d.Set("created_at", helper.StringFromMap(matched, "created_at"))

	return nil
}

func resourceENECSRouterPortDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":    helper.NormalizeRegion(d.Get("region").(string)),
		"id":        d.Get("router_id").(string),
		"subnet_id": d.Get("subnet_id").(string),
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/routers/remove_sub", req, &resp); err != nil {
		return diag.Errorf("failed to delete ECS router port: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS router port delete response: %s", err)
	}

	return nil
}
