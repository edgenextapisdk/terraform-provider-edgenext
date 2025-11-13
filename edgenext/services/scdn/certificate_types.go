package scdn

// ============================================================================
// Certificate Management Types
// ============================================================================

// CATextSaveRequest certificate text save request (add/update)
type CATextSaveRequest struct {
	ID          int    `json:"id,omitempty"`           // Certificate ID, omit for new certificate
	CAName      string `json:"ca_name"`                // Certificate name
	CACert      string `json:"ca_cert"`                // Certificate public key
	CAKey       string `json:"ca_key"`                 // Certificate private key
	ProductFlag string `json:"product_flag,omitempty"` // Product flag
}

// CATextSaveResponse certificate text save response
type CATextSaveResponse struct {
	Status Status         `json:"status"`
	Data   CATextSaveData `json:"data"`
}

// CATextSaveData certificate text save data
type CATextSaveData struct {
	ID   string `json:"id"`    // Certificate ID
	CASN string `json:"ca_sn"` // Certificate serial number
}

// CASelfListRequest certificate list request
type CASelfListRequest struct {
	Page          int    `json:"page,omitempty"`
	PerPage       int    `json:"per_page,omitempty"`
	Domain        string `json:"domain,omitempty"`
	ProductFlag   string `json:"product_flag,omitempty"`
	CAName        string `json:"ca_name,omitempty"`
	Binded        string `json:"binded,omitempty"`       // true/false
	ApplyStatus   string `json:"apply_status,omitempty"` // 1-申请中,2-已颁发,3-审核失败，4-上传托管
	Issuer        string `json:"issuer,omitempty"`
	ExpiryTime    string `json:"expiry_time,omitempty"`     // true/false/inno
	IsExactSearch string `json:"is_exact_search,omitempty"` // on/off
}

// CASelfListResponse certificate list response
type CASelfListResponse struct {
	Status Status         `json:"status"`
	Data   CASelfListData `json:"data"`
}

// CASelfListData certificate list data
type CASelfListData struct {
	Total      string            `json:"total"`
	IssuerList []string          `json:"issuer_list"`
	List       []CertificateInfo `json:"list"`
}

// CertificateInfo certificate information
type CertificateInfo struct {
	ID                              string      `json:"id"`
	MemberID                        string      `json:"member_id"`
	CAName                          string      `json:"ca_name"`
	Issuer                          string      `json:"issuer"`
	IssuerStartTime                 string      `json:"issuer_start_time"`
	IssuerExpiryTime                string      `json:"issuer_expiry_time"`
	IssuerExpiryTimeDesc            string      `json:"issuer_expiry_time_desc"`
	IssuerExpiryTimeAutoRenewStatus int         `json:"issuer_expiry_time_auto_renew_status"`
	RenewStatus                     string      `json:"renew_status"` // 1-默认, 2-续签中，3-续签失败, 4-续签成功
	Binded                          bool        `json:"binded"`
	CADomain                        []string    `json:"ca_domain"`
	ApplyStatus                     interface{} `json:"apply_status"`   // 1-申请中,2-已颁发,3-审核失败，4-上传托管
	CAType                          string      `json:"ca_type"`        // 1-上传,2-lets申购
	CATypeDomain                    string      `json:"ca_type_domain"` // 1-单域名,2-多域名,3-泛域名
	Code                            string      `json:"code"`
	Msg                             string      `json:"msg"`
}

// CASelfDetailRequest certificate detail request
type CASelfDetailRequest struct {
	ID int `json:"id"` // Certificate ID
}

// CASelfDetailResponse certificate detail response
type CASelfDetailResponse struct {
	Status Status                `json:"status"`
	Data   CertificateDetailInfo `json:"data"`
}

