package resource

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
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

// ResourceEdgenextScdnLogDownloadTask returns the SCDN log download task resource
func ResourceEdgenextScdnLogDownloadTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnLogDownloadTaskCreate,
		Read:   resourceScdnLogDownloadTaskRead,
		Update: resourceScdnLogDownloadTaskUpdate,
		Delete: resourceScdnLogDownloadTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"task_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Task name",
			},
			"is_use_template": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Whether to use template: 0-no, 1-yes",
			},
			"template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Template ID (required when is_use_template is 1)",
			},
			"data_source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data source: ng, cc, waf",
			},
			"download_fields": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Download fields",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"search_terms": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search conditions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Search key",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Search value",
						},
					},
				},
			},
			"file_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "File type: xls, csv, json",
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Start time (format: YYYY-MM-DD HH:MM:SS)",
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "End time (format: YYYY-MM-DD HH:MM:SS)",
			},
			"lang": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "zh_CN",
				Description: "Language: zh_CN, en_US, default: zh_CN",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the log download task",
			},
			"task_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The task ID",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task status: 0-not started, 1-in progress, 2-completed, 3-failed, 4-cancelled",
			},
			"download_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Download URL (available when task is completed)",
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
	}
}

func resourceScdnLogDownloadTaskCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request
	req := scdn.LogDownloadTaskAddRequest{
		TaskName:      d.Get("task_name").(string),
		IsUseTemplate: d.Get("is_use_template").(int),
		DataSource:    d.Get("data_source").(string),
		FileType:      d.Get("file_type").(string),
		StartTime:     d.Get("start_time").(string),
		EndTime:       d.Get("end_time").(string),
	}

	if v, ok := d.GetOk("template_id"); ok {
		req.TemplateID = v.(int)
	}

	if v, ok := d.GetOk("lang"); ok {
		req.Lang = v.(string)
	}

	// Download fields
	if v, ok := d.GetOk("download_fields"); ok {
		fieldsList := v.([]interface{})
		req.DownloadFields = make([]string, len(fieldsList))
		for i, item := range fieldsList {
			req.DownloadFields[i] = item.(string)
		}
	}

	// Search terms - convert from array format to map format
	if v, ok := d.GetOk("search_terms"); ok {
		termsList := v.([]interface{})
		req.SearchTerms = make(map[string]string)
		for _, item := range termsList {
			termMap := item.(map[string]interface{})
			key := termMap["key"].(string)
			value := termMap["value"].(string)
			// If multiple entries have the same key, the last one will be used
			req.SearchTerms[key] = value
		}
	}

	log.Printf("[INFO] Creating SCDN log download task: %s", req.TaskName)
	response, err := service.AddLogDownloadTask(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN log download task: %w", err)
	}

	log.Printf("[DEBUG] Log download task creation response: %+v", response)
	log.Printf("[DEBUG] Response Data: %+v", response.Data)
	log.Printf("[DEBUG] Task ID from response: %d", response.Data.TaskID)

	// If TaskID is 0, try to find the task by name
	taskID := response.Data.TaskID
	if taskID == 0 {
		log.Printf("[WARN] Task ID from creation response is 0, trying to find task by name: %s", req.TaskName)

		// Try to find the task by name
		listReq := scdn.LogDownloadTaskListRequest{
			Page:     1,
			PerPage:  100,
			Status:   -1,
			TaskName: req.TaskName,
		}

		listResp, err := service.ListLogDownloadTasks(listReq)
		if err == nil && len(listResp.Data.List) > 0 {
			// Find the most recent task with matching name
			for i := range listResp.Data.List {
				if listResp.Data.List[i].TaskName == req.TaskName {
					taskID = listResp.Data.List[i].TaskID
					log.Printf("[INFO] Found task by name: task_id=%d, task_name=%s", taskID, req.TaskName)
					break
				}
			}
		}

		// If still not found, return error
		if taskID == 0 {
			return fmt.Errorf("failed to create log download task: API returned task_id=0 and task not found by name")
		}
	}

	// Set ID from response or found task
	d.SetId(strconv.Itoa(taskID))
	log.Printf("[DEBUG] Set resource ID to: %d", taskID)

	// Set basic fields from creation request to avoid read issues immediately after creation
	// This prevents "inconsistent result" errors when task is just created
	if err := d.Set("task_id", taskID); err != nil {
		log.Printf("[WARN] Failed to set task_id: %v", err)
	}
	if err := d.Set("task_name", req.TaskName); err != nil {
		log.Printf("[WARN] Failed to set task_name: %v", err)
	}
	if err := d.Set("is_use_template", req.IsUseTemplate); err != nil {
		log.Printf("[WARN] Failed to set is_use_template: %v", err)
	}
	if req.TemplateID > 0 {
		if err := d.Set("template_id", req.TemplateID); err != nil {
			log.Printf("[WARN] Failed to set template_id: %v", err)
		}
	}
	if err := d.Set("data_source", req.DataSource); err != nil {
		log.Printf("[WARN] Failed to set data_source: %v", err)
	}
	if err := d.Set("download_fields", req.DownloadFields); err != nil {
		log.Printf("[WARN] Failed to set download_fields: %v", err)
	}
	if err := d.Set("file_type", req.FileType); err != nil {
		log.Printf("[WARN] Failed to set file_type: %v", err)
	}
	if err := d.Set("start_time", req.StartTime); err != nil {
		log.Printf("[WARN] Failed to set start_time: %v", err)
	}
	if err := d.Set("end_time", req.EndTime); err != nil {
		log.Printf("[WARN] Failed to set end_time: %v", err)
	}
	if req.Lang != "" {
		if err := d.Set("lang", req.Lang); err != nil {
			log.Printf("[WARN] Failed to set lang: %v", err)
		}
	}
	// Set default status
	if err := d.Set("status", "0"); err != nil {
		log.Printf("[WARN] Failed to set default status: %v", err)
	}
	// Set empty download_url initially
	if err := d.Set("download_url", ""); err != nil {
		log.Printf("[WARN] Failed to set empty download_url: %v", err)
	}
	// Set search_terms from request
	if len(req.SearchTerms) > 0 {
		searchTermsList := make([]map[string]interface{}, 0, len(req.SearchTerms))
		for key, value := range req.SearchTerms {
			searchTermsList = append(searchTermsList, map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}
		if err := d.Set("search_terms", searchTermsList); err != nil {
			log.Printf("[WARN] Failed to set search_terms: %v", err)
		}
	} else {
		if err := d.Set("search_terms", []interface{}{}); err != nil {
			log.Printf("[WARN] Failed to set empty search_terms: %v", err)
		}
	}

	log.Printf("[INFO] SCDN log download task created successfully: %s", d.Id())

	// Always call read to get full details from API
	// This ensures all fields are properly set from the API response
	return resourceScdnLogDownloadTaskRead(d, m)
}

func resourceScdnLogDownloadTaskRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	taskID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	// Query task list to find the task
	// Use task_name if available for more precise query
	req := scdn.LogDownloadTaskListRequest{
		Page:    1,
		PerPage: 100,
		Status:  -1, // -1 means don't filter by status (query all statuses)
	}

	// If task_name is available in state, use it for more precise query
	// Also use it if taskID is 0 (which might indicate the API didn't return the ID correctly)
	taskName := ""
	if taskNameVal, ok := d.GetOk("task_name"); ok && taskNameVal.(string) != "" {
		taskName = taskNameVal.(string)
		req.TaskName = taskName
		log.Printf("[DEBUG] Querying log download task by name: %s (task_id=%d)", req.TaskName, taskID)
	}

	log.Printf("[DEBUG] Querying log download tasks: page=%d, per_page=%d, status=%d, task_name=%s", req.Page, req.PerPage, req.Status, req.TaskName)
	response, err := service.ListLogDownloadTasks(req)
	if err != nil {
		return fmt.Errorf("failed to list log download tasks: %w", err)
	}

	log.Printf("[DEBUG] List response: total=%v, list_count=%d", response.Data.Total, len(response.Data.List))

	// Find the task by ID
	var task *scdn.LogDownloadTaskInfo
	for i := range response.Data.List {
		log.Printf("[DEBUG] Checking task: task_id=%d, task_name=%s", response.Data.List[i].TaskID, response.Data.List[i].TaskName)
		if response.Data.List[i].TaskID == taskID {
			task = &response.Data.List[i]
			log.Printf("[DEBUG] Found matching task: task_id=%d", taskID)
			break
		}
	}

	// If not found by ID, try searching by task_name if available
	if task == nil && taskName != "" {
		log.Printf("[DEBUG] Task not found by ID, trying to find by name: %s", taskName)
		for i := range response.Data.List {
			if response.Data.List[i].TaskName == taskName {
				task = &response.Data.List[i]
				log.Printf("[DEBUG] Found task by name: task_id=%d, task_name=%s", task.TaskID, task.TaskName)
				// Update the resource ID to the found task's ID
				if task.TaskID != taskID {
					log.Printf("[INFO] Updating resource ID from %d to %d based on task name match", taskID, task.TaskID)
					d.SetId(strconv.Itoa(task.TaskID))
					taskID = task.TaskID
				}
				break
			}
		}
	}

	// If still not found, try searching without task_name filter and across multiple pages
	if task == nil {
		log.Printf("[DEBUG] Task not found with current filters, trying without task_name filter")
		req.TaskName = ""
		req.Page = 1
		req.PerPage = 100

		// Search up to 5 pages
		for page := 1; page <= 5 && task == nil; page++ {
			req.Page = page
			response, err = service.ListLogDownloadTasks(req)
			if err != nil {
				log.Printf("[WARN] Failed to list tasks on page %d: %v", page, err)
				break
			}

			log.Printf("[DEBUG] Searching page %d: found %d tasks", page, len(response.Data.List))

			// Search by ID first
			for i := range response.Data.List {
				if response.Data.List[i].TaskID == taskID {
					task = &response.Data.List[i]
					log.Printf("[DEBUG] Found task by ID on page %d: task_id=%d", page, taskID)
					break
				}
			}

			// If not found by ID and we have task_name, search by name
			if task == nil && taskName != "" {
				for i := range response.Data.List {
					if response.Data.List[i].TaskName == taskName {
						task = &response.Data.List[i]
						log.Printf("[DEBUG] Found task by name on page %d: task_id=%d, task_name=%s", page, task.TaskID, task.TaskName)
						// Update the resource ID to the found task's ID
						if task.TaskID != taskID {
							log.Printf("[INFO] Updating resource ID from %d to %d based on task name match", taskID, task.TaskID)
							d.SetId(strconv.Itoa(task.TaskID))
							taskID = task.TaskID
						}
						break
					}
				}
			}

			// If no more tasks on this page, stop searching
			if len(response.Data.List) == 0 {
				break
			}
		}
	}

	if task == nil {
		// Task not found after exhaustive search
		log.Printf("[WARN] Log download task %d not found after searching multiple pages", taskID)

		// If taskID is 0, it means the resource was never properly created
		// Clear the ID to mark the resource for creation
		if taskID == 0 {
			log.Printf("[INFO] Task ID is 0, clearing resource ID to mark for creation")
			d.SetId("")
			return nil
		}

		// If taskID is not 0 but task not found, it might have been deleted
		// However, we should not clear the ID immediately - it might be a timing issue
		// Return an error so Terraform knows the resource state is inconsistent
		return fmt.Errorf("log download task %d not found - task may have been deleted or there is an API issue", taskID)
	}

	log.Printf("[DEBUG] Found log download task: task_id=%d, task_name=%s, status=%v", task.TaskID, task.TaskName, task.Status)

	// Set all fields from the task
	// Required fields
	if err := d.Set("task_id", task.TaskID); err != nil {
		return fmt.Errorf("error setting task_id: %w", err)
	}
	if err := d.Set("task_name", task.TaskName); err != nil {
		return fmt.Errorf("error setting task_name: %w", err)
	}

	// is_use_template is returned as string from API, convert to int
	isUseTemplate := 0
	if task.IsUseTemplate == "1" {
		isUseTemplate = 1
	} else if val, err := strconv.Atoi(task.IsUseTemplate); err == nil {
		isUseTemplate = val
	}
	if err := d.Set("is_use_template", isUseTemplate); err != nil {
		return fmt.Errorf("error setting is_use_template: %w", err)
	}

	// Optional fields
	if task.TemplateID > 0 {
		if err := d.Set("template_id", task.TemplateID); err != nil {
			log.Printf("[WARN] Failed to set template_id: %v", err)
		}
	}
	if err := d.Set("data_source", task.DataSource); err != nil {
		return fmt.Errorf("error setting data_source: %w", err)
	}
	if err := d.Set("download_fields", task.DownloadFields); err != nil {
		return fmt.Errorf("error setting download_fields: %w", err)
	}
	if err := d.Set("file_type", task.FileType); err != nil {
		return fmt.Errorf("error setting file_type: %w", err)
	}
	if err := d.Set("start_time", task.StartTime); err != nil {
		return fmt.Errorf("error setting start_time: %w", err)
	}
	if err := d.Set("end_time", task.EndTime); err != nil {
		return fmt.Errorf("error setting end_time: %w", err)
	}
	if task.Lang != "" {
		if err := d.Set("lang", task.Lang); err != nil {
			log.Printf("[WARN] Failed to set lang: %v", err)
		}
	}
	if err := d.Set("created_at", task.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set created_at: %v", err)
	}
	if err := d.Set("updated_at", task.UpdatedAt); err != nil {
		log.Printf("[WARN] Failed to set updated_at: %v", err)
	}

	// Status - convert to string if needed
	if task.Status != nil {
		var statusStr string
		if statusStrVal, ok := task.Status.(string); ok {
			statusStr = statusStrVal
		} else if statusInt, ok := task.Status.(int); ok {
			statusStr = fmt.Sprintf("%d", statusInt)
		} else if statusFloat, ok := task.Status.(float64); ok {
			statusStr = fmt.Sprintf("%.0f", statusFloat)
		} else {
			statusStr = fmt.Sprintf("%v", task.Status)
		}
		if err := d.Set("status", statusStr); err != nil {
			log.Printf("[WARN] Failed to set status: %v", err)
		}
	} else {
		// Set default status if nil
		if err := d.Set("status", "0"); err != nil {
			log.Printf("[WARN] Failed to set default status: %v", err)
		}
	}

	// Download URL - optional field
	if task.DownloadURL != nil {
		if url, ok := task.DownloadURL.(string); ok && url != "" {
			if err := d.Set("download_url", url); err != nil {
				log.Printf("[WARN] Failed to set download_url: %v", err)
			}
		} else {
			// Set empty string if download_url is null
			if err := d.Set("download_url", ""); err != nil {
				log.Printf("[WARN] Failed to set empty download_url: %v", err)
			}
		}
	} else {
		// Set empty string if download_url is nil
		if err := d.Set("download_url", ""); err != nil {
			log.Printf("[WARN] Failed to set empty download_url: %v", err)
		}
	}

	// Search terms - convert from map format to array format
	searchTerms := convertSearchTerms(task.SearchTerms)
	if len(searchTerms) > 0 {
		if err := d.Set("search_terms", searchTerms); err != nil {
			log.Printf("[WARN] Failed to set search_terms: %v", err)
		}
	} else {
		// Set empty list if search_terms is nil or empty
		if err := d.Set("search_terms", []interface{}{}); err != nil {
			log.Printf("[WARN] Failed to set empty search_terms: %v", err)
		}
	}

	log.Printf("[INFO] Log download task read successfully: task_id=%d", task.TaskID)
	return nil
}

func resourceScdnLogDownloadTaskUpdate(d *schema.ResourceData, m interface{}) error {
	// Log download tasks cannot be updated, only regenerated
	// For now, we'll delete and recreate
	return resourceScdnLogDownloadTaskDelete(d, m)
}

func resourceScdnLogDownloadTaskDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	taskID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid task ID: %w", err)
	}

	req := scdn.LogDownloadTaskDeleteRequest{
		TaskID: taskID,
	}

	log.Printf("[INFO] Deleting SCDN log download task: %d", taskID)
	_, err = service.DeleteLogDownloadTask(req)
	if err != nil {
		return fmt.Errorf("failed to delete SCDN log download task: %w", err)
	}

	log.Printf("[INFO] SCDN log download task deleted successfully: %d", taskID)
	d.SetId("")
	return nil
}
