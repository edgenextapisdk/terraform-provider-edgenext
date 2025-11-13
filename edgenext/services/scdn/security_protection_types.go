package scdn

// ============================================================================
// Security Protection Types
// ============================================================================

// ============================================================================
// DDoS Protection Types
// ============================================================================

// DdosProtectionGetConfigRequest get DDoS protection config request
type DdosProtectionGetConfigRequest struct {
	BusinessID int      `json:"business_id"`    // Business ID
	Keys       []string `json:"keys,omitempty"` // Specify config keys
}

// DdosProtectionGetConfigResponse get DDoS protection config response
type DdosProtectionGetConfigResponse struct {
	Status Status                      `json:"status"`
	Data   DdosProtectionGetConfigData `json:"data"`
}

// DdosProtectionGetConfigData get DDoS protection config data
type DdosProtectionGetConfigData struct {
	ApplicationDdosProtection *ApplicationDdosProtection `json:"application_ddos_protection,omitempty"`
	VisitorAuthentication     *VisitorAuthentication     `json:"visitor_authentication,omitempty"`
}

// ApplicationDdosProtection application layer DDoS protection config
type ApplicationDdosProtection struct {
	ID                  int    `json:"id,omitempty"`
	Status              string `json:"status"`                // on, off, keep
	AICCStatus          string `json:"ai_cc_status"`          // on, off
	Type                string `json:"type"`                  // default, normal, strict, captcha, keep
	NeedAttackDetection int    `json:"need_attack_detection"` // 0 or 1
	AIStatus            string `json:"ai_status"`             // on, off
}

// VisitorAuthentication visitor authentication config
type VisitorAuthentication struct {
	ID             int    `json:"id,omitempty"`
	Status         string `json:"status"`           // on, off
	AuthToken      string `json:"auth_token"`       // Authentication token
	PassStillCheck int    `json:"pass_still_check"` // 0 or 1
}

// DdosProtectionUpdateConfigRequest update DDoS protection config request
type DdosProtectionUpdateConfigRequest struct {
	BusinessID                int                        `json:"business_id"` // Business ID
	ApplicationDdosProtection *ApplicationDdosProtection `json:"application_ddos_protection,omitempty"`
	VisitorAuthentication     *VisitorAuthentication     `json:"visitor_authentication,omitempty"`
}

// DdosProtectionUpdateConfigResponse update DDoS protection config response
type DdosProtectionUpdateConfigResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// ============================================================================
// WAF Rule Config Types
// ============================================================================

// WafRuleConfigGetRequest get WAF rule config request
type WafRuleConfigGetRequest struct {
	BusinessID int      `json:"business_id"`    // Business ID
	Keys       []string `json:"keys,omitempty"` // Specify config keys
}

// WafRuleConfigGetResponse get WAF rule config response
type WafRuleConfigGetResponse struct {
	Status Status               `json:"status"`
	Data   WafRuleConfigGetData `json:"data"`
}

// WafRuleConfigGetData get WAF rule config data
type WafRuleConfigGetData struct {
	WafRuleConfig          *WafRuleConfig                `json:"waf_rule_config,omitempty"`
	WafInterceptPage       *WafInterceptPage             `json:"waf_intercept_page,omitempty"`
	ReplayAttackProtection *ReplayAttackProtectionConfig `json:"replay_attack_protection,omitempty"`
	CsrfProtection         *CSRFProtectionConfig         `json:"csrf_protection,omitempty"`
	WebShellProtection     *WebShellProtectionConfig     `json:"web_shell_protection,omitempty"`
}

// ReplayAttackProtectionConfig replay attack protection config
type ReplayAttackProtectionConfig struct {
	ID             int      `json:"id,omitempty"`
	Status         string   `json:"status"`                    // on, off, keep
	Action         string   `json:"action"`                    // captcha, deny, watch, keep
	Path           []string `json:"path"`                      // Path list
	IgnorePath     []string `json:"ignore_path"`               // Ignore path list
	ValidityPeriod int      `json:"validity_period,omitempty"` // Validity period
}

