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

// ResourceENECSSecurityGroupRule returns the resource schema for a single ECS security group rule.
func ResourceENECSSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENECSSecurityGroupRuleCreate,
		ReadContext:   resourceENECSSecurityGroupRuleRead,
		UpdateContext: resourceENECSSecurityGroupRuleUpdate,
		DeleteContext: resourceENECSSecurityGroupRuleDelete,
		CustomizeDiff: resourceENECSSecurityGroupRuleCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENECSSecurityGroupRuleImport,
		},
		Description: "Provides a single EdgeNext ECS security group rule. Arguments cannot be changed after creation.",
		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The security group ID this rule belongs to. Cannot be changed after creation.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protocol name (e.g. tcp, udp, icmp). Cannot be changed after creation.",
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Traffic direction: ingress or egress. Cannot be changed after creation.",
			},
			"ethertype": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP version (e.g. IPv4, IPv6). Cannot be changed after creation.",
			},
			"port_range_min": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Minimum port number. Cannot be changed after creation.",
			},
			"port_range_max": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum port number. Cannot be changed after creation.",
			},
			"remote_ip_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remote CIDR (e.g. 192.168.0.0/24). Cannot be changed after creation.",
			},
			"remote_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remote security group ID. Leave empty when using remote_ip_prefix only. Cannot be changed after creation.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule description. Cannot be changed after creation.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tenant ID of the rule.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Project ID of the rule.",
			},
			"revision_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Revision number.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Rule tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceENECSSecurityGroupRuleCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	// Skip this check during creation.
	if d.Id() == "" {
		return nil
	}
	immutableFields := []string{
		"security_group_id",
		"protocol",
		"direction",
		"ethertype",
		"port_range_min",
		"port_range_max",
		"remote_ip_prefix",
		"remote_group_id",
		"description",
	}
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

func resourceENECSSecurityGroupRuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import id as security_group_id/rule_id, got %q", d.Id())
	}
	if err := d.Set("security_group_id", parts[0]); err != nil {
		return nil, err
	}
	d.SetId(parts[1])

	if diags := resourceENECSSecurityGroupRuleRead(ctx, d, meta); diags.HasError() {
		errDiag := diags[0]
		if errDiag.Detail != "" {
			return nil, fmt.Errorf("%s: %s", errDiag.Summary, errDiag.Detail)
		}
		return nil, fmt.Errorf("%s", errDiag.Summary)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("security group rule %q not found under security group %q", parts[1], parts[0])
	}
	return []*schema.ResourceData{d}, nil
}

func ecsSecurityGroupDetailRules(ctx context.Context, ecsClient *connectivity.ECSClient, securityGroupID string) ([]map[string]interface{}, error) {
	req := map[string]interface{}{
		"id": securityGroupID,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/security_group/detail", req, &resp); err != nil {
		return nil, err
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, err
	}
	sg := helper.MapFromMap(payload, "security_group")
	if sg == nil {
		return nil, fmt.Errorf("missing security_group in detail response")
	}
	rulesRaw := helper.ListFromMap(sg, "security_group_rules")
	out := make([]map[string]interface{}, 0, len(rulesRaw))
	for _, raw := range rulesRaw {
		rule, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, rule)
	}
	return out, nil
}

func findSecurityGroupRuleByID(rules []map[string]interface{}, ruleID string) map[string]interface{} {
	for _, rule := range rules {
		if helper.StringFromMap(rule, "id") == ruleID {
			return rule
		}
	}
	return nil
}

func resourceENECSSecurityGroupRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	sgRule := map[string]interface{}{
		"security_group_id": d.Get("security_group_id").(string),
		"protocol":          d.Get("protocol").(string),
		"direction":         d.Get("direction").(string),
		"ethertype":         d.Get("ethertype").(string),
		"port_range_min":    d.Get("port_range_min").(int),
		"port_range_max":    d.Get("port_range_max").(int),
		"remote_ip_prefix":  d.Get("remote_ip_prefix").(string),
		"remote_group_id":   d.Get("remote_group_id").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		sgRule["description"] = v.(string)
	}

	req := map[string]interface{}{
		"security_group_rule": sgRule,
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/security_group_rule/add", req, &resp); err != nil {
		return diag.Errorf("failed to create ECS security_group_rule: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS security_group_rule create response: %s", err)
	}
	created := helper.MapFromMap(payload, "security_group_rule")
	if created == nil {
		return diag.Errorf("failed to parse ECS security_group_rule create response: missing security_group_rule")
	}
	createdID := helper.StringFromMap(created, "id")
	if createdID == "" {
		return diag.Errorf("failed to parse ECS security_group_rule create response: missing id")
	}
	d.SetId(createdID)
	return resourceENECSSecurityGroupRuleRead(ctx, d, m)
}

func resourceENECSSecurityGroupRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	sgID := d.Get("security_group_id").(string)
	rules, err := ecsSecurityGroupDetailRules(ctx, ecsClient, sgID)
	if err != nil {
		return diag.Errorf("failed to read ECS security_group_rule: %s", err)
	}
	rule := findSecurityGroupRuleByID(rules, d.Id())
	if rule == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("security_group_id", helper.StringFromMap(rule, "security_group_id"))
	_ = d.Set("protocol", helper.StringFromMap(rule, "protocol"))
	_ = d.Set("direction", helper.StringFromMap(rule, "direction"))
	_ = d.Set("ethertype", helper.StringFromMap(rule, "ethertype"))
	_ = d.Set("port_range_min", helper.IntFromMap(rule, "port_range_min"))
	_ = d.Set("port_range_max", helper.IntFromMap(rule, "port_range_max"))
	_ = d.Set("remote_ip_prefix", helper.StringFromMap(rule, "remote_ip_prefix"))
	_ = d.Set("remote_group_id", helper.StringFromMap(rule, "remote_group_id"))
	_ = d.Set("description", helper.StringFromMap(rule, "description"))
	_ = d.Set("tenant_id", helper.StringFromMap(rule, "tenant_id"))
	_ = d.Set("project_id", helper.StringFromMap(rule, "project_id"))
	_ = d.Set("revision_number", helper.IntFromMap(rule, "revision_number"))
	_ = d.Set("created_at", helper.StringFromMap(rule, "created_at"))
	_ = d.Set("updated_at", helper.StringFromMap(rule, "updated_at"))
	if err := d.Set("tags", helper.InterfaceToStringSlice(rule["tags"])); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceENECSSecurityGroupRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Defense in depth: CustomizeDiff blocks these at plan time; reject here if Update is still invoked.
	immutableFields := []string{
		"security_group_id",
		"protocol",
		"direction",
		"ethertype",
		"port_range_min",
		"port_range_max",
		"remote_ip_prefix",
		"remote_group_id",
		"description",
	}
	for _, field := range immutableFields {
		if d.HasChange(field) {
			return diag.Errorf("%s cannot be updated after creation", field)
		}
	}
	return resourceENECSSecurityGroupRuleRead(ctx, d, m)
}

func resourceENECSSecurityGroupRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	// API body: {"ids":["<rule_id>"]}; response data maps each id to status (e.g. "ok").
	req := map[string]interface{}{
		"ids": []string{d.Id()},
	}
	var resp map[string]interface{}
	if err := ecsClient.Post(ctx, "/ecs/openapi/v2/security_group_rule/delete", req, &resp); err != nil {
		return diag.Errorf("failed to delete ECS security_group_rule: %s", err)
	}
	payload, err := helper.ParseAPIResponsePayload(resp)
	if err != nil {
		return diag.Errorf("failed to parse ECS security_group_rule delete response: %s", err)
	}
	if m, ok := payload.(map[string]interface{}); ok {
		if status, ok := m[d.Id()].(string); !ok || status != "ok" {
			return diag.Errorf("ECS security_group_rule delete: unexpected status for id %q: %v", d.Id(), m[d.Id()])
		}
	}
	return nil
}
