package rds

import (
	"context"
	"fmt"
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

func rdsAccountPrivilegeComposeID(instanceID, userName string) string {
	return instanceID + rdsAccountPrivilegeIDSeparator + userName
}

func rdsAccountPrivilegeParseImportID(raw string) (instanceID, userName string, err error) {
	s := strings.TrimSpace(raw)
	parts := strings.SplitN(s, rdsAccountPrivilegeIDSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected import id as instance_id%suser_name, got %q", rdsAccountPrivilegeIDSeparator, raw)
	}
	return parts[0], parts[1], nil
}

func resourceENRDSAccountPrivilegeImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	instanceID, userName, err := rdsAccountPrivilegeParseImportID(d.Id())
	if err != nil {
		return nil, err
	}
	if err := d.Set("instance_id", instanceID); err != nil {
		return nil, err
	}
	if err := d.Set("user_name", userName); err != nil {
		return nil, err
	}
	d.SetId(rdsAccountPrivilegeComposeID(instanceID, userName))
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
	if id := d.Id(); id != "" && (instanceID == "" || userName == "") {
		i, u, perr := rdsAccountPrivilegeParseImportID(id)
		if perr == nil {
			instanceID, userName = i, u
			_ = d.Set("instance_id", instanceID)
			_ = d.Set("user_name", userName)
		}
	}
	if instanceID == "" || userName == "" {
		d.SetId("")
		return nil
	}

	rows, diags := rdsListDatabaseUsers(ctx, rdsClient, instanceID)
	if diags.HasError() {
		return diags
	}

	row, ok := rdsFindDatabaseUserByName(rows, userName)
	if !ok {
		d.SetId("")
		return nil
	}
	if err := d.Set("databases", rdsUserDatabasesFromMap(row)); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(rdsAccountPrivilegeComposeID(instanceID, userName))
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
	if diags := rdsAccountPrivilegeUpdateAccess(ctx, rdsClient, instanceID, userName, []string{}); diags.HasError() {
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
	databases := rdsStringSetToSortedSlice(d.Get("databases").(*schema.Set))
	if diags := rdsAccountPrivilegeUpdateAccess(ctx, rdsClient, instanceID, userName, databases); diags.HasError() {
		return diags
	}

	d.SetId(rdsAccountPrivilegeComposeID(instanceID, userName))
	return resourceENRDSAccountPrivilegeRead(ctx, d, m)
}

func rdsAccountPrivilegeUpdateAccess(ctx context.Context, rdsClient *connectivity.RDSClient, instanceID, userName string, databases []string) diag.Diagnostics {
	reqDatabases := make([]interface{}, 0, len(databases))
	for _, db := range databases {
		reqDatabases = append(reqDatabases, db)
	}
	req := map[string]interface{}{
		"instance_id": instanceID,
		"user_name":   userName,
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

func rdsFindDatabaseUserByName(rows []map[string]interface{}, userName string) (map[string]interface{}, bool) {
	for _, row := range rows {
		if helper.StringFromMap(row, "name") == userName {
			return row, true
		}
	}
	return nil, false
}