// CSRFProtectionConfig CSRF protection config
type CSRFProtectionConfig struct {
	ID         int      `json:"id,omitempty"`
	Status     string   `json:"status"`      // on, off, keep
	Action     string   `json:"action"`      // deny, watch, keep
	Path       []string `json:"path"`        // Path list
	IgnorePath []string `json:"ignore_path"` // Ignore path list
}

// WebShellProtectionConfig web shell protection config
type WebShellProtectionConfig struct {
	ID     int    `json:"id,omitempty"`
	Status string `json:"status"` // on, off
}

// BatchUpdateReplayAttackProtectionConfig batch update replay attack protection config
type BatchUpdateReplayAttackProtectionConfig struct {
	ID               int      `json:"id,omitempty"`
	Status           string   `json:"status"`                       // on, off, keep
	Action           string   `json:"action"`                       // captcha, deny, watch, keep
	Path             []string `json:"path"`                         // Path list
	PathAction       string   `json:"path_action,omitempty"`        // add, cover
	IgnorePath       []string `json:"ignore_path"`                  // Ignore path list
	IgnorePathAction string   `json:"ignore_path_action,omitempty"` // add, cover
	ValidityPeriod   int      `json:"validity_period,omitempty"`    // Validity period
}

// BatchUpdateCSRFProtectionConfig batch update CSRF protection config
type BatchUpdateCSRFProtectionConfig struct {
	ID               int      `json:"id,omitempty"`
	Status           string   `json:"status"`                       // on, off, keep
	Action           string   `json:"action"`                       // deny, watch, keep
	PathAction       string   `json:"path_action,omitempty"`        // add, cover
	Path             []string `json:"path"`                         // Path list
	IgnorePath       []string `json:"ignore_path"`                  // Ignore path list
	IgnorePathAction string   `json:"ignore_path_action,omitempty"` // add, cover
}

// BatchUpdateWafRuleConfigRequest batch update WAF rule config request
type BatchUpdateWafRuleConfigRequest struct {
	WafRuleConfig          *WafRuleConfig                           `json:"waf_rule_config,omitempty"`
	WafInterceptPage       *WafInterceptPage                        `json:"waf_intercept_page,omitempty"`
	ReplayAttackProtection *BatchUpdateReplayAttackProtectionConfig `json:"replay_attack_protection,omitempty"`
	CsrfProtection         *BatchUpdateCSRFProtectionConfig         `json:"csrf_protection,omitempty"`
	WebShellProtection     *WebShellProtectionConfig                `json:"web_shell_protection,omitempty"`
}

// WafRuleConfig WAF rule config
type WafRuleConfig struct {
	ID            int    `json:"id,omitempty"`
	Status        string `json:"status"`                    // on, off, keep
	AIStatus      string `json:"ai_status"`                 // on, off
	WafLevel      string `json:"waf_level"`                 // general, strict, keep
	WafMode       string `json:"waf_mode"`                  // off, active, block, ban, keep
	WafStrategyID int    `json:"waf_strategy_id,omitempty"` // WAF strategy ID
}

// WafInterceptPage WAF intercept page config
type WafInterceptPage struct {
	ID      int    `json:"id,omitempty"`
	Status  string `json:"status"`            // on, off
	Type    string `json:"type"`              // custom, default, keep
	Content string `json:"content,omitempty"` // Custom content
}

// WafRuleConfigUpdateRequest update WAF rule config request
type WafRuleConfigUpdateRequest struct {
	BusinessID             int                           `json:"business_id"` // Business ID
	WafRuleConfig          *WafRuleConfig                `json:"waf_rule_config,omitempty"`
	WafInterceptPage       *WafInterceptPage             `json:"waf_intercept_page,omitempty"`
	ReplayAttackProtection *ReplayAttackProtectionConfig `json:"replay_attack_protection,omitempty"`
	CsrfProtection         *CSRFProtectionConfig         `json:"csrf_protection,omitempty"`
	WebShellProtection     *WebShellProtectionConfig     `json:"web_shell_protection,omitempty"`
}

