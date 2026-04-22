package ecs

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceENECSVpc returns the resource schema for ECS vpc.
func ResourceENECSVpc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSVpcCreate,
		ReadContext:   resourceENECSVpcRead,
		UpdateContext: resourceENECSVpcUpdate,
		DeleteContext: resourceENECSVpcDelete,
		CustomizeDiff: resourceENECSVpcCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSVpcImport,
		},
		Description: "Provides an EdgeNext ECS vpc resource. subnet and its nested fields cannot be changed after creation.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionResourceSchema("region description"),
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name description",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description description",
			},
			"subnet": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Subnet configuration used when creating the VPC. Cannot be changed after creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet name. Cannot be changed after creation.",
						},
						"ip_version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     4,
							Description: "Subnet IP version. Cannot be changed after creation.",
						},
						"cidr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet CIDR. Cannot be changed after creation.",
						},
					},
				},
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The primary IPv4 CIDR returned by the API.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC status.",
			},
			"total_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of IPs in the VPC.",
			},
			"used_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Used number of IPs in the VPC.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project ID.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time.",
			},
		},
	}
}

func resourceENECSVpcCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip this check during creation.
	if d.Id() == "" {
		return nil
	}
	if d.HasChange("subnet") {
		oldRaw, newRaw := d.GetChange("subnet")
		if strings.TrimSpace(fmt.Sprintf("%v", oldRaw)) != strings.TrimSpace(fmt.Sprintf("%v", newRaw)) {
			return fmt.Errorf("subnet cannot be modified after creation")
		}
	}
	return nil
}

func resourceENECSVpcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	subnetCfg, err := vpcSubnetFromState(d)
	if err != nil {
		return diag.FromErr(err)
	}
	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"network": map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		},
		"subnet": map[string]interface{}{
			"name":       subnetCfg["name"].(string),
			"ip_version": subnetCfg["ip_version"].(int),
			"cidr":       subnetCfg["cidr"].(string),
		},
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/create", req, &resp)
	if err != nil {
		return diag.Errorf("failed to create ECS vpc: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS vpc create response: %s", err)
	}
	network := helper.MapFromMap(payload, "network")
	if network == nil {
		return diag.Errorf("failed to parse ECS vpc create response: missing network")
	}
	networkID := helper.StringFromMap(network, "id")
	if networkID == "" {
		return diag.Errorf("failed to parse ECS vpc create response: missing network.id")
	}
	d.SetId(networkID)

	return resourceENECSVpcRead(ctx, d, m)
}

func resourceENECSVpcImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import id as region/network_id, got %q", d.Id())
	}

	region := helper.NormalizeRegion(parts[0])
	if err := d.Set("region", region); err != nil {
		return nil, err
	}
	d.SetId(parts[1])

	if diags := resourceENECSVpcRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("vpc %q not found in region %q", parts[1], region)
	}

	return []*schema.ResourceData{d}, nil
}

func resourceENECSVpcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":     helper.NormalizeRegion(d.Get("region").(string)),
		"network_id": d.Id(),
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/detail", req, &resp)
	if err != nil {
		d.SetId("") // assume destroyed
		return nil
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS vpc detail response: %s", err)
	}

	// vpc/detail currently returns fields directly under data (not nested in "network").
	network := payload
	if name, ok := network["name"].(string); ok {
		_ = d.Set("name", name)
	}
	if val, ok := network["description"]; ok {
		_ = d.Set("description", val)
	}
	_ = d.Set("status", helper.StringFromMap(network, "status"))
	_ = d.Set("total_ips", helper.IntFromMap(network, "total_ips"))
	_ = d.Set("used_ips", helper.IntFromMap(network, "used_ips"))
	_ = d.Set("project_id", helper.StringFromMap(network, "project_id"))
	_ = d.Set("created_at", helper.StringFromMap(network, "created_at"))
	_ = d.Set("updated_at", helper.StringFromMap(network, "updated_at"))

	cidr := ""
	if ipv4CIDRs, ok := network["ipv4_cidrs"].([]interface{}); ok && len(ipv4CIDRs) > 0 {
		if first, ok := ipv4CIDRs[0].(string); ok {
			cidr = first
		}
	}
	if cidr != "" {
		_ = d.Set("cidr", cidr)
	}

	subnetCfg, err := vpcSubnetFromState(d)
	if err == nil {
		if cidr != "" {
			subnetCfg["cidr"] = cidr
		}
		if err := d.Set("subnet", []interface{}{subnetCfg}); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceENECSVpcUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	// Defense in depth: CustomizeDiff blocks this at plan time; reject here if Update is still invoked.
	if d.HasChange("subnet") {
		return diag.Errorf("subnet cannot be updated after creation")
	}

	if !d.HasChanges("name", "description") {
		return resourceENECSVpcRead(ctx, d, m)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"network": map[string]interface{}{
			"network_id":  d.Id(),
			"name":        d.Get("name"),
			"description": d.Get("description"),
		},
	}
	var resp map[string]interface{}
	err = ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/update", req, &resp)
	if err != nil {
		return diag.Errorf("failed to update ECS vpc: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS vpc update response: %s", err)
	}

	return resourceENECSVpcRead(ctx, d, m)
}

func resourceENECSVpcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":      helper.NormalizeRegion(d.Get("region").(string)),
		"network_ids": []string{d.Id()},
	}
	var resp map[string]interface{}

	err = ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/delete", req, &resp)
	if err != nil {
		return diag.Errorf("failed to delete ECS vpc: %s", err)
	}
	payload, err := helper.ParseAPIResponsePayload(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS vpc delete response: %s", err)
	}
	if m, ok := payload.(map[string]interface{}); ok {
		if status, ok := m[d.Id()].(string); !ok || status != "ok" {
			return diag.Errorf("ECS vpc delete: unexpected status for id %q: %v", d.Id(), m[d.Id()])
		}
	}

	return nil
}

func vpcSubnetFromState(d *schema.ResourceData) (map[string]interface{}, error) {
	rawSubnet, ok := d.GetOk("subnet")
	if !ok {
		return nil, fmt.Errorf("missing required subnet configuration")
	}
	subnetList, ok := rawSubnet.([]interface{})
	if !ok || len(subnetList) == 0 {
		return nil, fmt.Errorf("missing required subnet configuration")
	}
	subnetCfg, ok := subnetList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid subnet configuration")
	}
	return subnetCfg, nil
}
