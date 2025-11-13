package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionTemplates returns the SCDN security protection templates data source
func DataSourceEdgenextScdnSecurityProtectionTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionTemplatesRead,

		Schema: map[string]*schema.Schema{
			"tpl_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template type: global, only_domain, more_domain",
			},
			"search_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search type",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search keyword",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "Page size",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of templates",
			},
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Template list",
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
							Description: "Template type",
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
		},
	}
}

func dataSourceScdnSecurityProtectionTemplatesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.SecurityProtectionTemplateSearchRequest{
		TplType:  d.Get("tpl_type").(string),
		Page:     d.Get("page").(int),
		PageSize: d.Get("page_size").(int),
	}

	if searchType, ok := d.GetOk("search_type"); ok {
		req.SearchType = searchType.(string)
	}
	if searchKey, ok := d.GetOk("search_key"); ok {
		req.SearchKey = searchKey.(string)
	}

	log.Printf("[INFO] Querying SCDN security protection templates: tpl_type=%s", req.TplType)
	response, err := service.SearchSecurityProtectionTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to query templates: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("templates-%s-%d-%d", req.TplType, req.Page, req.PageSize))

	// Set total
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Convert templates to schema format
	templateList := make([]map[string]interface{}, 0, len(response.Data.Templates))
	for _, template := range response.Data.Templates {
		templateMap := map[string]interface{}{
			"id":         template.ID,
			"name":       template.Name,
			"type":       template.Type,
			"created_at": template.CreatedAt,
			"remark":     template.Remark,
		}
		templateList = append(templateList, templateMap)
	}

	if err := d.Set("templates", templateList); err != nil {
		return fmt.Errorf("error setting templates: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":     response.Data.Total,
			"templates": templateList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN security protection templates queried successfully: total=%d", response.Data.Total)
	return nil
}
