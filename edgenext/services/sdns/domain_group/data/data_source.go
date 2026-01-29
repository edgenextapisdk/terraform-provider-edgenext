package data

import (
	"fmt"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextDnsGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsGroupRead,

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by group name",
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of matched groups",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDnsGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	req := sdns.DnsGroupListRequest{
		PerPage: 1000,
	}
	if v, ok := d.GetOk("group_name"); ok {
		req.GroupName = v.(string)
	}

	resp, err := service.ListDnsGroups(req)
	if err != nil {
		return fmt.Errorf("failed to list DNS groups: %w", err)
	}

	groups := make([]map[string]interface{}, 0, len(resp.List))
	for _, info := range resp.List {
		groups = append(groups, map[string]interface{}{
			"id":         strconv.Itoa(info.ID),
			"group_name": info.GroupName,
			"remark":     info.Remark,
		})
	}

	if err := d.Set("groups", groups); err != nil {
		return fmt.Errorf("failed to set groups: %w", err)
	}

	d.SetId(strconv.Itoa(len(groups)))
	return nil
}
