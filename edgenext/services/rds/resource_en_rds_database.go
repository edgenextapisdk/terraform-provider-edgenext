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

const rdsDatabaseIDSeparator = "/"

// ResourceENRDSDatabase returns the resource schema for a single database on an RDS instance.
func ResourceENRDSDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENRDSDatabaseCreate,
		ReadContext:   resourceENRDSDatabaseRead,
		DeleteContext: resourceENRDSDatabaseDelete,
		CustomizeDiff: resourceENRDSDatabaseCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENRDSDatabaseImport,
		},
		Description: "Manages a database on an EdgeNext RDS instance (create and delete only; there is no update API).",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "RDS instance ID.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Database name.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"character_set": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          "utf8",
				ValidateDiagFunc: validateRDSMySQLCharacterSet,
				Description: "MySQL character set for the new database (create only). " +
					"Examples: utf8, utf8mb4, latin1. Default: utf8. " +
					"Matching is case-insensitive; the create request sends the canonical lowercase name.",
			},
			"collate": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "utf8_general_ci",
				Description: "MySQL collation for the selected character_set (create only). " +
					"Examples: utf8_general_ci (for utf8), utf8mb4_general_ci (for utf8mb4). " +
					"Default: utf8_general_ci. Matching is case-insensitive.",
			},
		},
	}
}

func rdsDatabaseComposeID(instanceID, dbName string) string {
	return instanceID + rdsDatabaseIDSeparator + dbName
}

func resourceENRDSDatabaseImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	raw := strings.TrimSpace(d.Id())
	parts := strings.SplitN(raw, rdsDatabaseIDSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("expected import id as instance_id%sdatabase_name, got %q", rdsDatabaseIDSeparator, raw)
	}
	if err := d.Set("instance_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("name", parts[1]); err != nil {
		return nil, err
	}
	d.SetId(rdsDatabaseComposeID(parts[0], parts[1]))
	diags := resourceENRDSDatabaseRead(ctx, d, meta)
	if diags.HasError() {
		return nil, diagToError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENRDSDatabaseCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	charset := "utf8"
	if raw, ok := d.Get("character_set").(string); ok {
		charset = strings.ToLower(strings.TrimSpace(raw))
	}
	if charset == "" {
		charset = "utf8"
	}
	collate := "utf8_general_ci"
	if raw, ok := d.Get("collate").(string); ok {
		collate = strings.TrimSpace(raw)
	}
	if collate == "" {
		collate = "utf8_general_ci"
	}
	if !rdsMySQLCollateValid(charset, collate) {
		allowed := rdsMySQLCharsetCollations[charset]
		hint := ""
		if len(allowed) > 0 {
			n := len(allowed)
			if n > 8 {
				n = 8
			}
			hint = fmt.Sprintf(" Example collations for %q: %s.", charset, strings.Join(allowed[:n], ", "))
			if len(allowed) > 8 {
				hint += " (list truncated)"
			}
		}
		return fmt.Errorf("collate %q is not valid for character_set %q.%s", collate, charset, hint)
	}
	return nil
}

func diagToError(diags diag.Diagnostics) error {
	if len(diags) == 0 {
		return fmt.Errorf("unknown error")
	}
	d := diags[0]
	if d.Detail != "" {
		return fmt.Errorf("%s: %s", d.Summary, d.Detail)
	}
	return fmt.Errorf("%s", d.Summary)
}

func resourceENRDSDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	characterSet := strings.ToLower(strings.TrimSpace(d.Get("character_set").(string)))
	if characterSet == "" {
		characterSet = "utf8"
	}
	collate := strings.TrimSpace(d.Get("collate").(string))
	if collate == "" {
		collate = "utf8_general_ci"
	}
	entry := map[string]interface{}{
		"name":          name,
		"character_set": characterSet,
		"collate":       collate,
	}
	req := map[string]interface{}{
		"instance_id": instanceID,
		"databases":   []map[string]interface{}{entry},
	}

	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databases/create", req, &resp); err != nil {
		return diag.Errorf("failed to create RDS database: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS database create response: %s", err)
	}

	d.SetId(rdsDatabaseComposeID(instanceID, name))
	return resourceENRDSDatabaseRead(ctx, d, m)
}

func resourceENRDSDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	if id := d.Id(); id != "" && (instanceID == "" || name == "") {
		parts := strings.SplitN(id, rdsDatabaseIDSeparator, 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			instanceID, name = parts[0], parts[1]
			_ = d.Set("instance_id", instanceID)
			_ = d.Set("name", name)
		}
	}
	if instanceID == "" || name == "" {
		d.SetId("")
		return nil
	}

	found, diags := rdsDatabaseExistsInList(ctx, rdsClient, instanceID, name)
	if diags.HasError() {
		return diags
	}
	if !found {
		d.SetId("")
		return nil
	}

	if v := d.Get("character_set").(string); v == "" {
		_ = d.Set("character_set", "utf8")
	}
	if v := d.Get("collate").(string); v == "" {
		_ = d.Set("collate", "utf8_general_ci")
	}

	d.SetId(rdsDatabaseComposeID(instanceID, name))
	return nil
}

func rdsDatabaseExistsInList(ctx context.Context, rdsClient *connectivity.RDSClient, instanceID, wantName string) (bool, diag.Diagnostics) {
	req := map[string]interface{}{
		"instance_id": instanceID,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databases/list", req, &resp); err != nil {
		return false, diag.Errorf("failed to list RDS databases: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return false, diag.Errorf("failed to parse RDS databases response: %s", err)
	}
	for _, raw := range helper.ListFromMap(payload, "databases") {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if strings.EqualFold(helper.StringFromMap(row, "name"), wantName) {
			return true, nil
		}
	}
	return false, nil
}

func resourceENRDSDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"instance_id":   d.Get("instance_id").(string),
		"database_name": d.Get("name").(string),
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databases/delete", req, &resp); err != nil {
		return diag.Errorf("failed to delete RDS database: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS database delete response: %s", err)
	}
	d.SetId("")
	return nil
}
