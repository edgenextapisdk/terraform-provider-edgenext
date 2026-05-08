package rds

import (
	"context"
	"sort"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const rdsAccountPrivilegeIDSeparator = "/"

// ResourceENRDSAccountPrivilege returns the resource schema for database access privileges of an RDS user.
func ResourceENRDSAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENRDSAccountPrivilegeCreate,
		ReadContext:   resourceENRDSAccountPrivilegeRead,
		UpdateContext: resourceENRDSAccountPrivilegeUpdate,
		DeleteContext: resourceENRDSAccountPrivilegeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENRDSAccountPrivilegeImport,
		},
		Description: "Manages database access privileges for an EdgeNext RDS database user via updateAccess API.",
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
				Description:  "Database user name to manage privileges for.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"host": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateRDSAccountHost,
				Description:      "Client host for the database user (same as edgenext_rds_account.host). Together with instance_id and user_name it identifies the user. Use % for any host.",
			},
			"databases": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Database names granted to the user.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},
		},
	}
}

func rdsAccountPrivilegeComposeID(instanceID, userName, host string) string {
	return strings.Join([]string{instanceID, userName, rdsNormalizeAccountHost(host)}, rdsAccountPrivilegeIDSeparator)
}

func resourceENRDSAccountPrivilegeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
	d.SetId(rdsAccountPrivilegeComposeID(instanceID, userName, host))
	diags := resourceENRDSAccountPrivilegeRead(ctx, d, meta)
	if diags.HasError() {
		return nil, diagToError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENRDSAccountPrivilegeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return rdsAccountPrivilegeApply(ctx, d, m)
}

func resourceENRDSAccountPrivilegeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return rdsAccountPrivilegeApply(ctx, d, m)
}

func resourceENRDSAccountPrivilegeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	userName := d.Get("user_name").(string)
	host := rdsNormalizeAccountHost(d.Get("host").(string))
	if id := d.Id(); id != "" && (instanceID == "" || userName == "" || host == "") {
		i, u, h, perr := rdsParseAccountUserHostImportID(id)
		if perr == nil {
			instanceID, userName, host = i, u, rdsNormalizeAccountHost(h)
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

	row, ok := rdsFindDatabaseUserByNameHost(rows, userName, host)
	if !ok {
		d.SetId("")
		return nil
	}
	if err := d.Set("databases", rdsUserDatabasesFromMap(row)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(rdsAccountPrivilegeComposeID(instanceID, userName, host))
	return nil
}

func resourceENRDSAccountPrivilegeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	userName := d.Get("user_name").(string)
	host := rdsNormalizeAccountHost(d.Get("host").(string))
	if diags := rdsAccountPrivilegeUpdateAccess(ctx, rdsClient, instanceID, userName, host, []string{}); diags.HasError() {
		return diags
	}

	d.SetId("")
	return nil
}

func rdsAccountPrivilegeApply(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	instanceID := d.Get("instance_id").(string)
	userName := d.Get("user_name").(string)
	host := rdsNormalizeAccountHost(d.Get("host").(string))
	databases := rdsStringSetToSortedSlice(d.Get("databases").(*schema.Set))
	if diags := rdsAccountPrivilegeUpdateAccess(ctx, rdsClient, instanceID, userName, host, databases); diags.HasError() {
		return diags
	}

	d.SetId(rdsAccountPrivilegeComposeID(instanceID, userName, host))
	return resourceENRDSAccountPrivilegeRead(ctx, d, m)
}

func rdsAccountPrivilegeUpdateAccess(ctx context.Context, rdsClient *connectivity.RDSClient, instanceID, userName, host string, databases []string) diag.Diagnostics {
	reqDatabases := make([]interface{}, 0, len(databases))
	for _, db := range databases {
		reqDatabases = append(reqDatabases, db)
	}
	req := map[string]interface{}{
		"instance_id": instanceID,
		"user_name":   userName,
		"host":        rdsNormalizeAccountHost(host),
		"databases":   reqDatabases,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databaseUsers/updateAccess", req, &resp); err != nil {
		return diag.Errorf("failed to update RDS database user privileges: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS database user privileges update response: %s", err)
	}
	return nil
}

func rdsStringSetToSortedSlice(set *schema.Set) []string {
	raw := set.List()
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		s, ok := v.(string)
		if !ok {
			continue
		}
		ss := strings.TrimSpace(s)
		if ss == "" {
			continue
		}
		out = append(out, ss)
	}
	sort.Strings(out)
	return out
}
