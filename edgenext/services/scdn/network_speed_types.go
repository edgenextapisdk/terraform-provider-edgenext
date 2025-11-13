package scdn

// ============================================================================
// Network Speed Management Types
// ============================================================================

// NetworkSpeedGetConfigRequest get template config request
type NetworkSpeedGetConfigRequest struct {
	BusinessID   int      `json:"business_id"`   // Business ID
	BusinessType string   `json:"business_type"` // Business type: "tpl" (template) or "global"
	ConfigGroups []string `json:"config_groups"` // Configuration groups to retrieve
}

// NetworkSpeedGetConfigResponse get template config response
type NetworkSpeedGetConfigResponse struct {
	Status Status                    `json:"status"`
	Data   NetworkSpeedGetConfigData `json:"data"`
}

// NetworkSpeedGetConfigData get template config data
type NetworkSpeedGetConfigData struct {
	BusinessType         string                `json:"business_type"`
	BusinessID           int                   `json:"business_id"`
	DomainProxyConf      *DomainProxyConf      `json:"domain_proxy_conf,omitempty"`
	UpstreamRedirect     *UpstreamRedirect     `json:"upstream_redirect,omitempty"`
	CustomizedReqHeaders *CustomizedReqHeaders `json:"customized_req_headers,omitempty"`
	RespHeaders          *RespHeaders          `json:"resp_headers,omitempty"`
	UpstreamURIChange    *UpstreamURIChange    `json:"upstream_uri_change,omitempty"`
	SourceSiteProtect    *SourceSiteProtect    `json:"source_site_protect,omitempty"`
	Slice                *Slice                `json:"slice,omitempty"`
	HTTPS                *NetworkSpeedHTTPS    `json:"https,omitempty"`
	PageGzip             *PageGzip             `json:"page_gzip,omitempty"`
	WebP                 *WebP                 `json:"webp,omitempty"`
	UploadFile           *UploadFile           `json:"upload_file,omitempty"`
	WebSocket            *WebSocket            `json:"websocket,omitempty"`
	MobileJump           *MobileJump           `json:"mobile_jump,omitempty"`
	CustomPage           *CustomPage           `json:"custom_page,omitempty"`
	UpstreamCheck        *UpstreamCheck        `json:"upstream_check,omitempty"`
}

// DomainProxyConf domain proxy configuration
type DomainProxyConf struct {
	ProxyConnectTimeout int `json:"proxy_connect_timeout"` // Connection timeout
	FailsTimeout        int `json:"fails_timeout"`         // Failure timeout
	KeepNewSrcTime      int `json:"keep_new_src_time"`     // Keep new source time
	MaxFails            int `json:"max_fails"`             // Max failures
	ProxyKeepalive      int `json:"proxy_keepalive"`       // Keepalive (0 or 1)
}

// UpstreamRedirect upstream redirect configuration
type UpstreamRedirect struct {
	Status string `json:"status"` // "on" or "off"
}

// CustomizedReqHeaders customized request headers configuration
type CustomizedReqHeaders struct {
	Status string `json:"status"` // "on" or "off"
}

// RespHeaders response headers configuration
type RespHeaders struct {
	Status string `json:"status"` // "on" or "off"
}

// UpstreamURIChange upstream URI change configuration
type UpstreamURIChange struct {
	Status string `json:"status"` // "on" or "off"
}

// SourceSiteProtect source site protection configuration
type SourceSiteProtect struct {
	Status string `json:"status"` // "on" or "off"
	Num    int    `json:"num"`    // Number of requests
	Second int    `json:"second"` // Time in seconds
}

// Slice range request configuration
type Slice struct {
	Status string `json:"status"` // "on" or "off"
}

// NetworkSpeedHTTPS HTTPS configuration
type NetworkSpeedHTTPS struct {
	Status                 string   `json:"status"`                   // "on" or "off"
	HTTP2HTTPS             string   `json:"http2https"`               // "off", "all", or "special"
	HTTP2HTTPSPort         int      `json:"http2https_port"`          // Redirect port
	HTTP2                  string   `json:"http2"`                    // "on" or "off"
	HSTS                   string   `json:"hsts"`                     // "on" or "off"
	OCSPStapling           string   `json:"ocsp_stapling"`            // "on" or "off"
	MinVersion             string   `json:"min_version"`              // "SSLv3", "TLSv1.0", "TLSv1.1", "TLSv1.2", "TLSv1.3"
	CiphersPreset          string   `json:"ciphers_preset"`           // "default", "strong", or "custom"
	CustomEncryptAlgorithm []string `json:"custom_encrypt_algorithm"` // Custom encryption algorithms
}

// PageGzip page gzip configuration
type PageGzip struct {
	Status string `json:"status"` // "on" or "off"
}

// WebP WebP format configuration
type WebP struct {
	Status string `json:"status"` // "on" or "off"
}

