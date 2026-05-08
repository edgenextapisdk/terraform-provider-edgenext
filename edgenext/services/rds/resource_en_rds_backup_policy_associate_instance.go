package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	rdsBackupPolicySourceTypeInstance = "instance"
	rdsBackupPolicySourceIDSeparator  = "/"
)

// ResourceENRDSBackupPolicyAssociateInstance manages one instance association on an RDS backup policy.
func ResourceENRDSBackupPolicyAssociateInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENRDSBackupPolicyAssociateInstanceCreate,
		ReadContext:   resourceENRDSBackupPolicyAssociateInstanceRead,
		DeleteContext: resourceENRDSBackupPolicyAssociateInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENRDSBackupPolicyAssociateInstanceImport,
		},
		Description: "Manages one instance association for an EdgeNext RDS backup policy.",
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Backup policy ID.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Instance ID.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance name returned from the API.",
			},
		},
	}
}

func rdsBackupPolicyAssociateInstanceComposeID(policyID, instanceID string) string {
	return strings.TrimSpace(policyID) + rdsBackupPolicySourceIDSeparator + strings.TrimSpace(instanceID)
}

func rdsBackupPolicyAssociateInstanceParseID(raw string) (policyID, instanceID string, err error) {
	s := strings.TrimSpace(raw)
	parts := strings.SplitN(s, rdsBackupPolicySourceIDSeparator, 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
		return "", "", fmt.Errorf("expected import id as policy_id%sinstance_id, got %q", rdsBackupPolicySourceIDSeparator, raw)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func resourceENRDSBackupPolicyAssociateInstanceImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	policyID, instanceID, err := rdsBackupPolicyAssociateInstanceParseID(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("policy_id", policyID); err != nil {
		return nil, err
	}
	if err := d.Set("instance_id", instanceID); err != nil {
		return nil, err
	}
	d.SetId(rdsBackupPolicyAssociateInstanceComposeID(policyID, instanceID))
	if diags := resourceENRDSBackupPolicyAssociateInstanceRead(ctx, d, m); diags.HasError() {
		return nil, diagToError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENRDSBackupPolicyAssociateInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Get("policy_id").(string))
	instanceID := strings.TrimSpace(d.Get("instance_id").(string))
	instanceName, diags := rdsResolveInstanceNameByID(ctx, rdsClient, instanceID)
	if diags.HasError() {
		return diags
	}

	sources := []map[string]interface{}{
		{
			"source_type": rdsBackupPolicySourceTypeInstance,
			"source_id":   instanceID,
			"source_name": instanceName,
		},
	}
	if diags := rdsBackupPolicySourcesAssociate(ctx, rdsClient, policyID, sources); diags.HasError() {
		return diags
	}

	d.SetId(rdsBackupPolicyAssociateInstanceComposeID(policyID, instanceID))
	return resourceENRDSBackupPolicyAssociateInstanceRead(ctx, d, m)
}

func resourceENRDSBackupPolicyAssociateInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Get("policy_id").(string))
	instanceID := strings.TrimSpace(d.Get("instance_id").(string))
	if id := strings.TrimSpace(d.Id()); id != "" && (policyID == "" || instanceID == "") {
		parsedPolicyID, parsedInstanceID, perr := rdsBackupPolicyAssociateInstanceParseID(id)
		if perr == nil {
			policyID = parsedPolicyID
			instanceID = parsedInstanceID
			_ = d.Set("policy_id", policyID)
			_ = d.Set("instance_id", instanceID)
		}
	}
	if policyID == "" || instanceID == "" {
		d.SetId("")
		return nil
	}

	sources, _, diags := rdsBackupPolicySourcesList(ctx, rdsClient, policyID)
	if diags.HasError() {
		return diags
	}
	matched, ok := rdsBackupPolicyFindSourceByID(sources, instanceID)
	if !ok {
		d.SetId("")
		return nil
	}

	_ = d.Set("policy_id", policyID)
	_ = d.Set("instance_id", helper.StringFromMap(matched, "source_id"))
	_ = d.Set("instance_name", helper.StringFromMap(matched, "source_name"))
	d.SetId(rdsBackupPolicyAssociateInstanceComposeID(policyID, instanceID))
	return nil
}

func resourceENRDSBackupPolicyAssociateInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := strings.TrimSpace(d.Get("policy_id").(string))
	instanceID := strings.TrimSpace(d.Get("instance_id").(string))

	removeSources := []map[string]interface{}{
		{
			"source_type": rdsBackupPolicySourceTypeInstance,
			"source_id":   instanceID,
		},
	}
	if diags := rdsBackupPolicySourcesRemove(ctx, rdsClient, policyID, removeSources); diags.HasError() {
		return diags
	}

	d.SetId("")
	return nil
}

func rdsBackupPolicyFindSourceByID(sources []map[string]interface{}, sourceID string) (map[string]interface{}, bool) {
	wantID := strings.TrimSpace(sourceID)
	for _, src := range sources {
		if strings.TrimSpace(helper.StringFromMap(src, "source_id")) == wantID {
			return src, true
		}
	}
	return nil, false
}

func rdsResolveInstanceNameByID(ctx context.Context, rdsClient *connectivity.RDSClient, instanceID string) (string, diag.Diagnostics) {
	req := map[string]interface{}{
		"id": instanceID,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/instances/get", req, &resp); err != nil {
		return "", diag.Errorf("failed to query RDS instance %q: %s", instanceID, err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return "", diag.Errorf("failed to parse RDS instance response: %s", err)
	}
	instance := helper.MapFromMap(payload, "instance")
	if instance == nil {
		return "", diag.Errorf("RDS instance %q not found", instanceID)
	}
	name := strings.TrimSpace(helper.StringFromMap(instance, "name"))
	if name == "" {
		return "", diag.Errorf("RDS instance %q returned empty name", instanceID)
	}
	return name, nil
}

func rdsBackupPolicySourcesList(ctx context.Context, rdsClient *connectivity.RDSClient, policyID string) ([]map[string]interface{}, int, diag.Diagnostics) {
	req := map[string]interface{}{
		"policy_id": policyID,
		"list_all":  true,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/sources/list", req, &resp); err != nil {
		return nil, 0, diag.Errorf("failed to list RDS backup policy sources: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, 0, diag.Errorf("failed to parse RDS backup policy sources list response: %s", err)
	}

	rawList := helper.ListFromMap(payload, "sources")
	out := make([]map[string]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if helper.StringFromMap(row, "source_type") != rdsBackupPolicySourceTypeInstance {
			continue
		}
		out = append(out, map[string]interface{}{
			"source_type": rdsBackupPolicySourceTypeInstance,
			"source_id":   helper.StringFromMap(row, "source_id"),
			"source_name": helper.StringFromMap(row, "source_name"),
		})
	}
	total := helper.IntFromMap(payload, "total")
	if total == 0 && len(out) > 0 {
		total = len(out)
	}
	return out, total, nil
}

func rdsBackupPolicySourcesAssociate(ctx context.Context, rdsClient *connectivity.RDSClient, policyID string, sources []map[string]interface{}) diag.Diagnostics {
	req := map[string]interface{}{
		"policy_id": policyID,
		"sources":   sources,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/sources/associate", req, &resp); err != nil {
		return diag.Errorf("failed to associate RDS backup policy sources: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS backup policy sources associate response: %s", err)
	}
	return nil
}

func rdsBackupPolicySourcesRemove(ctx context.Context, rdsClient *connectivity.RDSClient, policyID string, sources []map[string]interface{}) diag.Diagnostics {
	req := map[string]interface{}{
		"policy_id": policyID,
		"sources":   sources,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backup_policies/sources/remove", req, &resp); err != nil {
		return diag.Errorf("failed to remove RDS backup policy sources: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS backup policy sources remove response: %s", err)
	}
	return nil
}

func flattenRDSBackupPolicySources(sources []map[string]interface{}) []interface{} {
	out := make([]interface{}, 0, len(sources))
	for _, src := range sources {
		out = append(out, map[string]interface{}{
			"source_type": helper.StringFromMap(src, "source_type"),
			"source_id":   helper.StringFromMap(src, "source_id"),
			"source_name": helper.StringFromMap(src, "source_name"),
		})
	}
	return out
}
