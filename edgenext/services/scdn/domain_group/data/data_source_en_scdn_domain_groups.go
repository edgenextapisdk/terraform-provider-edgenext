package data

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/domain_group"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextScdnDomainGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnDomainGroupsRead,

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by group name",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by domain",
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
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of domain groups",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group ID",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group Name",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark",
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

func dataSourceScdnDomainGroupsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := domain_group.NewDomainGroupService(client)

	req := domain_group.DomainGroupListRequest{
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}
	if v, ok := d.GetOk("group_name"); ok {
		req.GroupName = v.(string)
	}
	if v, ok := d.GetOk("domain"); ok {
		req.Domain = v.(string)
	}

	log.Printf("[INFO] Listing SCDN Domain Groups: %+v", req)
	resp, err := service.ListDomainGroups(req)
	if err != nil {
		return fmt.Errorf("failed to list Domain Groups: %w", err)
	}

	if resp.Status.Code != 1 {
		return fmt.Errorf("failed to list Domain Groups: %s", resp.Status.Message)
	}

	items := make([]interface{}, 0, len(resp.Data.List))
	for _, item := range resp.Data.List {
		items = append(items, map[string]interface{}{
			"id":         item.ID,
			"group_name": item.GroupName,
			"remark":     item.Remark,
			"created_at": item.CreatedAt,
			"updated_at": item.UpdatedAt,
		})
	}

	d.Set("list", items)
	d.Set("total", resp.Data.Total)
	d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))

	return nil
}
