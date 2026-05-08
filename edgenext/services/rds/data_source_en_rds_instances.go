package rds

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const rdsInstancesListDatastoreTypeMySQL = "mysql"

// DataSourceENRDSInstances returns the data source schema for RDS instances.
func DataSourceENRDSInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENRDSInstancesRead,
		Description: "Data source to query EdgeNext RDS instances.",
		Schema: map[string]*schema.Schema{
			"page_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for the list request. The API request body uses the field name page_number.",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "Page size for the list request.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Filter by instance ID. Use empty string to omit the filter.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Filter by instance name. Use empty string to omit the filter.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of RDS instances returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance status.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance hostname.",
						},
						"flavor": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Flavor specification.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flavor ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flavor name.",
									},
								},
							},
						},
						"datastore": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Datastore engine and version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine type.",
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine version label.",
									},
									"version_number": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Engine version number.",
									},
								},
							},
						},
						"volume": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Primary storage volume.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Volume size in GB.",
									},
								},
							},
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region where the instance runs.",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time (RFC3339).",
						},
						"updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time (RFC3339).",
						},
						"operating_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating status of the instance.",
						},
						"access": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Network access settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_public": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the instance is exposed on a public network.",
									},
								},
							},
						},
						"ip": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP addresses attached to the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP address.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Address type (for example private).",
									},
									"network": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Network or subnet ID.",
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC name.",
									},
								},
							},
						},
						"root_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether root login is enabled.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of instances matching the query.",
			},
		},
	}
}

func dataSourceENRDSInstancesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"page_size":      d.Get("page_size").(int),
		"page_number":    d.Get("page_num").(int),
		"datastore_type": rdsInstancesListDatastoreTypeMySQL,
		"id":             d.Get("instance_id").(string),
		"name":           d.Get("name").(string),
	}

	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/instances/list", req, &resp); err != nil {
		return diag.Errorf("failed to list RDS instances: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse RDS instances response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "instances")
	instances := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		instances = append(instances, rdsInstanceAttrsFromMap(row))
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(instances) > 0 {
		total = len(instances)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("instances", instances); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "page_num", "page_size", "instance_id", "name")
	return nil
}

func rdsInstanceAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	flavor := make([]interface{}, 0)
	if f := helper.MapFromMap(m, "flavor"); f != nil {
		flavor = append(flavor, map[string]interface{}{
			"id":   helper.StringFromMap(f, "id"),
			"name": helper.StringFromMap(f, "name"),
		})
	}

	datastore := make([]interface{}, 0)
	if ds := helper.MapFromMap(m, "datastore"); ds != nil {
		datastore = append(datastore, map[string]interface{}{
			"type":           helper.StringFromMap(ds, "type"),
			"version":        helper.StringFromMap(ds, "version"),
			"version_number": helper.StringFromMap(ds, "version_number"),
		})
	}

	volume := make([]interface{}, 0)
	if vol := helper.MapFromMap(m, "volume"); vol != nil {
		volume = append(volume, map[string]interface{}{
			"size": helper.IntFromMap(vol, "size"),
		})
	}

	access := make([]interface{}, 0)
	if acc := helper.MapFromMap(m, "access"); acc != nil {
		access = append(access, map[string]interface{}{
			"is_public": boolFromInterface(acc["is_public"]),
		})
	}

	ipList := make([]interface{}, 0)
	for _, raw := range helper.ListFromMap(m, "ip") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		ipList = append(ipList, map[string]interface{}{
			"address":  helper.StringFromMap(row, "address"),
			"type":     helper.StringFromMap(row, "type"),
			"network":  helper.StringFromMap(row, "network"),
			"vpc_name": helper.StringFromMap(row, "vpc_name"),
		})
	}

	return map[string]interface{}{
		"id":               helper.StringFromMap(m, "id"),
		"name":             helper.StringFromMap(m, "name"),
		"status":           helper.StringFromMap(m, "status"),
		"hostname":         helper.StringFromMap(m, "hostname"),
		"flavor":           flavor,
		"datastore":        datastore,
		"volume":           volume,
		"region":           helper.StringFromMap(m, "region"),
		"created":          helper.StringFromMap(m, "created"),
		"updated":          helper.StringFromMap(m, "updated"),
		"operating_status": helper.StringFromMap(m, "operating_status"),
		"access":           access,
		"ip":               ipList,
		"root_enabled":     boolFromInterface(m["root_enabled"]),
	}
}

func boolFromInterface(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