// UploadFile upload file configuration
type UploadFile struct {
	UploadSize     int    `json:"upload_size"`      // Upload size
	UploadSizeUnit string `json:"upload_size_unit"` // Unit (e.g., "MB")
}

// WebSocket WebSocket configuration
type WebSocket struct {
	Status string `json:"status"` // "on" or "off"
}

// MobileJump mobile jump configuration
type MobileJump struct {
	Status  string `json:"status"`   // "on" or "off"
	JumpURL string `json:"jump_url"` // Jump URL
}

// CustomPage custom page configuration
type CustomPage struct {
	Status string `json:"status"` // "on" or "off"
}

// UpstreamCheck upstream check configuration
type UpstreamCheck struct {
	Status  string `json:"status"`         // "on" or "off"
	Fails   int    `json:"fails"`          // Consecutive unavailable times (1-10)
	Intval  int    `json:"intval"`         // Check interval in seconds (3-300)
	Rise    int    `json:"rise"`           // Consecutive available times (1-10)
	Timeout int    `json:"timeout"`        // TCP connection timeout in seconds (1-10)
	Type    string `json:"type"`           // "tcp" or "http"
	Op      string `json:"op,omitempty"`   // HTTP method: "HEAD", "GET", or "AUTO" (required when type is "http")
	Path    string `json:"path,omitempty"` // HTTP check path, must start with "/" (required when type is "http")
}

// NetworkSpeedUpdateConfigRequest update template config request
type NetworkSpeedUpdateConfigRequest struct {
	BusinessID           int                   `json:"business_id"`   // Business ID
	BusinessType         string                `json:"business_type"` // Business type: "tpl" or "global"
	DomainProxyConf      *DomainProxyConf      `json:"domain_proxy_conf,omitempty"`
	UpstreamRedirect     *UpstreamRedirect     `json:"upstream_redirect,omitempty"`
	CustomizedReqHeaders *CustomizedReqHeaders `json:"customized_req_headers,omitempty"`
	RespHeaders          *RespHeaders          `json:"resp_headers,omitempty"`
	UpstreamURIChange    *UpstreamURIChange    `json:"upstream_uri_change,omitempty"`
	SourceSiteProtect    *SourceSiteProtect    `json:"source_site_protect,omitempty"`
	Slice                *Slice                `json:"slice,omitempty"`
	HTTPS                *NetworkSpeedHTTPS    `json:"https,omitempty"`
	PageGzip             *PageGzip             `json:"page_gzip,omitempty"`
	WebP                 *WebP                 `json:"webp,omitempty"`
	UploadFile           *UploadFile           `json:"upload_file,omitempty"`
	WebSocket            *WebSocket            `json:"websocket,omitempty"`
	MobileJump           *MobileJump           `json:"mobile_jump,omitempty"`
	CustomPage           *CustomPage           `json:"custom_page,omitempty"`
	UpstreamCheck        *UpstreamCheck        `json:"upstream_check,omitempty"`
}

// NetworkSpeedUpdateConfigResponse update template config response
type NetworkSpeedUpdateConfigResponse struct {
	Status Status                       `json:"status"`
	Data   NetworkSpeedUpdateConfigData `json:"data"`
}

// NetworkSpeedUpdateConfigData update template config data
type NetworkSpeedUpdateConfigData struct {
	BusinessID   int    `json:"business_id"`
	BusinessType string `json:"business_type"`
	Updates      int    `json:"updates"` // Number of updated configs
	Adds         int    `json:"adds"`    // Number of added configs
}

// NetworkSpeedGetRulesRequest get rules request
type NetworkSpeedGetRulesRequest struct {
	BusinessID   int    `json:"business_id"`   // Business ID
	BusinessType string `json:"business_type"` // Business type
	ConfigGroup  string `json:"config_group"`  // Rule group: "custom_page", "upstream_uri_change_rule", "resp_headers_rule", "customized_req_headers_rule"
}

// NetworkSpeedGetRulesResponse get rules response
type NetworkSpeedGetRulesResponse struct {
	Status Status                   `json:"status"`
	Data   NetworkSpeedGetRulesData `json:"data"`
}

// NetworkSpeedGetRulesData get rules data
type NetworkSpeedGetRulesData struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Total    int                    `json:"total"`
	List     []NetworkSpeedRuleInfo `json:"list"`
}

// NetworkSpeedRuleInfo rule information
type NetworkSpeedRuleInfo struct {
	ID                       int                       `json:"id"`
	BusinessType             string                    `json:"business_type"`
	BusinessID               int                       `json:"business_id"`
	ConfigGroup              string                    `json:"config_group"`
	CustomPage               *CustomPageRule           `json:"custom_page,omitempty"`
	UpstreamURIChangeRule    *UpstreamURIChangeRule    `json:"upstream_uri_change_rule,omitempty"`
	RespHeadersRule          *RespHeadersRule          `json:"resp_headers_rule,omitempty"`
	CustomizedReqHeadersRule *CustomizedReqHeadersRule `json:"customized_req_headers_rule,omitempty"`
}

