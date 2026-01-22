package data

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnUserIps returns the SCDN User IP Lists data source
func DataSourceEdgenextScdnUserIps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnUserIpsRead,

		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number",
			},
			"per_page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Items per page",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of user IP lists",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark",
						},
						"item_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Item Num",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created At",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated At",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total count",
			},
		},
	}
}

func dataSourceScdnUserIpsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.UserIpListRequest{
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}

	log.Printf("[INFO] Listing SCDN User IPs: %+v", req)
	response, err := service.ListUserIps(req)
	if err != nil {
		return fmt.Errorf("failed to list User IPs: %w", err)
	}

	items := make([]interface{}, 0, len(response.Data.List))
	for _, item := range response.Data.List {
		items = append(items, map[string]interface{}{
			"id":         item.ID,
			"name":       item.Name,
			"remark":     item.Remark,
			"item_num":   item.ItemNum,
			"created_at": item.CreatedAt,
			"updated_at": item.UpdatedAt,
		})
	}

	if err := d.Set("items", items); err != nil {
		return fmt.Errorf("error setting items: %w", err)
	}
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Always generate a new ID to force refresh
	d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))

	return nil
}
