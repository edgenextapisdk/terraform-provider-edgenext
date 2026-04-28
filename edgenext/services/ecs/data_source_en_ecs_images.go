package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSImages returns the data source schema for ECS images.
func DataSourceENECSImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSImagesRead,
		Description: "Data source to query EdgeNext ECS images.",
		Schema: map[string]*schema.Schema{
			"visibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "public",
				Description: "Image visibility to filter by.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name to filter images.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Image status to filter by.",
			},
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for image listing.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Page size for image listing.",
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS images.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the image.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the image.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the image.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the image.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The image type.",
						},
						"visibility": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The visibility of the image.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the image in bytes.",
						},
						"min_ram": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum RAM required.",
						},
						"min_disk": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum disk required.",
						},
						"os_distro": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS distribution of the image.",
						},
						"os_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS version of the image.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the image.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time of the image.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of images.",
			},
		},
	}
}

func dataSourceENECSImagesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"visibility": d.Get("visibility").(string),
		"name":       d.Get("name").(string),
		"status":     d.Get("status").(string),
		"page_num":   d.Get("page_num").(int),
		"page_size":  d.Get("page_size").(int),
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/image/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS images: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS images response: %s", err)
	}
	imagesRaw := helper.ListFromMap(payload, "images")
	images := make([]interface{}, 0, len(imagesRaw))
	for _, raw := range imagesRaw {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		images = append(images, imageAttrsFromMap(row))
	}
	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(images) > 0 {
		total = len(images)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "visibility", "name", "status", "page_num", "page_size")
	if err := d.Set("images", images); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func imageAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          helper.StringFromMap(m, "id"),
		"name":        helper.StringFromMap(m, "name"),
		"description": helper.StringFromMap(m, "description"),
		"status":      helper.StringFromMap(m, "status"),
		"image_type":  helper.StringFromMap(m, "image_type"),
		"visibility":  helper.StringFromMap(m, "visibility"),
		"size":        helper.IntFromMap(m, "size"),
		"min_ram":     helper.IntFromMap(m, "min_ram"),
		"min_disk":    helper.IntFromMap(m, "min_disk"),
		"os_distro":   helper.StringFromMap(m, "os_distro"),
		"os_version":  helper.StringFromMap(m, "os_version"),
		"created_at":  helper.StringFromMap(m, "created_at"),
		"updated_at":  helper.StringFromMap(m, "updated_at"),
	}
}
