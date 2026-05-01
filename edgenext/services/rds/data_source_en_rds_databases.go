package rds

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceENRDSDatabases returns the data source schema for databases on an RDS instance.
func DataSourceENRDSDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENRDSDatabasesRead,
		Description: "Data source to list databases on an EdgeNext RDS instance.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "RDS instance ID to list databases for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"databases": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Databases returned by the API for the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
						},
					},
				},
			},
		},
	}
}

func dataSourceENRDSDatabasesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	req := map[string]interface{}{
		"instance_id": instanceID,
	}

	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databases/list", req, &resp); err != nil {
		return diag.Errorf("failed to list RDS databases: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse RDS databases response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "databases")
	databases := make([]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		databases = append(databases, map[string]interface{}{
			"name": helper.StringFromMap(row, "name"),
		})
	}

	if err := d.Set("databases", databases); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "instance_id")
	return nil
}
