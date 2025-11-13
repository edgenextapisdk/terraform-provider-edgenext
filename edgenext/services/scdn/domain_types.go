package scdn

// ============================================================================
// Domain Management Types
// ============================================================================

// DomainListRequest domain list request
type DomainListRequest struct {
	Page                int    `json:"page,omitempty"`
	PageSize            int    `json:"page_size,omitempty"`
	ID                  int    `json:"id,omitempty"`
	AccessProgress      string `json:"access_progress,omitempty"`
	GroupID             int    `json:"group_id,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Remark              string `json:"remark,omitempty"`
	OriginIP            string `json:"origin_ip,omitempty"`
	CAStatus            string `json:"ca_status,omitempty"`
	AccessMode          string `json:"access_mode,omitempty"`
	ProtectStatus       string `json:"protect_status,omitempty"`
	ExclusiveResourceID int    `json:"exclusive_resource_id,omitempty"`
}

// DomainListResponse domain list response
type DomainListResponse struct {
	Status Status         `json:"status"`
	Data   DomainListData `json:"data"`
}

// DomainListData domain list data
type DomainListData struct {
	Total int          `json:"total"`
	List  []DomainInfo `json:"list"`
}

// DomainInfo domain information
type DomainInfo struct {
	ID                  int       `json:"id"`
	Domain              string    `json:"domain"`
	Remark              string    `json:"remark"`
	AccessProgress      string    `json:"access_progress"`
	AccessMode          string    `json:"access_mode"`
	ProtectStatus       string    `json:"protect_status"`
	EIForwardStatus     string    `json:"ei_forward_status"`
	Cname               CnameInfo `json:"cname"`
	UseMyCname          int       `json:"use_my_cname"`
	UseMyDNS            int       `json:"use_my_dns"`
	CAStatus            string    `json:"ca_status"`
	ExclusiveResourceID int       `json:"exclusive_resource_id"`
	AccessProgressDesc  string    `json:"access_progress_desc"`
	HasOrigin           bool      `json:"has_origin"`
	CAID                int       `json:"ca_id"`
	CreatedAt           string    `json:"created_at"`
	UpdatedAt           string    `json:"updated_at"`
	PriDomain           string    `json:"pri_domain"`
}

// CnameInfo CNAME information
type CnameInfo struct {
	Master string   `json:"master"`
	Slaves []string `json:"slaves"`
}

// DomainSimpleListRequest simple domain list request
type DomainSimpleListRequest struct {
	Domain  string `json:"domain,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// DomainSimpleListResponse simple domain list response
type DomainSimpleListResponse struct {
	Status Status               `json:"status"`
	Data   DomainSimpleListData `json:"data"`
}

// DomainSimpleListData simple domain list data
type DomainSimpleListData struct {
	Total int                `json:"total"`
	List  []DomainSimpleInfo `json:"list"`
}

// DomainSimpleInfo simple domain information
type DomainSimpleInfo struct {
	ID       int    `json:"id"`
	Domain   string `json:"domain"`
	MemberID int    `json:"member_id"`
}

// DomainCreateRequest domain creation request
type DomainCreateRequest struct {
	Domain              string   `json:"domain"`
	GroupID             int      `json:"group_id,omitempty"`
	ExclusiveResourceID int      `json:"exclusive_resource_id,omitempty"`
	Remark              string   `json:"remark,omitempty"`
	TplID               int      `json:"tpl_id,omitempty"`
	Origins             []Origin `json:"origins"`
	ProtectStatus       string   `json:"protect_status,omitempty"`
	TplRecommend        string   `json:"tpl_recommend,omitempty"`
	AppType             string   `json:"app_type,omitempty"`
}

// Origin origin server configuration
type Origin struct {
	Id             int            `json:"id,omitempty"`
	Protocol       int            `json:"protocol"`
	ListenPorts    []int          `json:"listen_ports"`
	OriginProtocol int            `json:"origin_protocol"`
	LoadBalance    int            `json:"load_balance"`
	OriginType     int            `json:"origin_type"`
	Records        []OriginRecord `json:"records"`
}

// EditOrigin origin server configuration for editing
type EditOrigin struct {
	Id             int            `json:"id,omitempty"`
	Protocol       int            `json:"protocol"`
	ListenPort     int            `json:"listen_port"`
	OriginProtocol int            `json:"origin_protocol"`
	LoadBalance    int            `json:"load_balance"`
	OriginType     int            `json:"origin_type"`
	Records        []OriginRecord `json:"records"`
}

// OriginRecord origin record
type OriginRecord struct {
	View     string `json:"view"`
	Value    string `json:"value"`
	Port     int    `json:"port"`
	Priority int    `json:"priority"`
}

// DomainCreateResponse domain creation response
type DomainCreateResponse struct {
	Status Status           `json:"status"`
	Data   DomainCreateData `json:"data"`
}

// DomainCreateData domain creation data
type DomainCreateData struct {
	ID         int    `json:"id"`
	Record     string `json:"record"`
	Cname      string `json:"cname"`
	RecordType string `json:"record_type"`
	Domain     string `json:"domain"`
	PriDomain  string `json:"pri_domain"`
}

// DomainUpdateRequest domain update request
type DomainUpdateRequest struct {
	DomainID int    `json:"domain_id"`
	Remark   string `json:"remark,omitempty"`
}

// DomainUpdateResponse domain update response
type DomainUpdateResponse struct {
	Status Status           `json:"status"`
	Data   DomainUpdateData `json:"data"`
}

// DomainUpdateData domain update data
type DomainUpdateData struct {
	DomainID int `json:"domain_id"`
}

// DomainCertBindRequest domain certificate bind request
type DomainCertBindRequest struct {
	DomainID int `json:"domain_id"`
	CAID     int `json:"ca_id"`
}

// DomainCertBindResponse domain certificate bind response
type DomainCertBindResponse struct {
	Status Status             `json:"status"`
	Data   DomainCertBindData `json:"data"`
}

// DomainCertBindData domain certificate bind data
type DomainCertBindData struct {
	CAID int `json:"ca_id"`
}

// DomainCertUnbindRequest domain certificate unbind request
type DomainCertUnbindRequest struct {
	DomainID int `json:"domain_id"`
	CAID     int `json:"ca_id"`
}

// DomainCertUnbindResponse domain certificate unbind response
type DomainCertUnbindResponse struct {
	Status Status               `json:"status"`
	Data   DomainCertUnbindData `json:"data"`
}

// DomainCertUnbindData domain certificate unbind data
type DomainCertUnbindData struct {
	CAID int `json:"ca_id"`
}

// DomainDeleteRequest domain delete request
type DomainDeleteRequest struct {
	IDs []int `json:"ids"`
}

// DomainDeleteResponse domain delete response
type DomainDeleteResponse struct {
	Status Status           `json:"status"`
	Data   DomainDeleteData `json:"data"`
}

// DomainDeleteData domain delete data
type DomainDeleteData struct {
	IDs []int `json:"ids"`
}

// DomainDisableRequest domain disable request
type DomainDisableRequest struct {
	DomainIDs []int `json:"domain_ids"`
}

// DomainDisableResponse domain disable response
type DomainDisableResponse struct {
	Status Status            `json:"status"`
	Data   DomainDisableData `json:"data"`
}

// DomainDisableData domain disable data
type DomainDisableData struct {
	DomainIDs []int `json:"domain_ids"`
}

// DomainEnableRequest domain enable request
type DomainEnableRequest struct {
	DomainIDs []int `json:"domain_ids"`
}

// DomainEnableResponse domain enable response
type DomainEnableResponse struct {
	Status Status           `json:"status"`
	Data   DomainEnableData `json:"data"`
}

// DomainEnableData domain enable data
type DomainEnableData struct {
	DomainIDs []int `json:"domain_ids"`
}

// DomainAccessRefreshRequest domain access refresh request
type DomainAccessRefreshRequest struct {
	DomainIDs []int `json:"domain_ids"`
}

// DomainAccessRefreshResponse domain access refresh response
type DomainAccessRefreshResponse struct {
	Status Status                  `json:"status"`
	Data   DomainAccessRefreshData `json:"data"`
}

// DomainAccessRefreshData domain access refresh data
type DomainAccessRefreshData struct {
	DomainIDs []int `json:"domain_ids"`
}

// DomainExportRequest domain export request
type DomainExportRequest struct {
	DomainIDs           []int  `json:"domain_ids,omitempty"`
	AccessProgress      string `json:"access_progress,omitempty"`
	GroupID             int    `json:"group_id,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Remark              string `json:"remark,omitempty"`
	OriginIP            string `json:"origin_ip,omitempty"`
	CAStatus            string `json:"ca_status,omitempty"`
	AccessMode          string `json:"access_mode,omitempty"`
	ProtectStatus       string `json:"protect_status,omitempty"`
	ExclusiveResourceID int    `json:"exclusive_resource_id,omitempty"`
}

// DomainExportResponse domain export response
type DomainExportResponse struct {
	Status Status           `json:"status"`
	Data   DomainExportData `json:"data"`
}

// DomainExportData domain export data
type DomainExportData struct {
	Hash    string `json:"hash"`
	Key     string `json:"key"`
	RealURL string `json:"real_url"`
}

// OriginAddRequest origin add request
type OriginAddRequest struct {
	DomainID int      `json:"domain_id"`
	Origins  []Origin `json:"origins"`
}

// OriginAddResponse origin add response
type OriginAddResponse struct {
	Status Status        `json:"status"`
	Data   OriginAddData `json:"data"`
}

// OriginAddData origin add data
type OriginAddData struct {
	IDs []int `json:"ids"`
}

// OriginUpdateRequest origin update request
type OriginUpdateRequest struct {
	DomainID int          `json:"domain_id"`
	Origins  []EditOrigin `json:"origins"`
}

// OriginUpdateResponse origin update response
type OriginUpdateResponse struct {
	Status Status           `json:"status"`
	Data   OriginUpdateData `json:"data"`
}

// OriginUpdateData origin update data
type OriginUpdateData struct {
	ID int `json:"id"`
}

// OriginDeleteRequest origin delete request
type OriginDeleteRequest struct {
	IDs      []int `json:"ids"`
	DomainID int   `json:"domain_id"`
}

// OriginDeleteResponse origin delete response
type OriginDeleteResponse struct {
	Status Status           `json:"status"`
	Data   OriginDeleteData `json:"data"`
}

// OriginDeleteData origin delete data
type OriginDeleteData struct {
	IDs []int `json:"ids"`
}

// OriginListRequest origin list request
type OriginListRequest struct {
	DomainID int `json:"domain_id,omitempty"`
}

// OriginListResponse origin list response
type OriginListResponse struct {
	Status Status         `json:"status"`
	Data   OriginListData `json:"data"`
}

// OriginListData origin list data
type OriginListData struct {
	Total int          `json:"total"`
	List  []OriginInfo `json:"list"`
}

// OriginInfo origin information
type OriginInfo struct {
	ID             int            `json:"id"`
	DomainID       int            `json:"domain_id"`
	Protocol       int            `json:"protocol"`
	ListenPort     int            `json:"listen_port"`
	OriginProtocol int            `json:"origin_protocol"`
	LoadBalance    int            `json:"load_balance"`
	OriginType     int            `json:"origin_type"`
	Records        []OriginRecord `json:"records"`
}

// DomainNodeSwitchRequest domain node switch request
type DomainNodeSwitchRequest struct {
	DomainID            int    `json:"domain_id"`
	ProtectStatus       string `json:"protect_status"`
	ExclusiveResourceID int    `json:"exclusive_resource_id,omitempty"`
}

// DomainNodeSwitchResponse domain node switch response
type DomainNodeSwitchResponse struct {
	Status Status               `json:"status"`
	Data   DomainNodeSwitchData `json:"data"`
}

// DomainNodeSwitchData domain node switch data
type DomainNodeSwitchData struct {
	DomainID int `json:"domain_id"`
}

// DomainAccessModeSwitchRequest domain access mode switch request
type DomainAccessModeSwitchRequest struct {
	DomainID   int    `json:"domain_id"`
	AccessMode string `json:"access_mode"`
}

// DomainAccessModeSwitchResponse domain access mode switch response
type DomainAccessModeSwitchResponse struct {
	Status Status                     `json:"status"`
	Data   DomainAccessModeSwitchData `json:"data"`
}

// DomainAccessModeSwitchData domain access mode switch data
type DomainAccessModeSwitchData struct {
	DomainID int `json:"domain_id"`
}

// AccessProgressResponse access progress response
type AccessProgressResponse struct {
	Status Status             `json:"status"`
	Data   AccessProgressData `json:"data"`
}

// AccessProgressData access progress data
type AccessProgressData struct {
	Progress []ProgressInfo `json:"progress"`
}

// ProgressInfo progress information
type ProgressInfo struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// DomainBaseSettingsUpdateRequest domain base settings update request
type DomainBaseSettingsUpdateRequest struct {
	DomainID int                     `json:"domain_id"`
	Value    DomainBaseSettingsValue `json:"value"`
}

// DomainBaseSettingsValue domain base settings value
type DomainBaseSettingsValue struct {
	ProxyHost      *ProxyHostConfig      `json:"proxy_host,omitempty"`
	ProxySNI       *ProxySNIConfig       `json:"proxy_sni,omitempty"`
	DomainRedirect *DomainRedirectConfig `json:"domain_redirect,omitempty"`
}

// ProxyHostConfig proxy host configuration
type ProxyHostConfig struct {
	ProxyHost     string `json:"proxy_host,omitempty"`
	ProxyHostType string `json:"proxy_host_type,omitempty"`
}

// ProxySNIConfig proxy SNI configuration
type ProxySNIConfig struct {
	ProxySNI string `json:"proxy_sni,omitempty"`
	Status   string `json:"status,omitempty"`
}

// DomainRedirectConfig domain redirect configuration
type DomainRedirectConfig struct {
	Status   string `json:"status,omitempty"`
	JumpTo   string `json:"jump_to,omitempty"`
	JumpType string `json:"jump_type,omitempty"`
}

// DomainBaseSettingsUpdateResponse domain base settings update response
type DomainBaseSettingsUpdateResponse struct {
	Status Status                       `json:"status"`
	Data   DomainBaseSettingsUpdateData `json:"data"`
}

// DomainBaseSettingsUpdateData domain base settings update data
type DomainBaseSettingsUpdateData struct {
	DomainID int `json:"domain_id"`
}

// DomainBaseSettingsGetRequest domain base settings get request
type DomainBaseSettingsGetRequest struct {
	DomainID int `json:"domain_id"`
}

// DomainBaseSettingsGetResponse domain base settings get response
type DomainBaseSettingsGetResponse struct {
	Status Status                    `json:"status"`
	Data   DomainBaseSettingsGetData `json:"data"`
}

// DomainBaseSettingsGetData domain base settings get data
type DomainBaseSettingsGetData struct {
	DomainID       int                  `json:"domain_id"`
	ProxyHost      ProxyHostConfig      `json:"proxy_host"`
	ProxySNI       ProxySNIConfig       `json:"proxy_sni"`
	DomainRedirect DomainRedirectConfig `json:"domain_redirect"`
}

// BriefDomainListRequest brief domain list request
type BriefDomainListRequest struct {
	IDs []int `json:"ids,omitempty"`
}

// BriefDomainListResponse brief domain list response
type BriefDomainListResponse struct {
	Status Status              `json:"status"`
	Data   BriefDomainListData `json:"data"`
}

// BriefDomainListData brief domain list data
type BriefDomainListData struct {
	Total int               `json:"total"`
	List  []BriefDomainInfo `json:"list"`
}

// BriefDomainInfo brief domain information
type BriefDomainInfo struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
}

// DomainTemplatesRequest domain templates request
type DomainTemplatesRequest struct {
	DomainID int `json:"domain_id"`
}

// DomainTemplatesResponse domain templates response
type DomainTemplatesResponse struct {
	Status Status              `json:"status"`
	Data   DomainTemplatesData `json:"data"`
}

// DomainTemplatesData domain templates data
type DomainTemplatesData struct {
	DomainID        int                  `json:"domain_id"`
	BindedTemplates []BindedTemplateInfo `json:"binded_templates"`
}

// BindedTemplateInfo binded template information
type BindedTemplateInfo struct {
	BusinessID   int    `json:"business_id"`
	BusinessType string `json:"business_type"`
	AppType      string `json:"app_type"`
	Name         string `json:"name"`
}

// AccessInfoDownloadRequest access info download request
type AccessInfoDownloadRequest struct {
	DomainInfos []DomainInfoItem `json:"domain_infos"`
}

// DomainInfoItem domain info item
type DomainInfoItem struct {
	Domain     string `json:"domain"`
	DataKey    string `json:"data_key"`
	BizMainKey string `json:"biz_main_key"`
}

// AccessInfoDownloadResponse access info download response
type AccessInfoDownloadResponse struct {
	Status Status                 `json:"status"`
	Data   AccessInfoDownloadData `json:"data"`
}

// AccessInfoDownloadData access info download data
type AccessInfoDownloadData struct {
	Hash    string `json:"hash"`
	Key     string `json:"key"`
	RealURL string `json:"real_url"`
}
