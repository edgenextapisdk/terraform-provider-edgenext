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

// ResourceENECSNetworkInterfaceFloatingIPBinding binds a floating IP to a network interface.
// No UpdateContext: all arguments are ForceNew; SDK rejects a superfluous Update in that case.
func ResourceENECSNetworkInterfaceFloatingIPBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSNetworkInterfaceFloatingIPBindingCreate,
		ReadContext:   resourceENECSNetworkInterfaceFloatingIPBindingRead,
		DeleteContext: resourceENECSNetworkInterfaceFloatingIPBindingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSNetworkInterfaceFloatingIPBindingImport,
		},
		Description: "Provides an EdgeNext ECS network interface floating IP binding resource.",
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network interface ID. Changing this forces a new resource.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The floating IP address to bind. Changing this forces a new resource.",
			},
			"fixed_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fixed IP address used for floating IP binding.",
			},
		},
	}
}

func resourceENECSNetworkInterfaceFloatingIPBindingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import id as network_interface_id/floating_ip_address, got %q", d.Id())
	}
	portID := strings.TrimSpace(parts[0])
	floatingIP := strings.TrimSpace(parts[1])
	if portID == "" || floatingIP == "" {
		return nil, fmt.Errorf("expected import id as network_interface_id/floating_ip_address, got %q", d.Id())
	}
	_ = d.Set("network_interface_id", portID)
	_ = d.Set("floating_ip_address", floatingIP)
	d.SetId(fmt.Sprintf("%s/%s", portID, floatingIP))
	if diags := resourceENECSNetworkInterfaceFloatingIPBindingRead(ctx, d, meta); diags.HasError() {
		return nil, fmt.Errorf(diags[0].Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("network interface floating IP binding not found for import id %q", fmt.Sprintf("%s/%s", portID, floatingIP))
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSNetworkInterfaceFloatingIPBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}
	portID := strings.TrimSpace(d.Get("network_interface_id").(string))
	floatingIP := strings.TrimSpace(d.Get("floating_ip_address").(string))

	port, err := resourceENECSNetworkInterfacePortDetail(ctx, ecsClient, portID)
	if err != nil {
		return diag.Errorf("failed to query port detail before floating IP add: %s", err)
	}
	fixedIP := networkInterfaceFirstFixedIPAddress(port)
	if err := resourceENECSNetworkInterfaceRelationFloatingIPAdd(ctx, ecsClient, portID, floatingIP, fixedIP); err != nil {
		return diag.Errorf("failed to bind floating IP to ECS network_interface: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", portID, floatingIP))
	return resourceENECSNetworkInterfaceFloatingIPBindingRead(ctx, d, m)
}

func resourceENECSNetworkInterfaceFloatingIPBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}
	portID := strings.TrimSpace(d.Get("network_interface_id").(string))
	expectedFloatingIP := strings.TrimSpace(d.Get("floating_ip_address").(string))

	port, err := resourceENECSNetworkInterfacePortDetail(ctx, ecsClient, portID)
	if err != nil {
		d.SetId("")
		return nil
	}
	if err := resourceENECSNetworkInterfaceEnrichFixedIPs(ctx, ecsClient, portID, port); err != nil {
		return diag.Errorf("failed to enrich fixed IP details for ECS network_interface: %s", err)
	}
	fixedIPs := helper.InterfaceToList(port["fixed_ips"])
	for _, raw := range fixedIPs {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		floatingIP := strings.TrimSpace(helper.StringFromMap(item, "floating_ip"))
		if floatingIP != expectedFloatingIP {
			continue
		}
		_ = d.Set("fixed_ip_address", helper.StringFromMap(item, "ip_address"))
		return nil
	}
	d.SetId("")
	return nil
}

func resourceENECSNetworkInterfaceFloatingIPBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}
	floatingIP := strings.TrimSpace(d.Get("floating_ip_address").(string))
	if floatingIP != "" {
		if err := resourceENECSNetworkInterfaceRelationFloatingIPRemove(ctx, ecsClient, floatingIP); err != nil {
			return diag.Errorf("failed to unbind floating IP from ECS network_interface: %s", err)
		}
	}
	return nil
}

func resourceENECSNetworkInterfaceRelationFloatingIPAdd(ctx context.Context, ecsClient *connectivity.ECSClient, portID, floatingIPAddr, fixedIPAddr string) error {
	if fixedIPAddr == "" {
		return fmt.Errorf("missing fixed_ip_address for port %s", portID)
	}
	floatingID, err := ecsFloatingIPIDByAddress(ctx, ecsClient, floatingIPAddr)
	if err != nil {
		return err
	}
	req := map[string]interface{}{
		"action": "add",
		"floating_ip": map[string]interface{}{
			"port_id":          portID,
			"fixed_ip_address": fixedIPAddr,
			"id":               floatingID,
		},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/relation/floatingip", req, &resp); err != nil {
		return err
	}
	_, err = helper.ParseAPIResponsePayload(resp)
	return err
}

func resourceENECSNetworkInterfaceRelationFloatingIPRemove(ctx context.Context, ecsClient *connectivity.ECSClient, floatingIPAddr string) error {
	req := map[string]interface{}{
		"action": "remove",
		"floating_ip": map[string]interface{}{
			"floating_ip_address": floatingIPAddr,
		},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/relation/floatingip", req, &resp); err != nil {
		return err
	}
	_, err := helper.ParseAPIResponsePayload(resp)
	return err
}

func ecsFloatingIPIDByAddress(ctx context.Context, ecsClient *connectivity.ECSClient, floatingIPAddr string) (string, error) {
	req := map[string]interface{}{
		"floating_ip_address": floatingIPAddr,
		"limit":               50,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/floatingips/list", req, &resp); err != nil {
		return "", err
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return "", err
	}
	for _, raw := range helper.ListFromMap(payload, "floating_ip") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(row, "floating_ip_address") == floatingIPAddr {
			id := helper.StringFromMap(row, "id")
			if id == "" {
				break
			}
			return id, nil
		}
	}
	return "", fmt.Errorf("floating IP address %q not found", floatingIPAddr)
}

func networkInterfaceFirstFixedIPAddress(port map[string]interface{}) string {
	for _, raw := range helper.InterfaceToList(port["fixed_ips"]) {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if ip := helper.StringFromMap(m, "ip_address"); ip != "" {
			return ip
		}
	}
	return ""
}
