package domain

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnOrigin returns the SCDN origin data source
func DataSourceEdgenextScdnOrigin() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnOriginRead,

		Schema: map[string]*schema.Schema{
			"origin_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the origin to query",
			},
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the domain that owns the origin",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the origin",
			},
			"protocol": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS)",
			},
			"listen_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The listening port of the origin server",
			},
			"origin_protocol": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS), 2 (Follow)",
			},
			"load_balance": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The load balancing method. Valid values: 0 (IP hash), 1 (Round robin), 2 (Cookie)",
			},
			"origin_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The origin type. Valid values: 0 (IP), 1 (Domain)",
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
							Description: "The view of the record. Valid values: primary (primary line), backup (backup line)",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the record (IP address or domain)",
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
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The origin host, specifies the Host header when accessing the origin",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnOriginRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	originID := d.Get("origin_id").(int)
	domainID := d.Get("domain_id").(int)

	// Get origin information by listing origins for the domain
	req := scdn.OriginListRequest{
		DomainID: domainID,
	}

	log.Printf("[INFO] Querying SCDN origin: %d for domain %d", originID, domainID)
	response, err := service.ListOrigins(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN origin: %w", err)
	}

	var originInfo *scdn.OriginInfo
	for _, origin := range response.Data.List {
		if origin.ID == originID {
			originInfo = &origin
			break
		}
	}

	if originInfo == nil {
		return fmt.Errorf("SCDN origin not found: %d", originID)
	}

	// Set basic fields
	if err := d.Set("id", originInfo.ID); err != nil {
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

	// Set records
	records := make([]map[string]interface{}, len(originInfo.Records))
	for i, record := range originInfo.Records {
		records[i] = map[string]interface{}{
			"view":     record.View,
			"value":    record.Value,
			"port":     record.Port,
			"priority": record.Priority,
			"host":     record.Host,
		}
	}
	if err := d.Set("records", records); err != nil {
		return fmt.Errorf("error setting records: %w", err)
	}

	// Set the origin ID as the resource ID
	d.SetId(fmt.Sprintf("%d", originInfo.ID))

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"id":              originInfo.ID,
			"domain_id":       originInfo.DomainID,
			"protocol":        originInfo.Protocol,
			"listen_port":     originInfo.ListenPort,
			"origin_protocol": originInfo.OriginProtocol,
			"load_balance":    originInfo.LoadBalance,
			"origin_type":     originInfo.OriginType,
			"records":         records,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN origin queried successfully: %d", originInfo.ID)
	return nil
}

// DataSourceEdgenextScdnOrigins returns the SCDN origins data source
func DataSourceEdgenextScdnOrigins() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnOriginsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the domain to query origins for",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of origins",
			},
			"origins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of origins",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the origin",
						},
						"domain_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the domain",
						},
						"protocol": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS)",
						},
						"listen_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The listening port of the origin server",
						},
						"origin_protocol": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS), 2 (Follow)",
						},
						"load_balance": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The load balancing method. Valid values: 0 (IP hash), 1 (Round robin), 2 (Cookie)",
						},
						"origin_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The origin type. Valid values: 0 (IP), 1 (Domain)",
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
										Description: "The view of the record. Valid values: primary (primary line), backup (backup line)",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the record (IP address or domain)",
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
									"host": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The origin host, specifies the Host header when accessing the origin",
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

func dataSourceScdnOriginsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	// Get origins for the domain
	req := scdn.OriginListRequest{
		DomainID: domainID,
	}

	log.Printf("[INFO] Querying SCDN origins for domain: %d", domainID)
	response, err := service.ListOrigins(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN origins: %w", err)
	}

	// Convert origins to the format expected by Terraform
	origins := make([]map[string]interface{}, len(response.Data.List))
	ids := make([]string, len(response.Data.List))
	for i, origin := range response.Data.List {
		originMap := map[string]interface{}{
			"id":              origin.ID,
			"domain_id":       origin.DomainID,
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
				"host":     record.Host,
			}
		}
		originMap["records"] = records

		origins[i] = originMap
		ids[i] = fmt.Sprintf("%d", origin.ID)
	}

	// Set the resource ID
	d.SetId(helper.DataResourceIdsHash(ids))

	// Set the origins list
	if err := d.Set("origins", origins); err != nil {
		return fmt.Errorf("error setting origins: %w", err)
	}

	// Set the total count
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":   response.Data.Total,
			"origins": origins,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN origins queried successfully, %d origins found", len(response.Data.List))
	return nil
}
