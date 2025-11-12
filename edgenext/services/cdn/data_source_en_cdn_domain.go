package cdn

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextCdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDomainConfigRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain to query configuration for.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"config_item": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration items.",
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration area: mainland_china (mainland China), outside_mainland_china (outside mainland China), global (global).",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain type: page (web page), download (download), video_demand (video on demand), dynamic (dynamic).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain status.",
			},
			"icp_num": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ICP filing number.",
			},
			"icp_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ICP filing status.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CNAME.",
			},
			"https": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "HTTPS.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_master": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"default_slave": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"origin_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ori_https": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"origin_host": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cache_rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"pattern": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timeunit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ignore_no_cache": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ignore_expired": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ignore_query": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cache_rule_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_method": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pattern": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"case_ignore": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expire": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"expire_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ignore_no_cache_headers": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"follow_expired": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"query_params_op": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"cache_or_not": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"query_params_op_way": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"query_params_op_when": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"params": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"referer": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"allow_empty": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"ip_black_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"ip_white_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"add_response_head": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"cover": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"only_hit": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"add_back_source_head": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"head_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"head_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"write_when_exists": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"https": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"http2": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"force_https": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ocsp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_protocol": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"compress_response": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_type": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_size_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"speed_limit": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pattern": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"speed": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"begin_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"rate_limit": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_rate_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"leading_flow_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"leading_flow_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_rate_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cache_share": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"share_way": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"head_control": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"regex": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"head_op": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"head_direction": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"head": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"fun_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												// "order": {
												// 	Type:     schema.TypeInt,
												// 	Computed: true,
												// },
											},
										},
									},
								},
							},
						},
						"timeout": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"connect_timeout": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_connect_timeout": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"deny_url": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"urls": {
										Type:     schema.TypeList,
										Computed: true,
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

func dataSourceDomainConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	var configItem []string
	if configItemRaw := d.Get("config_item"); configItemRaw != nil {
		if configItemList, ok := configItemRaw.([]interface{}); ok && len(configItemList) > 0 {
			for _, item := range configItemList {
				if str, ok := item.(string); ok {
					configItem = append(configItem, str)
				}
			}
		}
	}

	// 1. Get domain information
	domain := d.Get("domain").(string)
	log.Printf("[INFO] Data source reading CDN domain: %s", domain)
	err := readDomain(d, service, domain)
	if err != nil {
		return fmt.Errorf("data source failed to read CDN domain: %w", err)
	}

	// 2. Get domain configuration
	log.Printf("[INFO] Data source reading domain configuration: %s", domain)
	err = readDomainConfig(d, service, domain, configItem)
	if err != nil {
		return fmt.Errorf("data source failed to read domain configuration: %w", err)
	}
	// Set data source ID
	d.SetId(domain)

	// Write result to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"domain":      d.Get("domain"),
			"area":        d.Get("area"),
			"type":        d.Get("type"),
			"status":      d.Get("status"),
			"icp_num":     d.Get("icp_num"),
			"icp_status":  d.Get("icp_status"),
			"cname":       d.Get("cname"),
			"https":       d.Get("https"),
			"create_time": d.Get("create_time"),
			"update_time": d.Get("update_time"),
			"config":      d.Get("config"),
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] Data source read domain and configuration successfully: %s", domain)
	return nil
}

// DataSourceEdgenextCdnDomains query multiple domain names
func DataSourceEdgenextCdnDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDomainsRead,

		Schema: map[string]*schema.Schema{
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				Description: "Get the page number. \n" +
					"Default value is 1 when not specified.",
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
				Description: "Page size, value range: 1-500. \n" +
					"Default value is 100 when not specified.",
			},
			"domain_status": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Specify the service status of the domain, support specifying multiple service status queries: \n" +
					"serving：Serving. When querying with serving, the domain whose \"status\" is \"deploying\" is in the configuration deployment state.\n" +
					"suspend：Suspended\n" +
					"deleted：Deleted\n" +
					"Default value is all status domain names when not specified.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Domain list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CDN domain name to query",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain type: page(web page), download(download), video_demand(video demand), dynamic(dynamic)",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Acceleration area: mainland_china, outside_mainland_china, global",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain status",
						},
						"icp_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ICP filing number",
						},
						"icp_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ICP filing status",
						},
						"cname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CNAME",
						},
						"https": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTPS",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time",
						},
					},
				},
			},
		},
	}
}

func dataSourceDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := NewCdnService(client)

	log.Printf("[INFO] Querying CDN domain list")

	response, err := service.ListDomains(DomainListRequest{
		PageNumber:   d.Get("page_number").(int),
		PageSize:     d.Get("page_size").(int),
		DomainStatus: d.Get("domain_status").(string),
	})
	if err != nil {
		return fmt.Errorf("failed to query CDN domain list: %w", err)
	}

	var list []map[string]interface{}
	ids := make([]string, 0)
	for _, elem := range response.Data.List {
		elemMap := map[string]interface{}{
			"id":          elem.ID,
			"domain":      elem.Domain,
			"type":        elem.Type,
			"area":        elem.Area,
			"status":      elem.Status,
			"cname":       elem.Cname,
			"https":       elem.Https,
			"create_time": elem.CreateTime,
			"update_time": elem.UpdateTime,
			"icp_num":     elem.IcpNum,
			"icp_status":  elem.IcpStatus,
		}
		list = append(list, elemMap)
		ids = append(ids, elem.ID)
	}
	// Set resource ID
	d.SetId(helper.DataResourceIdsHash(ids))
	// Set domain list
	err = d.Set("list", list)
	if err != nil {
		log.Printf("[ERROR] Failed to set domain list: %v", err)
		return err
	}

	// Write result to output file if specified
	if outputFile := d.Get("output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"page_number":   d.Get("page_number"),
			"page_size":     d.Get("page_size"),
			"domain_status": d.Get("domain_status"),
			"list":          list,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] CDN domain list queried successfully, %d domains", len(response.Data.List))
	return nil
}
