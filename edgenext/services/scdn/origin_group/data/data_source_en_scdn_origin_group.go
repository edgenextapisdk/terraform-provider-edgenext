package data

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnOriginGroup returns the SCDN origin group data source
func DataSourceEdgenextScdnOriginGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnOriginGroupRead,

		Schema: map[string]*schema.Schema{
			"origin_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Origin group ID",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Origin group name",
			},
			"remark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Remark",
			},
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
			"origins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Origin list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Origin ID",
						},
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

func dataSourceScdnOriginGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	originGroupID := d.Get("origin_group_id").(int)

	req := scdn.OriginGroupDetailRequest{
		ID: originGroupID,
	}

	log.Printf("[INFO] Querying SCDN origin group: origin_group_id=%d", originGroupID)
	response, err := service.GetOriginGroupDetail(req)
	if err != nil {
		return fmt.Errorf("failed to query origin group: %w", err)
	}

	originGroup := response.Data.OriginGroup
	d.SetId(strconv.Itoa(originGroupID))

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

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"origin_group_id": originGroupID,
			"name":            originGroup.Name,
			"remark":          originGroup.Remark,
			"member_id":       originGroup.MemberID,
			"username":        originGroup.Username,
			"origins":         originsList,
			"created_at":      originGroup.CreatedAt,
			"updated_at":      originGroup.UpdatedAt,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN origin group queried successfully: origin_group_id=%d", originGroupID)
	return nil
}
