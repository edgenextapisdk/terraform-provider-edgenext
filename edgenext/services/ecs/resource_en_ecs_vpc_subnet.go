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

// ResourceENECSVpcSubnet returns the resource schema for ECS vpc subnet.
func ResourceENECSVpcSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSVpcSubnetCreate,
		ReadContext:   resourceENECSVpcSubnetRead,
		UpdateContext: resourceENECSVpcSubnetUpdate,
		DeleteContext: resourceENECSVpcSubnetDelete,
		CustomizeDiff: resourceENECSVpcSubnetCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSVpcSubnetImport,
		},
		Description: "Provides an EdgeNext ECS vpc subnet resource. Except region, arguments cannot be changed after creation.",
		Schema: map[string]*schema.Schema{
			"region": helper.RegionResourceSchema("The region of the subnet."),
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC network ID. Cannot be changed after creation.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet name. Cannot be changed after creation.",
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4,
				Description: "IP version. Cannot be changed after creation.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet CIDR. Cannot be changed after creation.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant ID.",
			},
			"subnetpool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subnet pool ID.",
			},
			"enable_dhcp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether DHCP is enabled.",
			},
			"ipv6_ra_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv6 RA mode.",
			},
			"ipv6_address_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv6 address mode.",
			},
			"gateway_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway IP.",
			},
			"allocation_pools": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Allocation pools.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start IP.",
						},
						"end": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End IP.",
						},
					},
				},
			},
			"host_routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Host routes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route destination.",
						},
						"nexthop": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route next hop.",
						},
					},
				},
			},
			"dns_nameservers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DNS nameservers.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description.",
			},
			"service_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Service types.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
			"revision_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Revision number.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID.",
			},
			"used_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Used IP count.",
			},
			"total_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total IP count.",
			},
			"port_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Port count.",
			},
			"not_bind_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason if subnet is not bindable.",
			},
			"router_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Bound router ID.",
			},
		},
	}
}

func resourceENECSVpcSubnetCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip this check during creation.
	if d.Id() == "" {
		return nil
	}
	immutableFields := []string{"network_id", "name", "ip_version", "cidr"}
	for _, field := range immutableFields {
		if !d.HasChange(field) {
			continue
		}
		oldRaw, newRaw := d.GetChange(field)
		if strings.TrimSpace(fmt.Sprintf("%v", oldRaw)) != strings.TrimSpace(fmt.Sprintf("%v", newRaw)) {
			return fmt.Errorf("%s cannot be modified after creation", field)
		}
	}
	return nil
}

func resourceENECSVpcSubnetImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected import id as region/network_id/subnet_id, got %q", d.Id())
	}

	region := helper.NormalizeRegion(parts[0])
	if err := d.Set("region", region); err != nil {
		return nil, err
	}
	if err := d.Set("network_id", parts[1]); err != nil {
		return nil, err
	}
	d.SetId(parts[2])

	if diags := resourceENECSVpcSubnetRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("subnet %q not found under network %q in region %q", parts[2], parts[1], region)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSVpcSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region": helper.NormalizeRegion(d.Get("region").(string)),
		"network": map[string]interface{}{
			"network_id": d.Get("network_id").(string),
		},
		"subnet": map[string]interface{}{
			"name":       d.Get("name").(string),
			"ip_version": d.Get("ip_version").(int),
			"cidr":       d.Get("cidr").(string),
		},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/subnets_create", req, &resp); err != nil {
		return diag.Errorf("failed to create ECS vpc subnet: %s", err)
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS vpc subnet create response: %s", err)
	}
	subnet := helper.MapFromMap(payload, "subnet")
	if subnet == nil {
		return diag.Errorf("failed to parse ECS vpc subnet create response: missing subnet")
	}
	subnetID := helper.StringFromMap(subnet, "id")
	if subnetID == "" {
		return diag.Errorf("failed to parse ECS vpc subnet create response: missing subnet.id")
	}
	d.SetId(subnetID)

	return resourceENECSVpcSubnetRead(ctx, d, m)
}

func resourceENECSVpcSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	subnets, err := ecsVPCSubnetsList(ctx, ecsClient, d.Get("region").(string), d.Get("network_id").(string))
	if err != nil {
		return diag.Errorf("failed to read ECS vpc subnet: %s", err)
	}
	subnet := findVPCSubnetByID(subnets, d.Id())
	if subnet == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", helper.StringFromMap(subnet, "name"))
	_ = d.Set("tenant_id", helper.StringFromMap(subnet, "tenant_id"))
	_ = d.Set("network_id", helper.StringFromMap(subnet, "network_id"))
	_ = d.Set("ip_version", helper.IntFromMap(subnet, "ip_version"))
	_ = d.Set("subnetpool_id", helper.StringFromMap(subnet, "subnetpool_id"))
	_ = d.Set("enable_dhcp", subnetBool(subnet, "enable_dhcp"))
	_ = d.Set("ipv6_ra_mode", helper.StringFromMap(subnet, "ipv6_ra_mode"))
	_ = d.Set("ipv6_address_mode", helper.StringFromMap(subnet, "ipv6_address_mode"))
	_ = d.Set("gateway_ip", helper.StringFromMap(subnet, "gateway_ip"))
	_ = d.Set("cidr", helper.StringFromMap(subnet, "cidr"))
	_ = d.Set("description", helper.StringFromMap(subnet, "description"))
	_ = d.Set("service_types", helper.InterfaceToStringSlice(subnet["service_types"]))
	_ = d.Set("tags", helper.InterfaceToStringSlice(subnet["tags"]))
	_ = d.Set("created_at", helper.StringFromMap(subnet, "created_at"))
	_ = d.Set("updated_at", helper.StringFromMap(subnet, "updated_at"))
	_ = d.Set("revision_number", helper.IntFromMap(subnet, "revision_number"))
	_ = d.Set("project_id", helper.StringFromMap(subnet, "project_id"))
	_ = d.Set("used_ips", helper.IntFromMap(subnet, "used_ips"))
	_ = d.Set("total_ips", helper.IntFromMap(subnet, "total_ips"))
	_ = d.Set("port_num", helper.IntFromMap(subnet, "port_num"))
	_ = d.Set("not_bind_reason", helper.StringFromMap(subnet, "not_bind_reason"))
	_ = d.Set("router_id", helper.StringFromMap(subnet, "router_id"))
	if err := d.Set("dns_nameservers", helper.InterfaceToStringSlice(subnet["dns_nameservers"])); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allocation_pools", subnetAllocationPools(subnet)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_routes", subnetHostRoutes(subnet)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceENECSVpcSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Defense in depth: CustomizeDiff blocks these at plan time; reject here if Update is still invoked.
	immutableFields := []string{"network_id", "name", "ip_version", "cidr"}
	for _, field := range immutableFields {
		if d.HasChange(field) {
			return diag.Errorf("%s cannot be updated after creation", field)
		}
	}
	return resourceENECSVpcSubnetRead(ctx, d, m)
}

func resourceENECSVpcSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"region":     helper.NormalizeRegion(d.Get("region").(string)),
		"subnet_ids": []string{d.Id()},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/subnets_delete", req, &resp); err != nil {
		return diag.Errorf("failed to delete ECS vpc subnet: %s", err)
	}
	payload, err := helper.ParseAPIResponsePayload(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS vpc subnet delete response: %s", err)
	}
	if m, ok := payload.(map[string]interface{}); ok {
		if status, ok := m[d.Id()].(string); !ok || status != "ok" {
			return diag.Errorf("ECS vpc subnet delete: unexpected status for id %q: %v", d.Id(), m[d.Id()])
		}
	}

	return nil
}

func ecsVPCSubnetsList(ctx context.Context, ecsClient *connectivity.ECSClient, region, networkID string) ([]map[string]interface{}, error) {
	req := map[string]interface{}{
		"region":     helper.NormalizeRegion(region),
		"network_id": networkID,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/vpc/subnets_list", req, &resp); err != nil {
		return nil, err
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, err
	}
	rawSubnets := helper.ListFromMap(payload, "subnets")
	subnets := make([]map[string]interface{}, 0, len(rawSubnets))
	for _, raw := range rawSubnets {
		subnet, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		subnets = append(subnets, subnet)
	}
	return subnets, nil
}

func findVPCSubnetByID(subnets []map[string]interface{}, subnetID string) map[string]interface{} {
	for _, subnet := range subnets {
		if helper.StringFromMap(subnet, "id") == subnetID {
			return subnet
		}
	}
	return nil
}
