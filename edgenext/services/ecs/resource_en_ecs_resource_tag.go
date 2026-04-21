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

// ResourceENECSResourceTag manages tag bindings on a specific ECS resource.
func ResourceENECSResourceTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSResourceTagCreate,
		ReadContext:   resourceENECSResourceTagRead,
		UpdateContext: resourceENECSResourceTagUpdate,
		DeleteContext: resourceENECSResourceTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSResourceTagImport,
		},
		Description: "Provides an EdgeNext ECS resource tag binding resource.",
		Schema: map[string]*schema.Schema{
			"resource_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target resource UUID.",
			},
			"resource_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target resource name.",
			},
			"region": helper.RegionResourceSchema("region description"),
			"resource_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The target resource type code.",
			},
			"tag_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Tag IDs to bind to the target resource.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceENECSResourceTagImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf("expected import id as region/resource_uuid/resource_name/resource_type, got %q", d.Id())
	}
	region := helper.NormalizeRegion(parts[0])
	resourceUUID := strings.TrimSpace(parts[1])
	resourceName := strings.TrimSpace(parts[2])
	resourceTypeRaw := strings.TrimSpace(parts[3])
	if region == "" || resourceUUID == "" || resourceName == "" || resourceTypeRaw == "" {
		return nil, fmt.Errorf("expected import id as region/resource_uuid/resource_name/resource_type, got %q", d.Id())
	}
	resourceType, err := strconv.Atoi(resourceTypeRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid resource_type %q in import id %q", resourceTypeRaw, d.Id())
	}
	if err := d.Set("region", region); err != nil {
		return nil, err
	}
	if err := d.Set("resource_uuid", resourceUUID); err != nil {
		return nil, err
	}
	if err := d.Set("resource_name", resourceName); err != nil {
		return nil, err
	}
	if err := d.Set("resource_type", resourceType); err != nil {
		return nil, err
	}
	d.SetId(resourceUUID)
	if diags := resourceENECSResourceTagRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("resource tag binding for resource %q not found in region %q", resourceUUID, region)
	}
	return []*schema.ResourceData{d}, nil
}

// Create binds the provided tag IDs to the target resource.
func resourceENECSResourceTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := syncENECSResourceTags(ctx, ecsClient, d, expandIntList(d.Get("tag_ids").([]interface{}))); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("resource_uuid").(string))
	return resourceENECSResourceTagRead(ctx, d, m)
}

// Read refreshes current tag bindings from the resource query endpoint.
func resourceENECSResourceTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"resourceUuid": d.Get("resource_uuid").(string),
		"resourceName": d.Get("resource_name").(string),
		"region":       d.Get("region").(string),
		"resourceType": d.Get("resource_type").(int),
		"pageNum":      1,
		"pageSize":     100,
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
		return diag.Errorf("failed to parse ECS resource tag read response: %s", err)
	}
	list := helper.ListFromMap(payload, "list")
	uuid := d.Get("resource_uuid").(string)

	// Find the exact target resource and synchronize state fields from API output.
	for _, raw := range list {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(item, "resourceId") != uuid {
			continue
		}
		_ = d.Set("resource_name", helper.StringFromMap(item, "resourceName"))
		_ = d.Set("region", helper.StringFromMap(item, "region"))
		_ = d.Set("resource_type", d.Get("resource_type").(int))
		_ = d.Set("tag_ids", extractTagIDs(helper.ListFromMap(item, "tags")))
		return nil
	}
	// If the resource is not found, set the ID to empty and return nil.
	d.SetId("")
	return nil
}

// Update only syncs tag_ids changes; identity fields are ForceNew.
func resourceENECSResourceTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("tag_ids") {
		return resourceENECSResourceTagRead(ctx, d, m)
	}

	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := syncENECSResourceTags(ctx, ecsClient, d, expandIntList(d.Get("tag_ids").([]interface{}))); err != nil {
		return diag.FromErr(err)
	}

	return resourceENECSResourceTagRead(ctx, d, m)
}

// Delete clears tag bindings by syncing an empty tagIds list.
func resourceENECSResourceTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := syncENECSResourceTags(ctx, ecsClient, d, []int{}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// syncENECSResourceTags calls the tags/sync endpoint and validates API envelope.
func syncENECSResourceTags(ctx context.Context, ecsClient *connectivity.ECSClient, d *schema.ResourceData, tagIDs []int) error {
	req := map[string]interface{}{
		"resourceUuid": d.Get("resource_uuid").(string),
		"resourceName": d.Get("resource_name").(string),
		"region":       d.Get("region").(string),
		"resourceType": d.Get("resource_type").(int),
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

// expandIntList converts Terraform list values to []int for request payloads.
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
