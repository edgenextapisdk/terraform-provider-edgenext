package eip

import (
	"context"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	eipAssociationRelationPath = "/ecs/openapi/v2/floatingips/ports/fixed_ip/relation"
	eipFloatingIPListPath      = "/ecs/openapi/v2/floatingips/list"
	ecsInstanceListPortsPath   = "/ecs/openapi/v2/instance/list_ports"
	ecsPortInternalIPListPath  = "/ecs/openapi/v2/ports/internal_ip_list"
	elbLoadBalancerGetPath     = "/elb/openapi/v2/load_balancers/get"

	eipInstanceTypeEcs              = "EcsInstance"
	eipInstanceTypeElb              = "ElbInstance"
	eipInstanceTypeNetworkInterface = "NetworkInterface"
)

var eipInstanceTypes = []string{
	eipInstanceTypeEcs,
	eipInstanceTypeElb,
	eipInstanceTypeNetworkInterface,
}

// ResourceENEIPAssociation associates an EIP (floating IP) with an ECS instance, ELB, or network interface.
func ResourceENEIPAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENEIPAssociationCreate,
		ReadContext:   resourceENEIPAssociationRead,
		DeleteContext: resourceENEIPAssociationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Associates an EdgeNext EIP (floating IP) with an ECS instance, ELB load balancer, or network interface. " +
			"No in-place update; changing arguments forces replacement.",
		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "EIP (floating IP) ID.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Target instance ID. For EcsInstance and ElbInstance this is the server or load balancer ID; for NetworkInterface this is the network interface (port) ID.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      eipInstanceTypeEcs,
				ForceNew:     true,
				Description:  "Instance type: EcsInstance, ElbInstance, or NetworkInterface.",
				ValidateFunc: validation.StringInSlice(eipInstanceTypes, false),
			},
			"fixed_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Fixed (private) IP used for the association. Set on create when the target has multiple network interfaces or fixed IPs; otherwise the first available fixed IP is chosen. On read, populated from the API.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public floating IP address from the API.",
			},
		},
	}
}

func resourceENEIPAssociationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	allocationID := strings.TrimSpace(d.Get("allocation_id").(string))
	instanceID := strings.TrimSpace(d.Get("instance_id").(string))
	instanceType := strings.TrimSpace(d.Get("instance_type").(string))
	wantFixedIP := strings.TrimSpace(d.Get("fixed_ip_address").(string))

	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	if diags := eipEnsureFloatingIPAvailableForAssociation(ctx, ecsClient, allocationID); diags.HasError() {
		return diags
	}

	portID, fixedIP, diags := eipResolveAssociationTarget(ctx, client, instanceType, instanceID, wantFixedIP)
	if diags.HasError() {
		return diags
	}

	req := map[string]interface{}{
		"action": "add",
		"floating_ip": map[string]interface{}{
			"id":               allocationID,
			"port_id":          portID,
			"fixed_ip_address": fixedIP,
		},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, eipAssociationRelationPath, req, &resp); err != nil {
		return diag.Errorf("failed to associate EIP %q: %s", allocationID, err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse EIP association create response: %s", err)
	}

	d.SetId(allocationID)
	return resourceENEIPAssociationRead(ctx, d, m)
}

func resourceENEIPAssociationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	allocationID := strings.TrimSpace(d.Id())
	if allocationID == "" {
		d.SetId("")
		return nil
	}

	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	fip, err := eipFloatingIPByID(ctx, ecsClient, allocationID)
	if err != nil {
		return diag.FromErr(err)
	}
	if fip == nil {
		d.SetId("")
		return nil
	}

	if !eipFloatingIPIsAssociated(fip) {
		d.SetId("")
		return nil
	}
	fixedIP := helper.StringFromMap(fip, "fixed_ip_address")

	_ = d.Set("allocation_id", helper.StringFromMap(fip, "id"))
	_ = d.Set("floating_ip_address", helper.StringFromMap(fip, "floating_ip_address"))
	_ = d.Set("fixed_ip_address", fixedIP)
	return nil
}

func resourceENEIPAssociationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	allocationID := strings.TrimSpace(d.Get("allocation_id").(string))
	if allocationID == "" {
		allocationID = strings.TrimSpace(d.Id())
	}
	if allocationID == "" {
		return nil
	}

	fip, err := eipFloatingIPByID(ctx, ecsClient, allocationID)
	if err != nil {
		return diag.FromErr(err)
	}
	if fip == nil {
		return nil
	}
	if !eipFloatingIPIsAssociated(fip) {
		return nil
	}

	req := map[string]interface{}{
		"action": "remove",
		"floating_ip": map[string]interface{}{
			"id": allocationID,
		},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, eipAssociationRelationPath, req, &resp); err != nil {
		return diag.Errorf("failed to disassociate EIP: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse EIP association delete response: %s", err)
	}
	return nil
}

// eipResolveAssociationTarget resolves port_id and fixed_ip_address for the association API.
func eipResolveAssociationTarget(ctx context.Context, client *connectivity.EdgeNextClient, instanceType, instanceID, wantFixedIP string) (portID, fixedIP string, diags diag.Diagnostics) {
	switch instanceType {
	case eipInstanceTypeEcs:
		return eipResolveEcsInstanceTarget(ctx, client, instanceID, wantFixedIP)
	case eipInstanceTypeElb:
		return eipResolveElbInstanceTarget(ctx, client, instanceID, wantFixedIP)
	case eipInstanceTypeNetworkInterface:
		return eipResolveNetworkInterfaceTarget(ctx, client, instanceID, wantFixedIP)
	default:
		return "", "", diag.Errorf("unsupported instance_type %q", instanceType)
	}
}

