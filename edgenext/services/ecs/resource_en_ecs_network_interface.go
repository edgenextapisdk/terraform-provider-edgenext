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

// ResourceENECSNetworkInterface returns the resource schema for ECS network_interface (Neutron port / ENI).
func ResourceENECSNetworkInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSNetworkInterfaceCreate,
		ReadContext:   resourceENECSNetworkInterfaceRead,
		UpdateContext: resourceENECSNetworkInterfaceUpdate,
		DeleteContext: resourceENECSNetworkInterfaceDelete,
		CustomizeDiff: resourceENECSNetworkInterfaceCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSNetworkInterfaceImport,
		},
		Description: "Provides an EdgeNext ECS network interface (ENI). Updating name/description is supported in place. vpc_id and subnet_id cannot be changed after creation.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network interface name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Network interface description.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC network ID. Cannot be changed after creation.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subnet ID for the primary fixed IP. Cannot be changed after creation.",
			},
			"port_security_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether network interface security is enabled.",
			},
			"security_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security group IDs to apply to the network interface.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Attached instance ID to the network interface.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP address bound to this network interface.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant ID.",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Administrative state of the network interface.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network interface status.",
			},
			"instance_owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance owner (e.g. compute:nova).",
			},
			"fixed_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Fixed IP assignments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Fixed IP address.",
						},
						"floating_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated floating IP if present.",
						},
					},
				},
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID.",
			},
			"qos_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "QoS policy ID.",
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
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MAC address.",
			},
			"binding_vnic_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VNIC binding type.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resolved instance name when attached.",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resolved VPC name.",
			},
			"ipv4": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IPv4 addresses.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ipv6": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IPv6 addresses.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceENECSNetworkInterfaceCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip this check during creation.
	if d.Id() == "" {
		return nil
	}
	if d.HasChange("vpc_id") {
		oldRaw, newRaw := d.GetChange("vpc_id")
		if networkInterfaceNormalizeString(oldRaw) != networkInterfaceNormalizeString(newRaw) {
			return fmt.Errorf("vpc_id cannot be modified after creation")
		}
	}
	if d.HasChange("subnet_id") {
		oldRaw, newRaw := d.GetChange("subnet_id")
		if networkInterfaceNormalizeString(oldRaw) != networkInterfaceNormalizeString(newRaw) {
			return fmt.Errorf("subnet_id cannot be modified after creation")
		}
	}
	return nil
}

func parseNetworkInterfaceResourceID(id string) (networkInterfaceID string, err error) {
	networkInterfaceID = strings.TrimSpace(id)
	if networkInterfaceID == "" {
		return "", fmt.Errorf("expected id as network_interface_id, got %q", id)
	}
	return networkInterfaceID, nil
}

func resourceENECSNetworkInterfaceImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	portID, err := parseNetworkInterfaceResourceID(d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(portID)

	if diags := resourceENECSNetworkInterfaceRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("network interface %q not found", portID)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENECSNetworkInterfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"port": map[string]interface{}{
			"id":          "",
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"network_id":  d.Get("vpc_id").(string),
			"subnet_id":   d.Get("subnet_id").(string),
		},
	}
	var resp map[string]interface{}

	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/add", req, &resp); err != nil {
		return diag.Errorf("failed to create ECS network_interface: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS network_interface create response: %s", err)
	}
	port := helper.MapFromMap(payload, "port")
	if port == nil {
		return diag.Errorf("failed to parse ECS network_interface create response: missing data.port")
	}
	portID := helper.StringFromMap(port, "id")
	if portID == "" {
		return diag.Errorf("failed to parse ECS network_interface create response: missing port id")
	}
	d.SetId(portID)
	if relationEnabled, relationGroups, shouldCallRelation := networkInterfaceSecurityRelationInput(d); shouldCallRelation {
		if err := resourceENECSNetworkInterfaceRelationSecurityGroup(ctx, ecsClient, portID, relationEnabled, relationGroups); err != nil {
			return diag.Errorf("failed to set security relation on create for ECS network_interface: %s", err)
		}
	}

	return resourceENECSNetworkInterfaceRead(ctx, d, m)
}

func resourceENECSNetworkInterfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	portID, err := parseNetworkInterfaceResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"id": portID,
	}
	var resp map[string]interface{}

	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/detail", req, &resp); err != nil {
		return diag.Errorf("failed to read ECS network_interface: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS network_interface detail response: %s", err)
	}
	port := helper.MapFromMap(payload, "port")
	if port == nil {
		d.SetId("")
		return nil
	}
	if err := resourceENECSNetworkInterfaceEnrichFixedIPs(ctx, ecsClient, portID, port); err != nil {
		return diag.Errorf("failed to enrich ECS network_interface fixed IP details: %s", err)
	}

	if err := resourceENECSNetworkInterfaceApplyPortToState(d, port); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceENECSNetworkInterfaceApplyPortToState(d *schema.ResourceData, port map[string]interface{}) error {
	flat := flattenNetworkInterfacePort(port)
	delete(flat, "id")
	for k, v := range flat {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	if err := d.Set("subnet_id", networkInterfaceFirstSubnetID(port)); err != nil {
		return err
	}
	return d.Set("floating_ip_address", networkInterfaceFirstFloatingIPAddress(port))
}

func networkInterfaceFirstSubnetID(port map[string]interface{}) string {
	for _, raw := range helper.InterfaceToList(port["fixed_ips"]) {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if sid := helper.StringFromMap(m, "subnet_id"); sid != "" {
			return sid
		}
	}
	return ""
}

func networkInterfaceFirstFloatingIPAddress(port map[string]interface{}) string {
	for _, raw := range helper.InterfaceToList(port["fixed_ips"]) {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if ip := helper.StringFromMap(m, "floating_ip"); ip != "" {
			return ip
		}
	}
	return ""
}

func resourceENECSNetworkInterfaceEnrichFixedIPs(ctx context.Context, ecsClient *connectivity.ECSClient, portID string, port map[string]interface{}) error {
	req := map[string]interface{}{
		"id": portID,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/internal_ip_list", req, &resp); err != nil {
		return err
	}

	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return err
	}
	internalRows := helper.ListFromMap(payload, "data")
	if len(internalRows) == 0 {
		return nil
	}

	floatingByFixed := make(map[string]string, len(internalRows))
	for _, raw := range internalRows {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		fixedIP := helper.StringFromMap(row, "fixed_ip_address")
		if fixedIP == "" {
			continue
		}
		floatingByFixed[fixedIP] = helper.StringFromMap(row, "floating_ip_address")
	}
	if len(floatingByFixed) == 0 {
		return nil
	}

	fixedIPs := helper.InterfaceToList(port["fixed_ips"])
	updated := make([]interface{}, 0, len(fixedIPs))
	for _, raw := range fixedIPs {
		item, ok := raw.(map[string]interface{})
		if !ok {
			updated = append(updated, raw)
			continue
		}
		ipAddress := helper.StringFromMap(item, "ip_address")
		if ipAddress != "" {
			if floating, ok := floatingByFixed[ipAddress]; ok {
				item["floating_ip"] = floating
			}
		}
		updated = append(updated, item)
	}
	port["fixed_ips"] = updated
	return nil
}

func resourceENECSNetworkInterfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Defense in depth: CustomizeDiff blocks these at plan time; reject here if Update is still invoked.
	if d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		return diag.Errorf("vpc_id and subnet_id cannot be updated after creation")
	}

	nameOrDescChanged := d.HasChange("name") || d.HasChange("description")
	securityRelationChanged := d.HasChange("port_security_enabled") || d.HasChange("security_groups")
	if !nameOrDescChanged && !securityRelationChanged {
		return resourceENECSNetworkInterfaceRead(ctx, d, m)
	}

	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	portID, err := parseNetworkInterfaceResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if nameOrDescChanged {
		req := map[string]interface{}{
			"port": map[string]interface{}{
				"id":          portID,
				"name":        d.Get("name").(string),
				"description": d.Get("description").(string),
			},
		}
		var resp map[string]interface{}
		if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/add", req, &resp); err != nil {
			return diag.Errorf("failed to update ECS network_interface: %s", err)
		}
		if _, err := helper.ParseAPIResponseMap(resp); err != nil {
			return diag.Errorf("failed to parse ECS network_interface update response: %s", err)
		}
	}

	if securityRelationChanged {
		relationEnabled, relationGroups, _ := networkInterfaceSecurityRelationInput(d)
		if err := resourceENECSNetworkInterfaceRelationSecurityGroup(ctx, ecsClient, portID, relationEnabled, relationGroups); err != nil {
			return diag.Errorf("failed to update security relation for ECS network_interface: %s", err)
		}
	}

	return resourceENECSNetworkInterfaceRead(ctx, d, m)
}

func resourceENECSNetworkInterfaceRelationServer(ctx context.Context, ecsClient *connectivity.ECSClient, portID, action, serverID string) error {
	req := map[string]interface{}{
		"port_id":   portID,
		"action":    action,
		"server_id": serverID,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/relation/server", req, &resp); err != nil {
		return err
	}
	_, err := helper.ParseAPIResponsePayload(resp)
	return err
}

func resourceENECSNetworkInterfaceRelationSecurityGroup(ctx context.Context, ecsClient *connectivity.ECSClient, portID string, portSecurityEnabled bool, securityGroups []interface{}) error {
	req := map[string]interface{}{
		"port": map[string]interface{}{
			"id":                    portID,
			"port_security_enabled": portSecurityEnabled,
			"security_groups":       securityGroups,
		},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/relation/security_group", req, &resp); err != nil {
		return err
	}
	_, err := helper.ParseAPIResponsePayload(resp)
	return err
}

func networkInterfaceSecurityRelationInput(d *schema.ResourceData) (bool, []interface{}, bool) {
	enabledRaw, enabledSet := d.GetOkExists("port_security_enabled")
	groupsRaw, groupsSet := d.GetOk("security_groups")

	if !enabledSet && !groupsSet {
		return false, []interface{}{}, false
	}

	enabled := false
	if enabledSet {
		if b, ok := enabledRaw.(bool); ok {
			enabled = b
		}
	}
	return enabled, helper.InterfaceToStringSlice(groupsRaw), true
}

func resourceENECSNetworkInterfacePortDetail(ctx context.Context, ecsClient *connectivity.ECSClient, portID string) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"id": portID,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/detail", req, &resp); err != nil {
		return nil, err
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, err
	}
	port := helper.MapFromMap(payload, "port")
	if port == nil {
		return nil, fmt.Errorf("missing port in detail response")
	}
	return port, nil
}

func networkInterfaceNormalizeString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(t)
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", t))
	}
}

func resourceENECSNetworkInterfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	portID, err := parseNetworkInterfaceResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"ids": []interface{}{portID},
	}
	var resp map[string]interface{}

	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/ports/delete", req, &resp); err != nil {
		return diag.Errorf("failed to delete ECS network_interface: %s", err)
	}
	if _, err := helper.ParseAPIResponsePayload(resp); err != nil {
		return diag.Errorf("failed to parse ECS network_interface delete response: %s", err)
	}

	return nil
}
