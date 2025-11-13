package domain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnOrigin returns the SCDN origin resource
func ResourceEdgenextScdnOrigin() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnOriginCreate,
		Read:   resourceScdnOriginRead,
		Update: resourceScdnOriginUpdate,
		Delete: resourceScdnOriginDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain to add origins to",
			},
			"protocol": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS)",
			},
			"listen_ports": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The listening ports of the origin server",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"origin_protocol": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS), 2 (Follow)",
			},
			"load_balance": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The load balancing method. Valid values: 0 (IP hash), 1 (Round robin), 2 (Cookie)",
			},
			"origin_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The origin type. Valid values: 0 (IP), 1 (Domain)",
			},
			"records": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The origin records",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"view": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The view of the record",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the record (IP address or domain)",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The port of the record",
						},
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The priority of the record",
						},
					},
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the origin",
			},
			"listen_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The listening port of the origin server (single port from API)",
			},
		},
	}
}

func resourceScdnOriginCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Handle domain_id which might be string or int
	var domainID int
	switch v := d.Get("domain_id").(type) {
	case int:
		domainID = v
	case string:
		var err error
		domainID, err = strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid domain_id: %w", err)
		}
	default:
		return fmt.Errorf("domain_id must be int or string")
	}

	// Build origin
	origin := scdn.Origin{
		Protocol:       d.Get("protocol").(int),
		OriginProtocol: d.Get("origin_protocol").(int),
		LoadBalance:    d.Get("load_balance").(int),
		OriginType:     d.Get("origin_type").(int),
	}

	// Build listen ports
	if listenPortsRaw, ok := d.GetOk("listen_ports"); ok {
		listenPortsList := listenPortsRaw.([]interface{})
		listenPorts := make([]int, len(listenPortsList))
		for i, port := range listenPortsList {
			listenPorts[i] = port.(int)
		}
		origin.ListenPorts = listenPorts
	}

	// Build records
	if recordsRaw, ok := d.GetOk("records"); ok {
		recordsList := recordsRaw.([]interface{})
		records := make([]scdn.OriginRecord, len(recordsList))
		for i, recordRaw := range recordsList {
			recordMap := recordRaw.(map[string]interface{})
			records[i] = scdn.OriginRecord{
				View:     recordMap["view"].(string),
				Value:    recordMap["value"].(string),
				Port:     recordMap["port"].(int),
				Priority: recordMap["priority"].(int),
			}
		}
		origin.Records = records
	}

	req := scdn.OriginAddRequest{
		DomainID: domainID,
		Origins:  []scdn.Origin{origin},
	}

	log.Printf("[INFO] Creating SCDN origin: %+v", req)
	response, err := service.AddOrigins(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN origin: %w", err)
	}

	if len(response.Data.IDs) == 0 {
		return fmt.Errorf("no origin ID returned from API")
	}

	// Set the origin ID as the resource ID
	d.SetId(strconv.Itoa(response.Data.IDs[0]))

	log.Printf("[INFO] SCDN origin created successfully: %s", d.Id())
	return resourceScdnOriginRead(d, m)
}

func resourceScdnOriginRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	originID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid origin ID: %w", err)
	}

	domainID := d.Get("domain_id").(int)

	// Get origin information by listing origins for the domain
	req := scdn.OriginListRequest{
		DomainID: domainID,
	}

	response, err := service.ListOrigins(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN origin: %w", err)
	}

	var originInfo *scdn.OriginInfo
	for _, origin := range response.Data.List {
		if origin.ID == originID {
			originInfo = &origin
			break
		}
	}

	if originInfo == nil {
		log.Printf("[WARN] SCDN origin not found: %s", d.Id())
		d.SetId("")
		return nil
	}

	// Set basic fields
	// Note: ID is stored as string in Terraform but used as int in API
	if err := d.Set("id", strconv.Itoa(originInfo.ID)); err != nil {
		return fmt.Errorf("error setting id: %w", err)
	}
	if err := d.Set("domain_id", originInfo.DomainID); err != nil {
		return fmt.Errorf("error setting domain_id: %w", err)
	}
	if err := d.Set("protocol", originInfo.Protocol); err != nil {
		return fmt.Errorf("error setting protocol: %w", err)
	}
	if err := d.Set("listen_port", originInfo.ListenPort); err != nil {
		return fmt.Errorf("error setting listen_port: %w", err)
	}
	if err := d.Set("origin_protocol", originInfo.OriginProtocol); err != nil {
		return fmt.Errorf("error setting origin_protocol: %w", err)
	}
	if err := d.Set("load_balance", originInfo.LoadBalance); err != nil {
		return fmt.Errorf("error setting load_balance: %w", err)
	}
	if err := d.Set("origin_type", originInfo.OriginType); err != nil {
		return fmt.Errorf("error setting origin_type: %w", err)
	}

	// Set listen ports (convert single port to list)
	if err := d.Set("listen_ports", []int{originInfo.ListenPort}); err != nil {
		return fmt.Errorf("error setting listen_ports: %w", err)
	}

	// Set records
	records := make([]map[string]interface{}, len(originInfo.Records))
	for i, record := range originInfo.Records {
		records[i] = map[string]interface{}{
			"view":     record.View,
			"value":    record.Value,
			"port":     record.Port,
			"priority": record.Priority,
		}
	}
	if err := d.Set("records", records); err != nil {
		return fmt.Errorf("error setting records: %w", err)
	}

	log.Printf("[INFO] SCDN origin read successfully: %s", d.Id())
	return nil
}

func resourceScdnOriginUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid origin ID: %w", err)
	}

	domainID := d.Get("domain_id").(int)

	// Build updated origin
	origin := scdn.EditOrigin{
		Id:             d.Get("id").(int),
		Protocol:       d.Get("protocol").(int),
		ListenPort:     d.Get("listen_port").(int),
		OriginProtocol: d.Get("origin_protocol").(int),
		LoadBalance:    d.Get("load_balance").(int),
		OriginType:     d.Get("origin_type").(int),
	}

	// Build records
	if recordsRaw, ok := d.GetOk("records"); ok {
		recordsList := recordsRaw.([]interface{})
		records := make([]scdn.OriginRecord, len(recordsList))
		for i, recordRaw := range recordsList {
			recordMap := recordRaw.(map[string]interface{})
			records[i] = scdn.OriginRecord{
				View:     recordMap["view"].(string),
				Value:    recordMap["value"].(string),
				Port:     recordMap["port"].(int),
				Priority: recordMap["priority"].(int),
			}
		}
		origin.Records = records
	}

	req := scdn.OriginUpdateRequest{
		DomainID: domainID,
		Origins:  []scdn.EditOrigin{origin},
	}

	log.Printf("[INFO] Updating SCDN origin: %+v", req)
	_, err = service.UpdateOrigins(req)
	if err != nil {
		return fmt.Errorf("failed to update SCDN origin: %w", err)
	}

	log.Printf("[INFO] SCDN origin updated successfully: %s", d.Id())
	return resourceScdnOriginRead(d, m)
}

func resourceScdnOriginDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	originID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid origin ID: %w", err)
	}

	domainID := d.Get("domain_id").(int)

	req := scdn.OriginDeleteRequest{
		IDs:      []int{originID},
		DomainID: domainID,
	}

	log.Printf("[INFO] Deleting SCDN origin: %+v", req)
	_, err = service.DeleteOrigins(req)
	if err != nil {
		return fmt.Errorf("failed to delete SCDN origin: %w", err)
	}

	d.SetId("")
	log.Printf("[INFO] SCDN origin deleted successfully: %d", originID)
	return nil
}
