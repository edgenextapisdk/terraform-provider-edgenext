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

// ResourceENECSInstanceTag manages tag bindings on a specific ECS resource.
func ResourceENECSInstanceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSInstanceTagCreate,
		ReadContext:   resourceENECSInstanceTagRead,
		UpdateContext: resourceENECSInstanceTagUpdate,
		DeleteContext: resourceENECSInstanceTagDelete,
		CustomizeDiff: resourceENECSInstanceTagCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSInstanceTagImport,
		},
		Description: "Provides an EdgeNext ECS instance tag binding resource. instance_id and instance_name cannot be changed after creation.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target instance ID. Cannot be changed after creation.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The target instance name. Cannot be changed after creation.",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance type returned by query API.",
			},
			"tag_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of tags on this instance.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Detailed tag items for the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Tag item ID.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag item key.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tag item value.",
						},
					},
				},
			},
			"tag_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Tag IDs to bind to the target instance. Element order is ignored when comparing updates.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceENECSInstanceTagCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip this check during creation.
	if d.Id() == "" {
		return nil
	}
	if d.HasChange("instance_id") {
		oldRaw, newRaw := d.GetChange("instance_id")
		if strings.TrimSpace(fmt.Sprintf("%v", oldRaw)) != strings.TrimSpace(fmt.Sprintf("%v", newRaw)) {
			return fmt.Errorf("instance_id cannot be modified after creation")
		}
	}
	if d.HasChange("instance_name") {
		oldRaw, newRaw := d.GetChange("instance_name")
		if strings.TrimSpace(fmt.Sprintf("%v", oldRaw)) != strings.TrimSpace(fmt.Sprintf("%v", newRaw)) {
			return fmt.Errorf("instance_name cannot be modified after creation")
		}
	}
	return nil
}

func resourceENECSInstanceTagImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	instanceID := strings.TrimSpace(d.Id())
	if instanceID == "" {
		return nil, fmt.Errorf("expected import id as instance_id, got %q", d.Id())
	}
	if err := d.Set("instance_id", instanceID); err != nil {
		return nil, err
	}
	d.SetId(instanceID)
	if diags := resourceENECSInstanceTagRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("instance tag binding for instance %q not found", instanceID)
	}
	return []*schema.ResourceData{d}, nil
}

// Create binds the provided tag IDs to the target instance.
func resourceENECSInstanceTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := syncENECSInstanceTags(ctx, ecsClient, d, expandIntSet(d.Get("tag_ids").(*schema.Set))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("instance_id").(string))
	return resourceENECSInstanceTagRead(ctx, d, m)
}

// Read refreshes current tag bindings from the instance tag query endpoint.
func resourceENECSInstanceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":   ecsClient.Region(),
		"tagId":    0,
		"tagKey":   "",
		"tagValue": "",
		"pageNum":  1,
		"pageSize": 100,
	}
	var resp map[string]interface{}
	// The list endpoint is used as the authoritative source of current bindings.
	if err := ecsClient.Get(ctx, "/ecs/openapi/v2/resource/list", req, &resp); err != nil {
		// Treat read errors as not found to avoid hard-failing refresh.
		d.SetId("")
		return nil
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS instance tag read response: %s", err)
	}
	list := helper.ListFromMap(payload, "list")
	instanceID := d.Get("instance_id").(string)

	// Find the exact target instance and synchronize state fields from API output.
	for _, raw := range list {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(item, "resourceId") != instanceID {
			continue
		}
		_ = d.Set("instance_name", helper.StringFromMap(item, "resourceName"))
		_ = d.Set("instance_type", helper.StringFromMap(item, "productType"))
		_ = d.Set("tag_count", helper.IntFromMap(item, "tagCount"))
		_ = d.Set("tags", normalizeENECSInstanceTagItems(helper.ListFromMap(item, "tags")))
		_ = d.Set("tag_ids", extractTagIDs(helper.ListFromMap(item, "tags")))
		return nil
	}
	// If the instance is not found, set the ID to empty and return nil.
	d.SetId("")
	return nil
}

// Update only syncs tag_ids changes.
func resourceENECSInstanceTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Defense in depth: CustomizeDiff blocks these at plan time; reject here if Update is still invoked.
	if d.HasChange("instance_id") || d.HasChange("instance_name") {
		return diag.Errorf("instance_id and instance_name cannot be updated after creation")
	}

	if !d.HasChange("tag_ids") {
		return resourceENECSInstanceTagRead(ctx, d, m)
	}

	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := syncENECSInstanceTags(ctx, ecsClient, d, expandIntSet(d.Get("tag_ids").(*schema.Set))); err != nil {
		return diag.FromErr(err)
	}

	return resourceENECSInstanceTagRead(ctx, d, m)
}

// Delete clears tag bindings by syncing an empty tagIds list.
func resourceENECSInstanceTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := syncENECSInstanceTags(ctx, ecsClient, d, []int{}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// syncENECSInstanceTags calls the tags/sync endpoint and validates API envelope.
func syncENECSInstanceTags(ctx context.Context, ecsClient *connectivity.ECSClient, d *schema.ResourceData, tagIDs []int) error {
	req := map[string]interface{}{
		"region":       ecsClient.Region(),
		"resourceUuid": d.Get("instance_id").(string),
		"resourceName": d.Get("instance_name").(string),
		"resourceType": 1,
		"tagIds":       tagIDs,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/tags/sync", req, &resp); err != nil {
		return err
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return err
	}
	return nil
}

// expandIntList converts Terraform list/set values to []int for request payloads.
func expandIntList(raw []interface{}) []int {
	out := make([]int, 0, len(raw))
	for _, item := range raw {
		switch v := item.(type) {
		case int:
			out = append(out, v)
		case int32:
			out = append(out, int(v))
		case int64:
			out = append(out, int(v))
		case float64:
			out = append(out, int(v))
		}
	}
	return out
}

func expandIntSet(raw *schema.Set) []int {
	if raw == nil {
		return []int{}
	}
	return expandIntList(raw.List())
}

// extractTagIDs picks tag item IDs from the resource/list response shape.
func extractTagIDs(tags []interface{}) []interface{} {
	out := make([]interface{}, 0, len(tags))
	for _, raw := range tags {
		tag, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, helper.IntFromMap(tag, "id"))
	}
	return out
}
