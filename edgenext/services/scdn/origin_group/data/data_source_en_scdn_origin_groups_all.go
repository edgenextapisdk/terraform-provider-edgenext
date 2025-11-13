package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnOriginGroupsAll returns the SCDN all origin groups data source
func DataSourceEdgenextScdnOriginGroupsAll() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnOriginGroupsAllRead,

		Schema: map[string]*schema.Schema{
			"protect_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Protection status: scdn-shared nodes, exclusive-dedicated nodes",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of origin groups",
			},
			"origin_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Origin group list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Origin group ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin group name",
						},
						"origins": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Origin list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Origin type: 0-IP, 1-domain",
									},
									"records": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Origin record list",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Origin address",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Origin port",
												},
												"priority": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Weight",
												},
												"view": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Origin type: primary-backup, backup-backup",
												},
												"host": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Origin Host",
												},
											},
										},
									},
									"protocol_ports": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Protocol port mapping",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"protocol": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Protocol: 0-http, 1-https",
												},
												"listen_ports": {
													Type:        schema.TypeList,
													Computed:    true,
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
										Computed:    true,
										Description: "Origin protocol: 0-http, 1-https, 2-follow",
									},
									"load_balance": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Load balance strategy: 0-ip_hash, 1-round_robin, 2-cookie",
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

func dataSourceScdnOriginGroupsAllRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.OriginGroupAllRequest{
		ProtectStatus: d.Get("protect_status").(string),
	}

	log.Printf("[INFO] Querying SCDN all origin groups: protect_status=%s", req.ProtectStatus)
	response, err := service.GetAllOriginGroups(req)
	if err != nil {
		return fmt.Errorf("failed to query all origin groups: %w", err)
	}

	// Set resource ID
	d.SetId("origin-groups-all")

	// Set total
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Convert origin groups to schema format
	originGroupList := make([]map[string]interface{}, 0, len(response.Data.List))
	for _, og := range response.Data.List {
		ogMap := map[string]interface{}{
			"id":   og.ID,
			"name": og.Name,
		}

		// Set origins
		originsList := make([]map[string]interface{}, len(og.Origins))
		for i, origin := range og.Origins {
			originMap := map[string]interface{}{
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
		ogMap["origins"] = originsList

		originGroupList = append(originGroupList, ogMap)
	}

	if err := d.Set("origin_groups", originGroupList); err != nil {
		return fmt.Errorf("error setting origin_groups: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":         response.Data.Total,
			"origin_groups": originGroupList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN all origin groups queried successfully: total=%d", response.Data.Total)
	return nil
}
