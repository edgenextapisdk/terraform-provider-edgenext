package networkspeed

import (
	"fmt"
	"log"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnNetworkSpeedConfig returns the SCDN network speed config resource
func ResourceEdgenextScdnNetworkSpeedConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnNetworkSpeedConfigCreate,
		Read:   resourceScdnNetworkSpeedConfigRead,
		Update: resourceScdnNetworkSpeedConfigUpdate,
		Delete: resourceScdnNetworkSpeedConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Business ID (template ID for 'tpl' type, user ID for 'global' type)",
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business type: 'tpl' (template) or 'global'",
			},
			// Config groups - only include key ones, others follow similar pattern
			"domain_proxy_conf": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Domain proxy configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_connect_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Connection timeout",
						},
						"fails_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Failure timeout",
						},
						"keep_new_src_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Keep new source time",
						},
						"max_fails": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Max failures",
						},
						"proxy_keepalive": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Keepalive (0 or 1)",
						},
					},
				},
			},
			"upstream_redirect": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Upstream redirect configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"customized_req_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Customized request headers configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"resp_headers": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Response headers configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"upstream_uri_change": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Upstream URI change configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"source_site_protect": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Source site protection configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
						"num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of requests",
						},
						"second": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Time in seconds",
						},
					},
				},
			},
			"slice": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Range request configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"https": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "HTTPS configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
						"http2https": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP to HTTPS redirect: 'off', 'all', or 'special'",
						},
						"http2https_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Redirect port",
						},
						"http2": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP2: 'on' or 'off'",
						},
						"hsts": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HSTS: 'on' or 'off'",
						},
						"ocsp_stapling": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "OCSP Stapling: 'on' or 'off'",
						},
						"min_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Minimum TLS version",
						},
						"ciphers_preset": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ciphers preset: 'default', 'strong', or 'custom'",
						},
						"custom_encrypt_algorithm": {
							Type:        schema.TypeList,
							Optional:    true,
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
				Optional:    true,
				MaxItems:    1,
				Description: "Page Gzip configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"webp": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "WebP format configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"upload_file": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Upload file configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"upload_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Upload size",
						},
						"upload_size_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unit (e.g., 'MB')",
						},
					},
				},
			},
			"websocket": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "WebSocket configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"mobile_jump": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Mobile jump configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
						"jump_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Jump URL",
						},
					},
				},
			},
			"custom_page": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Custom page configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
					},
				},
			},
			"upstream_check": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Upstream check configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: 'on' or 'off'",
						},
						"fails": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Consecutive unavailable times (1-10)",
						},
						"intval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Check interval in seconds (3-300)",
						},
						"rise": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Consecutive available times (1-10)",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "TCP connection timeout in seconds (1-10)",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Check type: 'tcp' or 'http'",
						},
						"op": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP method: 'HEAD', 'GET', or 'AUTO' (required when type is 'http')",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP check path, must start with '/' (required when type is 'http')",
						},
					},
				},
			},
		},
	}
}

func resourceScdnNetworkSpeedConfigCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnNetworkSpeedConfigUpdate(d, m)
}

func resourceScdnNetworkSpeedConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	// Get all config groups
	configGroups := []string{
		"domain_proxy_conf", "upstream_redirect", "customized_req_headers",
		"source_site_protect", "slice", "https", "page_gzip", "webp",
		"upload_file", "websocket", "mobile_jump", "custom_page",
		"upstream_uri_change", "resp_headers", "upstream_check",
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

	// Set fields from response - all config groups
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

	log.Printf("[INFO] Network speed config read successfully")
	return nil
}

func resourceScdnNetworkSpeedConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	req := scdn.NetworkSpeedUpdateConfigRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
	}

	// Build domain_proxy_conf if provided
	if domainProxyConfList, ok := d.GetOk("domain_proxy_conf"); ok && len(domainProxyConfList.([]interface{})) > 0 {
		domainProxyConfMap := domainProxyConfList.([]interface{})[0].(map[string]interface{})
		req.DomainProxyConf = &scdn.DomainProxyConf{
			ProxyConnectTimeout: domainProxyConfMap["proxy_connect_timeout"].(int),
			FailsTimeout:        domainProxyConfMap["fails_timeout"].(int),
			KeepNewSrcTime:      domainProxyConfMap["keep_new_src_time"].(int),
			MaxFails:            domainProxyConfMap["max_fails"].(int),
			ProxyKeepalive:      domainProxyConfMap["proxy_keepalive"].(int),
		}
	}

	// Build upstream_redirect if provided
	if upstreamRedirectList, ok := d.GetOk("upstream_redirect"); ok && len(upstreamRedirectList.([]interface{})) > 0 {
		upstreamRedirectMap := upstreamRedirectList.([]interface{})[0].(map[string]interface{})
		req.UpstreamRedirect = &scdn.UpstreamRedirect{
			Status: upstreamRedirectMap["status"].(string),
		}
	}

	// Build customized_req_headers if provided
	if customizedReqHeadersList, ok := d.GetOk("customized_req_headers"); ok && len(customizedReqHeadersList.([]interface{})) > 0 {
		customizedReqHeadersMap := customizedReqHeadersList.([]interface{})[0].(map[string]interface{})
		req.CustomizedReqHeaders = &scdn.CustomizedReqHeaders{
			Status: customizedReqHeadersMap["status"].(string),
		}
	}

	// Build resp_headers if provided
	if respHeadersList, ok := d.GetOk("resp_headers"); ok && len(respHeadersList.([]interface{})) > 0 {
		respHeadersMap := respHeadersList.([]interface{})[0].(map[string]interface{})
		req.RespHeaders = &scdn.RespHeaders{
			Status: respHeadersMap["status"].(string),
		}
	}

	// Build upstream_uri_change if provided
	if upstreamURIChangeList, ok := d.GetOk("upstream_uri_change"); ok && len(upstreamURIChangeList.([]interface{})) > 0 {
		upstreamURIChangeMap := upstreamURIChangeList.([]interface{})[0].(map[string]interface{})
		req.UpstreamURIChange = &scdn.UpstreamURIChange{
			Status: upstreamURIChangeMap["status"].(string),
		}
	}

	// Build source_site_protect if provided
	if sourceSiteProtectList, ok := d.GetOk("source_site_protect"); ok && len(sourceSiteProtectList.([]interface{})) > 0 {
		sourceSiteProtectMap := sourceSiteProtectList.([]interface{})[0].(map[string]interface{})
		req.SourceSiteProtect = &scdn.SourceSiteProtect{
			Status: sourceSiteProtectMap["status"].(string),
			Num:    sourceSiteProtectMap["num"].(int),
			Second: sourceSiteProtectMap["second"].(int),
		}
	}

	// Build slice if provided
	if sliceList, ok := d.GetOk("slice"); ok && len(sliceList.([]interface{})) > 0 {
		sliceMap := sliceList.([]interface{})[0].(map[string]interface{})
		req.Slice = &scdn.Slice{
			Status: sliceMap["status"].(string),
		}
	}

	// Build https if provided
	if httpsList, ok := d.GetOk("https"); ok && len(httpsList.([]interface{})) > 0 {
		httpsMap := httpsList.([]interface{})[0].(map[string]interface{})
		https := &scdn.NetworkSpeedHTTPS{
			Status:         httpsMap["status"].(string),
			HTTP2HTTPS:     httpsMap["http2https"].(string),
			HTTP2HTTPSPort: httpsMap["http2https_port"].(int),
			HTTP2:          httpsMap["http2"].(string),
			HSTS:           httpsMap["hsts"].(string),
			OCSPStapling:   httpsMap["ocsp_stapling"].(string),
			MinVersion:     httpsMap["min_version"].(string),
			CiphersPreset:  httpsMap["ciphers_preset"].(string),
		}
		if customAlgo, ok := httpsMap["custom_encrypt_algorithm"]; ok && customAlgo != nil {
			algoList := customAlgo.([]interface{})
			https.CustomEncryptAlgorithm = make([]string, len(algoList))
			for i, v := range algoList {
				https.CustomEncryptAlgorithm[i] = v.(string)
			}
		}
		req.HTTPS = https
	}

	// Build page_gzip if provided
	if pageGzipList, ok := d.GetOk("page_gzip"); ok && len(pageGzipList.([]interface{})) > 0 {
		pageGzipMap := pageGzipList.([]interface{})[0].(map[string]interface{})
		req.PageGzip = &scdn.PageGzip{
			Status: pageGzipMap["status"].(string),
		}
	}

	// Build webp if provided
	if webpList, ok := d.GetOk("webp"); ok && len(webpList.([]interface{})) > 0 {
		webpMap := webpList.([]interface{})[0].(map[string]interface{})
		req.WebP = &scdn.WebP{
			Status: webpMap["status"].(string),
		}
	}

	// Build upload_file if provided
	if uploadFileList, ok := d.GetOk("upload_file"); ok && len(uploadFileList.([]interface{})) > 0 {
		uploadFileMap := uploadFileList.([]interface{})[0].(map[string]interface{})
		req.UploadFile = &scdn.UploadFile{
			UploadSize:     uploadFileMap["upload_size"].(int),
			UploadSizeUnit: uploadFileMap["upload_size_unit"].(string),
		}
	}

	// Build websocket if provided
	if websocketList, ok := d.GetOk("websocket"); ok && len(websocketList.([]interface{})) > 0 {
		websocketMap := websocketList.([]interface{})[0].(map[string]interface{})
		req.WebSocket = &scdn.WebSocket{
			Status: websocketMap["status"].(string),
		}
	}

	// Build mobile_jump if provided
	if mobileJumpList, ok := d.GetOk("mobile_jump"); ok && len(mobileJumpList.([]interface{})) > 0 {
		mobileJumpMap := mobileJumpList.([]interface{})[0].(map[string]interface{})
		status := mobileJumpMap["status"].(string)
		jumpURL := mobileJumpMap["jump_url"].(string)

		// Validate: if status is "on", jump_url must be provided and not empty
		if status == "on" && jumpURL == "" {
			return fmt.Errorf("jump_url is required when mobile_jump status is 'on'. Please provide a valid jump URL")
		}

		req.MobileJump = &scdn.MobileJump{
			Status:  status,
			JumpURL: jumpURL,
		}
	}

	// Build custom_page if provided
	if customPageList, ok := d.GetOk("custom_page"); ok && len(customPageList.([]interface{})) > 0 {
		customPageMap := customPageList.([]interface{})[0].(map[string]interface{})
		req.CustomPage = &scdn.CustomPage{
			Status: customPageMap["status"].(string),
		}
	}

	// Build upstream_check if provided
	if upstreamCheckList, ok := d.GetOk("upstream_check"); ok && len(upstreamCheckList.([]interface{})) > 0 {
		upstreamCheckMap := upstreamCheckList.([]interface{})[0].(map[string]interface{})
		checkType := upstreamCheckMap["type"].(string)

		upstreamCheck := &scdn.UpstreamCheck{
			Status:  upstreamCheckMap["status"].(string),
			Fails:   upstreamCheckMap["fails"].(int),
			Intval:  upstreamCheckMap["intval"].(int),
			Rise:    upstreamCheckMap["rise"].(int),
			Timeout: upstreamCheckMap["timeout"].(int),
			Type:    checkType,
		}

		// Validate and set HTTP-specific fields when type is "http"
		if checkType == "http" {
			var op, path string
			var opOk, pathOk bool

			if opVal, ok := upstreamCheckMap["op"]; ok && opVal != nil {
				op, opOk = opVal.(string)
			}
			if pathVal, ok := upstreamCheckMap["path"]; ok && pathVal != nil {
				path, pathOk = pathVal.(string)
			}

			if !opOk || op == "" {
				return fmt.Errorf("op is required when upstream_check type is 'http'. Valid values: 'HEAD', 'GET', 'AUTO'")
			}
			if op != "HEAD" && op != "GET" && op != "AUTO" {
				return fmt.Errorf("op must be one of 'HEAD', 'GET', or 'AUTO' when type is 'http', got: %s", op)
			}

			if !pathOk || path == "" {
				return fmt.Errorf("path is required when upstream_check type is 'http'")
			}
			if !strings.HasPrefix(path, "/") {
				return fmt.Errorf("path must start with '/' when type is 'http', got: %s", path)
			}

			upstreamCheck.Op = op
			upstreamCheck.Path = path
		}

		req.UpstreamCheck = upstreamCheck
	}

	log.Printf("[INFO] Updating SCDN network speed config: business_id=%d, business_type=%s", businessID, businessType)
	_, err := service.UpdateNetworkSpeedConfig(req)
	if err != nil {
		// Provide more specific error message for common errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "app not exist") || strings.Contains(errMsg, "code: 1010") {
			if businessType == "tpl" {
				return fmt.Errorf("template with business_id=%d does not exist. Please: 1) Verify the template ID is correct (you can list templates using the rule template data source), 2) Ensure the template exists in your account, 3) Check that you have permission to access it. Original error: %w", businessID, err)
			}
			return fmt.Errorf("business with business_id=%d and business_type=%s does not exist or is not accessible. Please verify the business ID is correct. Original error: %w", businessID, businessType, err)
		}
		if strings.Contains(errMsg, "JumpUrl is a required field") || strings.Contains(errMsg, "jump_url") {
			return fmt.Errorf("jump_url is required when mobile_jump status is 'on'. Please set mobile_jump.jump_url to a valid URL or set mobile_jump.status to 'off'. Original error: %w", err)
		}
		return fmt.Errorf("failed to update network speed config: %w", err)
	}

	return resourceScdnNetworkSpeedConfigRead(d, m)
}

func resourceScdnNetworkSpeedConfigDelete(d *schema.ResourceData, m interface{}) error {
	// Config deletion is typically a no-op - just remove from state
	log.Printf("[INFO] Network speed config resource deleted from state: %s", d.Id())
	return nil
}
