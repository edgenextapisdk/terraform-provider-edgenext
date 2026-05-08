package rds

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceENRDSAccounts returns the data source schema for database users on an RDS instance.
func DataSourceENRDSAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENRDSAccountsRead,
		Description: "Data source to list database users (accounts) on an EdgeNext RDS instance.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "RDS instance ID to list database users for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Database users returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User name.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host pattern the user is allowed to connect from (for example % or 127.0.0.1).",
						},
						"databases": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Database names granted to this user when returned by the API; may be empty.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceENRDSAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	rawList, diags := rdsListDatabaseUsers(ctx, rdsClient, instanceID)
	if diags.HasError() {
		return diags
	}

	users := make([]interface{}, 0, len(rawList))
	for _, row := range rawList {
		users = append(users, map[string]interface{}{
			"user_name": helper.StringFromMap(row, "name"),
			"host":      helper.StringFromMap(row, "host"),
			"databases": rdsUserDatabasesFromMap(row),
		})
	}

	if err := d.Set("users", users); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "instance_id")
	return nil
}

func rdsUserDatabasesFromMap(row map[string]interface{}) []interface{} {
	raw := helper.ListFromMap(row, "databases")
	out := make([]interface{}, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
