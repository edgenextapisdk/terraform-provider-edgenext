package data

import (
	"fmt"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextDnsDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsDomainRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by domain name",
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of matched domains",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDnsDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	req := sdns.DnsDomainListRequest{
		PerPage: 1000,
	}
	if v, ok := d.GetOk("domain"); ok {
		req.Domain = v.(string)
	}

	resp, err := service.ListDnsDomains(req)
	if err != nil {
		return fmt.Errorf("failed to list DNS domains: %w", err)
	}

	domains := make([]map[string]interface{}, 0, len(resp.List))
	for _, info := range resp.List {
		domains = append(domains, map[string]interface{}{
			"id":     strconv.Itoa(info.ID),
			"domain": info.Domain,
			"status": strconv.Itoa(info.Status),
		})
	}

	if err := d.Set("domains", domains); err != nil {
		return fmt.Errorf("failed to set domains: %w", err)
	}

	d.SetId(strconv.Itoa(len(domains)))
	return nil
}
