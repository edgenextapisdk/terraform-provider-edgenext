package domain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnDomain returns the SCDN domain resource
func ResourceEdgenextScdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnDomainCreate,
		Read:   resourceScdnDomainRead,
		Update: resourceScdnDomainUpdate,
		Delete: resourceScdnDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The domain name to be added to SCDN",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the domain group",
			},
			"exclusive_resource_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the exclusive resource package",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The remark for the domain",
			},
			"tpl_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The template ID to be applied to the domain",
			},
			"protect_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "scdn",
				Description: "The edge node type. Valid values: back_source, scdn, exclusive",
			},
			"tpl_recommend": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The template recommendation status",
			},
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application type",
			},
			"origins": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The origin server configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the domain",
			},
			"access_progress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access progress status",
			},
			"access_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access mode (ns or cname)",
			},
			"ei_forward_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The explicit/implicit forwarding status",
			},
			"cname": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The CNAME information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The master CNAME record",
						},
						"slaves": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The slave CNAME records",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"use_my_cname": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The CNAME resolution status",
			},
			"use_my_dns": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The DNS hosting status",
			},
			"ca_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The certificate binding status",
			},
			"access_progress_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the access progress status",
			},
			"has_origin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the domain has origin configuration",
			},
			"ca_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The certificate ID",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update timestamp",
			},
			"pri_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The primary domain",
			},
		},
	}
}

func resourceScdnDomainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build create request
	req := scdn.DomainCreateRequest{
		Domain:              d.Get("domain").(string),
		GroupID:             d.Get("group_id").(int),
		ExclusiveResourceID: d.Get("exclusive_resource_id").(int),
		Remark:              d.Get("remark").(string),
		TplID:               d.Get("tpl_id").(int),
		ProtectStatus:       d.Get("protect_status").(string),
		TplRecommend:        d.Get("tpl_recommend").(string),
		AppType:             d.Get("app_type").(string),
	}

	// Build origins
	if originsRaw, ok := d.GetOk("origins"); ok {
		originsList := originsRaw.([]interface{})
		origins := make([]scdn.Origin, len(originsList))
		for i, originRaw := range originsList {
			originMap := originRaw.(map[string]interface{})
			origin := scdn.Origin{
				Protocol:       originMap["protocol"].(int),
				OriginProtocol: originMap["origin_protocol"].(int),
				LoadBalance:    originMap["load_balance"].(int),
				OriginType:     originMap["origin_type"].(int),
			}

			// Build listen ports
			if listenPortsRaw, ok := originMap["listen_ports"]; ok {
				listenPortsList := listenPortsRaw.([]interface{})
				listenPorts := make([]int, len(listenPortsList))
				for j, port := range listenPortsList {
					listenPorts[j] = port.(int)
				}
				origin.ListenPorts = listenPorts
			}

			// Build records
			if recordsRaw, ok := originMap["records"]; ok {
				recordsList := recordsRaw.([]interface{})
				records := make([]scdn.OriginRecord, len(recordsList))
				for j, recordRaw := range recordsList {
					recordMap := recordRaw.(map[string]interface{})
					records[j] = scdn.OriginRecord{
						View:     recordMap["view"].(string),
						Value:    recordMap["value"].(string),
						Port:     recordMap["port"].(int),
						Priority: recordMap["priority"].(int),
					}
				}
				origin.Records = records
			}

			origins[i] = origin
		}
		req.Origins = origins
	}

	log.Printf("[INFO] Creating SCDN domain: %+v", req)
	response, err := service.CreateDomain(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN domain: %w", err)
	}

	log.Printf("[DEBUG] Domain creation response: %+v", response)
	log.Printf("[DEBUG] Domain ID from response: %d", response.Data.ID)

	// Set the domain ID as the resource ID
	d.SetId(strconv.Itoa(response.Data.ID))

	// Set basic fields directly from creation response to avoid read issues
	if err := d.Set("id", strconv.Itoa(response.Data.ID)); err != nil {
		log.Printf("[WARN] Failed to set domain id: %v", err)
	}
	if err := d.Set("domain", response.Data.Domain); err != nil {
		log.Printf("[WARN] Failed to set domain: %v", err)
	}
	if err := d.Set("remark", response.Data.Record); err != nil {
		log.Printf("[WARN] Failed to set remark: %v", err)
	}
	log.Printf("[INFO] SCDN domain created successfully: %s", d.Id())

	// Still call read to get full details
	return resourceScdnDomainRead(d, m)
}

func resourceScdnDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainName := d.Get("domain").(string)

	// Build request - support both domain name and ID queries
	req := scdn.DomainListRequest{
		Page:     1,
		PageSize: 100,
	}

	// If domain name is not set (e.g., during import), try to use ID
	if domainName == "" {
		// Try to get domain ID from resource ID
		domainID, err := strconv.Atoi(d.Id())
		if err == nil && domainID > 0 {
			log.Printf("[DEBUG] Domain name not set, using ID to query: %d", domainID)
			req.ID = domainID
		} else {
			log.Printf("[DEBUG] Domain name not set and invalid ID, skipping read operation")
			return nil
		}
	} else {
		// Use domain name to query
		req.Domain = domainName
	}

	if domainName != "" {
		log.Printf("[DEBUG] Reading SCDN domain by name: %s", domainName)
	} else {
		log.Printf("[DEBUG] Reading SCDN domain by ID: %d", req.ID)
	}

	response, err := service.ListDomains(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN domain: %w", err)
	}

	log.Printf("[DEBUG] Domain list response: %+v", response)
	log.Printf("[DEBUG] Found %d domains", len(response.Data.List))

	var domainInfo *scdn.DomainInfo
	if len(response.Data.List) > 0 {
		// If querying by ID, should return exactly one result
		// If querying by name, find matching domain
		if req.ID > 0 {
			// Query by ID - should return the domain with matching ID
			for _, domain := range response.Data.List {
				if domain.ID == req.ID {
					domainInfo = &domain
					break
				}
			}
		} else {
			// Query by name - find exact match
			for _, domain := range response.Data.List {
				log.Printf("[DEBUG] Checking domain: %s (ID: %d)", domain.Domain, domain.ID)
				if domain.Domain == domainName {
					domainInfo = &domain
					break
				}
			}
		}
	}

	if domainInfo == nil {
		if domainName != "" {
			log.Printf("[WARN] SCDN domain not found: %s, keeping current ID: %s", domainName, d.Id())
		} else {
			log.Printf("[WARN] SCDN domain not found by ID: %d, keeping current ID: %s", req.ID, d.Id())
		}
		// Don't clear the ID, keep the current one
		return nil
	}

	// Set basic fields
	// Note: ID is stored as string in Terraform but used as int in API
	if err := d.Set("id", strconv.Itoa(domainInfo.ID)); err != nil {
		return fmt.Errorf("error setting id: %w", err)
	}
	if err := d.Set("domain", domainInfo.Domain); err != nil {
		return fmt.Errorf("error setting domain: %w", err)
	}
	if err := d.Set("remark", domainInfo.Remark); err != nil {
		return fmt.Errorf("error setting remark: %w", err)
	}
	if err := d.Set("protect_status", domainInfo.ProtectStatus); err != nil {
		return fmt.Errorf("error setting protect_status: %w", err)
	}
	if err := d.Set("access_progress", domainInfo.AccessProgress); err != nil {
		return fmt.Errorf("error setting access_progress: %w", err)
	}
	if err := d.Set("access_mode", domainInfo.AccessMode); err != nil {
		return fmt.Errorf("error setting access_mode: %w", err)
	}
	if err := d.Set("ei_forward_status", domainInfo.EIForwardStatus); err != nil {
		return fmt.Errorf("error setting ei_forward_status: %w", err)
	}
	if err := d.Set("use_my_cname", domainInfo.UseMyCname); err != nil {
		return fmt.Errorf("error setting use_my_cname: %w", err)
	}
	if err := d.Set("use_my_dns", domainInfo.UseMyDNS); err != nil {
		return fmt.Errorf("error setting use_my_dns: %w", err)
	}
	if err := d.Set("ca_status", domainInfo.CAStatus); err != nil {
		return fmt.Errorf("error setting ca_status: %w", err)
	}
	if err := d.Set("exclusive_resource_id", domainInfo.ExclusiveResourceID); err != nil {
		return fmt.Errorf("error setting exclusive_resource_id: %w", err)
	}
	if err := d.Set("access_progress_desc", domainInfo.AccessProgressDesc); err != nil {
		return fmt.Errorf("error setting access_progress_desc: %w", err)
	}
	if err := d.Set("has_origin", domainInfo.HasOrigin); err != nil {
		return fmt.Errorf("error setting has_origin: %w", err)
	}
	if err := d.Set("ca_id", domainInfo.CAID); err != nil {
		return fmt.Errorf("error setting ca_id: %w", err)
	}
	if err := d.Set("created_at", domainInfo.CreatedAt); err != nil {
		return fmt.Errorf("error setting created_at: %w", err)
	}
	if err := d.Set("updated_at", domainInfo.UpdatedAt); err != nil {
		return fmt.Errorf("error setting updated_at: %w", err)
	}
	if err := d.Set("pri_domain", domainInfo.PriDomain); err != nil {
		return fmt.Errorf("error setting pri_domain: %w", err)
	}

	// Set CNAME information
	cnameInfo := map[string]interface{}{
		"master": domainInfo.Cname.Master,
		"slaves": domainInfo.Cname.Slaves,
	}
	if err := d.Set("cname", []map[string]interface{}{cnameInfo}); err != nil {
		return fmt.Errorf("error setting cname: %w", err)
	}

	// Get origins for this domain
	originReq := scdn.OriginListRequest{
		DomainID: domainInfo.ID,
	}
	originResponse, err := service.ListOrigins(originReq)
	if err != nil {
		log.Printf("[WARN] Failed to get origins for domain %d: %v", domainInfo.ID, err)
	} else {
		// Convert origins to the format expected by Terraform
		origins := make([]map[string]interface{}, len(originResponse.Data.List))
		for i, origin := range originResponse.Data.List {
			originMap := map[string]interface{}{
				"protocol":        origin.Protocol,
				"origin_protocol": origin.OriginProtocol,
				"load_balance":    origin.LoadBalance,
				"origin_type":     origin.OriginType,
				"listen_ports":    []int{origin.ListenPort},
			}

			// Convert records
			records := make([]map[string]interface{}, len(origin.Records))
			for j, record := range origin.Records {
				records[j] = map[string]interface{}{
					"view":     record.View,
					"value":    record.Value,
					"port":     record.Port,
					"priority": record.Priority,
				}
			}
			originMap["records"] = records

			origins[i] = originMap
		}
		if err := d.Set("origins", origins); err != nil {
			return fmt.Errorf("error setting origins: %w", err)
		}
	}

	log.Printf("[INFO] SCDN domain read successfully: %s", d.Id())
	return nil
}

func resourceScdnDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid domain ID: %w", err)
	}

	// Update domain basic information
	if d.HasChange("remark") {
		req := scdn.DomainUpdateRequest{
			DomainID: domainID,
			Remark:   d.Get("remark").(string),
		}

		log.Printf("[INFO] Updating SCDN domain: %+v", req)
		_, err := service.UpdateDomain(req)
		if err != nil {
			return fmt.Errorf("failed to update SCDN domain: %w", err)
		}
	}

	// Update origins if changed
	if d.HasChange("origins") {
		// First, get current origins to delete them
		originReq := scdn.OriginListRequest{
			DomainID: domainID,
		}
		originResponse, err := service.ListOrigins(originReq)
		if err == nil && len(originResponse.Data.List) > 0 {
			// Delete existing origins
			ids := make([]int, len(originResponse.Data.List))
			for i, origin := range originResponse.Data.List {
				ids[i] = origin.ID
			}
			deleteReq := scdn.OriginDeleteRequest{
				IDs:      ids,
				DomainID: domainID,
			}
			_, err := service.DeleteOrigins(deleteReq)
			if err != nil {
				log.Printf("[WARN] Failed to delete existing origins: %v", err)
			}
		}

		// Add new origins
		if originsRaw, ok := d.GetOk("origins"); ok {
			originsList := originsRaw.([]interface{})
			origins := make([]scdn.Origin, len(originsList))
			for i, originRaw := range originsList {
				originMap := originRaw.(map[string]interface{})
				origin := scdn.Origin{
					Protocol:       originMap["protocol"].(int),
					OriginProtocol: originMap["origin_protocol"].(int),
					LoadBalance:    originMap["load_balance"].(int),
					OriginType:     originMap["origin_type"].(int),
				}

				// Build listen ports
				if listenPortsRaw, ok := originMap["listen_ports"]; ok {
					listenPortsList := listenPortsRaw.([]interface{})
					listenPorts := make([]int, len(listenPortsList))
					for j, port := range listenPortsList {
						listenPorts[j] = port.(int)
					}
					origin.ListenPorts = listenPorts
				}

				// Build records
				if recordsRaw, ok := originMap["records"]; ok {
					recordsList := recordsRaw.([]interface{})
					records := make([]scdn.OriginRecord, len(recordsList))
					for j, recordRaw := range recordsList {
						recordMap := recordRaw.(map[string]interface{})
						records[j] = scdn.OriginRecord{
							View:     recordMap["view"].(string),
							Value:    recordMap["value"].(string),
							Port:     recordMap["port"].(int),
							Priority: recordMap["priority"].(int),
						}
					}
					origin.Records = records
				}

				origins[i] = origin
			}

			addReq := scdn.OriginAddRequest{
				DomainID: domainID,
				Origins:  origins,
			}
			_, err := service.AddOrigins(addReq)
			if err != nil {
				return fmt.Errorf("failed to add origins: %w", err)
			}
		}
	}

	// Update protect status if changed
	if d.HasChange("protect_status") {
		req := scdn.DomainNodeSwitchRequest{
			DomainID:      domainID,
			ProtectStatus: d.Get("protect_status").(string),
		}
		if exclusiveResourceID := d.Get("exclusive_resource_id").(int); exclusiveResourceID > 0 {
			req.ExclusiveResourceID = exclusiveResourceID
		}

		log.Printf("[INFO] Switching SCDN domain nodes: %+v", req)
		_, err := service.SwitchDomainNodes(req)
		if err != nil {
			return fmt.Errorf("failed to switch domain nodes: %w", err)
		}
	}

	log.Printf("[INFO] SCDN domain updated successfully: %s", d.Id())
	return resourceScdnDomainRead(d, m)
}

func resourceScdnDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid domain ID: %w", err)
	}

	req := scdn.DomainDeleteRequest{
		IDs: []int{domainID},
	}

	log.Printf("[INFO] Deleting SCDN domain: %+v", req)
	_, err = service.DeleteDomain(req)
	if err != nil {
		return fmt.Errorf("failed to delete SCDN domain: %w", err)
	}

	d.SetId("")
	log.Printf("[INFO] SCDN domain deleted successfully: %d", domainID)
	return nil
}
