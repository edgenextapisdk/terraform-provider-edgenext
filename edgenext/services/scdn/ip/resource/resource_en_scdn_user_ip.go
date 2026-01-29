package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnUserIp returns the SCDN User IP List resource
func ResourceEdgenextScdnUserIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnUserIpCreate,
		Read:   resourceScdnUserIpRead,
		Update: resourceScdnUserIpUpdate,
		Delete: resourceScdnUserIpDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the IP list",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The remark/description for the IP list",
			},
			"file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the file containing IP list to upload",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the IP list",
			},
			"item_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Number of IPs in the list",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time",
			},
		},
	}
}

func resourceScdnUserIpCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.UserIpAddRequest{
		Name:   d.Get("name").(string),
		Remark: d.Get("remark").(string),
	}

	log.Printf("[INFO] Creating SCDN User IP List: %s", req.Name)
	response, err := service.AddUserIp(req)
	if err != nil {
		return fmt.Errorf("failed to create User IP List: %w", err)
	}

	d.SetId(response.Data.ID)
	log.Printf("[INFO] SCDN User IP List created successfully: %s", d.Id())

	// Handle file upload if file_path is provided
	if filePath, ok := d.GetOk("file_path"); ok {
		log.Printf("[INFO] Uploading IP file: %s for User IP List: %s", filePath, d.Id())
		_, err := service.UploadUserIpFile(d.Id(), filePath.(string), d.Get("remark").(string))
		if err != nil {
			return fmt.Errorf("failed to upload IP file: %w", err)
		}
	}

	return resourceScdnUserIpRead(d, m)
}

func resourceScdnUserIpRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	userIpID := d.Id()

	// Since there is no GetDetail API, we list all and find by ID.
	// This might be slow if there are many lists.
	req := scdn.UserIpListRequest{
		Page:    1,
		PerPage: 1000, // Fetch large number to find the item
	}

	log.Printf("[DEBUG] Reading SCDN User IP List: %s", userIpID)

	// We might need to loop pages if > 1000.
	// For now assuming 1000 is enough or implementing simple loop.
	found := false
	var targetIpList scdn.UserIpInfo

	for {
		resp, err := service.ListUserIps(req)
		if err != nil {
			return fmt.Errorf("failed to list User IP Lists: %w", err)
		}

		if len(resp.Data.List) == 0 {
			break
		}

		for _, item := range resp.Data.List {
			if item.ID == userIpID {
				targetIpList = item
				found = true
				break
			}
		}

		if found {
			break
		}

		// Pagination check - if retrieved less than per_page, we are done
		if len(resp.Data.List) < req.PerPage {
			break
		}

		req.Page++
	}

	if !found {
		log.Printf("[WARN] User IP List not found: %s", userIpID)
		d.SetId("")
		return nil
	}

	// Update state
	if err := d.Set("name", targetIpList.Name); err != nil {
		log.Printf("[WARN] Failed to set name: %v", err)
	}
	if err := d.Set("remark", targetIpList.Remark); err != nil {
		log.Printf("[WARN] Failed to set remark: %v", err)
	}
	if err := d.Set("item_num", targetIpList.ItemNum); err != nil {
		log.Printf("[WARN] Failed to set item_num: %v", err)
	}
	if err := d.Set("created_at", targetIpList.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set created_at: %v", err)
	}
	if err := d.Set("updated_at", targetIpList.UpdatedAt); err != nil {
		log.Printf("[WARN] Failed to set updated_at: %v", err)
	}

	return nil
}

func resourceScdnUserIpUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	userIpID := d.Id()

	req := scdn.UserIpSaveRequest{
		ID:     userIpID,
		Name:   d.Get("name").(string),
		Remark: d.Get("remark").(string),
	}

	log.Printf("[INFO] Updating SCDN User IP List: %s", userIpID)
	_, err := service.UpdateUserIp(req)
	if err != nil {
		return fmt.Errorf("failed to update User IP List: %w", err)
	}

	// Handle file upload if file_path changed
	if d.HasChange("file_path") {
		filePath := d.Get("file_path").(string)
		if filePath != "" {
			log.Printf("[INFO] Uploading new IP file: %s for User IP List: %s", filePath, userIpID)
			_, err := service.UploadUserIpFile(userIpID, filePath, d.Get("remark").(string))
			if err != nil {
				return fmt.Errorf("failed to upload IP file: %w", err)
			}
		}
	}

	return resourceScdnUserIpRead(d, m)
}

func resourceScdnUserIpDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	userIpID := d.Id()

	req := scdn.UserIpDelRequest{
		IDs: []string{userIpID},
	}

	log.Printf("[INFO] Deleting SCDN User IP List: %s", userIpID)
	_, err := service.DeleteUserIp(req)
	if err != nil {
		return fmt.Errorf("failed to delete User IP List: %w", err)
	}

	d.SetId("")
	return nil
}
