package template

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnRuleTemplates returns the SCDN rule templates list data source
func DataSourceEdgenextScdnRuleTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnRuleTemplatesRead,

		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for pagination, default: 1",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "Items per page, max: 1000, default: 1000",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by rule template name",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by associated domain",
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by application type",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of rule templates",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of rule templates",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule template ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule template name",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule template description",
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application type",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template creation timestamp",
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
				},
			},
		},
	}
}

func dataSourceScdnRuleTemplatesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.RuleTemplateListRequest{
		Page:     d.Get("page").(int),
		PageSize: d.Get("page_size").(int),
	}

	if name, ok := d.GetOk("name"); ok {
		req.Name = name.(string)
	}
	if domain, ok := d.GetOk("domain"); ok {
		req.Domain = domain.(string)
	}
	if appType, ok := d.GetOk("app_type"); ok {
		req.AppType = appType.(string)
	}

	log.Printf("[INFO] Listing SCDN rule templates")
	response, err := service.ListRuleTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to list rule templates: %w", err)
	}

	// Set the resource ID (use a hash of the request parameters)
	d.SetId(fmt.Sprintf("templates-%d-%d-%s-%s-%s", req.Page, req.PageSize, req.Name, req.Domain, req.AppType))

	// Set computed fields
	if err := d.Set("total", response.Data.Total); err != nil {
		log.Printf("[WARN] Failed to set total: %v", err)
	}

	// Convert list to schema format
	templateList := make([]map[string]interface{}, 0, len(response.Data.List))
	for _, tpl := range response.Data.List {
		bindDomains := make([]map[string]interface{}, 0, len(tpl.BindDomains))
		for _, bd := range tpl.BindDomains {
			bindDomains = append(bindDomains, map[string]interface{}{
				"domain_id":  bd.DomainID,
				"domain":     bd.Domain,
				"created_at": bd.CreatedAt,
			})
		}

		templateList = append(templateList, map[string]interface{}{
			"id":           tpl.ID,
			"name":         tpl.Name,
			"description":  tpl.Description,
			"app_type":     tpl.AppType,
			"created_at":   tpl.CreatedAt,
			"bind_domains": bindDomains,
		})
	}

	if err := d.Set("list", templateList); err != nil {
		log.Printf("[WARN] Failed to set list: %v", err)
	}

	// Save to file if requested
	if _, ok := d.GetOk("result_output_file"); ok {
		output := map[string]interface{}{
			"total": response.Data.Total,
			"list":  templateList,
		}
		if err := helper.WriteToFile(d, output); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
	}

	log.Printf("[INFO] SCDN rule templates listed successfully: %d templates", response.Data.Total)
	return nil
}
