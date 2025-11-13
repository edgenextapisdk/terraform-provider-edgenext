package networkspeed

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnNetworkSpeedConfig returns the SCDN network speed config data source
func DataSourceEdgenextScdnNetworkSpeedConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnNetworkSpeedConfigRead,

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Business ID (template ID for 'tpl' type, user ID for 'global' type)",
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Business type: 'tpl' (template) or 'global'",
			},
			"config_groups": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration groups to retrieve",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields - all config groups
			"domain_proxy_conf": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Domain proxy configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_connect_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Connection timeout",
						},
						"fails_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Failure timeout",
						},
						"keep_new_src_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Keep new source time",
						},
						"max_fails": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max failures",
						},
						"proxy_keepalive": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Keepalive (0 or 1)",
						},
					},
				},
			},
			"upstream_redirect": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Upstream redirect configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"customized_req_headers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Customized request headers configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"resp_headers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Response headers configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"upstream_uri_change": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Upstream URI change configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"source_site_protect": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source site protection configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
						"num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of requests",
						},
						"second": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time in seconds",
						},
					},
				},
			},
			"slice": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Range request configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"https": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "HTTPS configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
						"http2https": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP to HTTPS redirect: 'off', 'all', or 'special'",
						},
						"http2https_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Redirect port",
						},
						"http2": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP2: 'on' or 'off'",
						},
						"hsts": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HSTS: 'on' or 'off'",
						},
						"ocsp_stapling": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OCSP Stapling: 'on' or 'off'",
						},
						"min_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Minimum TLS version",
						},
						"ciphers_preset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ciphers preset: 'default', 'strong', or 'custom'",
						},
						"custom_encrypt_algorithm": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Custom encryption algorithms",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"page_gzip": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Page Gzip configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"webp": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "WebP format configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"upload_file": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Upload file configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upload_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Upload size",
						},
						"upload_size_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unit (e.g., 'MB')",
						},
					},
				},
			},
			"websocket": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "WebSocket configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"mobile_jump": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Mobile jump configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
						"jump_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Jump URL",
						},
					},
				},
			},
			"custom_page": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Custom page configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"upstream_check": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Upstream check configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: 'on' or 'off'",
						},
						"fails": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Consecutive unavailable times (1-10)",
						},
						"intval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Check interval in seconds (3-300)",
						},
						"rise": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Consecutive available times (1-10)",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TCP connection timeout in seconds (1-10)",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Check type: 'tcp' or 'http'",
						},
						"op": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP method: 'HEAD', 'GET', or 'AUTO' (when type is 'http')",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP check path (when type is 'http')",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnNetworkSpeedConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	// Get config_groups if provided
	var configGroups []string
	if groups, ok := d.GetOk("config_groups"); ok {
		groupsList := groups.([]interface{})
		configGroups = make([]string, len(groupsList))
		for i, v := range groupsList {
			configGroups[i] = v.(string)
		}
	} else {
		// Default: get all config groups
		configGroups = []string{
			"domain_proxy_conf",
			"upstream_redirect",
			"customized_req_headers",
			"source_site_protect",
			"slice",
			"https",
			"page_gzip",
			"webp",
			"upload_file",
			"websocket",
			"mobile_jump",
			"custom_page",
			"upstream_uri_change",
			"resp_headers",
			"upstream_check",
		}
	}

	req := scdn.NetworkSpeedGetConfigRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ConfigGroups: configGroups,
	}

	log.Printf("[INFO] Reading SCDN network speed config: business_id=%d, business_type=%s", businessID, businessType)
	response, err := service.GetNetworkSpeedConfig(req)
	if err != nil {
		return fmt.Errorf("failed to read network speed config: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s", businessID, businessType))

	// Set computed fields
	if response.Data.DomainProxyConf != nil {
		domainProxyConfMap := map[string]interface{}{
			"proxy_connect_timeout": response.Data.DomainProxyConf.ProxyConnectTimeout,
			"fails_timeout":         response.Data.DomainProxyConf.FailsTimeout,
			"keep_new_src_time":     response.Data.DomainProxyConf.KeepNewSrcTime,
			"max_fails":             response.Data.DomainProxyConf.MaxFails,
			"proxy_keepalive":       response.Data.DomainProxyConf.ProxyKeepalive,
		}
		if err := d.Set("domain_proxy_conf", []interface{}{domainProxyConfMap}); err != nil {
			log.Printf("[WARN] Failed to set domain_proxy_conf: %v", err)
		}
	}

	if response.Data.UpstreamRedirect != nil {
		upstreamRedirectMap := map[string]interface{}{
			"status": response.Data.UpstreamRedirect.Status,
		}
		if err := d.Set("upstream_redirect", []interface{}{upstreamRedirectMap}); err != nil {
			log.Printf("[WARN] Failed to set upstream_redirect: %v", err)
		}
	}

	if response.Data.CustomizedReqHeaders != nil {
		customizedReqHeadersMap := map[string]interface{}{
			"status": response.Data.CustomizedReqHeaders.Status,
		}
		if err := d.Set("customized_req_headers", []interface{}{customizedReqHeadersMap}); err != nil {
			log.Printf("[WARN] Failed to set customized_req_headers: %v", err)
		}
	}

	if response.Data.RespHeaders != nil {
		respHeadersMap := map[string]interface{}{
			"status": response.Data.RespHeaders.Status,
		}
		if err := d.Set("resp_headers", []interface{}{respHeadersMap}); err != nil {
			log.Printf("[WARN] Failed to set resp_headers: %v", err)
		}
	}

	if response.Data.UpstreamURIChange != nil {
		upstreamURIChangeMap := map[string]interface{}{
			"status": response.Data.UpstreamURIChange.Status,
		}
		if err := d.Set("upstream_uri_change", []interface{}{upstreamURIChangeMap}); err != nil {
			log.Printf("[WARN] Failed to set upstream_uri_change: %v", err)
		}
	}

	if response.Data.SourceSiteProtect != nil {
		sourceSiteProtectMap := map[string]interface{}{
			"status": response.Data.SourceSiteProtect.Status,
			"num":    response.Data.SourceSiteProtect.Num,
			"second": response.Data.SourceSiteProtect.Second,
		}
		if err := d.Set("source_site_protect", []interface{}{sourceSiteProtectMap}); err != nil {
			log.Printf("[WARN] Failed to set source_site_protect: %v", err)
		}
	}

	if response.Data.Slice != nil {
		sliceMap := map[string]interface{}{
			"status": response.Data.Slice.Status,
		}
		if err := d.Set("slice", []interface{}{sliceMap}); err != nil {
			log.Printf("[WARN] Failed to set slice: %v", err)
		}
	}

	if response.Data.HTTPS != nil {
		httpsMap := map[string]interface{}{
			"status":                   response.Data.HTTPS.Status,
			"http2https":               response.Data.HTTPS.HTTP2HTTPS,
			"http2https_port":          response.Data.HTTPS.HTTP2HTTPSPort,
			"http2":                    response.Data.HTTPS.HTTP2,
			"hsts":                     response.Data.HTTPS.HSTS,
			"ocsp_stapling":            response.Data.HTTPS.OCSPStapling,
			"min_version":              response.Data.HTTPS.MinVersion,
			"ciphers_preset":           response.Data.HTTPS.CiphersPreset,
			"custom_encrypt_algorithm": response.Data.HTTPS.CustomEncryptAlgorithm,
		}
		if err := d.Set("https", []interface{}{httpsMap}); err != nil {
			log.Printf("[WARN] Failed to set https: %v", err)
		}
	}

	if response.Data.PageGzip != nil {
		pageGzipMap := map[string]interface{}{
			"status": response.Data.PageGzip.Status,
		}
		if err := d.Set("page_gzip", []interface{}{pageGzipMap}); err != nil {
			log.Printf("[WARN] Failed to set page_gzip: %v", err)
		}
	}

	if response.Data.WebP != nil {
		webpMap := map[string]interface{}{
			"status": response.Data.WebP.Status,
		}
		if err := d.Set("webp", []interface{}{webpMap}); err != nil {
			log.Printf("[WARN] Failed to set webp: %v", err)
		}
	}

	if response.Data.UploadFile != nil {
		uploadFileMap := map[string]interface{}{
			"upload_size":      response.Data.UploadFile.UploadSize,
			"upload_size_unit": response.Data.UploadFile.UploadSizeUnit,
		}
		if err := d.Set("upload_file", []interface{}{uploadFileMap}); err != nil {
			log.Printf("[WARN] Failed to set upload_file: %v", err)
		}
	}

	if response.Data.WebSocket != nil {
		websocketMap := map[string]interface{}{
			"status": response.Data.WebSocket.Status,
		}
		if err := d.Set("websocket", []interface{}{websocketMap}); err != nil {
			log.Printf("[WARN] Failed to set websocket: %v", err)
		}
	}

	if response.Data.MobileJump != nil {
		mobileJumpMap := map[string]interface{}{
			"status":   response.Data.MobileJump.Status,
			"jump_url": response.Data.MobileJump.JumpURL,
		}
		if err := d.Set("mobile_jump", []interface{}{mobileJumpMap}); err != nil {
			log.Printf("[WARN] Failed to set mobile_jump: %v", err)
		}
	}

	if response.Data.CustomPage != nil {
		customPageMap := map[string]interface{}{
			"status": response.Data.CustomPage.Status,
		}
		if err := d.Set("custom_page", []interface{}{customPageMap}); err != nil {
			log.Printf("[WARN] Failed to set custom_page: %v", err)
		}
	}

	if response.Data.UpstreamCheck != nil {
		upstreamCheckMap := map[string]interface{}{
			"status":  response.Data.UpstreamCheck.Status,
			"fails":   response.Data.UpstreamCheck.Fails,
			"intval":  response.Data.UpstreamCheck.Intval,
			"rise":    response.Data.UpstreamCheck.Rise,
			"timeout": response.Data.UpstreamCheck.Timeout,
			"type":    response.Data.UpstreamCheck.Type,
			"op":      response.Data.UpstreamCheck.Op,
			"path":    response.Data.UpstreamCheck.Path,
		}
		if err := d.Set("upstream_check", []interface{}{upstreamCheckMap}); err != nil {
			log.Printf("[WARN] Failed to set upstream_check: %v", err)
		}
	}

	// Handle result_output_file if provided
	if _, ok := d.GetOk("result_output_file"); ok {
		outputData := map[string]interface{}{
			"business_id":   businessID,
			"business_type": businessType,
			"config":        response.Data,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
	}

	log.Printf("[INFO] Network speed config read successfully")
	return nil
}