// WafRuleConfigUpdateResponse update WAF rule config response
type WafRuleConfigUpdateResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// ============================================================================
// Security Protection Template Types
// ============================================================================

// SecurityProtectionTemplateGetMemberGlobalResponse get member global template response
type SecurityProtectionTemplateGetMemberGlobalResponse struct {
	Status Status                                        `json:"status"`
	Data   SecurityProtectionTemplateGetMemberGlobalData `json:"data"`
}

// SecurityProtectionTemplateGetMemberGlobalData get member global template data
type SecurityProtectionTemplateGetMemberGlobalData struct {
	Template        *SecurityProtectionTemplateInfo `json:"template,omitempty"`
	BindDomainCount int                             `json:"bind_domain_count,omitempty"`
}

// SecurityProtectionTemplateInfo security protection template information
type SecurityProtectionTemplateInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // domain, template, global
	CreatedAt   string `json:"created_at"`
	Remark      string `json:"remark,omitempty"`
	SubMemberID int    `json:"sub_member_id,omitempty"`
	DomainID    int    `json:"domain_id,omitempty"` // When querying domain template, this field has value
	DomainCount int    `json:"domain_count,omitempty"`
}

// SecurityProtectionTemplateCreateRequest create template request
type SecurityProtectionTemplateCreateRequest struct {
	Name             string   `json:"name"`                         // Template name
	Remark           string   `json:"remark,omitempty"`             // Remark
	TemplateSourceID int      `json:"template_source_id,omitempty"` // Source template ID
	DomainIDs        []int    `json:"domain_ids,omitempty"`         // Domain ID list
	GroupIDs         []int    `json:"group_ids,omitempty"`          // Group ID list
	Domains          []string `json:"domains,omitempty"`            // Domain list
	BindAll          bool     `json:"bind_all,omitempty"`           // Bind all domains
}

// SecurityProtectionTemplateCreateResponse create template response
type SecurityProtectionTemplateCreateResponse struct {
	Status Status                               `json:"status"`
	Data   SecurityProtectionTemplateCreateData `json:"data"`
}

// SecurityProtectionTemplateCreateData create template data
type SecurityProtectionTemplateCreateData struct {
	BusinessID  int               `json:"business_id"`            // Created template ID
	FailDomains map[string]string `json:"fail_domains,omitempty"` // Failed domains
}

// SecurityProtectionTemplateCreateDomainRequest create domain template request
type SecurityProtectionTemplateCreateDomainRequest struct {
	DomainIDs        []int `json:"domain_ids"`         // Domain ID list
	TemplateSourceID int   `json:"template_source_id"` // Source template ID
}

// SecurityProtectionTemplateCreateDomainResponse create domain template response
type SecurityProtectionTemplateCreateDomainResponse struct {
	Status Status                                     `json:"status"`
	Data   SecurityProtectionTemplateCreateDomainData `json:"data"`
}

// SecurityProtectionTemplateCreateDomainData create domain template data
type SecurityProtectionTemplateCreateDomainData struct {
	FailDomains map[string]string `json:"fail_domains,omitempty"` // Failed domains
}

// SecurityProtectionTemplateSearchRequest search template list request
type SecurityProtectionTemplateSearchRequest struct {
	TplType    string `json:"tpl_type"`              // global, only_domain, more_domain
	SearchType string `json:"search_type,omitempty"` // Search type
	SearchKey  string `json:"search_key,omitempty"`  // Search keyword
	Page       int    `json:"page"`                  // Page number
	PageSize   int    `json:"page_size"`             // Page size
}

// SecurityProtectionTemplateSearchResponse search template list response
type SecurityProtectionTemplateSearchResponse struct {
	Status Status                               `json:"status"`
	Data   SecurityProtectionTemplateSearchData `json:"data"`
}

