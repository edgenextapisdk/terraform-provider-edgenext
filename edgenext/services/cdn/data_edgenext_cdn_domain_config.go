package cdn

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceEdgenextCdnDomainConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDomainConfigRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain to query configuration for",
			},
			"config_item": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration items",
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration area: mainland_china (mainland China), outside_mainland_china (outside mainland China), global (global)",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Domain type: page (web page), download (download), video_demand (video on demand), dynamic (dynamic)",
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
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Configuration",
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
										Type:     schema.TypeString,
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
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
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
									"mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"add_response_head": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
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
								},
							},
						},
						"visit_timestamp": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pattern": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time_format": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"deadtime": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"req_uri_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"origin_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"style": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"timename": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"keyname": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"forbid_http_x": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cache_error_code": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"bcache": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cache_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"cache_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"video_drag": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mp4": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flv": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end": {
										Type:     schema.TypeString,
										Computed: true,
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
										Type:     schema.TypeString,
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
								},
							},
						},
						"visit_deny_whitelist": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"turn_on": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"range_back_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extend": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"squid": {
										Type:     schema.TypeString,
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
						"new_origin": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"level": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"weight_level": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"isp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"connect_time": {
										Type:     schema.TypeInt,
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
												"order": {
													Type:     schema.TypeInt,
													Computed: true,
												},
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
						"qiniu_origin_auth": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"match_method": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pattern": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"forward_status": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"codes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
									},
								},
							},
						},
						"error_page_rewrite": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_status_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"redirect_status_code": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"redirect_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"post_upload_size_limit": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"source_url_rewrite": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"combined_ban": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"configs": {
										Type:     schema.TypeList,
										Computed: true,

										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"method": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"case_insensitive": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"method_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
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
												"ip": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"case_insensitive": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"method_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
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
												"referer": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"case_insensitive": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"method_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
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
												"url": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"case_insensitive": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"method_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
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
												"ua": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_match": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"case_insensitive": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"method_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
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
											},
										},
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
						"tos_origin": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"isp": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ips": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"group_sort": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"origin_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"host_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_secret_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_access_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_bucket_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"auth_expire": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"auth_cdn_tag": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parse_priority": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"client_real_ip": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"head": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"user_agent": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ua_list": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url_pattern": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url_case_insensitive": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"allow_empty_ua": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ua_case_insensitive": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url_match_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ua_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ua_match_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"visit_areas_limit": {
							Type:     schema.TypeList,
							Computed: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"limit_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"country_list": {
										Type:     schema.TypeString,
										Computed: true,
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
	client := m.(*connectivity.Client)
	service := NewCdnService(client)

	domain := d.Get("domain").(string)
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

	log.Printf("[INFO] Reading CDN domain: %s", domain)

	response, err := service.GetDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to read CDN domain: %w", err)
	}

	if len(response.Data) == 0 {
		log.Printf("[WARN] Domain does not exist: %s", domain)
		d.SetId("")
		return nil
	}

	domainData := response.Data[0]

	// Set all computed fields
	if err := d.Set("id", domainData.ID); err != nil {
		return fmt.Errorf("error setting id: %w", err)
	}
	if err := d.Set("domain", domainData.Domain); err != nil {
		return fmt.Errorf("error setting domain: %w", err)
	}
	if err := d.Set("type", domainData.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("status", domainData.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if err := d.Set("icp_num", domainData.IcpNum); err != nil {
		return fmt.Errorf("error setting icp_num: %w", err)
	}
	if err := d.Set("icp_status", domainData.IcpStatus); err != nil {
		return fmt.Errorf("error setting icp_status: %w", err)
	}
	if err := d.Set("area", domainData.Area); err != nil {
		return fmt.Errorf("error setting area: %w", err)
	}
	if err := d.Set("cname", domainData.Cname); err != nil {
		return fmt.Errorf("error setting cname: %w", err)
	}
	if err := d.Set("https", domainData.Https); err != nil {
		return fmt.Errorf("error setting https: %w", err)
	}
	if err := d.Set("create_time", domainData.CreateTime); err != nil {
		return fmt.Errorf("error setting create_time: %w", err)
	}
	if err := d.Set("update_time", domainData.UpdateTime); err != nil {
		return fmt.Errorf("error setting update_time: %w", err)
	}

	// 2. Get domain configuration
	log.Printf("[INFO] Reading domain configuration: %s", domain)
	response2, err := service.GetDomainConfig(domain, configItem)
	if err != nil {
		return fmt.Errorf("failed to read domain configuration: %w", err)
	}

	if len(response2.Data) == 0 {
		log.Printf("[WARN] Domain configuration does not exist: %s", domain)
		d.SetId("")
		return nil
	}

	// Set resource ID
	d.SetId(domain)
	// Need to convert API returned configuration to Terraform schema expected format
	apiConfig := response2.Data[0].Config
	terraformConfig := convertAPIConfigToTerraform(apiConfig)
	configList := []map[string]interface{}{terraformConfig}
	if err := d.Set("config", configList); err != nil {
		return fmt.Errorf("error setting config: %w", err)
	}

	log.Printf("[INFO] Domain and configuration read successfully: %s", domain)
	return nil
}
