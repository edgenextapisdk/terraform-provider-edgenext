package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCacheCleanTasks returns the SCDN cache clean tasks data source
func DataSourceEdgenextScdnCacheCleanTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCacheCleanTasksRead,

		Schema: map[string]*schema.Schema{
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
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start time, format: YYYY-MM-DD HH:II:SS",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "End time, format: YYYY-MM-DD HH:II:SS",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status: 1-executing, 2-completed",
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
						"user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User ID",
						},
						"sub_user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sub user ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status (can be null): Failed, Finished, etc.",
						},
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID",
						},
						"sub_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task type: Directory, SubDomain, URL",
						},
						"total": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Total number of nodes",
						},
						"succeed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of successful nodes",
						},
						"failed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of failed nodes",
						},
						"ongoing": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Number of executing nodes",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time, ISO 8601 format",
						},
						"operator_user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operator user name",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnCacheCleanTasksRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.CacheCleanTaskListRequest{
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}

	if v, ok := d.GetOk("start_time"); ok {
		req.StartTime = v.(string)
	}
	if v, ok := d.GetOk("end_time"); ok {
		req.EndTime = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		req.Status = v.(string)
	}

	log.Printf("[INFO] Querying SCDN cache clean tasks")
	response, err := service.GetCacheCleanTaskList(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN cache clean tasks: %w", err)
	}

	// Set total
	if err := d.Set("total", response.Data.Total.String()); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Set tasks
	tasks := make([]map[string]interface{}, len(response.Data.List))
	for i, task := range response.Data.List {
		taskMap := map[string]interface{}{
			"user_id":            task.UserID,
			"sub_user_id":        task.SubUserID,
			"task_id":            task.TaskID,
			"sub_type":           task.SubType,
			"total":              task.Total.String(),
			"succeed":            task.Succeed.String(),
			"failed":             task.Failed.String(),
			"ongoing":            task.Ongoing.String(),
			"created_at":         task.CreatedAt,
			"operator_user_name": task.OperatorUserName,
		}
		// Handle nullable status
		if task.Status != nil {
			taskMap["status"] = *task.Status
		} else {
			taskMap["status"] = ""
		}
		tasks[i] = taskMap
	}
	if err := d.Set("tasks", tasks); err != nil {
		return fmt.Errorf("error setting tasks: %w", err)
	}

	// Set ID
	d.SetId(fmt.Sprintf("cache-clean-tasks-%d-%d", req.Page, req.PerPage))

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total": response.Data.Total.String(),
			"tasks": tasks,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN cache clean tasks queried successfully")
	return nil
}
