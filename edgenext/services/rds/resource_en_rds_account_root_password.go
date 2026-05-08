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

// ResourceENRDSAccountRootPassword manages root password enable/rotation on an RDS instance.
func ResourceENRDSAccountRootPassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENRDSAccountRootPasswordCreate,
		ReadContext:   resourceENRDSAccountRootPasswordRead,
		UpdateContext: resourceENRDSAccountRootPasswordUpdate,
		DeleteContext: resourceENRDSAccountRootPasswordDelete,
		Description:   "Manages RDS root password via enableRoot API. Delete only removes Terraform state.",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "RDS instance ID.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				ValidateDiagFunc: validateRDSAccountPassword,
				Description:      "Root password. Please enter 8-20 characters, must include all four: uppercase letters, lowercase letters, numbers, and special characters from ()~!@#$%^&*_-+=|{}[]:;'<>,.?/. ",
			},
		},
	}
}

func resourceENRDSAccountRootPasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	instanceID := strings.TrimSpace(d.Get("instance_id").(string))
	d.SetId(instanceID + "/root")
	return resourceENRDSAccountRootPasswordApply(ctx, d, m)
}

func resourceENRDSAccountRootPasswordRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No query endpoint is available for this action resource.
	return nil
}

func resourceENRDSAccountRootPasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if !d.HasChange("password") {
		return resourceENRDSAccountRootPasswordRead(ctx, d, m)
	}
	return resourceENRDSAccountRootPasswordApply(ctx, d, m)
}

func resourceENRDSAccountRootPasswordDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No delete endpoint; deleting this resource only removes Terraform state.
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "State-only delete for RDS root password resource",
			Detail:   fmt.Sprintf("Deleting edgenext_rds_account_root_password for instance %q only removes Terraform state and does not disable root access.", d.Get("instance_id").(string)),
		},
	}
}

func resourceENRDSAccountRootPasswordApply(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"instance_id": d.Get("instance_id").(string),
		"password":    d.Get("password").(string),
	}

	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/databaseUsers/enableRoot", req, &resp); err != nil {
		return diag.Errorf("failed to set RDS root password: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS root password response: %s", err)
	}
	return resourceENRDSAccountRootPasswordRead(ctx, d, m)
}