// CertificateDetailInfo certificate detail information
type CertificateDetailInfo struct {
	ID                              string   `json:"id"`
	CAID                            string   `json:"ca_id"`
	MemberID                        string   `json:"member_id"`
	CAName                          string   `json:"ca_name"`
	Issuer                          string   `json:"issuer"`
	IssuerStartTime                 string   `json:"issuer_start_time"`
	IssuerExpiryTime                string   `json:"issuer_expiry_time"`
	IssuerExpiryTimeDesc            string   `json:"issuer_expiry_time_desc"`
	IssuerExpiryTimeAutoRenewStatus int      `json:"issuer_expiry_time_auto_renew_status"`
	RenewStatus                     string   `json:"renew_status"`
	Binded                          bool     `json:"binded"`
	CADomain                        []string `json:"ca_domain"`
	ApplyStatus                     string   `json:"apply_status"`
	CAType                          string   `json:"ca_type"`
	CATypeDomain                    string   `json:"ca_type_domain"`
	Code                            string   `json:"code"`
	Msg                             string   `json:"msg"`
	CreatedAt                       string   `json:"created_at"`
	UpdatedAt                       string   `json:"updated_at"`
	IssuerOrganization              string   `json:"issuer_organization"`
	IssuerOrganizationElement       string   `json:"issuer_organization_element"`
	SerialNumber                    string   `json:"serial_number"`
	IssuerObject                    string   `json:"issuer_object"`
	UseOrganization                 string   `json:"use_organization"`
	UseOrganizationElement          string   `json:"use_organization_element"`
	City                            string   `json:"city"`
	Province                        string   `json:"province"`
	Country                         string   `json:"country"`
	AuthenticationUsableDomain      string   `json:"authentication_usable_domain"`
}

// CASelfDeleteRequest certificate delete request
type CASelfDeleteRequest struct {
	IDs         string `json:"ids"` // Certificate IDs, comma separated
	ProductFlag string `json:"product_flag,omitempty"`
}

// CASelfDeleteResponse certificate delete response
type CASelfDeleteResponse struct {
	Status Status           `json:"status"`
	Data   CASelfDeleteData `json:"data"`
}

// CASelfDeleteData certificate delete data
type CASelfDeleteData struct {
	Info string `json:"info"`
}

// CAEditNameRequest certificate edit name request
type CAEditNameRequest struct {
	ID          int    `json:"id"`
	CAName      string `json:"ca_name"`
	ProductFlag string `json:"product_flag,omitempty"`
}

// CAEditNameResponse certificate edit name response
type CAEditNameResponse struct {
	Status Status         `json:"status"`
	Data   CAEditNameData `json:"data"`
}

// CAEditNameData certificate edit name data
type CAEditNameData struct {
	Info string `json:"info"`
}

// CABatchListRequest batch certificate list by domains request
type CABatchListRequest struct {
	Domains []string `json:"domains"`
}

// CABatchListResponse batch certificate list response
type CABatchListResponse struct {
	Status Status            `json:"status"`
	Data   []CertificateInfo `json:"data"`
}

// CAApplyAddRequest certificate apply request
type CAApplyAddRequest struct {
	Domain []string `json:"domain"`
}

// CAApplyAddResponse certificate apply response
type CAApplyAddResponse struct {
	Status Status         `json:"status"`
	Data   CAApplyAddData `json:"data"`
}

// CAApplyAddData certificate apply data
type CAApplyAddData struct {
	CAIDDomains map[string]string `json:"ca_id_domains"` // domain_id: domain
	CAIDNames   map[string]string `json:"ca_id_names"`   // ca_id: ca_name
}

// CASelfExportRequest certificate export request
type CASelfExportRequest struct {
	ID          interface{} `json:"id"` // Certificate ID(s), can be string or array
	ProductFlag string      `json:"product_flag,omitempty"`
}

// CASelfExportResponse certificate export response
type CASelfExportResponse struct {
	Status Status             `json:"status"`
	Data   []CASelfExportData `json:"data"`
}

// CASelfExportData certificate export data
type CASelfExportData struct {
	Hash    string `json:"hash"`
	Key     string `json:"key"`
	RealURL string `json:"real_url"`
}
