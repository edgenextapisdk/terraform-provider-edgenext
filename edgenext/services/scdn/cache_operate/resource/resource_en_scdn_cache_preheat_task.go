package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCachePreheatTask returns the SCDN cache preheat task resource
func ResourceEdgenextScdnCachePreheatTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCachePreheatTaskCreate,
		Read:   resourceScdnCachePreheatTaskRead,
		Update: resourceScdnCachePreheatTaskUpdate,
		Delete: resourceScdnCachePreheatTaskDelete,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Group ID, can refresh cache by group",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protocol: http/https; only valid when refreshing by group",
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Website port, only needed for special ports; only valid when refreshing by group",
			},
			"preheat_url": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Preheat URLs",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the preheat task (generated timestamp)",
			},
			"error_url": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of URLs with preheat errors",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceScdnCachePreheatTaskCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request
	req := scdn.CachePreheatSaveRequest{}

	if v, ok := d.GetOk("group_id"); ok {
		req.GroupID = v.(int)
	}
	if v, ok := d.GetOk("protocol"); ok {
		req.Protocol = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		req.Port = v.(string)
	}
	if v, ok := d.GetOk("preheat_url"); ok {
		preheatURLList := v.([]interface{})
		req.PreheatURL = make([]string, len(preheatURLList))
		for i, item := range preheatURLList {
			req.PreheatURL[i] = item.(string)
		}
	}

	// Validate that preheat_url is provided
	if len(req.PreheatURL) == 0 {
		return fmt.Errorf("preheat_url is required")
	}

	log.Printf("[INFO] Creating SCDN cache preheat task")
	response, err := service.SaveCachePreheatTask(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN cache preheat task: %w", err)
	}

	log.Printf("[DEBUG] Cache preheat task creation response: %+v", response)

	// Use timestamp as ID since the API doesn't return a task ID immediately
	// The actual task ID will be available in the task list
	// Generate a simple ID based on the request content
	d.SetId("cache-preheat-task")

	// Set error_url if present
	errorURLs := response.Data.GetErrorURLs()
	if len(errorURLs) > 0 {
		if err := d.Set("error_url", errorURLs); err != nil {
			log.Printf("[WARN] Failed to set error_url: %v", err)
		}
	}

	log.Printf("[INFO] SCDN cache preheat task created successfully: %s", d.Id())
	return resourceScdnCachePreheatTaskRead(d, m)
}

func resourceScdnCachePreheatTaskRead(d *schema.ResourceData, m interface{}) error {
	// Cache preheat tasks are one-time operations, so read is a no-op
	// The task details can be queried via data source
	log.Printf("[DEBUG] Reading SCDN cache preheat task: %s", d.Id())
	return nil
}

func resourceScdnCachePreheatTaskUpdate(d *schema.ResourceData, m interface{}) error {
	// For cache preheat tasks, update means creating a new task
	return resourceScdnCachePreheatTaskCreate(d, m)
}

func resourceScdnCachePreheatTaskDelete(d *schema.ResourceData, m interface{}) error {
	// Cache preheat tasks cannot be deleted, they are one-time operations
	log.Printf("[INFO] Cache preheat task %s cannot be deleted (one-time operation)", d.Id())
	d.SetId("")
	return nil
}
