package scdn

// ============================================================================
// Origin Group Management Types
// ============================================================================

// OriginGroupListRequest origin group list request
type OriginGroupListRequest struct {
	Page     int    `json:"page,omitempty"`      // Page number
	PageSize int    `json:"page_size,omitempty"` // Page size
	Name     string `json:"name,omitempty"`      // Origin group name filter
}

// OriginGroupListResponse origin group list response
type OriginGroupListResponse struct {
	Status Status              `json:"status"`
	Data   OriginGroupListData `json:"data"`
}

// OriginGroupListData origin group list data
type OriginGroupListData struct {
	Total int               `json:"total"` // Total number of origin groups
	List  []OriginGroupInfo `json:"list"`  // List of origin groups
}

// OriginGroupDetailRequest origin group detail request
type OriginGroupDetailRequest struct {
	ID int `json:"id"` // Origin group ID
}

// OriginGroupDetailResponse origin group detail response
type OriginGroupDetailResponse struct {
	Status Status                `json:"status"`
	Data   OriginGroupDetailData `json:"data"`
}

// OriginGroupDetailData origin group detail data
type OriginGroupDetailData struct {
	OriginGroup OriginGroupInfo `json:"origin_group"` // Origin group information
}

// OriginGroupCreateRequest origin group create request
type OriginGroupCreateRequest struct {
	Name    string              `json:"name"`             // Origin group name (2-16 characters)
	Remark  string              `json:"remark,omitempty"` // Remark (2-64 characters)
	Origins []OriginGroupOrigin `json:"origins"`          // Origin list (at least 1)
}

// OriginGroupCreateResponse origin group create response
type OriginGroupCreateResponse struct {
	Status Status                `json:"status"`
	Data   OriginGroupCreateData `json:"data"`
}

// OriginGroupCreateData origin group create data
type OriginGroupCreateData struct {
	ID int `json:"id"` // Created origin group ID
}

// OriginGroupUpdateRequest origin group update request
type OriginGroupUpdateRequest struct {
	ID      int                 `json:"id"`               // Origin group ID
	Name    string              `json:"name"`             // Origin group name (2-16 characters)
	Remark  string              `json:"remark,omitempty"` // Remark (2-64 characters)
	Origins []OriginGroupOrigin `json:"origins"`          // Origin list (at least 1)
}

// OriginGroupUpdateResponse origin group update response
type OriginGroupUpdateResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// OriginGroupDeleteRequest origin group delete request
type OriginGroupDeleteRequest struct {
	IDs []int `json:"ids"` // Origin group ID array (at least 1)
}

// OriginGroupDeleteResponse origin group delete response
type OriginGroupDeleteResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// OriginGroupBindDomainsRequest bind origin group to domains request
type OriginGroupBindDomainsRequest struct {
	OriginGroupID  int      `json:"origin_group_id"`            // Origin group ID
	DomainIDs      []int    `json:"domain_ids,omitempty"`       // Domain ID array
	DomainGroupIDs []int    `json:"domain_group_ids,omitempty"` // Domain group ID array
	Domains        []string `json:"domains,omitempty"`          // Domain array
}

// OriginGroupBindDomainsResponse bind origin group to domains response
type OriginGroupBindDomainsResponse struct {
	Status Status                     `json:"status"`
	Data   OriginGroupBindDomainsData `json:"data"`
}

// OriginGroupBindDomainsData bind origin group to domains data
type OriginGroupBindDomainsData struct {
	JobID string `json:"job_id"` // Batch job ID
}

// OriginGroupAllRequest get all origin groups request
type OriginGroupAllRequest struct {
	ProtectStatus string `json:"protect_status"` // Protection status: scdn-shared nodes, exclusive-dedicated nodes
}

// OriginGroupAllResponse get all origin groups response
type OriginGroupAllResponse struct {
	Status Status              `json:"status"`
	Data   OriginGroupListData `json:"data"`
}

// OriginGroupCopyRequest copy origin group to domain request
type OriginGroupCopyRequest struct {
	OriginGroupID int `json:"origin_group_id"` // Origin group ID
	DomainID      int `json:"domain_id"`       // Domain ID
}

// OriginGroupCopyResponse copy origin group to domain response
type OriginGroupCopyResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// OriginGroupBindHistoryRequest get latest bind history request
type OriginGroupBindHistoryRequest struct {
	OriginGroupID int `json:"origin_group_id"` // Origin group ID
}

// OriginGroupBindHistoryResponse get latest bind history response
type OriginGroupBindHistoryResponse struct {
	Status Status                     `json:"status"`
	Data   OriginGroupBindHistoryData `json:"data"`
}

// OriginGroupBindHistoryData bind history data
type OriginGroupBindHistoryData struct {
	History OriginGroupBindHistory `json:"history"` // Bind history record
}

// OriginGroupBindHistory bind history record
type OriginGroupBindHistory struct {
	ID            int                            `json:"id"`              // History record ID
	OriginGroupID int                            `json:"origin_group_id"` // Origin group ID
	MemberID      int                            `json:"member_id"`       // Member ID
	Domains       []OriginGroupBindHistoryDomain `json:"domains"`         // Bound domain list
	CreatedAt     string                         `json:"created_at"`      // Creation time
	UpdatedAt     string                         `json:"updated_at"`      // Update time
}

// OriginGroupBindHistoryDomain bound domain in history
type OriginGroupBindHistoryDomain struct {
	DomainID   int    `json:"domain_id"`   // Domain ID
	DomainName string `json:"domain_name"` // Domain name
}

// OriginGroupInfo origin group information
type OriginGroupInfo struct {
	ID        int                 `json:"id"`         // Origin group ID
	Name      string              `json:"name"`       // Origin group name
	Remark    string              `json:"remark"`     // Remark
	MemberID  int                 `json:"member_id"`  // Member ID
	Username  string              `json:"username"`   // Username
	Origins   []OriginGroupOrigin `json:"origins"`    // Origin list
	CreatedAt string              `json:"created_at"` // Creation time
	UpdatedAt string              `json:"updated_at"` // Update time
}

// OriginGroupOrigin origin configuration
type OriginGroupOrigin struct {
	ID             int                       `json:"id,omitempty"`    // Origin ID (0 for new, >0 for update)
	OriginType     int                       `json:"origin_type"`     // Origin type: 0-IP, 1-domain
	Records        []OriginGroupRecord       `json:"records"`         // Origin record list (at least 1)
	ProtocolPorts  []OriginGroupProtocolPort `json:"protocol_ports"`  // Protocol port mapping (at least 1)
	OriginProtocol int                       `json:"origin_protocol"` // Origin protocol: 0-http, 1-https, 2-follow
	LoadBalance    int                       `json:"load_balance"`    // Load balance strategy: 0-ip_hash, 1-round_robin, 2-cookie
}

// OriginGroupRecord origin record
type OriginGroupRecord struct {
	Value    string `json:"value"`          // Origin address
	Port     int    `json:"port"`           // Origin port (1-65535)
	Priority int    `json:"priority"`       // Weight (1-100)
	View     string `json:"view"`           // Origin type: primary-backup, backup-backup
	Host     string `json:"host,omitempty"` // Origin Host
}

// OriginGroupProtocolPort protocol port mapping
type OriginGroupProtocolPort struct {
	Protocol    int   `json:"protocol"`     // Protocol: 0-http, 1-https
	ListenPorts []int `json:"listen_ports"` // Listen port list
}
