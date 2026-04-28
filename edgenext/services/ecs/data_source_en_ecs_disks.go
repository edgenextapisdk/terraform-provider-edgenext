package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSDisks returns the data source schema for ECS disks.
func DataSourceENECSDisks() *schema.Resource {
	attachmentElem := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the instance this disk is attached to.",
			},
			"device": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Device path on the instance (e.g. /dev/vda).",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the instance this disk is attached to.",
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceENECSDisksRead,
		Description: "Data source to query EdgeNext ECS disks via GET /ecs/openapi/v2/volume/list.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Disk name filter (empty string lists all names).",
			},
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for listing.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Page size for listing.",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of disks reported by the API.",
			},
			"disks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Disks returned for the current page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk name.",
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size in GB.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "API status field.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk product type (API field type), e.g. Quick Disk.",
						},
						"disk_label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk label, e.g. System Disk.",
						},
						"billing_model": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing model code.",
						},
						"billing_model_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Billing model display name.",
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time if applicable.",
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated server name from API (may be empty).",
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image name associated with the disk.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk description.",
						},
						"disk_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk status code.",
						},
						"status_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Human-readable status, e.g. in-use, available.",
						},
						"volume_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume type code (API field volumeType).",
						},
						"policy_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Backup or policy names attached to the disk.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp.",
						},
						"attachment": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Attachment records when the disk is mounted on an instance.",
							Elem:        attachmentElem,
						},
					},
				},
			},
		},
	}
}

func dataSourceENECSDisksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"name":      d.Get("name").(string),
		"page_num":  d.Get("page_num").(int),
		"page_size": d.Get("page_size").(int),
	}
	var resp map[string]interface{}

	err = ecsClient.Get(ctx, "/ecs/openapi/v2/volume/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS disks: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS disks response: %s", err)
	}

	dataList := helper.ListFromMap(payload, "list")
	items := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		items = append(items, diskAttrsFromMap(row))
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(items) > 0 {
		total = len(items)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "name", "page_num", "page_size")
	if err := d.Set("disks", items); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func diskAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                 helper.StringFromMap(m, "id"),
		"name":               helper.StringFromMap(m, "name"),
		"size":               helper.IntFromMap(m, "size"),
		"status":             helper.IntFromMap(m, "status"),
		"disk_type":          helper.StringFromMap(m, "type"),
		"disk_label":         helper.StringFromMap(m, "disk_label"),
		"billing_model":      helper.IntFromMap(m, "billing_model"),
		"billing_model_name": helper.StringFromMap(m, "billing_model_name"),
		"expiration_time":    helper.StringFromMap(m, "expiration_time"),
		"server_name":        helper.StringFromMap(m, "server_name"),
		"image_name":         helper.StringFromMap(m, "image_name"),
		"description":        helper.StringFromMap(m, "description"),
		"disk_status":        helper.IntFromMap(m, "disk_status"),
		"status_name":        helper.StringFromMap(m, "status_name"),
		"volume_type":        helper.IntFromMap(m, "volumeType"),
		"policy_names":       policyNamesFromDisk(m),
		"created_at":         helper.StringFromMap(m, "created_at"),
		"attachment":         normalizeDiskAttachments(helper.ListFromMap(m, "attachment")),
	}
}

func policyNamesFromDisk(m map[string]interface{}) []interface{} {
	raw := helper.ListFromMap(m, "policy_names")
	out := make([]interface{}, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func normalizeDiskAttachments(items []interface{}) []interface{} {
	out := make([]interface{}, 0, len(items))
	for _, raw := range items {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, map[string]interface{}{
			"instance_name": helper.StringFromMap(row, "instance_name"),
			"device":        helper.StringFromMap(row, "device"),
			"instance_id":   helper.StringFromMap(row, "instance_id"),
		})
	}
	return out
}
