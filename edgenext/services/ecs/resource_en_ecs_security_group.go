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
	securityGroupID := strings.TrimSpace(d.Id())
	if securityGroupID == "" {
		return nil, fmt.Errorf("expected import id as security_group_id, got %q", d.Id())
	}
	d.SetId(securityGroupID)
	if diags := resourceENECSSecurityGroupRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("security group %q not found", securityGroupID)
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
		"id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/detail", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS security_group detail response: %s", err)
	}
	securityGroup := helper.MapFromMap(payload, "security_group")
	if securityGroup == nil {
		d.SetId("")
		return nil
	}
	currentID := helper.StringFromMap(securityGroup, "id")
	if currentID == "" {
		d.SetId("")
		return nil
	}
	d.SetId(currentID)
	_ = d.Set("name", helper.StringFromMap(securityGroup, "name"))
	_ = d.Set("description", helper.StringFromMap(securityGroup, "description"))
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
		"id": d.Id(),
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
		"ids": []string{d.Id()},
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
