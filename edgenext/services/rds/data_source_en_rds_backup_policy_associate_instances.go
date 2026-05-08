package rds

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceENRDSBackupPolicyAssociateInstances queries associated instances of a backup policy.
func DataSourceENRDSBackupPolicyAssociateInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENRDSBackupPolicyAssociateInstancesRead,
		Description: "Data source to query instance associations for an EdgeNext RDS backup policy.",
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Backup policy ID.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Associated instances returned by the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of associated sources returned by the list API.",
			},
		},
	}
}

func dataSourceENRDSBackupPolicyAssociateInstancesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := d.Get("policy_id").(string)
	sources, total, diags := rdsBackupPolicySourcesList(ctx, rdsClient, policyID)
	if diags.HasError() {
		return diags
	}

	if err := d.Set("instances", flattenRDSBackupPolicyInstances(sources)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "policy_id")
	return nil
}

func flattenRDSBackupPolicyInstances(sources []map[string]interface{}) []interface{} {
	out := make([]interface{}, 0, len(sources))
	for _, src := range sources {
		out = append(out, map[string]interface{}{
			"instance_type": helper.StringFromMap(src, "source_type"),
			"instance_id":   helper.StringFromMap(src, "source_id"),
			"instance_name": helper.StringFromMap(src, "source_name"),
		})
	}
	return out
}
