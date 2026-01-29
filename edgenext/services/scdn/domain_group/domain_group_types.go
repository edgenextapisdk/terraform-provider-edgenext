package domain_group

import "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"

// Status represents the standard API status response
type Status = scdn.Status

// DomainGroupSaveRequest represents the request to add or update a domain group
type DomainGroupSaveRequest struct {
	GroupID   int    `json:"group_id,omitempty"`   // Optional for add, required for update/save
	GroupName string `json:"group_name,omitempty"` // Required for add
	Remark    string `json:"remark,omitempty"`
}

// DomainGroupSaveResponse represents the response when saving a domain group
type DomainGroupSaveResponse struct {
	Status Status              `json:"status"`
	Data   DomainGroupSaveData `json:"data"`
}

type DomainGroupSaveData struct {
	ID string `json:"id"`
}

// DomainGroupDelRequest represents the request to delete a domain group
type DomainGroupDelRequest struct {
	GroupID int `json:"group_id"`
}

// DomainGroupDelResponse represents the response when deleting a domain group
type DomainGroupDelResponse struct {
	Status Status `json:"status"`
}

// DomainGroupListRequest represents the request to list domain groups
type DomainGroupListRequest struct {
	GroupName    string `json:"group_name,omitempty"`
	Domain       string `json:"domain,omitempty"`
	BindedDomain int    `json:"binded_domain,omitempty"` // 0: unbind, 1: bind, other: no limit
	Page         int    `json:"page,omitempty"`
	PerPage      int    `json:"per_page,omitempty"`
}

// DomainGroupListResponse represents the list of domain groups
type DomainGroupListResponse struct {
	Status Status              `json:"status"`
	Data   DomainGroupListData `json:"data"`
}

type DomainGroupListData struct {
	Total string        `json:"total"`
	List  []DomainGroup `json:"list"`
}

type DomainGroup struct {
	ID              string `json:"id"`
	MemberID        string `json:"member_id"`
	ProductFlag     string `json:"product_flag"`
	GroupName       string `json:"group_name"`
	Remark          string `json:"remark"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	FormatCreatedAt string `json:"format_created_at"` // Normalized field name in Terraform usually
	FormatUpdatedAt string `json:"format_updated_at"`
}

// DomainGroupDomainSaveRequest represents the request to bind/unbind domains to a group
type DomainGroupDomainSaveRequest struct {
	GroupID   int      `json:"group_id"`
	DomainIDs []string `json:"domain_ids,omitempty"`
	Domains   []string `json:"domains,omitempty"`
	Action    string   `json:"action"` // "add" or "del"
}

// DomainGroupDomainSaveResponse represents the response for domain binding
type DomainGroupDomainSaveResponse struct {
	Status Status `json:"status"`
}

// DomainGroupDomainListRequest represents the request to list domains in a group
type DomainGroupDomainListRequest struct {
	GroupID int    `json:"group_id"`
	Domain  string `json:"domain,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// DomainGroupDomainListResponse represents the list of domains in a group
type DomainGroupDomainListResponse struct {
	Status Status                    `json:"status"`
	Data   DomainGroupDomainListData `json:"data"`
}

type DomainGroupDomainListData struct {
	Total string              `json:"total"`
	Ports []string            `json:"ports"`
	List  []DomainGroupDomain `json:"list"`
}

type DomainGroupDomain struct {
	DomainID string `json:"domain_id"`
	Domain   string `json:"domain"`
}

// DomainGroupMoveDomainRequest represents the request to move domains between groups
type DomainGroupMoveDomainRequest struct {
	FromGroupID int   `json:"from_group_id"`
	ToGroupID   int   `json:"to_group_id"`
	DomainIDs   []int `json:"domain_ids"`
}

// DomainGroupMoveDomainResponse represents the response for moving domains
type DomainGroupMoveDomainResponse struct {
	Status Status `json:"status"`
}
