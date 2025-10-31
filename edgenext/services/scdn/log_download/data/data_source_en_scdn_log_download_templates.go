package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// convertSearchTermsForTemplate converts search_terms from map[string]string or map[string][]string to []map[string]string
func convertSearchTermsForTemplate(searchTerms interface{}) []map[string]interface{} {
	if searchTerms == nil {
		return nil
	}

	// Try to convert from map format (API response)
	if termMap, ok := searchTerms.(map[string]interface{}); ok {
		result := make([]map[string]interface{}, 0)
		for key, value := range termMap {
			// First try map[string]string format (API document format)
			if strValue, ok := value.(string); ok {
				result = append(result, map[string]interface{}{
					"key":   key,
					"value": strValue,
				})
				continue
			}

			// Then try map[string][]string format (array format, for backward compatibility)
			var values []string
			if strSlice, ok := value.([]string); ok {
				values = strSlice
			} else if ifaceSlice, ok := value.([]interface{}); ok {
				values = make([]string, len(ifaceSlice))
				for i, v := range ifaceSlice {
					if str, ok := v.(string); ok {
						values[i] = str
					}
				}
			}
			// Create one entry per value in the array
			for _, val := range values {
				result = append(result, map[string]interface{}{
					"key":   key,
					"value": val,
				})
			}
		}
		return result
	}

	// Try to convert from []LogDownloadSearchTerm format (if already converted)
	if termSlice, ok := searchTerms.([]scdn.LogDownloadSearchTerm); ok {
		result := make([]map[string]interface{}, len(termSlice))
		for i, term := range termSlice {
			result[i] = map[string]interface{}{
				"key":   term.Key,
				"value": term.Value,
			}
		}
		return result
	}

	return nil
}

// DataSourceEdgenextScdnLogDownloadTemplates returns the SCDN log download templates data source
func DataSourceEdgenextScdnLogDownloadTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnLogDownloadTemplatesRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Status: 1-enabled, 0-disabled",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Group ID",
			},
			"template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template name",
			},
			"data_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data source: ng, cc, waf",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number",
			},
			"per_page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "Items per page",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"total": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Total number of templates",
			},
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Template list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Template ID",
						},
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template name",
						},
						"member_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member ID",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group ID",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group name",
						},
						"data_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status",
						},
						"download_fields": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Download fields",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"search_terms": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Search conditions",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Search key",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Search value",
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update timestamp",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnLogDownloadTemplatesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.LogDownloadTemplateListRequest{
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}

	if v, ok := d.GetOk("status"); ok {
		req.Status = v.(int)
	}
	if v, ok := d.GetOk("group_id"); ok {
		req.GroupID = v.(int)
	}
	if v, ok := d.GetOk("template_name"); ok {
		req.TemplateName = v.(string)
	}
	if v, ok := d.GetOk("data_source"); ok {
		req.DataSource = v.(string)
	}

	log.Printf("[DEBUG] Reading SCDN log download templates with request: %+v", req)
	response, err := service.ListLogDownloadTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to list log download templates: %w", err)
	}

	// Set total - convert interface{} to string
	totalStr := convertTotalToString(response.Data.Total)
	d.Set("total", totalStr)

	// Set templates
	templates := make([]map[string]interface{}, len(response.Data.List))
	for i, template := range response.Data.List {
		templateMap := map[string]interface{}{
			"template_id":     template.TemplateID,
			"template_name":   template.TemplateName,
			"member_id":       template.MemberID,
			"group_id":        template.GroupID,
			"group_name":      template.GroupName,
			"data_source":     template.DataSource,
			"status":          template.Status,
			"download_fields": template.DownloadFields,
			"created_at":      template.CreatedAt,
			"updated_at":      template.UpdatedAt,
		}

		// Search terms - convert from map format to array format
		searchTerms := convertSearchTermsForTemplate(template.SearchTerms)
		if searchTerms != nil {
			templateMap["search_terms"] = searchTerms
		}

		templates[i] = templateMap
	}
	d.Set("templates", templates)

	// Set ID
	d.SetId(fmt.Sprintf("log-download-templates-%d-%d", req.Page, req.PerPage))

	// Save to file if specified
	if v, ok := d.GetOk("result_output_file"); ok {
		outputFile := v.(string)
		outputData := map[string]interface{}{
			"total":     totalStr,
			"templates": templates,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
		_ = outputFile // Suppress unused variable warning
	}

	return nil
}
