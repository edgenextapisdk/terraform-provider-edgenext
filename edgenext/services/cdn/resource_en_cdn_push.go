package cdn

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEdgenextCdnPush() *schema.Resource {
	return &schema.Resource{
		Create: resourcePushCreate,
		Read:   resourcePushRead,
		Update: nil, // Cache refresh does not support updates
		Delete: resourcePushDelete,

		Schema: map[string]*schema.Schema{
			"urls": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true, // Need to recreate task when urls list is updated
				Description: "List of URLs/directories to refresh, maximum 500 URLs per request",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Need to recreate task when push type is updated
				Description: "URL type for push: dir(directory), url(URL)",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task ID for this submission",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of successfully submitted URLs/directories",
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
							Description: "URL/Directory",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL type",
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

func resourcePushCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	urlsList := d.Get("urls").([]interface{})
	refreshType := d.Get("type").(string)

	log.Printf("[INFO] Creating cache refresh task: type=%s, URL count=%d", refreshType, len(urlsList))

	// Convert URL list
	var urls []string
	for _, url := range urlsList {
		urls = append(urls, url.(string))
	}

	// Call cache refresh API
	response, err := service.CacheRefresh(urls, refreshType)
	if err != nil {
		return fmt.Errorf("failed to create cache refresh task: %w", err)
	}

	// Set resource ID
	d.SetId(response.Data.TaskID)

	log.Printf("[INFO] Cache refresh task created successfully: %s", response.Data.TaskID)
	return resourcePushRead(d, m)
}

func resourcePushRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	taskID := d.Id()

	log.Printf("[INFO] Reading cache refresh task: %s", taskID)

	// Query refresh status by task ID
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", taskID)
	}
	response, err := service.QueryCacheRefreshByTaskID(taskIDInt)
	if err != nil {
		return fmt.Errorf("failed to read cache refresh task: %w", err)
	}
	if len(response.Data.List) == 0 {
		log.Printf("[WARN] Cache refresh task does not exist: %s", taskID)
		d.SetId("")
		return nil
	}

	// Set resource ID
	d.SetId(taskID)
	// Set response data
	if err := d.Set("task_id", taskID); err != nil {
		return fmt.Errorf("error setting task_id: %w", err)
	}
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}
	var list []map[string]interface{}
	for _, elem := range response.Data.List {
		elemMap := map[string]interface{}{
			"id":            elem.ID,
			"url":           elem.URL,
			"type":          elem.Type,
			"status":        elem.Status,
			"create_time":   elem.CreateTime,
			"complete_time": elem.CompleteTime,
		}
		list = append(list, elemMap)
	}
	// Set the list of successfully submitted URLs/directories
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[ERROR] Failed to set successfully submitted URL/directory list: %v", err)
		return err
	}

	return nil
}

func resourcePushDelete(d *schema.ResourceData, m interface{}) error {
	// API does not support deletion, can only no-op
	log.Printf("[WARN] Cache refresh task %s cannot be deleted (API limitation)", d.Id())
	d.SetId("") // Remove from state, Terraform considers it deleted
	return nil
}
