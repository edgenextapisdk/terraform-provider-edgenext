package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSTags returns the data source schema for ECS tags.
func DataSourceENECSTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSTagsRead,
		Description: "Data source to query EdgeNext ECS tags.",
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The tag key to filter tags.",
			},
			"tag_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The tag value to filter tags.",
			},
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for tag listing.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Page size for tag listing.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the tag.",
						},
						"tag_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key of the tag.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the tag.",
						},
						"resource_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The count of resources using this tag.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of tags.",
			},
		},
	}
}

func dataSourceENECSTagsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"tagKey":   d.Get("tag_key").(string),
		"tagValue": d.Get("tag_value").(string),
		"pageNum":  d.Get("page_num").(int),
		"pageSize": d.Get("page_size").(int),
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Get(ctx, "/ecs/openapi/v2/tags/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS tags: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS tags response: %s", err)
	}
	dataList := helper.ListFromMap(payload, "list")
	tags := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		tags = append(tags, map[string]interface{}{
			"id":             helper.IntFromMap(row, "id"),
			"tag_key":        helper.StringFromMap(row, "tagKey"),
			"tag_value":      helper.StringFromMap(row, "tagValue"),
			"resource_count": helper.IntFromMap(row, "resourceCount"),
		})
	}
	if err := d.Set("total", helper.IntFromMap(payload, "total")); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "tag_key", "tag_value", "page_num", "page_size")
	if err := d.Set("tags", tags); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
