package resource

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnLogDownloadTemplateStatus returns the SCDN log download template status resource
func ResourceEdgenextScdnLogDownloadTemplateStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnLogDownloadTemplateStatusCreate,
		Read:   resourceScdnLogDownloadTemplateStatusRead,
		Update: resourceScdnLogDownloadTemplateStatusUpdate,
		Delete: resourceScdnLogDownloadTemplateStatusDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Template ID",
			},
			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Status: 1-enabled, 0-disabled",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the template status (same as template_id)",
			},
		},
	}
}

func resourceScdnLogDownloadTemplateStatusCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnLogDownloadTemplateStatusUpdate(d, m)
}

func resourceScdnLogDownloadTemplateStatusRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	// Query template list to find the template
	req := scdn.LogDownloadTemplateListRequest{
		Page:    1,
		PerPage: 100,
	}

	response, err := service.ListLogDownloadTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to list log download templates: %w", err)
	}

	// Find the template by ID
	var template *scdn.LogDownloadTemplateInfo
	for i := range response.Data.List {
		if response.Data.List[i].TemplateID == templateID {
			template = &response.Data.List[i]
			break
		}
	}

	if template == nil {
		log.Printf("[WARN] Log download template %d not found, removing from state", templateID)
		d.SetId("")
		return nil
	}

	// Set fields
	d.Set("template_id", template.TemplateID)
	d.Set("status", template.Status)

	return nil
}

func resourceScdnLogDownloadTemplateStatusUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID := d.Get("template_id").(int)
	status := d.Get("status").(int)

	req := scdn.LogDownloadTemplateChangeStatusRequest{
		TemplateID: templateID,
		Status:     status,
	}

	log.Printf("[INFO] Changing SCDN log download template status: template_id=%d, status=%d", templateID, status)
	_, err := service.ChangeLogDownloadTemplateStatus(req)
	if err != nil {
		return fmt.Errorf("failed to change template status: %w", err)
	}

	// Set ID
	id := strconv.Itoa(templateID)
	d.SetId(id)

	log.Printf("[INFO] SCDN log download template status changed successfully: %d", templateID)
	return resourceScdnLogDownloadTemplateStatusRead(d, m)
}

func resourceScdnLogDownloadTemplateStatusDelete(d *schema.ResourceData, m interface{}) error {
	// Setting status to 0 (disabled) instead of deleting
	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.LogDownloadTemplateChangeStatusRequest{
		TemplateID: templateID,
		Status:     0, // Disable
	}

	log.Printf("[INFO] Disabling SCDN log download template: %d", templateID)
	_, err = service.ChangeLogDownloadTemplateStatus(req)
	if err != nil {
		return fmt.Errorf("failed to disable template: %w", err)
	}

	log.Printf("[INFO] SCDN log download template disabled successfully: %d", templateID)
	d.SetId("")
	return nil
}
