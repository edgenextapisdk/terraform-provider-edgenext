package scdn

// ============================================================================
// Rule Template Management Types
// ============================================================================

// RuleTemplateCreateRequest rule template creation request
type RuleTemplateCreateRequest struct {
	Name        string                  `json:"name"`                  // Template name
	Description string                  `json:"description,omitempty"` // Template description
	AppType     string                  `json:"app_type"`              // Application type (e.g., "network_speed")
	TplType     string                  `json:"tpl_type,omitempty"`    // Template type
	FromTplID   int                     `json:"from_tpl_id,omitempty"` // Existing template ID to copy from
	BindDomain  *RuleTemplateBindDomain `json:"bind_domain,omitempty"` // Domain binding information
}

// RuleTemplateBindDomain domain binding information for template
type RuleTemplateBindDomain struct {
	AllDomain      bool     `json:"all_domain"`                 // If true, bind to all domains
	DomainIDs      []int    `json:"domain_ids,omitempty"`       // List of domain IDs to bind
	DomainGroupIDs []int    `json:"domain_group_ids,omitempty"` // List of domain group IDs to bind
	Domains        []string `json:"domains,omitempty"`          // List of domain names to bind
	IsBind         bool     `json:"is_bind"`                    // Whether to bind domains
}

// RuleTemplateCreateResponse rule template creation response
type RuleTemplateCreateResponse struct {
	Status Status                 `json:"status"`
	Data   RuleTemplateCreateData `json:"data"`
}

// RuleTemplateCreateData rule template creation data
type RuleTemplateCreateData struct {
	ID int `json:"id"` // Newly created rule template ID
}

// RuleTemplateUpdateRequest rule template update request
type RuleTemplateUpdateRequest struct {
	ID          int    `json:"id"`                    // Rule template ID to update
	Name        string `json:"name,omitempty"`        // New name for the rule template
	Description string `json:"description,omitempty"` // New description for the rule template
}

// RuleTemplateUpdateResponse rule template update response
type RuleTemplateUpdateResponse struct {
	Status Status                 `json:"status"`
	Data   RuleTemplateUpdateData `json:"data"`
}

// RuleTemplateUpdateData rule template update data
type RuleTemplateUpdateData struct {
	ID int `json:"id"` // Updated rule template ID
}

// RuleTemplateDeleteRequest rule template delete request
type RuleTemplateDeleteRequest struct {
	ID int `json:"id"` // Rule template ID to delete
}

// RuleTemplateDeleteResponse rule template delete response
type RuleTemplateDeleteResponse struct {
	Status Status                 `json:"status"`
	Data   RuleTemplateDeleteData `json:"data"`
}

// RuleTemplateDeleteData rule template delete data
type RuleTemplateDeleteData struct {
	ID int `json:"id"` // Deleted rule template ID
}

// RuleTemplateListRequest rule template list request
type RuleTemplateListRequest struct {
	Page     int    `json:"page,omitempty"`      // Page number for pagination, default: 1
	PageSize int    `json:"page_size,omitempty"` // Items per page, max: 1000, default: 1000
	Name     string `json:"name,omitempty"`      // Filter by rule template name
	Domain   string `json:"domain,omitempty"`    // Filter by associated domain
	AppType  string `json:"app_type,omitempty"`  // Filter by application type
}

// RuleTemplateListResponse rule template list response
type RuleTemplateListResponse struct {
	Status Status               `json:"status"`
	Data   RuleTemplateListData `json:"data"`
}

// RuleTemplateListData rule template list data
type RuleTemplateListData struct {
	Total int                `json:"total"` // Total number of rule templates
	List  []RuleTemplateInfo `json:"list"`  // List of rule templates
}

// RuleTemplateInfo rule template information
type RuleTemplateInfo struct {
	ID          int                          `json:"id"`           // Rule template ID
	Name        string                       `json:"name"`         // Rule template name
	Description string                       `json:"description"`  // Rule template description
	AppType     string                       `json:"app_type"`     // Application type
	BindDomains []RuleTemplateBindDomainInfo `json:"bind_domains"` // List of domains bound to this template
	CreatedAt   string                       `json:"created_at"`   // Template creation timestamp
}

// RuleTemplateBindDomainInfo domain binding information
type RuleTemplateBindDomainInfo struct {
	DomainID  int    `json:"domain_id"`  // Bound domain ID
	Domain    string `json:"domain"`     // Bound domain name
	CreatedAt string `json:"created_at"` // Domain binding timestamp
}

// RuleTemplateUnbindDomainRequest rule template unbind domain request
type RuleTemplateUnbindDomainRequest struct {
	ID        int   `json:"id"`         // Rule template ID
	DomainIDs []int `json:"domain_ids"` // List of domain IDs to unbind from template
}

// RuleTemplateUnbindDomainResponse rule template unbind domain response
type RuleTemplateUnbindDomainResponse struct {
	Status Status                       `json:"status"`
	Data   RuleTemplateUnbindDomainData `json:"data"`
}

// RuleTemplateUnbindDomainData rule template unbind domain data
type RuleTemplateUnbindDomainData struct {
	ID int `json:"id"` // Rule template ID
}

// RuleTemplateBindDomainRequest rule template bind domain request
type RuleTemplateBindDomainRequest struct {
	ID        int   `json:"id"`         // Rule template ID
	DomainIDs []int `json:"domain_ids"` // List of domain IDs to bind to template
}

// RuleTemplateBindDomainResponse rule template bind domain response
type RuleTemplateBindDomainResponse struct {
	Status Status                     `json:"status"`
	Data   RuleTemplateBindDomainData `json:"data"`
}

// RuleTemplateBindDomainData rule template bind domain data
type RuleTemplateBindDomainData struct {
	ID int `json:"id"` // Rule template ID
}

// RuleTemplateListDomainsRequest rule template list domains request
type RuleTemplateListDomainsRequest struct {
	ID       int    `json:"id"`                  // Rule template ID
	Page     int    `json:"page,omitempty"`      // Page number for pagination
	PageSize int    `json:"page_size,omitempty"` // Items per page
	AppType  string `json:"app_type"`            // Application type
	Domain   string `json:"domain,omitempty"`    // Filter by domain name
}

// RuleTemplateListDomainsResponse rule template list domains response
type RuleTemplateListDomainsResponse struct {
	Status Status                      `json:"status"`
	Data   RuleTemplateListDomainsData `json:"data"`
}

// RuleTemplateListDomainsData rule template list domains data
type RuleTemplateListDomainsData struct {
	Total int                      `json:"total"` // Total number of domains bound to template
	List  []RuleTemplateDomainInfo `json:"list"`  // List of domain information
}

// RuleTemplateDomainInfo domain information bound to template
type RuleTemplateDomainInfo struct {
	ID        int    `json:"id"`         // Domain ID
	Domain    string `json:"domain"`     // Domain name
	CreatedAt string `json:"created_at"` // Domain binding timestamp
}
