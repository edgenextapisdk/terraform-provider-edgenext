package template

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnRuleTemplate returns the SCDN rule template data source
func DataSourceEdgenextScdnRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnRuleTemplateRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The rule template ID",
			},
			"app_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application type (e.g., 'network_speed')",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule template name",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule template description",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template creation timestamp",
			},
			"bind_domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of domains bound to this template",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bound domain ID",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bound domain name",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain binding timestamp",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnRuleTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateIDStr := d.Get("id").(string)
	templateID, err := strconv.Atoi(templateIDStr)
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	appType := d.Get("app_type").(string)

	// List templates with filter to find this specific one
	req := scdn.RuleTemplateListRequest{
		Page:     1,
		PageSize: 1000,
		AppType:  appType,
	}

	log.Printf("[INFO] Reading SCDN rule template: %d", templateID)
	response, err := service.ListRuleTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to list rule templates: %w", err)
	}

	// Find the template with matching ID
	var template *scdn.RuleTemplateInfo
	for _, tpl := range response.Data.List {
		if tpl.ID == templateID {
			template = &tpl
			break
		}
	}

	if template == nil {
		return fmt.Errorf("rule template not found: %d", templateID)
	}

	// Set the resource ID
	d.SetId(templateIDStr)

	// Set computed fields
	if err := d.Set("name", template.Name); err != nil {
		log.Printf("[WARN] Failed to set template name: %v", err)
	}
	if err := d.Set("description", template.Description); err != nil {
		log.Printf("[WARN] Failed to set template description: %v", err)
	}
	if err := d.Set("created_at", template.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set template created_at: %v", err)
	}

	// Set bind_domains
	bindDomains := make([]map[string]interface{}, 0, len(template.BindDomains))
	for _, bd := range template.BindDomains {
		bindDomains = append(bindDomains, map[string]interface{}{
			"domain_id":  bd.DomainID,
			"domain":     bd.Domain,
			"created_at": bd.CreatedAt,
		})
	}
	if err := d.Set("bind_domains", bindDomains); err != nil {
		log.Printf("[WARN] Failed to set template bind_domains: %v", err)
	}

	// Save to file if requested
	if _, ok := d.GetOk("result_output_file"); ok {
		output := map[string]interface{}{
			"id":           template.ID,
			"name":         template.Name,
			"description":  template.Description,
			"app_type":     template.AppType,
			"created_at":   template.CreatedAt,
			"bind_domains": bindDomains,
		}
		if err := helper.WriteToFile(d, output); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
	}

	log.Printf("[INFO] SCDN rule template read successfully: %s", d.Id())
	return nil
}
