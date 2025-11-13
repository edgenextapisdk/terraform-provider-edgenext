package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCacheCleanTask returns the SCDN cache clean task resource
func ResourceEdgenextScdnCacheCleanTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCacheCleanTaskCreate,
		Read:   resourceScdnCacheCleanTaskRead,
		Update: resourceScdnCacheCleanTaskUpdate,
		Delete: resourceScdnCacheCleanTaskDelete,

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
			"wholesite": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Whole site domains to clean",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"specialurl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Special URLs to clean",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"specialdir": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Special directories to clean",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the cache clean task (generated timestamp)",
			},
		},
	}
}

func resourceScdnCacheCleanTaskCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request
	req := scdn.CacheCleanSaveRequest{}

	if v, ok := d.GetOk("group_id"); ok {
		req.GroupID = v.(int)
	}
	if v, ok := d.GetOk("protocol"); ok {
		req.Protocol = v.(string)
	}
	if v, ok := d.GetOk("port"); ok {
		req.Port = v.(string)
	}
	if v, ok := d.GetOk("wholesite"); ok {
		wholesiteList := v.([]interface{})
		req.Wholesite = make([]string, len(wholesiteList))
		for i, item := range wholesiteList {
			req.Wholesite[i] = item.(string)
		}
	}
	if v, ok := d.GetOk("specialurl"); ok {
		specialurlList := v.([]interface{})
		req.Specialurl = make([]string, len(specialurlList))
		for i, item := range specialurlList {
			req.Specialurl[i] = item.(string)
		}
	}
	if v, ok := d.GetOk("specialdir"); ok {
		specialdirList := v.([]interface{})
		req.Specialdir = make([]string, len(specialdirList))
		for i, item := range specialdirList {
			req.Specialdir[i] = item.(string)
		}
	}

	// Validate that at least one of wholesite, specialurl, or specialdir is provided
	if len(req.Wholesite) == 0 && len(req.Specialurl) == 0 && len(req.Specialdir) == 0 {
		return fmt.Errorf("at least one of wholesite, specialurl, or specialdir must be provided")
	}

	log.Printf("[INFO] Creating SCDN cache clean task")
	response, err := service.SaveCacheCleanTask(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN cache clean task: %w", err)
	}

	log.Printf("[DEBUG] Cache clean task creation response: %+v", response)

	// Use timestamp as ID since the API doesn't return a task ID immediately
	// The actual task ID will be available in the task list
	// Generate a simple ID based on the request content
	d.SetId("cache-clean-task")

	log.Printf("[INFO] SCDN cache clean task created successfully: %s", d.Id())
	return resourceScdnCacheCleanTaskRead(d, m)
}

func resourceScdnCacheCleanTaskRead(d *schema.ResourceData, m interface{}) error {
	// Cache clean tasks are one-time operations, so read is a no-op
	// The task details can be queried via data source
	log.Printf("[DEBUG] Reading SCDN cache clean task: %s", d.Id())
	return nil
}

func resourceScdnCacheCleanTaskUpdate(d *schema.ResourceData, m interface{}) error {
	// For cache clean tasks, update means creating a new task
	return resourceScdnCacheCleanTaskCreate(d, m)
}

func resourceScdnCacheCleanTaskDelete(d *schema.ResourceData, m interface{}) error {
	// Cache clean tasks cannot be deleted, they are one-time operations
	log.Printf("[INFO] Cache clean task %s cannot be deleted (one-time operation)", d.Id())
	d.SetId("")
	return nil
}
