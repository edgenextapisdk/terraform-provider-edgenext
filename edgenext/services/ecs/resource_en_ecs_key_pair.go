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

// ResourceENECSKeyPair returns the resource schema for ECS key_pair.
func ResourceENECSKeyPair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSKeyPairCreate,
		ReadContext:   resourceENECSKeyPairRead,
		DeleteContext: resourceENECSKeyPairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSKeyPairImport,
		},
		Description: "Provides an EdgeNext ECS key_pair resource.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionResourceSchema("region description"),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "name description",
			},
			"public_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "public_key description",
			},
			"private_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "private_key description",
			},
		},
	}
}

func resourceENECSKeyPairImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import id as region/name, got %q", d.Id())
	}
	region := helper.NormalizeRegion(parts[0])
	name := strings.TrimSpace(parts[1])
	if region == "" || name == "" {
		return nil, fmt.Errorf("expected import id as region/name, got %q", d.Id())
	}
	if err := d.Set("region", region); err != nil {
		return nil, err
	}
	if err := d.Set("name", name); err != nil {
		return nil, err
	}
	d.SetId(name)
	if diags := resourceENECSKeyPairRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("key pair %q not found in region %q", name, region)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSKeyPairCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":     d.Get("region").(string),
		"public_key": d.Get("public_key").(string),
	}
	if n, ok := d.GetOk("name"); ok {
		req["name"] = n.(string)
	}
	var resp map[string]interface{}

	path := "/ecs/openapi/v2/keypair/create"

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS key_pair: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS key_pair create response: %s", err)
	}
	fallbackID := d.Get("name").(string)
	d.SetId(helper.ExtractIDFromPayload(payload, fallbackID))

	return resourceENECSKeyPairRead(ctx, d, m)
}

func resourceENECSKeyPairRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"name":   d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/keypair/detail", req, &resp)
	if err != nil {
		// Only clear state on explicit not-found; other errors should fail apply/refresh.
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "404") {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read ECS key_pair: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS key_pair detail response: %s", err)
	}

	if name, ok := payload["name"].(string); ok {
		d.Set("name", name)
	}
	if val, ok := payload["public_key"]; ok {
		d.Set("public_key", val)
	}
	if val, ok := payload["private_key"]; ok {
		d.Set("private_key", val)
	}

	return nil
}

func resourceENECSKeyPairDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"names":  []string{d.Id()},
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/keypair/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS key_pair: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS key_pair delete response: %s", err)
	}

	return nil
}
