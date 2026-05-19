package elb

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const elbListenersListPath = "/elb/openapi/v2/listeners/list"

// DataSourceENELBListeners returns the data source schema for ELB listeners on a load balancer.
func DataSourceENELBListeners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENELBListenersRead,
		Description: "Data source to list EdgeNext ELB listeners for a load balancer.",
		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Load balancer ID to list listeners for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "Maximum number of listeners to return in one request.",
			},
			"listeners": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Listeners returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Listener ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Listener name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Listener protocol (for example TERMINATED_HTTPS, HTTP).",
						},
						"protocol_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Protocol port.",
						},
						"connection_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Connection limit (-1 if unlimited).",
						},
						"provisioning_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Provisioning status.",
						},
						"operating_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating status.",
						},
						"created_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation time as Unix timestamp (seconds).",
						},
						"default_target_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default target group ID when set.",
						},
						"insert_headers": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Insert headers configuration (string key to string value, for example X-Forwarded-For, X-Forwarded-Port, X-Forwarded-Proto).",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of listeners matching the query (from the API when present).",
			},
			"hasmore": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether more results exist beyond this page.",
			},
		},
	}
}

func dataSourceENELBListenersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"loadbalancer_id": strings.TrimSpace(d.Get("loadbalancer_id").(string)),
		"limit":           d.Get("limit").(int),
	}

	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbListenersListPath, req, &resp); err != nil {
		return diag.Errorf("failed to list ELB listeners: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB listeners list response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "listeners")
	listeners := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		listeners = append(listeners, elbListenerAttrsFromMap(row))
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(listeners) > 0 {
		total = len(listeners)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hasmore", boolFromInterface(payload["hasmore"])); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("listeners", listeners); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "loadbalancer_id", "limit")
	return nil
}

func elbListenerAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                      helper.StringFromMap(m, "id"),
		"name":                    helper.StringFromMap(m, "name"),
		"description":             helper.StringFromMap(m, "description"),
		"protocol":                helper.StringFromMap(m, "protocol"),
		"protocol_port":           helper.IntFromMap(m, "protocol_port"),
		"connection_limit":        helper.IntFromMap(m, "connection_limit"),
		"provisioning_status":     helper.StringFromMap(m, "provisioning_status"),
		"operating_status":        helper.StringFromMap(m, "operating_status"),
		"created_at":              helper.IntFromMap(m, "created_at"),
		"default_target_group_id": helper.StringFromMap(m, "default_pool_id"),
		"insert_headers":          elbInsertHeadersStringMap(m["insert_headers"]),
	}
}

func elbInsertHeadersStringMap(v interface{}) map[string]interface{} {
	if v == nil {
		return map[string]interface{}{}
	}
	raw, ok := v.(map[string]interface{})
	if !ok || len(raw) == 0 {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(raw))
	for k, val := range raw {
		switch t := val.(type) {
		case string:
			out[k] = t
		case bool:
			if t {
				out[k] = "true"
			} else {
				out[k] = "false"
			}
		default:
			out[k] = fmt.Sprintf("%v", t)
		}
	}
	return out
}
