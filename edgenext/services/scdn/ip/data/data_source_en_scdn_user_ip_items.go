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

// DataSourceEdgenextScdnUserIpItems returns the SCDN User IP Items data source
func DataSourceEdgenextScdnUserIpItems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnUserIpItemsRead,

		Schema: map[string]*schema.Schema{
			"user_ip_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "User IP List ID",
			},
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
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by IP",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of user IP items",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark",
						},
						"user_ip_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User IP ID",
						},
						"format_created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Created At",
						},
						"format_updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Updated At",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count",
			},
		},
	}
}

func dataSourceScdnUserIpItemsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.UserIpItemListRequest{
		UserIpID: d.Get("user_ip_id").(int),
		Page:     d.Get("page").(int),
		PerPage:  d.Get("per_page").(int),
		IP:       d.Get("ip").(string),
	}

	log.Printf("[INFO] Listing SCDN User IP Items: %+v", req)
	response, err := service.ListUserIpItems(req)
	if err != nil {
		return fmt.Errorf("failed to list User IP Items: %w", err)
	}

	items := make([]interface{}, 0, len(response.Data.List))
	for _, item := range response.Data.List {
		items = append(items, map[string]interface{}{
			"id":                item.ID,
			"ip":                item.IP,
			"remark":            item.Remark,
			"user_ip_id":        item.UserIpId,
			"format_created_at": item.FormatCreatedAt,
			"format_updated_at": item.FormatUpdatedAt,
		})
	}

	if err := d.Set("items", items); err != nil {
		return fmt.Errorf("error setting items: %w", err)
	}
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))

	return nil
}
