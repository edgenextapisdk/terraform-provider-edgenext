package ecs

import (
	"context"
	"strconv"

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
		DeleteContext: resourceENECSTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Provides an EdgeNext ECS tag resource.",
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "key description",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "value description",
			},
		},
	}
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
				"key":   d.Get("key").(string),
				"value": d.Get("value").(string),
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
		_ = d.Set("key", key)
	}
	if value := helper.StringFromMap(first, "value"); value != "" {
		_ = d.Set("value", value)
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
		"tagKey":   d.Get("key").(string),
		"tagValue": d.Get("value").(string),
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
		if currentID == idInt || (currentKey == d.Get("key").(string) && currentValue == d.Get("value").(string)) {
			found = true
			_ = d.Set("key", currentKey)
			_ = d.Set("value", currentValue)
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