// SecurityProtectionTemplateSearchData search template list data
type SecurityProtectionTemplateSearchData struct {
	Templates        []SecurityProtectionTemplateInfo `json:"templates"`
	Total            int                              `json:"total"`
	TotalDomainCount int                              `json:"total_domain_count,omitempty"`
}

// SecurityProtectionTemplateBindDomainSearchRequest get template bind domain list request
type SecurityProtectionTemplateBindDomainSearchRequest struct {
	BusinessID int    `json:"business_id"`        // Business ID
	Page       int    `json:"page"`               // Page number
	PageSize   int    `json:"page_size"`          // Page size
	Domain     string `json:"domain,omitempty"`   // Domain
	TplType    string `json:"tpl_type,omitempty"` // global, only_domain, more_domain
}

// SecurityProtectionTemplateBindDomainSearchResponse get template bind domain list response
type SecurityProtectionTemplateBindDomainSearchResponse struct {
	Status Status                                         `json:"status"`
	Data   SecurityProtectionTemplateBindDomainSearchData `json:"data"`
}

// SecurityProtectionTemplateBindDomainSearchData get template bind domain list data
type SecurityProtectionTemplateBindDomainSearchData struct {
	Domains []SecurityProtectionTemplateDomainInfo `json:"domains"`
	Total   int                                    `json:"total"`
}