// CustomPageRule custom page rule
type CustomPageRule struct {
	StatusCode  int    `json:"status_code"`  // Status code
	PageType    string `json:"page_type"`    // Page type
	PageContent string `json:"page_content"` // Page content
}

// UpstreamURIChangeRule upstream URI change rule
type UpstreamURIChangeRule struct {
	Type   string `json:"typ"`    // Type
	Action string `json:"action"` // Action
	Match  string `json:"match"`  // Match value
	Target string `json:"target"` // Target value
}

// RespHeadersRule response headers rule
type RespHeadersRule struct {
	Type    string `json:"type"`    // Type
	Content string `json:"content"` // Content
	Remark  string `json:"remark"`  // Remark
}

// CustomizedReqHeadersRule customized request headers rule
type CustomizedReqHeadersRule struct {
	Type    string `json:"type"`    // Type
	Content string `json:"content"` // Content
	Remark  string `json:"remark"`  // Remark
}

// NetworkSpeedCreateRuleRequest create rule request
type NetworkSpeedCreateRuleRequest struct {
	BusinessID               int                       `json:"business_id"`   // Business ID
	BusinessType             string                    `json:"business_type"` // Business type
	ConfigGroup              string                    `json:"config_group"`  // Rule group
	CustomPage               *CustomPageRule           `json:"custom_page,omitempty"`
	UpstreamURIChangeRule    *UpstreamURIChangeRule    `json:"upstream_uri_change_rule,omitempty"`
	RespHeadersRule          *RespHeadersRule          `json:"resp_headers_rule,omitempty"`
	CustomizedReqHeadersRule *CustomizedReqHeadersRule `json:"customized_req_headers_rule,omitempty"`
}

// NetworkSpeedCreateRuleResponse create rule response
type NetworkSpeedCreateRuleResponse struct {
	Status Status                     `json:"status"`
	Data   NetworkSpeedCreateRuleData `json:"data"`
}

// NetworkSpeedCreateRuleData create rule data
type NetworkSpeedCreateRuleData struct {
	ID int `json:"id"` // Rule ID
}

// NetworkSpeedDeleteRuleRequest delete rule request
type NetworkSpeedDeleteRuleRequest struct {
	BusinessID   int    `json:"business_id"`   // Business ID
	BusinessType string `json:"business_type"` // Business type
	ConfigGroup  string `json:"config_group"`  // Rule group
	IDs          []int  `json:"ids"`           // Rule IDs to delete
}

// NetworkSpeedDeleteRuleResponse delete rule response
type NetworkSpeedDeleteRuleResponse struct {
	Status Status                     `json:"status"`
	Data   NetworkSpeedDeleteRuleData `json:"data"`
}

// NetworkSpeedDeleteRuleData delete rule data
type NetworkSpeedDeleteRuleData struct {
	IDs []int `json:"ids"` // Deleted rule IDs
}

// NetworkSpeedSortRulesRequest sort rules request
type NetworkSpeedSortRulesRequest struct {
	BusinessID   int    `json:"business_id"`   // Business ID
	BusinessType string `json:"business_type"` // Business type
	ConfigGroup  string `json:"config_group"`  // Rule group
	IDs          []int  `json:"ids"`           // Sorted rule IDs
}

// NetworkSpeedSortRulesResponse sort rules response
type NetworkSpeedSortRulesResponse struct {
	Status Status                    `json:"status"`
	Data   NetworkSpeedSortRulesData `json:"data"`
}

// NetworkSpeedSortRulesData sort rules data
type NetworkSpeedSortRulesData struct {
	IDs []int `json:"ids"` // Sorted rule IDs
}

// NetworkSpeedUpdateRuleRequest update rule request
type NetworkSpeedUpdateRuleRequest struct {
	ID                       int                       `json:"id"`           // Rule ID
	ConfigGroup              string                    `json:"config_group"` // Rule group
	CustomPage               *CustomPageRule           `json:"custom_page,omitempty"`
	UpstreamURIChangeRule    *UpstreamURIChangeRule    `json:"upstream_uri_change_rule,omitempty"`
	RespHeadersRule          *RespHeadersRule          `json:"resp_headers_rule,omitempty"`
	CustomizedReqHeadersRule *CustomizedReqHeadersRule `json:"customized_req_headers_rule,omitempty"`
}

// NetworkSpeedUpdateRuleResponse update rule response
type NetworkSpeedUpdateRuleResponse struct {
	Status Status                     `json:"status"`
	Data   NetworkSpeedUpdateRuleData `json:"data"`
}

// NetworkSpeedUpdateRuleData update rule data
type NetworkSpeedUpdateRuleData struct {
	ID int `json:"id"` // Rule ID
}
