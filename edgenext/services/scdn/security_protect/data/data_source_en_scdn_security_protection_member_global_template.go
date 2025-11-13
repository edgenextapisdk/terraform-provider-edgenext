package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionMemberGlobalTemplate returns the SCDN security protection member global template data source
func DataSourceEdgenextScdnSecurityProtectionMemberGlobalTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionMemberGlobalTemplateRead,

		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"template": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Global template information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Template ID",
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
					},
				},
			},
			"bind_domain_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Bind domain count",
			},
		},
	}
}

func dataSourceScdnSecurityProtectionMemberGlobalTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	log.Printf("[INFO] Querying SCDN security protection member global template")
	response, err := service.GetMemberGlobalTemplate()
	if err != nil {
		return fmt.Errorf("failed to get member global template: %w", err)
	}

	// Set resource ID
	d.SetId("member-global-template")

	// Set bind_domain_count
	if err := d.Set("bind_domain_count", response.Data.BindDomainCount); err != nil {
		return fmt.Errorf("error setting bind_domain_count: %w", err)
	}

	// Set template if exists
	if response.Data.Template != nil {
		templateList := []map[string]interface{}{
			{
				"id":         response.Data.Template.ID,
				"name":       response.Data.Template.Name,
				"type":       response.Data.Template.Type,
				"created_at": response.Data.Template.CreatedAt,
				"remark":     response.Data.Template.Remark,
			},
		}
		if err := d.Set("template", templateList); err != nil {
			return fmt.Errorf("error setting template: %w", err)
		}
	} else {
		// Set empty list if template is nil to avoid null in output
		if err := d.Set("template", []map[string]interface{}{}); err != nil {
			log.Printf("[WARN] Failed to set empty template list: %v", err)
		}
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"bind_domain_count": response.Data.BindDomainCount,
		}
		if response.Data.Template != nil {
			outputData["template"] = map[string]interface{}{
				"id":         response.Data.Template.ID,
				"name":       response.Data.Template.Name,
				"type":       response.Data.Template.Type,
				"created_at": response.Data.Template.CreatedAt,
				"remark":     response.Data.Template.Remark,
			}
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN security protection member global template queried successfully")
	return nil
}
