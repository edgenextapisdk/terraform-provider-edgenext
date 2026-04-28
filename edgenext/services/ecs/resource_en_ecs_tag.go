package ecs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceENECSTag returns the resource schema for ECS tag.
func ResourceENECSTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSTagCreate,
		ReadContext:   resourceENECSTagRead,
		UpdateContext: resourceENECSTagUpdate,
		DeleteContext: resourceENECSTagDelete,
		CustomizeDiff: resourceENECSTagCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSTagImport,
		},
		Description: "Provides an EdgeNext ECS tag resource. key and value cannot be changed after creation.",
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag key. Cannot be changed after creation.",
			},
			"tag_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag value. Cannot be changed after creation.",
			},
		},
	}
}

func resourceENECSTagCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip this check during creation.
	if d.Id() == "" {
		return nil
	}
	if d.HasChange("tag_key") {
		oldRaw, newRaw := d.GetChange("tag_key")
		if strings.TrimSpace(fmt.Sprintf("%v", oldRaw)) != strings.TrimSpace(fmt.Sprintf("%v", newRaw)) {
			return fmt.Errorf("tag_key cannot be modified after creation")
		}
	}
	if d.HasChange("tag_value") {
		oldRaw, newRaw := d.GetChange("tag_value")
		if strings.TrimSpace(fmt.Sprintf("%v", oldRaw)) != strings.TrimSpace(fmt.Sprintf("%v", newRaw)) {
			return fmt.Errorf("tag_value cannot be modified after creation")
		}
	}
	return nil
}

func resourceENECSTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// No update API; UpdateContext exists so key/value do not require ForceNew in the schema (SDK rule).
	if d.HasChange("tag_key") || d.HasChange("tag_value") {
		return diag.Errorf("ECS tag does not support updating tag_key or tag_value in place; destroy and recreate the resource if you need different tag_key/tag_value")
	}
	return resourceENECSTagRead(ctx, d, m)
}

func resourceENECSTagImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected import id as tag_id/tag_key/tag_value, got %q", d.Id())
	}
	tagID := strings.TrimSpace(parts[0])
	key := strings.TrimSpace(parts[1])
	value := strings.TrimSpace(parts[2])
	if tagID == "" || key == "" {
		return nil, fmt.Errorf("expected import id as tag_id/tag_key/tag_value, got %q", d.Id())
	}
	if _, err := strconv.Atoi(tagID); err != nil {
		return nil, fmt.Errorf("invalid tag_id %q in import id %q", tagID, d.Id())
	}
	if err := d.Set("tag_key", key); err != nil {
		return nil, err
	}
	if err := d.Set("tag_value", value); err != nil {
		return nil, err
	}
	d.SetId(tagID)
	if diags := resourceENECSTagRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("tag %q not found", tagID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"tags": []map[string]interface{}{
			{
				"key":   d.Get("tag_key").(string),
				"value": d.Get("tag_value").(string),
			},
		},
	}
	var resp map[string]interface{}
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/tags/create", req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS tag: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS tag create response: %s", err)
	}
	createdTags := helper.ListFromMap(payload, "createdTags")
	if len(createdTags) == 0 {
		return diag.Errorf("failed to parse ECS tag create response: missing createdTags")
	}
	first, ok := createdTags[0].(map[string]interface{})
	if !ok {
		return diag.Errorf("failed to parse ECS tag create response: invalid createdTags item type %T", createdTags[0])
	}
	createdID := helper.IntFromMap(first, "id")
	if createdID <= 0 {
		return diag.Errorf("failed to parse ECS tag create response: invalid created tag id")
	}
	d.SetId(strconv.Itoa(createdID))
	if key := helper.StringFromMap(first, "key"); key != "" {
		_ = d.Set("tag_key", key)
	}
	if value := helper.StringFromMap(first, "value"); value != "" {
		_ = d.Set("tag_value", value)
	}

	return resourceENECSTagRead(ctx, d, m)
}

func resourceENECSTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"tagKey":   d.Get("tag_key").(string),
		"tagValue": d.Get("tag_value").(string),
		"pageNum":  1,
		"pageSize": 100,
	}
	var resp map[string]interface{}
	err = ecsClient.Get(ctx, "/ecs/openapi/v2/tags/list", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS tag list response: %s", err)
	}
	list := helper.ListFromMap(payload, "list")
	idInt, _ := strconv.Atoi(d.Id())
	found := false
	for _, raw := range list {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		currentID := helper.IntFromMap(item, "id")
		currentKey := helper.StringFromMap(item, "tagKey")
		currentValue := helper.StringFromMap(item, "tagValue")
		if currentID == idInt || (currentKey == d.Get("tag_key").(string) && currentValue == d.Get("tag_value").(string)) {
			found = true
			_ = d.Set("tag_key", currentKey)
			_ = d.Set("tag_value", currentValue)
			break
		}
	}
	if !found {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceENECSTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	tagID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse ECS tag ID: %s", err)
	}
	req := map[string]interface{}{"tagIds": []int{tagID}}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/tags/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS tag: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS tag delete response: %s", err)
	}

	return nil
}
