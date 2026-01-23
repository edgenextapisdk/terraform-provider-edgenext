package cdn

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEdgenextCdnPurge() *schema.Resource {
	return &schema.Resource{
		Create: resourcePurgeCreate,
		Read:   resourcePurgeRead,
		Update: nil, // Cache purge does not support updates
		Delete: resourcePurgeDelete,

		Schema: map[string]*schema.Schema{
			"urls": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true, // Need to recreate task when urls list is updated
				Description: "List of URLs/directories to purge, maximum 500 URLs per request",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true, // Need to recreate task when purge type is updated
				Description: "URL type for purge: dir(directory), url(URL)",
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

func resourcePurgeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	urlsList := d.Get("urls").([]interface{})
	purgeType := d.Get("type").(string)

	log.Printf("[INFO] Creating cache purge task: type=%s, URL count=%d", purgeType, len(urlsList))

	// Convert URL list
	var urls []string
	for _, url := range urlsList {
		urls = append(urls, url.(string))
	}

	// Call cache refresh API
	response, err := service.CacheRefresh(urls, purgeType)
	if err != nil {
		return fmt.Errorf("failed to create cache purge task: %w", err)
	}

	// Set resource ID
	d.SetId(response.Data.TaskID)

	log.Printf("[INFO] Cache purge task created successfully: %s", response.Data.TaskID)
	return resourcePurgeRead(d, m)
}

func resourcePurgeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	taskID := d.Id()

	log.Printf("[INFO] Reading cache purge task: %s", taskID)

	// Query purge status by task ID
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", taskID)
	}
	response, err := service.QueryCacheRefreshByTaskID(taskIDInt)
	if err != nil {
		return fmt.Errorf("failed to read cache purge task: %w", err)
	}
	if len(response.Data.List) == 0 {
		log.Printf("[WARN] Cache purge task does not exist: %s", taskID)
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

func resourcePurgeDelete(d *schema.ResourceData, m interface{}) error {
	// API does not support deletion, can only no-op
	log.Printf("[WARN] Cache purge task %s cannot be deleted (API limitation)", d.Id())
	d.SetId("") // Remove from state, Terraform considers it deleted
	return nil
}
