package data

import (
	"fmt"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextDnsRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsRecordRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Domain ID to list records for",
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of records in the domain",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"view": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mx": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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

func dataSourceDnsRecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	domainID := d.Get("domain_id").(int)
	req := sdns.DnsRecordListRequest{
		DomainID: domainID,
		PerPage:  1000,
	}

	resp, err := service.ListDnsRecords(req)
	if err != nil {
		return fmt.Errorf("failed to list DNS records: %w", err)
	}

	records := make([]map[string]interface{}, 0, len(resp.List))
	for _, info := range resp.List {
		records = append(records, map[string]interface{}{
			"id":     strconv.Itoa(info.ID),
			"name":   info.Name,
			"type":   info.Type,
			"view":   info.View,
			"value":  info.Value,
			"mx":     strconv.Itoa(info.MX),
			"ttl":    strconv.Itoa(info.TTL),
			"status": strconv.Itoa(info.Status),
			"remark": info.Remark,
		})
	}

	if err := d.Set("records", records); err != nil {
		return fmt.Errorf("failed to set records: %w", err)
	}

	d.SetId(strconv.Itoa(domainID))
	return nil
}
