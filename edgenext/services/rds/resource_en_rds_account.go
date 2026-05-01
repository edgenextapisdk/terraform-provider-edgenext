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

const rdsAccountIDSeparator = "/"

// --- databaseUsers/list helpers (also used by data_source_en_rds_accounts.go)

// rdsNormalizeAccountHost treats an empty host from the API as "%".
func rdsNormalizeAccountHost(v string) string {
	s := strings.TrimSpace(v)
	if s == "" {
		return "%"
	}
	return s
}

// rdsListDatabaseUsers returns data.users rows from POST /rds/openapi/v2/databaseUsers/list.
func rdsListDatabaseUsers(ctx context.Context, rdsClient *connectivity.RDSClient, instanceID string) ([]map[string]interface{}, diag.Diagnostics) {
	req := map[string]interface{}{
		"instance_id": instanceID,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databaseUsers/list", req, &resp); err != nil {
		return nil, diag.Errorf("failed to list RDS database users: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return nil, diag.Errorf("failed to parse RDS database users response: %s", err)
	}
	rawList := helper.ListFromMap(payload, "users")
	out := make([]map[string]interface{}, 0, len(rawList))
	for _, raw := range rawList {
		row, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		out = append(out, row)
	}
	return out, nil
}

// rdsFindDatabaseUserByNameHost finds the first list row matching name and host (after normalization).
func rdsFindDatabaseUserByNameHost(rows []map[string]interface{}, name, host string) (map[string]interface{}, bool) {
	wantHost := rdsNormalizeAccountHost(host)
	for _, row := range rows {
		if helper.StringFromMap(row, "name") != name {
			continue
		}
		if rdsNormalizeAccountHost(helper.StringFromMap(row, "host")) == wantHost {
			return row, true
		}
	}
	return nil, false
}

// ResourceENRDSAccount returns the resource schema for a database user on an RDS instance.
func ResourceENRDSAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENRDSAccountCreate,
		ReadContext:   resourceENRDSAccountRead,
		UpdateContext: resourceENRDSAccountUpdate,
		DeleteContext: resourceENRDSAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENRDSAccountImport,
		},
		Description: "Manages a database user on an EdgeNext RDS instance (create, update password or host, delete).",
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
				Description:  "Database user name (maps to API user_name).",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "%",
				Description: "Host pattern for the user (MySQL-style). Sent via update after create when not %.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "User password. Required on create; omit after import unless rotating (use lifecycle ignore_changes when not managing password).",
			},
		},
	}
}

func rdsAccountComposeID(instanceID, userName, host string) string {
	return strings.Join([]string{instanceID, userName, rdsNormalizeAccountHost(host)}, rdsAccountIDSeparator)
}

func rdsAccountParseImportID(raw string) (instanceID, userName, host string, err error) {
	s := strings.TrimSpace(raw)
	parts := strings.SplitN(s, rdsAccountIDSeparator, 3)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", "", fmt.Errorf("expected import id as instance_id%sname or instance_id%sname%shost, got %q",
			rdsAccountIDSeparator, rdsAccountIDSeparator, rdsAccountIDSeparator, raw)
	}
	instanceID, userName = parts[0], parts[1]
	host = "%"
	if len(parts) == 3 && strings.TrimSpace(parts[2]) != "" {
		host = parts[2]
	}
	return instanceID, userName, host, nil
}

func resourceENRDSAccountImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	instanceID, userName, host, err := rdsAccountParseImportID(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("instance_id", instanceID); err != nil {
		return nil, err
	}
	if err := d.Set("name", userName); err != nil {
		return nil, err
	}
	if err := d.Set("host", host); err != nil {
		return nil, err
	}
	d.SetId(rdsAccountComposeID(instanceID, userName, host))
	diags := resourceENRDSAccountRead(ctx, d, meta)
	if diags.HasError() {
		return nil, diagToError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENRDSAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	password := strings.TrimSpace(d.Get("password").(string))
	if password == "" {
		return diag.Errorf("password is required when creating an RDS database user")
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	wantHost := rdsNormalizeAccountHost(d.Get("host").(string))

	req := map[string]interface{}{
		"instance_id": instanceID,
		"users": []map[string]interface{}{
			{
				"name":     name,
				"password": password,
			},
		},
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databaseUsers/create", req, &resp); err != nil {
		return diag.Errorf("failed to create RDS database user: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS database user create response: %s", err)
	}

	if wantHost != "%" {
		if diags := rdsAccountPostUpdate(ctx, rdsClient, instanceID, name, map[string]interface{}{"host": wantHost}); diags.HasError() {
			return diags
		}
	}

	d.SetId(rdsAccountComposeID(instanceID, name, wantHost))
	return resourceENRDSAccountRead(ctx, d, m)
}

func rdsAccountPostUpdate(ctx context.Context, rdsClient *connectivity.RDSClient, instanceID, userName string, user map[string]interface{}) diag.Diagnostics {
	req := map[string]interface{}{
		"instance_id": instanceID,
		"user_name":   userName,
		"user":        user,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databaseUsers/update", req, &resp); err != nil {
		return diag.Errorf("failed to update RDS database user: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS database user update response: %s", err)
	}
	return nil
}

func resourceENRDSAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	host := rdsNormalizeAccountHost(d.Get("host").(string))

	if id := d.Id(); id != "" && (instanceID == "" || name == "") {
		i, n, h, perr := rdsAccountParseImportID(id)
		if perr == nil {
			instanceID, name, host = i, n, rdsNormalizeAccountHost(h)
			_ = d.Set("instance_id", instanceID)
			_ = d.Set("name", name)
			_ = d.Set("host", host)
		}
	}
	if instanceID == "" || name == "" {
		d.SetId("")
		return nil
	}

	rows, diags := rdsListDatabaseUsers(ctx, rdsClient, instanceID)
	if diags.HasError() {
		return diags
	}
	found, ok := rdsFindDatabaseUserByNameHost(rows, name, host)
	if !ok {
		d.SetId("")
		return nil
	}

	apiHost := strings.TrimSpace(helper.StringFromMap(found, "host"))
	if apiHost == "" {
		apiHost = "%"
	}
	_ = d.Set("host", apiHost)

	d.SetId(rdsAccountComposeID(instanceID, name, apiHost))
	return nil
}

func resourceENRDSAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	userPatch := map[string]interface{}{}
	if d.HasChange("password") {
		if pw := strings.TrimSpace(d.Get("password").(string)); pw != "" {
			userPatch["password"] = pw
		}
	}
	if d.HasChange("host") {
		userPatch["host"] = rdsNormalizeAccountHost(d.Get("host").(string))
	}
	if len(userPatch) == 0 {
		return resourceENRDSAccountRead(ctx, d, m)
	}

	instanceID := d.Get("instance_id").(string)
	userName := d.Get("name").(string)
	if diags := rdsAccountPostUpdate(ctx, rdsClient, instanceID, userName, userPatch); diags.HasError() {
		return diags
	}
	return resourceENRDSAccountRead(ctx, d, m)
}

func resourceENRDSAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"instance_id": d.Get("instance_id").(string),
		"user_name":   d.Get("name").(string),
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databaseUsers/delete", req, &resp); err != nil {
		return diag.Errorf("failed to delete RDS database user: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS database user delete response: %s", err)
	}
	d.SetId("")
	return nil
}
