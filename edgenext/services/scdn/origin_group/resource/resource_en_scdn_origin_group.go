package resource

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnOriginGroup returns the SCDN origin group resource
func ResourceEdgenextScdnOriginGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnOriginGroupCreate,
		Read:   resourceScdnOriginGroupRead,
		Update: resourceScdnOriginGroupUpdate,
		Delete: resourceScdnOriginGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: resourceScdnOriginGroupCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"origin_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Origin group ID. Required for update/delete, computed for create. If provided during create, will update existing origin group instead.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Origin group name (2-16 characters)",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark (2-64 characters)",
			},
			"origins": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Origin list (at least 1)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Origin ID (0 for new, >0 for update)",
						},
						"origin_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Origin type: 0-IP, 1-domain",
						},
						"records": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Origin record list (at least 1)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Origin address",
									},
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Origin port (1-65535)",
									},
									"priority": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Weight (1-100)",
									},
									"view": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Origin type: primary-backup, backup-backup",
									},
									"host": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Origin Host",
									},
								},
							},
						},
						"protocol_ports": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Protocol port mapping (at least 1)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Protocol: 0-http, 1-https",
									},
									"listen_ports": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Listen port list",
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"origin_protocol": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Origin protocol: 0-http, 1-https, 2-follow",
						},
						"load_balance": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Load balance strategy: 0-ip_hash, 1-round_robin, 2-cookie",
						},
					},
				},
			},
			// Computed fields
			"member_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Member ID",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username",
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
		},
	}
}

func resourceScdnOriginGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Check if origin_group_id is provided (for updating existing origin group)
	// Only do this if resource doesn't already exist (d.Id() is empty)
	if d.Id() == "" {
		if originGroupID, hasOriginGroupID := d.GetOk("origin_group_id"); hasOriginGroupID {
			// Update existing origin group
			originGroupIDInt := originGroupID.(int)
			if originGroupIDInt > 0 {
				log.Printf("[INFO] origin_group_id provided (%d) for new resource, updating existing origin group instead of creating", originGroupIDInt)
				// Set the ID so Update function can use it
				d.SetId(strconv.Itoa(originGroupIDInt))
				return resourceScdnOriginGroupUpdate(d, m)
			}
		}
	}

	req := scdn.OriginGroupCreateRequest{
		Name: d.Get("name").(string),
	}

	if remark, ok := d.GetOk("remark"); ok {
		req.Remark = remark.(string)
	}

	// Build origins
	origins := d.Get("origins").([]interface{})
	req.Origins = make([]scdn.OriginGroupOrigin, len(origins))
	for i, origin := range origins {
		originMap := origin.(map[string]interface{})
		originCfg := scdn.OriginGroupOrigin{
			OriginType:     originMap["origin_type"].(int),
			OriginProtocol: originMap["origin_protocol"].(int),
			LoadBalance:    originMap["load_balance"].(int),
		}

		if id, ok := originMap["id"].(int); ok {
			originCfg.ID = id
		}

		// Build records
		records := originMap["records"].([]interface{})
		originCfg.Records = make([]scdn.OriginGroupRecord, len(records))
		for j, record := range records {
			recordMap := record.(map[string]interface{})
			originCfg.Records[j] = scdn.OriginGroupRecord{
				Value:    recordMap["value"].(string),
				Port:     recordMap["port"].(int),
				Priority: recordMap["priority"].(int),
				View:     recordMap["view"].(string),
			}
			if host, ok := recordMap["host"].(string); ok && host != "" {
				originCfg.Records[j].Host = host
			}
		}

		// Build protocol_ports
		protocolPorts := originMap["protocol_ports"].([]interface{})
		originCfg.ProtocolPorts = make([]scdn.OriginGroupProtocolPort, len(protocolPorts))
		for j, pp := range protocolPorts {
			ppMap := pp.(map[string]interface{})
			listenPorts := ppMap["listen_ports"].([]interface{})
			originCfg.ProtocolPorts[j] = scdn.OriginGroupProtocolPort{
				Protocol:    ppMap["protocol"].(int),
				ListenPorts: make([]int, len(listenPorts)),
			}
			for k, lp := range listenPorts {
				originCfg.ProtocolPorts[j].ListenPorts[k] = lp.(int)
			}
		}

		req.Origins[i] = originCfg
	}

	log.Printf("[INFO] Creating SCDN origin group: name=%s", req.Name)
	response, err := service.CreateOriginGroup(req)
	if err != nil {
		return fmt.Errorf("failed to create origin group: %w", err)
	}

	originGroupID := response.Data.ID
	d.SetId(strconv.Itoa(originGroupID))
	if err := d.Set("origin_group_id", originGroupID); err != nil {
		return fmt.Errorf("error setting origin_group_id: %w", err)
	}

	return resourceScdnOriginGroupRead(d, m)
}

func resourceScdnOriginGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	originGroupID := d.Get("origin_group_id").(int)
	if originGroupID == 0 {
		// Try to parse from resource ID
		idStr := d.Id()
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				originGroupID = id
				if err := d.Set("origin_group_id", originGroupID); err != nil {
					return fmt.Errorf("error setting origin_group_id: %w", err)
				}
			} else {
				return fmt.Errorf("origin_group_id is required for reading origin group")
			}
		} else {
			return fmt.Errorf("origin_group_id is required for reading origin group")
		}
	}

	req := scdn.OriginGroupDetailRequest{
		ID: originGroupID,
	}

	log.Printf("[INFO] Reading SCDN origin group: origin_group_id=%d", originGroupID)
	response, err := service.GetOriginGroupDetail(req)
	if err != nil {
		return fmt.Errorf("failed to read origin group: %w", err)
	}

	originGroup := response.Data.OriginGroup

	// Set fields
	if err := d.Set("name", originGroup.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("remark", originGroup.Remark); err != nil {
		log.Printf("[WARN] Failed to set remark: %v", err)
	}
	if err := d.Set("member_id", originGroup.MemberID); err != nil {
		log.Printf("[WARN] Failed to set member_id: %v", err)
	}
	if err := d.Set("username", originGroup.Username); err != nil {
		log.Printf("[WARN] Failed to set username: %v", err)
	}
	if err := d.Set("created_at", originGroup.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set created_at: %v", err)
	}
	if err := d.Set("updated_at", originGroup.UpdatedAt); err != nil {
		log.Printf("[WARN] Failed to set updated_at: %v", err)
	}

	// Set origins
	originsList := make([]map[string]interface{}, len(originGroup.Origins))
	for i, origin := range originGroup.Origins {
		originMap := map[string]interface{}{
			"id":              origin.ID,
			"origin_type":     origin.OriginType,
			"origin_protocol": origin.OriginProtocol,
			"load_balance":    origin.LoadBalance,
		}

		// Set records
		recordsList := make([]map[string]interface{}, len(origin.Records))
		for j, record := range origin.Records {
			recordsList[j] = map[string]interface{}{
				"value":    record.Value,
				"port":     record.Port,
				"priority": record.Priority,
				"view":     record.View,
				"host":     record.Host,
			}
		}
		originMap["records"] = recordsList

		// Set protocol_ports
		protocolPortsList := make([]map[string]interface{}, len(origin.ProtocolPorts))
		for j, pp := range origin.ProtocolPorts {
			listenPortsList := make([]int, len(pp.ListenPorts))
			copy(listenPortsList, pp.ListenPorts)
			protocolPortsList[j] = map[string]interface{}{
				"protocol":     pp.Protocol,
				"listen_ports": listenPortsList,
			}
		}
		originMap["protocol_ports"] = protocolPortsList

		originsList[i] = originMap
	}
	if err := d.Set("origins", originsList); err != nil {
		return fmt.Errorf("error setting origins: %w", err)
	}

	return nil
}

func resourceScdnOriginGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// For update, always use the resource ID from state, not from config
	// This ensures we update the existing resource, not try to change its ID
	var originGroupID int
	idStr := d.Id()
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			originGroupID = id
		} else {
			return fmt.Errorf("invalid resource ID: %s", idStr)
		}
	} else {
		// Fallback to config if no ID in state (shouldn't happen in normal flow)
		originGroupID = d.Get("origin_group_id").(int)
		if originGroupID == 0 {
			return fmt.Errorf("origin_group_id is required for updating origin group")
		}
		// Set the ID from config
		d.SetId(strconv.Itoa(originGroupID))
	}

	// Ensure origin_group_id in state matches the resource ID
	if err := d.Set("origin_group_id", originGroupID); err != nil {
		log.Printf("[WARN] Failed to set origin_group_id: %v", err)
	}

	req := scdn.OriginGroupUpdateRequest{
		ID:   originGroupID,
		Name: d.Get("name").(string),
	}

	if remark, ok := d.GetOk("remark"); ok {
		req.Remark = remark.(string)
	}

	// Build origins (same as create)
	origins := d.Get("origins").([]interface{})
	req.Origins = make([]scdn.OriginGroupOrigin, len(origins))
	for i, origin := range origins {
		originMap := origin.(map[string]interface{})
		originCfg := scdn.OriginGroupOrigin{
			OriginType:     originMap["origin_type"].(int),
			OriginProtocol: originMap["origin_protocol"].(int),
			LoadBalance:    originMap["load_balance"].(int),
		}

		if id, ok := originMap["id"].(int); ok {
			originCfg.ID = id
		}

		// Build records
		records := originMap["records"].([]interface{})
		originCfg.Records = make([]scdn.OriginGroupRecord, len(records))
		for j, record := range records {
			recordMap := record.(map[string]interface{})
			originCfg.Records[j] = scdn.OriginGroupRecord{
				Value:    recordMap["value"].(string),
				Port:     recordMap["port"].(int),
				Priority: recordMap["priority"].(int),
				View:     recordMap["view"].(string),
			}
			if host, ok := recordMap["host"].(string); ok && host != "" {
				originCfg.Records[j].Host = host
			}
		}

		// Build protocol_ports
		protocolPorts := originMap["protocol_ports"].([]interface{})
		originCfg.ProtocolPorts = make([]scdn.OriginGroupProtocolPort, len(protocolPorts))
		for j, pp := range protocolPorts {
			ppMap := pp.(map[string]interface{})
			listenPorts := ppMap["listen_ports"].([]interface{})
			originCfg.ProtocolPorts[j] = scdn.OriginGroupProtocolPort{
				Protocol:    ppMap["protocol"].(int),
				ListenPorts: make([]int, len(listenPorts)),
			}
			for k, lp := range listenPorts {
				originCfg.ProtocolPorts[j].ListenPorts[k] = lp.(int)
			}
		}

		req.Origins[i] = originCfg
	}

	log.Printf("[INFO] Updating SCDN origin group: origin_group_id=%d", originGroupID)
	_, err := service.UpdateOriginGroup(req)
	if err != nil {
		return fmt.Errorf("failed to update origin group: %w", err)
	}

	return resourceScdnOriginGroupRead(d, m)
}

func resourceScdnOriginGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	originGroupID := d.Get("origin_group_id").(int)
	if originGroupID == 0 {
		// Try to parse from resource ID
		idStr := d.Id()
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				originGroupID = id
			}
		}
	}

	if originGroupID == 0 {
		return fmt.Errorf("origin_group_id is required for deleting origin group")
	}

	req := scdn.OriginGroupDeleteRequest{
		IDs: []int{originGroupID},
	}

	log.Printf("[INFO] Deleting SCDN origin group: origin_group_id=%d", originGroupID)
	_, err := service.DeleteOriginGroups(req)
	if err != nil {
		return fmt.Errorf("failed to delete origin group: %w", err)
	}

	d.SetId("")
	return nil
}

// resourceScdnOriginGroupCustomizeDiff customizes the diff to ignore origin_group_id changes
// when the resource already exists, since origin_group_id should be immutable
func resourceScdnOriginGroupCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	// If resource already exists (has ID), ignore changes to origin_group_id
	// The origin_group_id should always match the resource ID
	if d.Id() != "" {
		// Get the current origin_group_id from state
		oldOriginGroupID, _ := d.GetChange("origin_group_id")
		// Get the new origin_group_id from config
		newOriginGroupID := d.Get("origin_group_id")

		// If they differ, clear the diff since we'll use the ID from state
		if oldOriginGroupID != newOriginGroupID {
			// Parse ID from resource ID
			idStr := d.Id()
			if idStr != "" {
				if id, err := strconv.Atoi(idStr); err == nil {
					// Force origin_group_id to match resource ID
					if err := d.SetNew("origin_group_id", id); err != nil {
						log.Printf("[WARN] Failed to set origin_group_id in diff: %v", err)
					}
				}
			}
		}
	}
	return nil
}
