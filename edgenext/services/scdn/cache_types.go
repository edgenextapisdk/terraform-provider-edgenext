package scdn

// ============================================================================
// Cache Rule Management Types
// ============================================================================

// CacheRuleGetRulesRequest get cache rules request
type CacheRuleGetRulesRequest struct {
	BusinessID   int    `json:"business_id"`         // Business ID
	BusinessType string `json:"business_type"`       // Business type: "domain" or "tpl"
	Page         int    `json:"page,omitempty"`      // Page number
	PageSize     int    `json:"page_size,omitempty"` // Page size
	ID           int    `json:"id,omitempty"`        // Rule ID (optional, converted to id query parameter)
}

// CacheRuleGetRulesResponse get cache rules response
type CacheRuleGetRulesResponse struct {
	Status Status                `json:"status"`
	Data   CacheRuleGetRulesData `json:"data"`
}

// CacheRuleGetRulesData get cache rules data
type CacheRuleGetRulesData struct {
	Page     int             `json:"page"`      // Current page
	PageSize int             `json:"page_size"` // Page size
	Total    int             `json:"total"`     // Total records
	List     []CacheRuleInfo `json:"list"`      // Rule list
}

// CacheRuleInfo cache rule information
type CacheRuleInfo struct {
	ID     int            `json:"id"`     // Rule ID
	Name   string         `json:"name"`   // Rule name
	Remark string         `json:"remark"` // Remark
	Status int            `json:"status"` // Status (1: enabled, 2: disabled)
	Weight int            `json:"weight"` // Weight
	Mode   int            `json:"mode"`   // Mode (can be ignored)
	Expr   string         `json:"expr"`   // Wirefilter rule
	Type   string         `json:"type"`   // Type: "domain", "tpl", or "global"
	Conf   *CacheRuleConf `json:"conf"`   // Cache configuration
}

// CacheRuleConf cache rule configuration
type CacheRuleConf struct {
	NoCache          bool              `json:"nocache"`                      // Cache eligibility (true: bypass cache, false: cache)
	CacheRule        *CacheRule        `json:"cache_rule,omitempty"`         // Edge TTL cache
	BrowserCacheRule *BrowserCacheRule `json:"browser_cache_rule,omitempty"` // Browser cache
	CacheErrStatus   []CacheErrStatus  `json:"cache_errstatus,omitempty"`    // Status code cache config
	CacheURLRewrite  *CacheURLRewrite  `json:"cache_url_rewrite,omitempty"`  // Custom cache key
	CacheShare       *CacheShare       `json:"cache_share,omitempty"`        // Cache sharing
}

// CacheRule edge TTL cache configuration
type CacheRule struct {
	CacheTime           int    `json:"cachetime"`                       // Cache time
	IgnoreCacheTime     bool   `json:"ignore_cache_time,omitempty"`     // Ignore source cache time
	IgnoreNoCacheHeader bool   `json:"ignore_nocache_header,omitempty"` // Ignore no-cache header
	NoCacheControlOp    string `json:"no_cache_control_op,omitempty"`   // No cache control operation
	Action              string `json:"action,omitempty"`                // Cache action: "default", "nocache", "cachetime", "force"
}

// BrowserCacheRule browser cache configuration
type BrowserCacheRule struct {
	CacheTime       int  `json:"cachetime"`         // Cache time
	IgnoreCacheTime bool `json:"ignore_cache_time"` // Ignore source cache time (cache-control)
	NoCache         bool `json:"nocache"`           // Whether to cache
}

// CacheErrStatus status code cache configuration
type CacheErrStatus struct {
	CacheTime int   `json:"cachetime"`  // Status code cache time
	ErrStatus []int `json:"err_status"` // Status code array
}

// CacheURLRewrite custom cache key configuration
type CacheURLRewrite struct {
	SortArgs   bool                    `json:"sort_args"`         // Parameter sorting
	IgnoreCase bool                    `json:"ignore_case"`       // Ignore case
	Queries    *CacheURLRewriteQueries `json:"queries,omitempty"` // Query string processing
	Cookies    *CacheURLRewriteCookies `json:"cookies,omitempty"` // Cookie processing
}

// CacheURLRewriteQueries query string processing
type CacheURLRewriteQueries struct {
	ArgsMethod string   `json:"args_method"` // Action: "SAVE", "DEL", "IGNORE", "CUT"
	Items      []string `json:"items"`       // Parameter keys
}

// CacheURLRewriteCookies cookie processing
type CacheURLRewriteCookies struct {
	ArgsMethod string   `json:"args_method"` // Action: "SAVE", "DEL", "IGNORE", "CUT"
	Items      []string `json:"items"`       // Cookie keys
}

// CacheShare cache sharing configuration
type CacheShare struct {
	Scheme string `json:"scheme"` // HTTP/HTTPS cache sharing method: "http" or "https"
}

