package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSInstanceTags returns the data source schema for ECS instances by tag filters.
func DataSourceENECSInstanceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSInstanceTagsRead,
		Description: "Data source to query EdgeNext ECS instances by tag filters.",
		Schema: map[string]*schema.Schema{
			"tag_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The tag ID to filter instances.",
			},
			"tag_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The tag key to filter instances.",
			},
			"tag_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The tag value to filter instances.",
			},
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for instance tag listing.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Page size for instance tag listing.",
			},
			"instance_tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of instances matched by tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The record ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance name.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance type, e.g. ECS.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance region.",
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
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched instances.",
			},
		},
	}
}

func dataSourceENECSInstanceTagsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":   ecsClient.Region(),
		"tagId":    d.Get("tag_id").(int),
		"tagKey":   d.Get("tag_key").(string),
		"tagValue": d.Get("tag_value").(string),
		"pageNum":  d.Get("page_num").(int),
		"pageSize": d.Get("page_size").(int),
	}
	var resp map[string]interface{}

	err = ecsClient.Get(ctx, "/ecs/openapi/v2/resource/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS instance tags: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS instance tags response: %s", err)
	}
	dataList := helper.ListFromMap(payload, "list")
	items := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, map[string]interface{}{
			"id":            helper.IntFromMap(row, "id"),
			"instance_id":   helper.StringFromMap(row, "resourceId"),
			"instance_name": helper.StringFromMap(row, "resourceName"),
			"instance_type": helper.StringFromMap(row, "productType"),
			"region":        helper.StringFromMap(row, "region"),
			"tag_count":     helper.IntFromMap(row, "tagCount"),
			"tags":          normalizeENECSInstanceTagItems(helper.ListFromMap(row, "tags")),
		})
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(items) > 0 {
		total = len(items)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("instance_tags", items); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "tag_id", "tag_key", "tag_value", "page_num", "page_size")

	return nil
}

func normalizeENECSInstanceTagItems(items []interface{}) []interface{} {
	out := make([]interface{}, 0, len(items))
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":    helper.IntFromMap(item, "id"),
			"key":   helper.StringFromMap(item, "key"),
			"value": helper.StringFromMap(item, "value"),
		})
	}
	return out
}