// SecurityProtectionTemplateDomainInfo template domain information
type SecurityProtectionTemplateDomainInfo struct {
	ID        int    `json:"id"`
	Domain    string `json:"domain"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Remark    string `json:"remark,omitempty"`
}

// SecurityProtectionTemplateBindDomainRequest bind template domain request
type SecurityProtectionTemplateBindDomainRequest struct {
	BusinessID      int   `json:"business_id"`                 // Business ID
	DomainIDs       []int `json:"domain_ids,omitempty"`        // Domain ID list
	BindBusinessIDs []int `json:"bind_business_ids,omitempty"` // Bind business ID list
	GroupIDs        []int `json:"group_ids,omitempty"`         // Group ID list
}

// SecurityProtectionTemplateBindDomainResponse bind template domain response
type SecurityProtectionTemplateBindDomainResponse struct {
	Status Status                                   `json:"status"`
	Data   SecurityProtectionTemplateBindDomainData `json:"data"`
}

// SecurityProtectionTemplateBindDomainData bind template domain data
type SecurityProtectionTemplateBindDomainData struct {
	FailDomains map[string]string `json:"fail_domains,omitempty"` // Failed domains
}

// SecurityProtectionTemplateDeleteRequest delete template request
type SecurityProtectionTemplateDeleteRequest struct {
	BusinessID int `json:"business_id"` // Business ID
}

// SecurityProtectionTemplateDeleteResponse delete template response
type SecurityProtectionTemplateDeleteResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// SecurityProtectionTemplateBatchConfigRequest batch config template request
type SecurityProtectionTemplateBatchConfigRequest struct {
	TemplateIDs                []int                                    `json:"template_ids"`                            // Template ID list
	DdosConfig                 *DdosProtectionGetConfigData             `json:"ddos_config,omitempty"`                   // DDoS config (GetDdosProtectionConfigResponse in proto)
	PreciseAccessControlConfig *UpdatePreciseAccessControlConfigRequest `json:"precise_access_control_config,omitempty"` // Precise access control config
	WafRuleConfig              *BatchUpdateWafRuleConfigRequest         `json:"waf_rule_config,omitempty"`               // WAF rule config
	BotManagementConfig        *UpdateBotManagementConfigRequest        `json:"bot_management_config,omitempty"`         // Bot management config
	All                        int                                      `json:"all,omitempty"`                           // All flag
	Domains                    []string                                 `json:"domains,omitempty"`                       // Domain list
	DomainIDs                  []int                                    `json:"domain_ids,omitempty"`                    // Domain ID list
}

// SecurityProtectionTemplateBatchConfigResponse batch config template response
type SecurityProtectionTemplateBatchConfigResponse struct {
	Status Status                                    `json:"status"`
	Data   SecurityProtectionTemplateBatchConfigData `json:"data"`
}

// SecurityProtectionTemplateBatchConfigData batch config template data
type SecurityProtectionTemplateBatchConfigData struct {
	FailTemplates map[string]string `json:"fail_templates,omitempty"` // Failed templates
}

// SecurityProtectionTemplateUnboundDomainSearchRequest get unbound template domain list request
type SecurityProtectionTemplateUnboundDomainSearchRequest struct {
	Domain   string `json:"domain,omitempty"`    // Domain
	Page     int    `json:"page"`                // Page number
	PageSize int    `json:"page_size"`           // Page size
	MemberID int    `json:"member_id,omitempty"` // Member ID
}

// SecurityProtectionTemplateUnboundDomainSearchResponse get unbound template domain list response
type SecurityProtectionTemplateUnboundDomainSearchResponse struct {
	Status Status                                            `json:"status"`
	Data   SecurityProtectionTemplateUnboundDomainSearchData `json:"data"`
}

// SecurityProtectionTemplateUnboundDomainSearchData get unbound template domain list data
type SecurityProtectionTemplateUnboundDomainSearchData struct {
	Domains []SecurityProtectionTemplateDomainInfo `json:"domains"`
	Total   int                                    `json:"total"`
}

// SecurityProtectionTemplateEditRequest edit template request
type SecurityProtectionTemplateEditRequest struct {
	BusinessID int    `json:"business_id"`      // Business ID
	Name       string `json:"name"`             // Template name
	Remark     string `json:"remark,omitempty"` // Remark
}

// SecurityProtectionTemplateEditResponse edit template response
type SecurityProtectionTemplateEditResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// ============================================================================
// Security Protection Iota Types
// ============================================================================

// SecurityProtectionIotaResponse get iota response
type SecurityProtectionIotaResponse struct {
	Status Status                     `json:"status"`
	Data   SecurityProtectionIotaData `json:"data"`
}

// SecurityProtectionIotaData get iota data
type SecurityProtectionIotaData struct {
	Iota map[string]string `json:"iota"` // Enum key-value pairs
}

// ============================================================================
// Member Package Change Types
// ============================================================================

// MemberPackageChangeRequest member package change request
type MemberPackageChangeRequest struct {
	MemberID    int    `json:"member_id"`     // Member ID
	NewMealFlag string `json:"new_meal_flag"` // New meal flag: YD-WEBAQJS-TY, YD-WEBAQJS-SY, YD-WEBAQJS-GJ, YD-WEBAQJS-QJ
}

// MemberPackageChangeResponse member package change response
type MemberPackageChangeResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// ============================================================================
// Domain Gray Types
// ============================================================================

// DomainGrayRequest domain gray request
type DomainGrayRequest struct {
	DomainID int `json:"domain_id"` // Domain ID
	MemberID int `json:"member_id"` // Member ID
}

// DomainGrayResponse domain gray response
type DomainGrayResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// ============================================================================
// Precise Access Control Types
// ============================================================================

// UpdatePreciseAccessControlConfigRequest update precise access control config request
type UpdatePreciseAccessControlConfigRequest struct {
	Action   string                       `json:"action"`   // add, cover
	Policies []PreciseAccessControlPolicy `json:"policies"` // Policy list
}

// PreciseAccessControlPolicy precise access control policy
type PreciseAccessControlPolicy struct {
	Type       string                   `json:"type"`        // Policy type
	Action     string                   `json:"action"`      // Policy action
	ActionData map[string]interface{}   `json:"action_data"` // Action data
	Rules      []map[string]interface{} `json:"rules"`       // Rules list
	From       string                   `json:"from"`        // From source
	Status     int                      `json:"status"`      // Status
}

// ============================================================================
// Bot Management Types
// ============================================================================

// UpdateBotManagementConfigRequest update bot management config request
type UpdateBotManagementConfigRequest struct {
	BusinessID int         `json:"business_id"`           // Business ID
	IDs        []int       `json:"ids"`                   // ID list
	DataAction interface{} `json:"data_action,omitempty"` // Data action
}
