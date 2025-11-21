package domain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnDomain returns the SCDN domain data source
func DataSourceEdgenextScdnDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnDomainRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The domain name to query (either domain, id, or domain_id must be provided)",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the domain to query (either domain, id, or domain_id must be provided). Also returned as the computed ID.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the domain to query (deprecated, use id instead)",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the domain group",
			},
			"exclusive_resource_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the exclusive resource package",
			},
			"remark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The remark for the domain",
			},
			"tpl_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The template ID applied to the domain",
			},
			"protect_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The edge node type",
			},
			"tpl_recommend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The template recommendation status",
			},
			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The application type",
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
			"origins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The origin server configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the origin",
						},
						"protocol": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The origin protocol",
						},
						"listen_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The listening port of the origin server",
						},
						"origin_protocol": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The origin protocol",
						},
						"load_balance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The load balancing method",
						},
						"origin_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The origin type",
						},
						"records": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The origin records",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"view": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The view of the record",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the record",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The port of the record",
									},
									"priority": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The priority of the record",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Get domain, id, and domain_id parameters
	domain, domainOk := d.GetOk("domain")
	idStr, idOk := d.GetOk("id")
	domainIDStr, domainIDOk := d.GetOk("domain_id")

	// Get the ID value (prefer id over domain_id for backward compatibility)
	var finalIDStr string
	if idOk && idStr.(string) != "" {
		finalIDStr = idStr.(string)
	} else if domainIDOk && domainIDStr.(string) != "" {
		finalIDStr = domainIDStr.(string)
	}

	// Validate that at least one of domain, id, or domain_id is provided
	if !domainOk && finalIDStr == "" {
		return fmt.Errorf("either domain, id, or domain_id must be provided")
	}

	// Build request
	req := scdn.DomainListRequest{
		Page:     1,
		PageSize: 100,
	}

	// Priority: If ID is provided, use it (more precise)
	// Otherwise, use domain name
	if finalIDStr != "" {
		domainID, err := strconv.Atoi(finalIDStr)
		if err != nil {
			return fmt.Errorf("invalid domain ID: %w", err)
		}
		req.ID = domainID
		log.Printf("[INFO] Querying SCDN domain by ID: %d", domainID)
	} else if domainOk && domain.(string) != "" {
		// Use domain name for query
		req.Domain = domain.(string)
		log.Printf("[INFO] Querying SCDN domain by name: %s", domain.(string))
	} else {
		return fmt.Errorf("either domain, id, or domain_id must be provided with a non-empty value")
	}

	response, err := service.ListDomains(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN domain: %w", err)
	}

	var domainInfo *scdn.DomainInfo
	if len(response.Data.List) > 0 {
		// If querying by ID, should return exactly one result
		// If querying by name, find matching domain
		if req.ID > 0 {
			// Query by ID - should return the domain with matching ID
			for _, domainItem := range response.Data.List {
				if domainItem.ID == req.ID {
					domainInfo = &domainItem
					break
				}
			}
		} else {
			// Query by name - find exact match
			domainName := req.Domain
			for _, domainItem := range response.Data.List {
				if domainItem.Domain == domainName {
					domainInfo = &domainItem
					break
				}
			}
		}
	}

	if domainInfo == nil {
		if req.ID > 0 {
			return fmt.Errorf("SCDN domain not found by ID: %d", req.ID)
		}
		return fmt.Errorf("SCDN domain not found: %s", req.Domain)
	}

	// Set basic fields
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
				"id":              origin.ID,
				"protocol":        origin.Protocol,
				"listen_port":     origin.ListenPort,
				"origin_protocol": origin.OriginProtocol,
				"load_balance":    origin.LoadBalance,
				"origin_type":     origin.OriginType,
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

	// Set the domain ID as the resource ID
	d.SetId(strconv.Itoa(domainInfo.ID))

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"id":                    domainInfo.ID,
			"domain":                domainInfo.Domain,
			"remark":                domainInfo.Remark,
			"protect_status":        domainInfo.ProtectStatus,
			"access_progress":       domainInfo.AccessProgress,
			"access_mode":           domainInfo.AccessMode,
			"ei_forward_status":     domainInfo.EIForwardStatus,
			"use_my_cname":          domainInfo.UseMyCname,
			"use_my_dns":            domainInfo.UseMyDNS,
			"ca_status":             domainInfo.CAStatus,
			"exclusive_resource_id": domainInfo.ExclusiveResourceID,
			"access_progress_desc":  domainInfo.AccessProgressDesc,
			"has_origin":            domainInfo.HasOrigin,
			"ca_id":                 domainInfo.CAID,
			"created_at":            domainInfo.CreatedAt,
			"updated_at":            domainInfo.UpdatedAt,
			"pri_domain":            domainInfo.PriDomain,
			"cname":                 cnameInfo,
			"origins":               d.Get("origins"),
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN domain queried successfully: %s", domain)
	return nil
}

