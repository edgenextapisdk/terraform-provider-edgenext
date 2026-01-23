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

func DataSourceEdgenextScdnDomainGroupDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnDomainGroupDomainsRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Group ID to list domains for",
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
				Description: "List of domains",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain ID",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total count",
			},
			"ports": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Common ports",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceScdnDomainGroupDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := domain_group.NewDomainGroupService(client)

	req := domain_group.DomainGroupDomainListRequest{
		GroupID: d.Get("group_id").(int),
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}
	if v, ok := d.GetOk("domain"); ok {
		req.Domain = v.(string)
	}

	log.Printf("[INFO] Listing SCDN Domain Group Domains: %+v", req)
	resp, err := service.ListGroupDomains(req)
	if err != nil {
		return fmt.Errorf("failed to list Domain Group Domains: %w", err)
	}

	if resp.Status.Code != 1 {
		return fmt.Errorf("failed to list Domain Group Domains: %s", resp.Status.Message)
	}

	items := make([]interface{}, 0, len(resp.Data.List))
	for _, item := range resp.Data.List {
		items = append(items, map[string]interface{}{
			"domain_id": item.DomainID,
			"domain":    item.Domain,
		})
	}

	d.Set("list", items)
	d.Set("total", resp.Data.Total)
	d.Set("ports", resp.Data.Ports)
	d.SetId(strconv.FormatInt(time.Now().UnixNano(), 10))

	return nil
}
