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

// ResourceENRDSBackup returns the resource schema for a single RDS backup.
func ResourceENRDSBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENRDSBackupCreate,
		ReadContext:   resourceENRDSBackupRead,
		DeleteContext: resourceENRDSBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENRDSBackupImport,
		},
		Description: "Manages an EdgeNext RDS backup (create and delete only; get API is used for read).",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "RDS instance ID to back up.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Backup name passed to the create API.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backup status from the API.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RDS instance name from the get API.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backup type (for example full).",
			},
			"size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Backup size in GB from the get API.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time (RFC3339).",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time (RFC3339).",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region from the get API.",
			},
			"datastore": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Datastore engine and version.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine type.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine version.",
						},
					},
				},
			},
		},
	}
}

func resourceENRDSBackupImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(d.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import id as backup_id, got empty string")
	}
	d.SetId(id)
	diags := resourceENRDSBackupRead(ctx, d, meta)
	if diags.HasError() {
		return nil, diagToError(diags)
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("backup %q not found", id)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceENRDSBackupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	req := map[string]interface{}{
		"instance_id": d.Get("instance_id").(string),
		"type":        rdsBackupListTypeFull,
		"name":        d.Get("name").(string),
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backups/create", req, &resp); err != nil {
		return diag.Errorf("failed to create RDS backup: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse RDS backup create response: %s", err)
	}

	backupID := helper.StringFromMap(payload, "id")
	if backupID == "" {
		return diag.Errorf("RDS backup create response missing id")
	}
	d.SetId(backupID)
	if v := helper.StringFromMap(payload, "status"); v != "" {
		_ = d.Set("status", v)
	}
	if v := helper.StringFromMap(payload, "name"); v != "" {
		_ = d.Set("name", v)
	}
	return resourceENRDSBackupRead(ctx, d, m)
}

func resourceENRDSBackupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	backupID := strings.TrimSpace(d.Id())
	if backupID == "" {
		return nil
	}

	req := map[string]interface{}{
		"id": backupID,
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backups/get", req, &resp); err != nil {
		return diag.Errorf("failed to get RDS backup: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse RDS backup get response: %s", err)
	}

	row := helper.MapFromMap(payload, "backup")
	if row == nil {
		d.SetId("")
		return nil
	}

	attrs := rdsBackupAttrsFromMap(row)
	for k, v := range attrs {
		if k == "id" {
			continue
		}
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(backupID)
	return nil
}

func resourceENRDSBackupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	rdsClient, err := client.RDSClient()
	if err != nil {
		return diag.FromErr(err)
	}

	backupID := strings.TrimSpace(d.Id())
	if backupID == "" {
		return nil
	}

	req := map[string]interface{}{
		"ids": []string{backupID},
	}
	var resp map[string]interface{}
	if err := rdsClient.Post(ctx, "/rds/openapi/v2/backups/delete", req, &resp); err != nil {
		return diag.Errorf("failed to delete RDS backup: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		return diag.Errorf("failed to parse RDS backup delete response: %s", err)
	}
	d.SetId("")
	return nil
}