// DataSourceEdgenextScdnDomains returns the SCDN domains data source
func DataSourceEdgenextScdnDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnDomainsRead,

		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The page number for pagination",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "The page size for pagination",
			},
			"access_progress": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by access progress status",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter by domain group ID",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by domain name (fuzzy search)",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by remark (fuzzy search)",
			},
			"origin_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by origin IP",
			},
			"ca_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by certificate binding status",
			},
			"access_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by access mode",
			},
			"protect_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by edge node type",
			},
			"exclusive_resource_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Filter by exclusive resource package ID",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of domains",
			},
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of domains",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the domain",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain name",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The remark for the domain",
						},
						"access_progress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access progress status",
						},
						"access_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access mode",
						},
						"protect_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The edge node type",
						},
						"ei_forward_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The explicit/implicit forwarding status",
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
						"exclusive_resource_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The exclusive resource package ID",
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
					},
				},
			},
		},
	}
}

func dataSourceScdnDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request
	req := scdn.DomainListRequest{
		Page:                d.Get("page").(int),
		PageSize:            d.Get("page_size").(int),
		AccessProgress:      d.Get("access_progress").(string),
		GroupID:             d.Get("group_id").(int),
		Domain:              d.Get("domain").(string),
		Remark:              d.Get("remark").(string),
		OriginIP:            d.Get("origin_ip").(string),
		CAStatus:            d.Get("ca_status").(string),
		AccessMode:          d.Get("access_mode").(string),
		ProtectStatus:       d.Get("protect_status").(string),
		ExclusiveResourceID: d.Get("exclusive_resource_id").(int),
	}

	log.Printf("[INFO] Querying SCDN domains with filters: %+v", req)
	response, err := service.ListDomains(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN domains: %w", err)
	}

	// Convert domains to the format expected by Terraform
	domains := make([]map[string]interface{}, len(response.Data.List))
	ids := make([]string, len(response.Data.List))
	for i, domain := range response.Data.List {
		domainMap := map[string]interface{}{
			"id":                    domain.ID,
			"domain":                domain.Domain,
			"remark":                domain.Remark,
			"access_progress":       domain.AccessProgress,
			"access_mode":           domain.AccessMode,
			"protect_status":        domain.ProtectStatus,
			"ei_forward_status":     domain.EIForwardStatus,
			"use_my_cname":          domain.UseMyCname,
			"use_my_dns":            domain.UseMyDNS,
			"ca_status":             domain.CAStatus,
			"exclusive_resource_id": domain.ExclusiveResourceID,
			"access_progress_desc":  domain.AccessProgressDesc,
			"has_origin":            domain.HasOrigin,
			"ca_id":                 domain.CAID,
			"created_at":            domain.CreatedAt,
			"updated_at":            domain.UpdatedAt,
			"pri_domain":            domain.PriDomain,
		}

		// Set CNAME information
		cnameInfo := map[string]interface{}{
			"master": domain.Cname.Master,
			"slaves": domain.Cname.Slaves,
		}
		domainMap["cname"] = []map[string]interface{}{cnameInfo}

		domains[i] = domainMap
		ids[i] = fmt.Sprintf("%d", domain.ID)
	}

	// Set the resource ID
	d.SetId(helper.DataResourceIdsHash(ids))

	// Set the domains list
	if err := d.Set("domains", domains); err != nil {
		return fmt.Errorf("error setting domains: %w", err)
	}

	// Set the total count
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":   response.Data.Total,
			"domains": domains,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN domains queried successfully, %d domains found", len(response.Data.List))
	return nil
}
