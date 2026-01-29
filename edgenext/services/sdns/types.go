package sdns

// Status represents the standard API status response
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ============================================================================
// DNS Domain Types
// ============================================================================

type DnsDomainInfo struct {
	ID              int    `json:"id"`
	MemberID        int    `json:"member_id"`
	Domain          string `json:"domain"`
	TrustStatus     int    `json:"trust_status"`
	TrustStatusDesc string `json:"trust_status_desc"`
	Status          int    `json:"status"` // 1正常
}

type DnsDomainListRequest struct {
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	GroupID int    `json:"group_id,omitempty"`
	Domain  string `json:"domain,omitempty"`
	Id      int    `json:"id,omitempty"`
}

type DnsDomainListResponse struct {
	Status Status            `json:"status"`
	Data   DnsDomainListData `json:"data"`
}

type DnsDomainListData struct {
	Total int             `json:"total"`
	List  []DnsDomainInfo `json:"list"`
}

type DnsDomainAddRequest struct {
	Domain string `json:"domain"`
}

type DnsDomainAddResponse struct {
	Status Status           `json:"status"`
	Data   DnsDomainAddData `json:"data"`
}

type DnsDomainAddData struct {
	ID int `json:"id"`
}

type DnsDomainDeleteRequest struct {
	DomainIDs []int `json:"domain_ids"`
}

type DnsDomainDeleteResponse struct {
	Status Status `json:"status"`
}

// ============================================================================
// DNS Domain Group Types
// ============================================================================

type DnsGroupListRequest struct {
	Page      int    `json:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty"`
	GroupName string `json:"group_name,omitempty"`
	Domain    string `json:"domain,omitempty"`
	Id        int    `json:"id,omitempty"`
}

type DnsGroupListResponse struct {
	Status Status           `json:"status"`
	Data   DnsGroupListData `json:"data"`
}

type DnsGroupListData struct {
	Total int        `json:"total"`
	List  []DnsGroup `json:"list"`
}

type DnsGroup struct {
	ID        int    `json:"id"`
	MemberID  int    `json:"member_id"`
	GroupName string `json:"group_name"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DnsGroupAddRequest struct {
	GroupName string `json:"group_name"`
	Remark    string `json:"remark,omitempty"`
	DomainIDs []int  `json:"domain_ids,omitempty"`
}

type DnsGroupAddResponse struct {
	Status Status `json:"status"`
	Data   struct {
		ID int `json:"id"`
	} `json:"data"`
}

type DnsGroupSaveRequest struct {
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
	Remark    string `json:"remark,omitempty"`
	DomainIDs []int  `json:"domain_ids,omitempty"`
}

type DnsGroupDelRequest struct {
	GroupID int `json:"group_id"`
}

type DnsGroupDomainListRequest struct {
	GroupID int    `json:"group_id"`
	Domain  string `json:"domain,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

type DnsGroupDomainListResponse struct {
	Status Status                 `json:"status"`
	Data   DnsGroupDomainListData `json:"data"`
}

type DnsGroupDomainListData struct {
	Total int              `json:"total"`
	List  []DnsGroupDomain `json:"list"`
}

type DnsGroupDomain struct {
	DomainID interface{} `json:"domain_id"` // Sometimes string, sometimes int in documentation
	Domain   string      `json:"domain"`
}

type DnsGroupDomainSaveRequest struct {
	GroupID   int    `json:"group_id"`
	DomainIDs []int  `json:"domain_ids,omitempty"`
	Action    string `json:"action"` // add, del
}

// ============================================================================
// DNS Domain Record Types
// ============================================================================

type DnsRecordListRequest struct {
	DomainID   int    `json:"domain_id"`
	Page       int    `json:"page,omitempty"`
	PerPage    int    `json:"per_page,omitempty"`
	RecordType string `json:"record_type,omitempty"`
	RecordName string `json:"record_name,omitempty"`
	GroupID    int    `json:"group_id,omitempty"`
}

type DnsRecordListResponse struct {
	Status Status            `json:"status"`
	Data   DnsRecordListData `json:"data"`
}

type DnsRecordListData struct {
	Total int         `json:"total"`
	List  []DnsRecord `json:"list"`
}

type DnsRecord struct {
	ID        int    `json:"id"`
	DomainID  int    `json:"domain_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	View      string `json:"view"`
	Value     string `json:"value"`
	MX        int    `json:"mx"`
	TTL       int    `json:"ttl"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	IsSyncCDN int    `json:"is_sync_cdn"`
	Remark    string `json:"remark"`
	Locked    bool   `json:"locked"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DnsRecordAddRequest struct {
	DomainID     int    `json:"domain_id"`
	RecordName   string `json:"record_name"`
	RecordType   string `json:"record_type"`
	RecordView   string `json:"record_view"`
	RecordValue  string `json:"record_value"`
	RecordMX     int    `json:"record_mx,omitempty"`
	RecordTTL    int    `json:"record_ttl,omitempty"`
	RecordRemark string `json:"record_remark,omitempty"`
}

type DnsRecordEditRequest struct {
	RecordID     int    `json:"record_id"`
	DomainID     int    `json:"domain_id"`
	RecordName   string `json:"record_name"`
	RecordType   string `json:"record_type"`
	RecordView   string `json:"record_view"`
	RecordValue  string `json:"record_value"`
	RecordMX     int    `json:"record_mx,omitempty"`
	RecordTTL    int    `json:"record_ttl,omitempty"`
	RecordRemark string `json:"record_remark,omitempty"`
}

type DnsRecordResponse struct {
	Status Status `json:"status"`
	Data   struct {
		ID int `json:"id"`
	} `json:"data"`
}

type DnsRecordDeleteRequest struct {
	RecordID int `json:"record_id"`
	DomainID int `json:"domain_id"`
}
