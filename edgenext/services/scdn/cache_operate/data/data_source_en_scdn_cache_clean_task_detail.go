package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnCacheCleanTaskDetail returns the SCDN cache clean task detail data source
func DataSourceEdgenextScdnCacheCleanTaskDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnCacheCleanTaskDetailRead,

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Task ID",
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
			"result": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Result filter: 1-success, 2-failed, 3-executing",
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
			"details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task detail list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution result",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution message",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time",
						},
						"directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Directory (present when this task type)",
						},
						"subdomain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subdomain (present when this task type)",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL (present when this task type)",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnCacheCleanTaskDetailRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.CacheCleanTaskDetailRequest{
		TaskID:  d.Get("task_id").(int),
		Page:    d.Get("page").(int),
		PerPage: d.Get("per_page").(int),
	}

	if v, ok := d.GetOk("result"); ok {
		req.Result = v.(int)
	}

	log.Printf("[INFO] Querying SCDN cache clean task detail: %d", req.TaskID)
	response, err := service.GetCacheCleanTaskDetail(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN cache clean task detail: %w", err)
	}

	// Set total
	if err := d.Set("total", response.Data.Total.String()); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Set details
	details := make([]map[string]interface{}, len(response.Data.List))
	for i, detail := range response.Data.List {
		detailMap := map[string]interface{}{
			"result":     detail.Result,
			"message":    detail.Message,
			"created_at": detail.CreatedAt,
			"updated_at": detail.UpdatedAt,
		}
		if detail.Directory != "" {
			detailMap["directory"] = detail.Directory
		}
		if detail.Subdomain != "" {
			detailMap["subdomain"] = detail.Subdomain
		}
		if detail.URL != "" {
			detailMap["url"] = detail.URL
		}
		details[i] = detailMap
	}
	if err := d.Set("details", details); err != nil {
		return fmt.Errorf("error setting details: %w", err)
	}

	// Set ID
	d.SetId(fmt.Sprintf("%d", req.TaskID))

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":   response.Data.Total.String(),
			"details": details,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN cache clean task detail queried successfully: %d", req.TaskID)
	return nil
}
