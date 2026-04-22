package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECSInstances returns the data source schema for ECS instances.
func DataSourceENECSInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceENECSInstancesRead,
		Description: "Data source to query EdgeNext ECS instances.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionDataSchema("region description"),
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name to filter instances.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of instances to return.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the instance.",
						},
						"flavor": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The flavor name of the instance.",
						},
						"flavor_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Flavor detail information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vcpus": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of vCPUs.",
									},
									"ram": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The RAM size in MB.",
									},
								},
							},
						},
						"image_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The image name of the instance.",
						},
						"fixed_addresses": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of fixed IP addresses.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"floating_addresses": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of floating IP addresses.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the instance.",
						},
						"instance_cost_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance billing and expiration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_cost_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance billing type.",
									},
									"network_cost_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The network billing type.",
									},
									"instance_expiration_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance expiration time.",
									},
									"billing_model": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The billing model code.",
									},
								},
							},
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of tag names.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of matched instances.",
			},
		},
	}
}

func dataSourceENECSInstancesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	region := d.Get("region").(string)
	req := map[string]interface{}{
		"region": region,
	}
	if name, ok := d.GetOk("name"); ok {
		req["name"] = name.(string)
	}
	if limit, ok := d.GetOk("limit"); ok {
		req["limit"] = limit.(int)
	}
	var resp map[string]interface{}

	// List action
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/instance/list", req, &resp)
	if err != nil {
		return diag.Errorf("failed to read ECS instances: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS instances response: %s", err)
	}
	serversRaw := helper.ListFromMap(payload, "servers")
	instances := make([]interface{}, 0, len(serversRaw))
	for _, raw := range serversRaw {
		server, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		instances = append(instances, instanceAttrsFromMap(server))
	}
	if countRaw, ok := payload["count"]; ok {
		count := helper.IntFromMap(map[string]interface{}{"count": countRaw}, "count")
		if err := d.Set("total", count); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := d.Set("total", len(instances)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("instances", instances); err != nil {
		return diag.FromErr(err)
	}

	helper.SetDataSourceStableID(d, "region", "name", "limit")

	return nil
}

func instanceAttrsFromMap(m map[string]interface{}) map[string]interface{} {
	flavorInfo := []interface{}{}
	if raw := helper.MapFromMap(m, "flavor_info"); raw != nil {
		flavorInfo = append(flavorInfo, map[string]interface{}{
			"vcpus": helper.IntFromMap(raw, "vcpus"),
			"ram":   helper.IntFromMap(raw, "ram"),
		})
	}
	instanceCostInfo := []interface{}{}
	if raw := helper.MapFromMap(m, "instance_cost_info"); raw != nil {
		instanceCostInfo = append(instanceCostInfo, map[string]interface{}{
			"instance_cost_type":       helper.StringFromMap(raw, "instance_cost_type"),
			"network_cost_type":        helper.StringFromMap(raw, "network_cost_type"),
			"instance_expiration_time": helper.StringFromMap(raw, "instance_expiration_time"),
			"billing_model":            helper.IntFromMap(raw, "billing_model"),
		})
	}

	return map[string]interface{}{
		"id":                 helper.StringFromMap(m, "id"),
		"name":               helper.StringFromMap(m, "name"),
		"status":             helper.StringFromMap(m, "status"),
		"flavor":             helper.StringFromMap(m, "flavor"),
		"flavor_info":        flavorInfo,
		"image_name":         helper.StringFromMap(m, "image_name"),
		"fixed_addresses":    helper.InterfaceToStringSlice(m["fixed_addresses"]),
		"floating_addresses": helper.InterfaceToStringSlice(m["floating_addresses"]),
		"created_at":         helper.StringFromMap(m, "created_at"),
		"instance_cost_info": instanceCostInfo,
		"tags":               helper.InterfaceToStringSlice(m["tags"]),
	}
}