func eipResolveEcsInstanceTarget(ctx context.Context, client *connectivity.EdgeNextClient, serverID, wantFixedIP string) (string, string, diag.Diagnostics) {
	ecsClient, err := client.ECSClient()
	if err != nil {
		return "", "", diag.FromErr(err)
	}
	req := map[string]interface{}{
		"server_id": strings.TrimSpace(serverID),
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, ecsInstanceListPortsPath, req, &resp); err != nil {
		return "", "", diag.Errorf("failed to list ports for ECS instance %q: %s", serverID, err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return "", "", diag.Errorf("failed to parse instance list_ports response: %s", err)
	}
	ports := helper.ListFromMap(payload, "ports")
	if len(ports) == 0 {
		return "", "", diag.Errorf("ECS instance %q has no network interfaces", serverID)
	}

	wantFixedIP = strings.TrimSpace(wantFixedIP)
	if wantFixedIP != "" {
		for _, raw := range ports {
			port, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			pid := helper.StringFromMap(port, "id")
			if ip := eipFixedIPOnPort(port, wantFixedIP); ip != "" {
				return pid, ip, nil
			}
		}
		return "", "", diag.Errorf("fixed_ip_address %q not found on any network interface of ECS instance %q", wantFixedIP, serverID)
	}

	for _, raw := range ports {
		port, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		pid := helper.StringFromMap(port, "id")
		if ip := eipFirstFixedIPOnPort(port); ip != "" {
			return pid, ip, nil
		}
	}
	return "", "", diag.Errorf("ECS instance %q has no fixed IP addresses on its network interfaces", serverID)
}

func eipResolveElbInstanceTarget(ctx context.Context, client *connectivity.EdgeNextClient, loadBalancerID, wantFixedIP string) (string, string, diag.Diagnostics) {
	elbClient, err := client.ELBClient()
	if err != nil {
		return "", "", diag.FromErr(err)
	}
	req := map[string]interface{}{
		"loadbalancer_id": strings.TrimSpace(loadBalancerID),
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbLoadBalancerGetPath, req, &resp); err != nil {
		return "", "", diag.Errorf("failed to get ELB load balancer %q: %s", loadBalancerID, err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return "", "", diag.Errorf("failed to parse load balancer get response: %s", err)
	}
	lb := helper.MapFromMap(payload, "loadbalancer")
	if lb == nil {
		return "", "", diag.Errorf("load balancer %q not found in API response", loadBalancerID)
	}
	portID := helper.StringFromMap(lb, "vip_port_id")
	fixedIP := helper.StringFromMap(lb, "vip_address")
	if portID == "" || fixedIP == "" {
		return "", "", diag.Errorf("load balancer %q is missing vip_port_id or vip_address", loadBalancerID)
	}
	wantFixedIP = strings.TrimSpace(wantFixedIP)
	if wantFixedIP != "" && wantFixedIP != fixedIP {
		return "", "", diag.Errorf("fixed_ip_address %q does not match load balancer VIP %q", wantFixedIP, fixedIP)
	}
	return portID, fixedIP, nil
}

func eipResolveNetworkInterfaceTarget(ctx context.Context, client *connectivity.EdgeNextClient, portID, wantFixedIP string) (string, string, diag.Diagnostics) {
	ecsClient, err := client.ECSClient()
	if err != nil {
		return "", "", diag.FromErr(err)
	}
	portID = strings.TrimSpace(portID)
	req := map[string]interface{}{
		"id": portID,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, ecsPortInternalIPListPath, req, &resp); err != nil {
		return "", "", diag.Errorf("failed to list internal IPs for network interface %q: %s", portID, err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return "", "", diag.Errorf("failed to parse ports internal_ip_list response: %s", err)
	}
	rows := helper.ListFromMap(payload, "data")
	if len(rows) == 0 {
		return "", "", diag.Errorf("network interface %q has no fixed IP addresses", portID)
	}

	wantFixedIP = strings.TrimSpace(wantFixedIP)
	if wantFixedIP != "" {
		for _, raw := range rows {
			row, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			if helper.StringFromMap(row, "fixed_ip_address") == wantFixedIP {
				return portID, wantFixedIP, nil
			}
		}
		return "", "", diag.Errorf("fixed_ip_address %q not found on network interface %q", wantFixedIP, portID)
	}

	row, ok := rows[0].(map[string]interface{})
	if !ok {
		return "", "", diag.Errorf("invalid internal_ip_list entry for network interface %q", portID)
	}
	fixedIP := helper.StringFromMap(row, "fixed_ip_address")
	if fixedIP == "" {
		return "", "", diag.Errorf("network interface %q has no fixed IP addresses", portID)
	}
	return portID, fixedIP, nil
}

func eipFixedIPOnPort(port map[string]interface{}, want string) string {
	for _, raw := range helper.InterfaceToList(port["fixed_ips"]) {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(item, "ip_address") == want {
			return want
		}
	}
	for _, raw := range helper.InterfaceToList(port["ipv4"]) {
		if s, ok := raw.(string); ok && strings.TrimSpace(s) == want {
			return want
		}
	}
	return ""
}

func eipFirstFixedIPOnPort(port map[string]interface{}) string {
	for _, raw := range helper.InterfaceToList(port["fixed_ips"]) {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if ip := helper.StringFromMap(item, "ip_address"); ip != "" {
			return ip
		}
	}
	for _, raw := range helper.InterfaceToList(port["ipv4"]) {
		if s, ok := raw.(string); ok {
			if ip := strings.TrimSpace(s); ip != "" {
				return ip
			}
		}
	}
	return ""
}

func eipFloatingIPIsAssociated(fip map[string]interface{}) bool {
	return helper.StringFromMap(fip, "port_id") != "" && helper.StringFromMap(fip, "fixed_ip_address") != ""
}

// eipEnsureFloatingIPAvailableForAssociation rejects create when the EIP is missing or already bound to a port.
func eipEnsureFloatingIPAvailableForAssociation(ctx context.Context, ecsClient *connectivity.ECSClient, allocationID string) diag.Diagnostics {
	fip, err := eipFloatingIPByID(ctx, ecsClient, allocationID)
	if err != nil {
		return diag.FromErr(err)
	}
	if fip == nil {
		return diag.Errorf("EIP %q not found", allocationID)
	}
	if !eipFloatingIPIsAssociated(fip) {
		return nil
	}
	return diag.Errorf(
		"EIP %q is already associated (floating_ip_address=%q, port_id=%q, fixed_ip_address=%q); disassociate it before creating edgenext_eip_association",
		allocationID,
		helper.StringFromMap(fip, "floating_ip_address"),
		helper.StringFromMap(fip, "port_id"),
		helper.StringFromMap(fip, "fixed_ip_address"),
	)
}

func eipFloatingIPByID(ctx context.Context, ecsClient *connectivity.ECSClient, allocationID string) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"id":    strings.TrimSpace(allocationID),
		"limit": 100,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, eipFloatingIPListPath, req, &resp); err != nil {
		return nil, err
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, err
	}
	for _, raw := range helper.ListFromMap(payload, "floating_ip") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(row, "id") == allocationID {
			return row, nil
		}
	}
	return nil, nil
}
