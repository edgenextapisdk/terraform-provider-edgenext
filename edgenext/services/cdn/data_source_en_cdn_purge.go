package cdn

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextCdnPurge() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePurgeRead,

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task ID for querying the purge status of a specific task",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of records",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of successfully submitted URLs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL ID",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"complete_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Completion time",
						},
					},
				},
			},
		},
	}
}

func dataSourcePurgeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	taskID := d.Get("task_id").(string)

	log.Printf("[INFO] Querying file purge task, task_id: %s", taskID)

	var response *FilePurgeQueryResponse
	var err error

	// Query by task ID
	log.Printf("[INFO] Querying by task ID: %s", taskID)
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", taskID)
	}
	response, err = service.QueryFilePurgeByTaskID(taskIDInt)
	if err != nil {
		return fmt.Errorf("failed to query by task ID: %w", err)
	}
	// Set resource ID
	d.SetId(taskID)

	// Set response data
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}
	// Set the list of successfully submitted URLs
	var list []map[string]interface{}
	for _, elem := range response.Data.List {
		elemMap := map[string]interface{}{
			"id":            elem.ID,
			"url":           elem.URL,
			"status":        elem.Status,
			"create_time":   elem.CreateTime,
			"complete_time": elem.CompleteTime,
		}
		list = append(list, elemMap)
	}
	if err := d.Set("list", list); err != nil {
		return fmt.Errorf("error setting list: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"task_id": taskID,
			"total":   response.Data.Total,
			"list":    list,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] File purge task query successful, total %d records", response.Data.Total)
	return nil
}

// DataSourceEdgenextCdnPurges data source for querying multiple file purge tasks
func DataSourceEdgenextCdnPurges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePurgesRead,

		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time, format: YYYY-MM-DD",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time, format: YYYY-MM-DD",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL",
			},
			"page_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1",
				Description: "Page number to retrieve, default 1",
			},
			"page_size": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "50",
				Description: "Page size, default 50, range 1-500",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of records",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of successfully submitted URLs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL ID",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"complete_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Completion time",
						},
					},
				},
			},
		},
	}
}

func dataSourcePurgesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	startTime := d.Get("start_time").(string)
	endTime := d.Get("end_time").(string)
	url := d.Get("url").(string)
	pageNumber := d.Get("page_number").(string)
	pageSize := d.Get("page_size").(string)

	log.Printf("[INFO] Querying multiple file purge tasks: %s to %s", startTime, endTime)

	// Query by time range
	response, err := service.QueryFilePurgeByTimeRange(startTime, endTime, url, pageNumber, pageSize)
	if err != nil {
		return fmt.Errorf("failed to query file purge tasks: %w", err)
	}

	// Set response data
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}
	var list []map[string]interface{}
	ids := make([]string, 0)
	for _, elem := range response.Data.List {
		elemMap := map[string]interface{}{
			"id":            elem.ID,
			"url":           elem.URL,
			"status":        elem.Status,
			"create_time":   elem.CreateTime,
			"complete_time": elem.CompleteTime,
		}
		list = append(list, elemMap)
		ids = append(ids, elem.ID)
	}
	// Set resource ID
	d.SetId(helper.DataResourceIdsHash(ids))
	// Set successfully submitted URL list
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[ERROR] Failed to set successfully submitted URL list: %v", err)
		return err
	}

	// Write result to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"start_time":  startTime,
			"end_time":    endTime,
			"url":         url,
			"page_number": pageNumber,
			"page_size":   pageSize,
			"total":       response.Data.Total,
			"list":        list,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] Multiple file purge tasks query successful, total %d records", len(list))
	return nil
}
