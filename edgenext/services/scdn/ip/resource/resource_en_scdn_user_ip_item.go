package resource

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnUserIpItem returns the SCDN User IP Item resource
func ResourceEdgenextScdnUserIpItem() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnUserIpItemCreate,
		Read:   resourceScdnUserIpItemRead,
		Update: resourceScdnUserIpItemUpdate,
		Delete: resourceScdnUserIpItemDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceScdnUserIpItemImport,
		},

		Schema: map[string]*schema.Schema{
			"user_ip_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the IP list to which this item belongs",
			},
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IP address or CIDR",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark for the IP item",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID (UUID) of the IP item",
			},
			"format_created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time",
			},
			"format_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time",
			},
		},
	}
}

func resourceScdnUserIpItemCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	userIpID := d.Get("user_ip_id").(int)
	ip := d.Get("ip").(string)
	remark := d.Get("remark").(string)

	req := scdn.UserIpItemAddRequest{
		UserIpID: strconv.Itoa(userIpID),
		IP:       ip,
		Remark:   remark,
	}

	log.Printf("[INFO] Creating SCDN User IP Item: %s in List %d", req.IP, userIpID)
	response, err := service.AddUserIpItem(req)
	if err != nil {
		return fmt.Errorf("failed to create User IP Item: %w", err)
	}

	// Response Data.IDs is an array. We assume we only added ONE item (one IP in req).
	// But API allows text blob. Here we send single IP.
	// If multiple IDs returned, we take the first one?
	// The API doc example says "ids": ["uuid_1"].
	if len(response.Data.IDs) == 0 {
		return fmt.Errorf("no ID returned after creating User IP Item")
	}

	d.SetId(response.Data.IDs[0])
	log.Printf("[INFO] SCDN User IP Item created successfully: %s", d.Id())

	return resourceScdnUserIpItemRead(d, m)
}

func resourceScdnUserIpItemRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	ipItemID := d.Id()
	userIpID := d.Get("user_ip_id").(int)
	// We might need "ip" to filter efficiently, but if we are reading executing `terraform import` or just refresh,
	// `ip` might be in state. If not, we list all items in the list to find this ID.
	targetIP := d.Get("ip").(string) // May be empty if not in state yet? No, d.Get returns from state.

	req := scdn.UserIpItemListRequest{
		UserIpID: userIpID,
		// If we have IP, use it filter to speed up
		IP:      targetIP,
		Page:    1,
		PerPage: 100,
	}

	log.Printf("[DEBUG] Reading SCDN User IP Item: %s", ipItemID)

	found := false
	var targetItem scdn.UserIpItem

	// Pagination loop
	for {
		resp, err := service.ListUserIpItems(req)
		if err != nil {
			return fmt.Errorf("failed to list User IP Items: %w", err)
		}

		if len(resp.Data.List) == 0 {
			break
		}

		for _, item := range resp.Data.List {
			if item.ID == ipItemID {
				targetItem = item
				found = true
				break
			}
		}

		if found {
			break
		}

		if len(resp.Data.List) < req.PerPage {
			break
		}

		req.Page++
	}

	// If not found with IP filter, maybe IP changed or something?
	// Try searching WITHOUT IP filter if we used it and didn't find it?
	// But `ip` is part of uniqueness often.
	// If ID is UUID, it should be unique.
	// Let's assume if not found, it's gone.
	if !found {
		log.Printf("[WARN] User IP Item not found: %s", ipItemID)
		d.SetId("")
		return nil
	}

	d.Set("ip", targetItem.IP)
	d.Set("remark", targetItem.Remark)
	d.Set("user_ip_id", targetItem.UserIpId)
	d.Set("format_created_at", targetItem.FormatCreatedAt)
	d.Set("format_updated_at", targetItem.FormatUpdatedAt)

	return nil
}

func resourceScdnUserIpItemUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	ipItemID := d.Id()
	userIpID := d.Get("user_ip_id").(int)

	req := scdn.UserIpItemEditRequest{
		ID:       ipItemID,
		UserIpID: strconv.Itoa(userIpID),
		IP:       d.Get("ip").(string),
		Remark:   d.Get("remark").(string),
	}

	log.Printf("[INFO] Updating SCDN User IP Item: %s", ipItemID)
	_, err := service.UpdateUserIpItem(req)
	if err != nil {
		return fmt.Errorf("failed to update User IP Item: %w", err)
	}

	return resourceScdnUserIpItemRead(d, m)
}

func resourceScdnUserIpItemDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	ipItemID := d.Id()
	userIpID := d.Get("user_ip_id").(int)

	req := scdn.UserIpItemDelRequest{
		IDs:      []string{ipItemID},
		UserIpID: strconv.Itoa(userIpID),
	}

	log.Printf("[INFO] Deleting SCDN User IP Item: %s", ipItemID)
	_, err := service.DeleteUserIpItem(req)
	if err != nil {
		return fmt.Errorf("failed to delete User IP Item: %w", err)
	}

	d.SetId("")
	return nil
}

func resourceScdnUserIpItemImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Expected format: user_ip_id:item_id
	if !strings.Contains(d.Id(), ":") {
		return nil, fmt.Errorf("invalid import import id. Expected format: user_ip_id:item_id")
	}

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import format. Expected format: user_ip_id:item_id")
	}

	userIpID, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid user_ip_id: %s", parts[0])
	}

	d.Set("user_ip_id", userIpID)
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
