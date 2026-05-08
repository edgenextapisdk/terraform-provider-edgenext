package rds

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const rdsAccountIDSeparator = "/"

var (
	rdsPasswordUpperRe   = regexp.MustCompile(`[A-Z]`)
	rdsPasswordLowerRe   = regexp.MustCompile(`[a-z]`)
	rdsPasswordDigitRe   = regexp.MustCompile(`[0-9]`)
	rdsPasswordSpecialRe = regexp.MustCompile(`[()~!@#$%^&*_\-+=|{}\[\]:;'<>,.?/]`)
	rdsPasswordAllowedRe = regexp.MustCompile(`^[A-Za-z0-9()~!@#$%^&*_\-+=|{}\[\]:;'<>,.?/]+$`)
)

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
			"user_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "User name.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"host": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateRDSAccountHost,
				Description:      "Client host for the database user: use % for any host, or a literal IPv4 address. Sent on update and delete.",
			},
			"password": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				ValidateDiagFunc: validateRDSAccountPassword,
				Description: "User password. Required on create; omit after import unless rotating. " +
					"Please enter 8-20 characters, must include all four: uppercase letters, lowercase letters, numbers, and special characters from ()~!@#$%^&*_-+=|{}[]:;'<>,.?/. ",
			},
		},
	}
}

func rdsAccountComposeID(instanceID, userName, host string) string {
	return strings.Join([]string{instanceID, userName, rdsNormalizeAccountHost(host)}, rdsAccountIDSeparator)
}

func rdsParseAccountUserHostImportID(raw string) (instanceID, userName, host string, err error) {
	s := strings.TrimSpace(raw)
	parts := strings.SplitN(s, rdsAccountIDSeparator, 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || strings.TrimSpace(parts[2]) == "" {
		return "", "", "", fmt.Errorf("expected import id as instance_id%suser_name%shost, got %q",
			rdsAccountIDSeparator, rdsAccountIDSeparator, raw)
	}
	instanceID, userName = parts[0], parts[1]
	host = strings.TrimSpace(parts[2])
	return instanceID, userName, host, nil
}

func resourceENRDSAccountImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	instanceID, userName, host, err := rdsParseAccountUserHostImportID(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("instance_id", instanceID); err != nil {
		return nil, err
	}
	if err := d.Set("user_name", userName); err != nil {
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
	userName := d.Get("user_name").(string)
	wantHost := rdsNormalizeAccountHost(d.Get("host").(string))

	req := map[string]interface{}{
		"instance_id": instanceID,
		"users": []map[string]interface{}{
			{
				"name":     userName,
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
		if diags := rdsAccountPostUpdate(ctx, rdsClient, instanceID, userName, "%", map[string]interface{}{"host": wantHost}); diags.HasError() {
			return diags
		}
	}

	d.SetId(rdsAccountComposeID(instanceID, userName, wantHost))
	return resourceENRDSAccountRead(ctx, d, m)
}

func rdsAccountPostUpdate(ctx context.Context, rdsClient *connectivity.RDSClient, instanceID, userName, host string, user map[string]interface{}) diag.Diagnostics {
	h := rdsNormalizeAccountHost(host)
	req := map[string]interface{}{
		"instance_id": instanceID,
		"user_name":   userName,
		"host":        h,
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
	userName := d.Get("user_name").(string)
	host := rdsNormalizeAccountHost(d.Get("host").(string))

	if id := d.Id(); id != "" && (instanceID == "" || userName == "" || host == "") {
		i, n, h, perr := rdsParseAccountUserHostImportID(id)
		if perr == nil {
			instanceID, userName, host = i, n, rdsNormalizeAccountHost(h)
			_ = d.Set("instance_id", instanceID)
			_ = d.Set("user_name", userName)
			_ = d.Set("host", host)
		}
	}
	if instanceID == "" || userName == "" || host == "" {
		d.SetId("")
		return nil
	}

	rows, diags := rdsListDatabaseUsers(ctx, rdsClient, instanceID)
	if diags.HasError() {
		return diags
	}
	_, ok := rdsFindDatabaseUserByNameHost(rows, userName, host)
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(rdsAccountComposeID(instanceID, userName, host))
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
	userName := d.Get("user_name").(string)
	currentHost := rdsNormalizeAccountHost(d.Get("host").(string))
	if d.HasChange("host") {
		oldRaw, _ := d.GetChange("host")
		if oldHost, ok := oldRaw.(string); ok {
			currentHost = rdsNormalizeAccountHost(oldHost)
		}
	}
	if diags := rdsAccountPostUpdate(ctx, rdsClient, instanceID, userName, currentHost, userPatch); diags.HasError() {
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
		"user_name":   d.Get("user_name").(string),
		"host":        rdsNormalizeAccountHost(d.Get("host").(string)),
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

func validateRDSAccountPassword(v interface{}, _ cty.Path) diag.Diagnostics {
	if v == nil {
		return nil
	}
	password, ok := v.(string)
	if !ok {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid password",
			Detail:   fmt.Sprintf("Expected string, got %T", v),
		}}
	}
	// Keep password optional at schema level for import/non-rotation workflows.
	if password == "" {
		return nil
	}
	if len(password) < 8 || len(password) > 20 {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid password",
			Detail:   "Please enter 8-20 characters.",
		}}
	}
	if !rdsPasswordAllowedRe.MatchString(password) {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid password",
			Detail:   "Password contains unsupported characters. Allowed special characters are: ()~!@#$%^&*_-+=|{}[]:;'<>,.?/",
		}}
	}
	if !rdsPasswordUpperRe.MatchString(password) ||
		!rdsPasswordLowerRe.MatchString(password) ||
		!rdsPasswordDigitRe.MatchString(password) ||
		!rdsPasswordSpecialRe.MatchString(password) {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid password",
			Detail:   "Password must include all four: uppercase letters, lowercase letters, numbers, and special characters from ()~!@#$%^&*_-+=|{}[]:;'<>,.?/.",
		}}
	}
	return nil
}

func validateRDSAccountHost(v interface{}, _ cty.Path) diag.Diagnostics {
	if v == nil {
		return nil
	}
	raw, ok := v.(string)
	if !ok {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid host",
			Detail:   fmt.Sprintf("Expected string, got %T", v),
		}}
	}
	h := strings.TrimSpace(raw)
	if h == "%" {
		return nil
	}
	ip := net.ParseIP(h)
	if ip != nil && ip.To4() != nil {
		return nil
	}
	return diag.Diagnostics{{
		Severity: diag.Error,
		Summary:  "Invalid host",
		Detail:   "host must be % (any host) or a valid IPv4 address (dotted decimal).",
	}}
}

// --- databaseUsers/list helpers (also used by data_source_en_rds_accounts.go)

// rdsNormalizeAccountHost trims surrounding whitespace.
func rdsNormalizeAccountHost(v string) string {
	return strings.TrimSpace(v)
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
