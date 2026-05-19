package elb

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const elbCertificatesListPath = "/elb/openapi/v2/certificates/list"

// DataSourceENELBCertificates returns the data source schema for ELB certificate containers (list).
func DataSourceENELBCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENELBCertificatesRead,
		Description: "Data source to list EdgeNext ELB certificates (Barbican containers).",
		Schema: map[string]*schema.Schema{
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for the list request.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Page size for the list request.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Filter by certificate name. Use empty string to omit the filter.",
			},
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Certificates returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate ID (Barbican container ID).",
						},
						"certificate_ref": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate reference URL (Barbican container reference).",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Container type (for example certificate).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate status.",
						},
						"created": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation time as Unix timestamp (seconds).",
						},
						"updated": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Last update time as Unix timestamp (seconds).",
						},
						"expiration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Expiration time as Unix timestamp (seconds).",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of certificates matching the query.",
			},
		},
	}
}

func dataSourceENELBCertificatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"page_num":  d.Get("page_num").(int),
		"page_size": d.Get("page_size").(int),
		"name":      d.Get("name").(string),
	}

	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbCertificatesListPath, req, &resp); err != nil {
		return diag.Errorf("failed to list ELB certificates: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB certificates list response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "certificates")
	certs := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		certs = append(certs, elbCertificateListItemFromMap(row))
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(certs) > 0 {
		total = len(certs)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("certificates", certs); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "page_num", "page_size", "name")
	return nil
}

func elbCertificateListItemFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"certificate_id":  helper.StringFromMap(m, "container_id"),
		"certificate_ref": helper.StringFromMap(m, "container_ref"),
		"name":            helper.StringFromMap(m, "name"),
		"type":            helper.StringFromMap(m, "type"),
		"status":          helper.StringFromMap(m, "status"),
		"created":         helper.IntFromMap(m, "created"),
		"updated":         helper.IntFromMap(m, "updated"),
		"expiration":      helper.IntFromMap(m, "expiration"),
	}
}
