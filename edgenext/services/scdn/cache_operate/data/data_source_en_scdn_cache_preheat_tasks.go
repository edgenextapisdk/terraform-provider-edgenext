package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCachePreheatTasks returns the SCDN cache preheat tasks data source
func DataSourceEdgenextScdnCachePreheatTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCachePreheatTasksRead,

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
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status filter",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL filter",
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
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID",
						},
						"user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User ID",
						},
						"time_create": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"time_update": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time",
						},
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID",
						},
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Domain ID",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status: 1-Prefetch waiting, 2-Prefetch pending, 3-Prefetch successful, 4-Prefetch failed",
						},
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Weight",
						},
						"sub_user_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Sub user ID",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User name",
						},
						"strategy_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy ID",
						},
						"strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Strategy",
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

func dataSourceScdnCachePreheatTasksRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.CachePreheatTaskListRequest{
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}

	if v, ok := d.GetOk("status"); ok {
		req.Status = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		req.URL = v.(string)
	}

	log.Printf("[INFO] Querying SCDN cache preheat tasks")
	response, err := service.GetCachePreheatTaskList(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN cache preheat tasks: %w", err)
	}

	// Set total
	if err := d.Set("total", response.Data.Total.String()); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Set tasks
	tasks := make([]map[string]interface{}, len(response.Data.List))
	for i, task := range response.Data.List {
		taskMap := map[string]interface{}{
			"id":                 task.ID,
			"user_id":            task.UserID,
			"time_create":        task.TimeCreate,
			"time_update":        task.TimeUpdate,
			"task_id":            task.TaskID,
			"domain_id":          task.DomainID,
			"url":                task.URL,
			"status":             task.Status,
			"total":              task.Total,
			"weight":             task.Weight,
			"sub_user_id":        task.SubUserID,
			"user_name":          task.UserName,
			"strategy_id":        task.StrategyID,
			"strategy":           task.Strategy,
			"operator_user_name": task.OperatorUserName,
		}
		tasks[i] = taskMap
	}
	if err := d.Set("tasks", tasks); err != nil {
		return fmt.Errorf("error setting tasks: %w", err)
	}

	// Set ID
	d.SetId(fmt.Sprintf("cache-preheat-tasks-%d-%d", req.Page, req.PerPage))

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

	log.Printf("[INFO] SCDN cache preheat tasks queried successfully")
	return nil
}
