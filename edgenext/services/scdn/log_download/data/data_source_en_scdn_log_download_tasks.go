package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// convertSearchTerms converts search_terms from map[string]string or map[string][]string to []map[string]string
func convertSearchTerms(searchTerms interface{}) []map[string]interface{} {
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

// DataSourceEdgenextScdnLogDownloadTasks returns the SCDN log download tasks data source
func DataSourceEdgenextScdnLogDownloadTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnLogDownloadTasksRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Task status: 0-not started, 1-in progress, 2-completed, 3-failed, 4-cancelled",
			},
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Task name",
			},
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File type: xls, csv, json",
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
				Description: "Total number of tasks",
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID",
						},
						"task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task name",
						},
						"member_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member ID",
						},
						"is_use_template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether to use template",
						},
						"template_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Template ID",
						},
						"data_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source",
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
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time",
						},
						"file_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "File type",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status",
						},
						"download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Download URL",
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

func dataSourceScdnLogDownloadTasksRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.LogDownloadTaskListRequest{
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}

	if v, ok := d.GetOk("status"); ok {
		req.Status = v.(int)
	}
	if v, ok := d.GetOk("task_name"); ok {
		req.TaskName = v.(string)
	}
	if v, ok := d.GetOk("file_type"); ok {
		req.FileType = v.(string)
	}
	if v, ok := d.GetOk("data_source"); ok {
		req.DataSource = v.(string)
	}

	log.Printf("[DEBUG] Reading SCDN log download tasks with request: %+v", req)
	response, err := service.ListLogDownloadTasks(req)
	if err != nil {
		return fmt.Errorf("failed to list log download tasks: %w", err)
	}

	// Set total - convert interface{} to string
	totalStr := convertTotalToString(response.Data.Total)
	d.Set("total", totalStr)

	// Set tasks
	tasks := make([]map[string]interface{}, len(response.Data.List))
	for i, task := range response.Data.List {
		taskMap := map[string]interface{}{
			"task_id":         task.TaskID,
			"task_name":       task.TaskName,
			"member_id":       task.MemberID,
			"is_use_template": task.IsUseTemplate,
			"template_id":     task.TemplateID,
			"data_source":     task.DataSource,
			"download_fields": task.DownloadFields,
			"start_time":      task.StartTime,
			"end_time":        task.EndTime,
			"file_type":       task.FileType,
			"status":          task.Status,
			"created_at":      task.CreatedAt,
			"updated_at":      task.UpdatedAt,
		}

		// Search terms - convert from map format to array format
		searchTerms := convertSearchTerms(task.SearchTerms)
		if searchTerms != nil {
			taskMap["search_terms"] = searchTerms
		}

		// Status - convert to string if needed
		if task.Status != nil {
			if statusStr, ok := task.Status.(string); ok {
				taskMap["status"] = statusStr
			} else if statusInt, ok := task.Status.(int); ok {
				taskMap["status"] = fmt.Sprintf("%d", statusInt)
			} else if statusFloat, ok := task.Status.(float64); ok {
				taskMap["status"] = fmt.Sprintf("%.0f", statusFloat)
			} else {
				taskMap["status"] = fmt.Sprintf("%v", task.Status)
			}
		}

		// Download URL
		if task.DownloadURL != nil {
			if url, ok := task.DownloadURL.(string); ok {
				taskMap["download_url"] = url
			}
		}

		tasks[i] = taskMap
	}
	d.Set("tasks", tasks)

	// Set ID
	d.SetId(fmt.Sprintf("log-download-tasks-%d-%d", req.Page, req.PerPage))

	// Save to file if specified
	if v, ok := d.GetOk("result_output_file"); ok {
		outputFile := v.(string)
		outputData := map[string]interface{}{
			"total": totalStr,
			"tasks": tasks,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
		_ = outputFile // Suppress unused variable warning
	}

	return nil
}
