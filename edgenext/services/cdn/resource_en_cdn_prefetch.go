package cdn

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEdgenextCdnPrefetch() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrefetchCreate,
		Read:   resourcePrefetchRead,
		Update: nil, // Prefetch does not support updates
		Delete: resourcePrefetchDelete,

		Schema: map[string]*schema.Schema{
			"urls": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true, // Need to recreate task when urls list is updated
				Description: "List of URLs to prefetch, maximum 500 URLs per request",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task ID for this submission",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of successfully submitted URLs",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of successfully submitted URLs",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
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

func resourcePrefetchCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	urlsList := d.Get("urls").([]interface{})

	log.Printf("[INFO] Creating file prefetch task: URL count=%d", len(urlsList))

	// Convert URL list
	var urls []string
	for _, url := range urlsList {
		urls = append(urls, url.(string))
	}

	// Call file prefetch API
	response, err := service.FilePrefetch(urls)
	if err != nil {
		return fmt.Errorf("failed to create file prefetch task: %w", err)
	}

	// Set resource ID
	d.SetId(response.Data.TaskID)

	log.Printf("[INFO] File prefetch task created successfully: %s", response.Data.TaskID)
	return resourcePrefetchRead(d, m)
}

func resourcePrefetchRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	taskID := d.Id()

	log.Printf("[INFO] Reading file prefetch task: %s", taskID)

	// Query prefetch status by task ID
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", taskID)
	}
	response, err := service.QueryFilePrefetchByTaskID(taskIDInt)
	if err != nil {
		return fmt.Errorf("failed to read file prefetch task: %w", err)
	}
	if len(response.Data.List) == 0 {
		log.Printf("[WARN] File prefetch task does not exist: %s", taskID)
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
			"status":        elem.Status,
			"create_time":   elem.CreateTime,
			"complete_time": elem.CompleteTime,
		}
		list = append(list, elemMap)
	}
	// Set the list of successfully submitted URLs
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[ERROR] Failed to set successfully submitted URL list: %v", err)
		return err
	}

	return nil
}

func resourcePrefetchDelete(d *schema.ResourceData, m interface{}) error {
	// API does not support deletion, can only no-op
	log.Printf("[WARN] File prefetch task %s cannot be deleted (API limitation)", d.Id())
	d.SetId("") // Remove from state, Terraform considers it deleted
	return nil
}
