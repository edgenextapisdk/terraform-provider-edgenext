package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionTemplate returns the SCDN security protection template data source
func DataSourceEdgenextScdnSecurityProtectionTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionTemplateRead,

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Business ID (template ID)",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template name",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template type: global, only_domain, more_domain",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time",
			},
			"remark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template remark",
			},
			"bind_domain_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Bind domain count",
			},
		},
	}
}

func dataSourceScdnSecurityProtectionTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	// Try to get global template first
	globalResp, err := service.GetMemberGlobalTemplate()
	if err == nil && globalResp.Data.Template != nil && globalResp.Data.Template.ID == businessID {
		template := globalResp.Data.Template
		d.SetId(fmt.Sprintf("template-%d", businessID))

		if err := d.Set("name", template.Name); err != nil {
			return fmt.Errorf("error setting name: %w", err)
		}
		if err := d.Set("type", template.Type); err != nil {
			return fmt.Errorf("error setting type: %w", err)
		}
		if err := d.Set("created_at", template.CreatedAt); err != nil {
			log.Printf("[WARN] Failed to set created_at: %v", err)
		}
		if err := d.Set("remark", template.Remark); err != nil {
			log.Printf("[WARN] Failed to set remark: %v", err)
		}
		if err := d.Set("bind_domain_count", globalResp.Data.BindDomainCount); err != nil {
			log.Printf("[WARN] Failed to set bind_domain_count: %v", err)
		}

		// Write result to output file if specified
		if outputFile := d.Get("result_output_file").(string); outputFile != "" {
			outputData := map[string]interface{}{
				"business_id":       businessID,
				"name":              template.Name,
				"type":              template.Type,
				"created_at":        template.CreatedAt,
				"remark":            template.Remark,
				"bind_domain_count": globalResp.Data.BindDomainCount,
			}
			if err := helper.WriteToFile(d, outputData); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
		}

		return nil
	}

	// Search for template by ID
	req := scdn.SecurityProtectionTemplateSearchRequest{
		TplType:  "global",
		Page:     1,
		PageSize: 100,
	}

	response, err := service.SearchSecurityProtectionTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to search templates: %w", err)
	}

	// Find the template by business_id
	var foundTemplate *scdn.SecurityProtectionTemplateInfo
	for _, template := range response.Data.Templates {
		if template.ID == businessID {
			foundTemplate = &template
			break
		}
	}

	if foundTemplate == nil {
		return fmt.Errorf("template not found: business_id=%d", businessID)
	}

	d.SetId(fmt.Sprintf("template-%d", businessID))

	if err := d.Set("name", foundTemplate.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("type", foundTemplate.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("created_at", foundTemplate.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set created_at: %v", err)
	}
	if err := d.Set("remark", foundTemplate.Remark); err != nil {
		log.Printf("[WARN] Failed to set remark: %v", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"business_id": businessID,
			"name":        foundTemplate.Name,
			"type":        foundTemplate.Type,
			"created_at":  foundTemplate.CreatedAt,
			"remark":      foundTemplate.Remark,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN security protection template queried successfully: business_id=%d", businessID)
	return nil
}
