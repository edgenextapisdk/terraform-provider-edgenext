package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSKeyPairs returns the data source schema for ECS key_pairs.
func DataSourceENECSKeyPairs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSKeyPairsRead,
		Description: "Data source to query EdgeNext ECS key_pairs.",
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of key_pairs to return.",
			},
			"key_pairs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS key_pairs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the key_pair.",
						},
						"fingerprint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fingerprint of the key_pair.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public key material.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key type (e.g. ssh).",
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of matched key_pairs.",
			},
		},
	}
}

func dataSourceENECSKeyPairsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{}
	if limit, ok := d.GetOk("limit"); ok {
		req["limit"] = limit.(int)
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/keypair/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS key_pairs: %s", err)
	}

	dataList, err := helper.ParseAPIResponseList(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS key_pairs response: %s", err)
	}
	total := 0
	if payload, err := helper.ParseAPIResponseMap(resp); err == nil {
		total = helper.IntFromMap(payload, "count")
		if total == 0 {
			total = helper.IntFromMap(payload, "total")
		}
	}
	flat := make([]interface{}, 0, len(dataList))
	for _, raw := range dataList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if inner, ok := row["keypair"].(map[string]interface{}); ok {
			flat = append(flat, keyPairAttrsFromMap(inner))
			continue
		}
		flat = append(flat, keyPairAttrsFromMap(row))
	}
	if total == 0 && len(flat) > 0 {
		total = len(flat)
	}
	if err := d.Set("total", total); err != nil {
		return diag.FromErr(err)
	}
	helper.SetDataSourceStableID(d, "limit")
	if err := d.Set("key_pairs", flat); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func keyPairAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":        helper.StringFromMap(m, "name"),
		"fingerprint": helper.StringFromMap(m, "fingerprint"),
		"public_key":  helper.StringFromMap(m, "public_key"),
		"type":        helper.StringFromMap(m, "type"),
	}
}
