package scdn

// ============================================================================
// User IP Intelligence Management Types
// ============================================================================

// UserIpListRequest IP list request
type UserIpListRequest struct {
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	Domain  string `json:"domain,omitempty"`
}

// UserIpListResponse IP list response
type UserIpListResponse struct {
	Status Status         `json:"status"`
	Data   UserIpListData `json:"data"`
}

// UserIpListData IP list data
type UserIpListData struct {
	Total string       `json:"total"`
	List  []UserIpInfo `json:"list"`
}

// UserIpInfo IP list info
type UserIpInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberID    string `json:"member_id"`
	SubUserID   string `json:"sub_user_id"`
	ItemNum     string `json:"item_num"`
	Remark      string `json:"remark"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	FileKey     string `json:"file_key"`
	FileVersion string `json:"file_version"`
	WriteMmdb   string `json:"write_mmdb"` // 1:Not updated, 2:Updated, 3:Failed
	FileError   string `json:"file_error"`
	MqTime      string `json:"mq_time"`
	OwnerName   string `json:"owner_name"`
}

// UserIpAddRequest Add IP list request
type UserIpAddRequest struct {
	Name   string `json:"name"`
	Remark string `json:"remark,omitempty"`
}

// UserIpAddResponse Add IP list response
type UserIpAddResponse struct {
	Status Status        `json:"status"`
	Data   UserIpAddData `json:"data"`
}

// UserIpAddData Add IP list data
type UserIpAddData struct {
	ID string `json:"id"`
}

// UserIpSaveRequest Update IP list request
type UserIpSaveRequest struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark,omitempty"`
}

// UserIpSaveResponse Update IP list response
type UserIpSaveResponse struct {
	Status Status         `json:"status"`
	Data   UserIpSaveData `json:"data"` // ID number in doc response example, but checking consistency
}

// UserIpSaveData Update IP list data
type UserIpSaveData struct {
	ID int `json:"id"` // Doc says number
}

// UserIpDelRequest Delete IP list request
type UserIpDelRequest struct {
	IDs []string `json:"ids"`
}

// UserIpDelResponse Delete IP list response
type UserIpDelResponse struct {
	Status Status `json:"status"`
}

// UserIpItemListRequest IP item list request
type UserIpItemListRequest struct {
	UserIpID int    `json:"user_ip_id"`
	IP       string `json:"ip,omitempty"`
	Page     int    `json:"page,omitempty"`
	PerPage  int    `json:"per_page,omitempty"`
}

// UserIpItemListResponse IP item list response
type UserIpItemListResponse struct {
	Status Status             `json:"status"`
	Data   UserIpItemListData `json:"data"`
}

// UserIpItemListData IP item list data
type UserIpItemListData struct {
	Total int          `json:"total"`
	List  []UserIpItem `json:"list"`
}

// UserIpItem IP item info
type UserIpItem struct {
	ID              string `json:"_id"`
	IP              string `json:"ip"`
	Remark          string `json:"remark"`
	UserIpId        int    `json:"user_ip_id"`
	FormatCreatedAt string `json:"format_created_at"`
	FormatUpdatedAt string `json:"format_updated_at"`
}

// UserIpItemAddRequest Add IP item request
type UserIpItemAddRequest struct {
	UserIpID string `json:"user_ip_id"`
	IP       string `json:"ip"`
	Remark   string `json:"remark,omitempty"`
}

// UserIpItemAddResponse Add IP item response
type UserIpItemAddResponse struct {
	Status Status            `json:"status"`
	Data   UserIpItemAddData `json:"data"`
}

// UserIpItemAddData Add IP item data
type UserIpItemAddData struct {
	IDs []string `json:"ids"`
}

// UserIpItemEditRequest Edit IP item request
type UserIpItemEditRequest struct {
	ID       string `json:"_id"`
	UserIpID string `json:"user_ip_id"`
	IP       string `json:"ip"`
	Remark   string `json:"remark,omitempty"`
}

// UserIpItemEditResponse Edit IP item response
type UserIpItemEditResponse struct {
	Status Status             `json:"status"`
	Data   UserIpItemEditData `json:"data"`
}

// UserIpItemEditData Edit IP item data
type UserIpItemEditData struct {
	ID string `json:"id"`
}

// UserIpItemDelRequest Delete IP item request
type UserIpItemDelRequest struct {
	IDs      []string `json:"ids"`
	UserIpID string   `json:"user_ip_id"`
}

// UserIpItemDelResponse Delete IP item response
type UserIpItemDelResponse struct {
	Status Status            `json:"status"`
	Data   UserIpItemDelData `json:"data"`
}

// UserIpItemDelData Delete IP item data
type UserIpItemDelData struct {
	ID int `json:"id"`
}

// UserIpItemDelAllRequest Delete all IP items request
type UserIpItemDelAllRequest struct {
	UserIpID string `json:"user_ip_id"`
}

// UserIpItemDelAllResponse Delete all IP items response
type UserIpItemDelAllResponse struct {
	Status Status `json:"status"`
	Data   struct {
		ID int `json:"id"`
	} `json:"data"`
}

// UserIpCopyRequest Copy IP list request
type UserIpCopyRequest struct {
	UserIpID string `json:"user_ip_id"`
	Name     string `json:"name"`
	Remark   string `json:"remark,omitempty"`
}

// UserIpCopyResponse Copy IP list response
type UserIpCopyResponse struct {
	Status Status `json:"status"`
	Data   struct {
		ID int `json:"id"`
	} `json:"data"`
}

// UserIpItemFileSaveResponse File save upload response
type UserIpItemFileSaveResponse struct {
	Status Status                 `json:"status"`
	Data   UserIpItemFileSaveData `json:"data"`
}

// UserIpItemFileSaveData File save upload data
type UserIpItemFileSaveData struct {
	IDs []string `json:"ids"`
}
