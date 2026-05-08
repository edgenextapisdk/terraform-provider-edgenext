package rds

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const rdsBackupListTypeFull = "full"

// DataSourceENRDSBackups returns the data source schema for RDS backups.
func DataSourceENRDSBackups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENRDSBackupsRead,
		Description: "Data source to query EdgeNext RDS backups (list API always uses type full).",
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
			"backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by backup ID. Maps to the API field id. Omit or leave empty to list without this filter.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by backup name. Omit or leave empty to match the API empty-name filter.",
			},
			"backups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Backups returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup name.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RDS instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RDS instance name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup type (for example full).",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Backup status.",
						},
						"size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Backup size in GB (as returned by the API).",
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
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
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
										Description: "Engine version.",
									},
								},
							},
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of backups matching the query.",
			},
		},
	}
}

// rdsBackupsListPost calls POST /rds/openapi/v2/backups/list with type fixed to full.
func rdsBackupsListPost(ctx context.Context, rdsClient *connectivity.RDSClient, pageNum, pageSize int, idFilter, nameFilter string) (map[string]interface{}, diag.Diagnostics) {
	req := map[string]interface{}{
		"page_size":   pageSize,
		"page_number": pageNum,
		"type":        rdsBackupListTypeFull,
		"id":          idFilter,
		"name":        nameFilter,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backups/list", req, &resp); err != nil {
		return nil, diag.Errorf("failed to list RDS backups: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, diag.Errorf("failed to parse RDS backups response: %s", err)
	}
	return payload, nil
}

func dataSourceENRDSBackupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	payload, diags := rdsBackupsListPost(ctx, rdsClient, d.Get("page_num").(int), d.Get("page_size").(int), d.Get("backup_id").(string), d.Get("name").(string))
	if diags.HasError() {
		return diags
	}

	rawList := helper.ListFromMap(payload, "backups")
	backups := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		backups = append(backups, rdsBackupAttrsFromMap(row))
	}

	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(backups) > 0 {
		total = len(backups)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("backups", backups); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "page_num", "page_size", "backup_id", "name")
	return nil
}

func rdsBackupAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	datastore := make([]interface{}, 0)
	if ds := helper.MapFromMap(m, "datastore"); ds != nil {
		datastore = append(datastore, map[string]interface{}{
			"type":    helper.StringFromMap(ds, "type"),
			"version": helper.StringFromMap(ds, "version"),
		})
	}

	return map[string]interface{}{
		"id":            helper.StringFromMap(m, "id"),
		"name":          helper.StringFromMap(m, "name"),
		"instance_id":   helper.StringFromMap(m, "instance_id"),
		"instance_name": helper.StringFromMap(m, "instance_name"),
		"type":          helper.StringFromMap(m, "type"),
		"status":        helper.StringFromMap(m, "status"),
		"size":          floatFromInterface(m["size"]),
		"created":       helper.StringFromMap(m, "created"),
		"updated":       helper.StringFromMap(m, "updated"),
		"region":        helper.StringFromMap(m, "region"),
		"datastore":     datastore,
	}
}

func floatFromInterface(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}