// CacheRuleCreateRequest create cache rule request
type CacheRuleCreateRequest struct {
	BusinessID   int            `json:"business_id"`      // Business ID
	BusinessType string         `json:"business_type"`    // Business type: "domain" or "tpl"
	Name         string         `json:"name"`             // Rule name
	Expr         string         `json:"expr"`             // Wirefilter rule
	Remark       string         `json:"remark,omitempty"` // Rule remark
	Conf         *CacheRuleConf `json:"conf"`             // Configuration
}

// CacheRuleCreateResponse create cache rule response
type CacheRuleCreateResponse struct {
	Status Status              `json:"status"`
	Data   CacheRuleCreateData `json:"data"`
}

// CacheRuleCreateData create cache rule data
type CacheRuleCreateData struct {
	ID int `json:"id"` // Rule ID
}

// CacheRuleUpdateRequest update cache rule (name/remark) request
type CacheRuleUpdateRequest struct {
	ID     int    `json:"id"`               // Rule ID
	Name   string `json:"name,omitempty"`   // Rule name
	Remark string `json:"remark,omitempty"` // Rule remark
}

// CacheRuleUpdateResponse update cache rule response
type CacheRuleUpdateResponse struct {
	Status Status              `json:"status"`
	Data   CacheRuleUpdateData `json:"data"`
}

// CacheRuleUpdateData update cache rule data
type CacheRuleUpdateData struct {
	ID int `json:"id"` // Rule ID
}

// CacheRuleUpdateConfigRequest update cache rule configuration request
type CacheRuleUpdateConfigRequest struct {
	ID           int            `json:"id"`               // Rule ID
	Name         string         `json:"name,omitempty"`   // Rule name
	Remark       string         `json:"remark,omitempty"` // Rule remark
	Expr         string         `json:"expr,omitempty"`   // Wirefilter rule
	Conf         *CacheRuleConf `json:"conf"`             // Configuration
	BusinessID   int            `json:"business_id"`      // Business ID
	BusinessType string         `json:"business_type"`    // Business type
}

// CacheRuleUpdateConfigResponse update cache rule configuration response
type CacheRuleUpdateConfigResponse struct {
	Status Status                    `json:"status"`
	Data   CacheRuleUpdateConfigData `json:"data"`
}

// CacheRuleUpdateConfigData update cache rule configuration data
type CacheRuleUpdateConfigData struct {
	ID int `json:"id"` // Rule ID
}

// CacheRuleUpdateStatusRequest update cache rule status request
type CacheRuleUpdateStatusRequest struct {
	BusinessID   int    `json:"business_id"`   // Business ID
	BusinessType string `json:"business_type"` // Business type: "domain" or "tpl"
	IDs          []int  `json:"ids"`           // Rule IDs array
	Status       int    `json:"status"`        // Status: 1 (enabled) or 2 (disabled)
}

// CacheRuleUpdateStatusResponse update cache rule status response
type CacheRuleUpdateStatusResponse struct {
	Status Status                    `json:"status"`
	Data   CacheRuleUpdateStatusData `json:"data"`
}

// CacheRuleUpdateStatusData update cache rule status data
type CacheRuleUpdateStatusData struct {
	IDs []int `json:"ids"` // Rule IDs that were updated
}

// CacheRuleSortRequest sort cache rules request
type CacheRuleSortRequest struct {
	BusinessID   int    `json:"business_id"`   // Business ID
	BusinessType string `json:"business_type"` // Business type: "domain" or "tpl"
	IDs          []int  `json:"ids"`           // Sorted rule IDs array
}

// CacheRuleSortResponse sort cache rules response
type CacheRuleSortResponse struct {
	Status Status            `json:"status"`
	Data   CacheRuleSortData `json:"data"`
}

// CacheRuleSortData sort cache rules data
type CacheRuleSortData struct {
	IDs []int `json:"ids"` // Sorted rule IDs
}

// CacheRuleDeleteRequest delete cache rule request
type CacheRuleDeleteRequest struct {
	BusinessID   int    `json:"business_id"`   // Business ID
	BusinessType string `json:"business_type"` // Business type: "domain" or "tpl"
	IDs          []int  `json:"ids"`           // Rule IDs to delete
}

// CacheRuleDeleteResponse delete cache rule response
type CacheRuleDeleteResponse struct {
	Status Status              `json:"status"`
	Data   CacheRuleDeleteData `json:"data"`
}

// CacheRuleDeleteData delete cache rule data
type CacheRuleDeleteData struct {
	IDs []int `json:"ids"` // Deleted rule IDs
}

// CacheGlobalConfigGetResponse get global cache config response
type CacheGlobalConfigGetResponse struct {
	Status Status                   `json:"status"`
	Data   CacheGlobalConfigGetData `json:"data"`
}

// CacheGlobalConfigGetData get global cache config data
type CacheGlobalConfigGetData struct {
	ID   int            `json:"id"`   // Rule ID
	Name string         `json:"name"` // Rule name
	Conf *CacheRuleConf `json:"conf"` // Cache configuration
}
