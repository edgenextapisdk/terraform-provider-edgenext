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

// ResourceENECSSecurityGroup returns the resource schema for ECS security_group.
func ResourceENECSSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSSecurityGroupCreate,
		ReadContext:   resourceENECSSecurityGroupRead,
		UpdateContext: resourceENECSSecurityGroupUpdate,
		DeleteContext: resourceENECSSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSSecurityGroupImport,
		},
		Description: "Provides an EdgeNext ECS security_group resource.",
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
		},
	}
}

func resourceENECSSecurityGroupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
	// Read uses name to query and then sets canonical id.
	d.SetId(name)
	if diags := resourceENECSSecurityGroupRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("security group %q not found in region %q", name, region)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"security_group": map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		},
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/add", req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS security_group: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS security_group create response: %s", err)
	}
	created := helper.MapFromMap(payload, "security_group")
	if created == nil {
		return diag.Errorf("failed to parse ECS security_group create response: missing security_group")
	}
	createdID := helper.StringFromMap(created, "id")
	if createdID == "" {
		return diag.Errorf("failed to parse ECS security_group create response: missing security_group.id")
	}
	d.SetId(createdID)

	return resourceENECSSecurityGroupRead(ctx, d, m)
}

func resourceENECSSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"name":   d.Get("name").(string),
		"limit":  10,
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/list", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS security_group list response: %s", err)
	}
	securityGroups := helper.ListFromMap(payload, "security_groups")
	if len(securityGroups) == 0 {
		d.SetId("")
		return nil
	}
	first, ok := securityGroups[0].(map[string]interface{})
	if !ok {
		d.SetId("")
		return nil
	}
	currentID := helper.StringFromMap(first, "id")
	if currentID == "" {
		d.SetId("")
		return nil
	}
	d.SetId(currentID)
	_ = d.Set("name", helper.StringFromMap(first, "name"))
	_ = d.Set("description", helper.StringFromMap(first, "description"))
	return nil
}

func resourceENECSSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if !d.HasChanges("name", "description") {
		return resourceENECSSecurityGroupRead(ctx, d, m)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"id":     d.Id(),
		"security_group": map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		},
	}
	var resp map[string]interface{}
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/add", req, &resp)
	if err != nil {
		return diag.Errorf("failed to update ECS security_group: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse ECS security_group update response: %s", err)
	}

	return resourceENECSSecurityGroupRead(ctx, d, m)
}

func resourceENECSSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS security_group: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS security_group delete response: %s", err)
	}

	return nil
}
