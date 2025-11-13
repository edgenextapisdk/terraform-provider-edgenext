package scdn

// ============================================================================
// Log Download Management Types
// ============================================================================

// LogDownloadTaskListRequest log download task list request
type LogDownloadTaskListRequest struct {
	Status     int    `json:"status,omitempty"`      // Task status: 0-not started, 1-in progress, 2-completed, 3-failed, 4-cancelled
	TaskName   string `json:"task_name,omitempty"`   // Task name
	FileType   string `json:"file_type,omitempty"`   // File type: xls, csv, json
	DataSource string `json:"data_source,omitempty"` // Data source: ng, cc, waf
	Page       int    `json:"page,omitempty"`        // Page number, default: 1
	PerPage    int    `json:"per_page,omitempty"`    // Page size, default: 20
}

// LogDownloadTaskListResponse log download task list response
type LogDownloadTaskListResponse struct {
	Status Status                  `json:"status"`
	Data   LogDownloadTaskListData `json:"data"`
}

// LogDownloadTaskListData log download task list data
type LogDownloadTaskListData struct {
	Total interface{}           `json:"total"` // Can be string or number
	List  []LogDownloadTaskInfo `json:"list"`
}

// LogDownloadTaskInfo log download task information
type LogDownloadTaskInfo struct {
	MemberID         string      `json:"member_id"`
	SubUserID        string      `json:"sub_user_id,omitempty"`
	TaskName         string      `json:"task_name"`
	IsUseTemplate    string      `json:"is_use_template"`
	TemplateID       int         `json:"template_id"`
	DataSource       string      `json:"data_source"`
	DownloadFields   []string    `json:"download_fields"`
	SearchTerms      interface{} `json:"search_terms"` // Can be map[string][]string (response) or []LogDownloadSearchTerm (request)
	StartTime        string      `json:"start_time"`
	EndTime          string      `json:"end_time"`
	FileType         string      `json:"file_type"`
	Lang             string      `json:"lang,omitempty"`
	Status           interface{} `json:"status"` // Can be string or int
	Progress         int         `json:"progress,omitempty"`
	Result           interface{} `json:"result"`
	DownloadURL      interface{} `json:"download_url"`
	ResetDownloadURL string      `json:"reset_download_url,omitempty"`
	TaskStartTime    interface{} `json:"task_start_time"`
	TaskEndTime      interface{} `json:"task_end_time"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
	DataSourceName   string      `json:"data_source_name,omitempty"`
	TaskID           int         `json:"task_id"`
	TaskExpireDesc   string      `json:"task_expire_desc,omitempty"`
	UTCStartTime     string      `json:"utc_start_time,omitempty"`
	UTCEndTime       string      `json:"utc_end_time,omitempty"`
	UTCTaskStartTime string      `json:"utc_task_start_time,omitempty"`
	UTCTaskEndTime   string      `json:"utc_task_end_time,omitempty"`
	TemplateName     string      `json:"template_name,omitempty"`
}

// LogDownloadSearchTerm search term for log download
type LogDownloadSearchTerm struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// LogDownloadTaskAddRequest log download task add request
type LogDownloadTaskAddRequest struct {
	TaskName       string            `json:"task_name"`             // Task name
	IsUseTemplate  int               `json:"is_use_template"`       // Whether to use template: 0-no, 1-yes
	TemplateID     int               `json:"template_id,omitempty"` // Template ID
	DataSource     string            `json:"data_source"`           // Data source: ng, cc, waf
	DownloadFields []string          `json:"download_fields"`       // Download fields
	SearchTerms    map[string]string `json:"search_terms"`          // Search conditions (map format: key -> value)
	FileType       string            `json:"file_type"`             // File type: xls, csv, json
	StartTime      string            `json:"start_time"`            // Start time
	EndTime        string            `json:"end_time"`              // End time
	Lang           string            `json:"lang,omitempty"`        // Language: zh_CN, en_US, default: zh_CN
}

// LogDownloadTaskAddResponse log download task add response
type LogDownloadTaskAddResponse struct {
	Status Status                 `json:"status"`
	Data   LogDownloadTaskAddData `json:"data"`
}

// LogDownloadTaskAddData log download task add data
type LogDownloadTaskAddData struct {
	TaskID int `json:"task_id"`
}

// LogDownloadTaskCancelRequest log download task cancel request
type LogDownloadTaskCancelRequest struct {
	TaskID int `json:"task_id"` // Task ID
}

// LogDownloadTaskCancelResponse log download task cancel response
type LogDownloadTaskCancelResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTaskBatchCancelRequest log download task batch cancel request
type LogDownloadTaskBatchCancelRequest struct {
	TaskIDs []int `json:"task_ids"` // Task IDs
}

// LogDownloadTaskBatchCancelResponse log download task batch cancel response
type LogDownloadTaskBatchCancelResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTaskDeleteRequest log download task delete request
type LogDownloadTaskDeleteRequest struct {
	TaskID int `json:"task_id"` // Task ID
}

// LogDownloadTaskDeleteResponse log download task delete response
type LogDownloadTaskDeleteResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTaskBatchDeleteRequest log download task batch delete request
type LogDownloadTaskBatchDeleteRequest struct {
	TaskIDs []int `json:"task_ids"` // Task IDs
}

// LogDownloadTaskBatchDeleteResponse log download task batch delete response
type LogDownloadTaskBatchDeleteResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTaskRegenerateRequest log download task regenerate request
type LogDownloadTaskRegenerateRequest struct {
	TaskID int `json:"task_id"` // Task ID
}

// LogDownloadTaskRegenerateResponse log download task regenerate response
type LogDownloadTaskRegenerateResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadFieldsResponse log download fields response
type LogDownloadFieldsResponse struct {
	Status Status                            `json:"status"`
	Data   map[string]LogDownloadFieldConfig `json:"data"`
}

// LogDownloadFieldConfig log download field configuration for a data source
type LogDownloadFieldConfig struct {
	Name           string   `json:"name"`            // Data source name
	DownloadFields []string `json:"download_fields"` // Available download fields
	SearchTerms    []string `json:"search_terms"`    // Available search terms
}

// LogDownloadTemplateListRequest log download template list request
type LogDownloadTemplateListRequest struct {
	Status       int    `json:"status,omitempty"`        // Status: 1-enabled, 0-disabled
	GroupID      int    `json:"group_id,omitempty"`      // Group ID
	TemplateName string `json:"template_name,omitempty"` // Template name
	DataSource   string `json:"data_source,omitempty"`   // Data source: ng, cc, waf
	Page         int    `json:"page,omitempty"`          // Page number, default: 1
	PerPage      int    `json:"per_page,omitempty"`      // Page size, default: 20
}

// LogDownloadTemplateListResponse log download template list response
type LogDownloadTemplateListResponse struct {
	Status Status                      `json:"status"`
	Data   LogDownloadTemplateListData `json:"data"`
}

// LogDownloadTemplateListData log download template list data
type LogDownloadTemplateListData struct {
	Total interface{}               `json:"total"` // Can be string or number
	List  []LogDownloadTemplateInfo `json:"list"`
}

// LogDownloadTemplateInfo log download template information
type LogDownloadTemplateInfo struct {
	ID             string      `json:"id"`
	TemplateName   string      `json:"template_name"`
	MemberID       string      `json:"member_id"`
	GroupID        string      `json:"group_id"`
	DataSource     string      `json:"data_source"`
	Status         int         `json:"status"`
	DownloadFields []string    `json:"download_fields"`
	SearchTerms    interface{} `json:"search_terms"` // Can be map[string]string (response) or []LogDownloadSearchTerm (request)
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	TemplateID     int         `json:"template_id"`
	GroupName      string      `json:"group_name"`
}

// LogDownloadTemplateDomainListRequest log download template domain list request
type LogDownloadTemplateDomainListRequest struct {
	Domain string `json:"domain,omitempty"` // Domain name
}

// LogDownloadTemplateDomainListResponse log download template domain list response
type LogDownloadTemplateDomainListResponse struct {
	Status Status   `json:"status"`
	Total  int      `json:"total"`
	Data   []string `json:"data"`
}

// LogDownloadTemplateAddRequest log download template add request
type LogDownloadTemplateAddRequest struct {
	TemplateName     string            `json:"template_name"`      // Template name
	GroupName        string            `json:"group_name"`         // Group name
	GroupID          int               `json:"group_id"`           // Group ID
	DataSource       string            `json:"data_source"`        // Data source: ng, cc, waf
	Status           int               `json:"status"`             // Status: 1-enabled, 0-disabled, default: 1 (always included)
	DownloadFields   []string          `json:"download_fields"`    // Download fields
	SearchTerms      map[string]string `json:"search_terms"`       // Search conditions (map format: key -> value)
	DomainSelectType int               `json:"domain_select_type"` // Domain select type: 0-partial, 1-all
}

// LogDownloadTemplateAddResponse log download template add response
type LogDownloadTemplateAddResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTemplateSaveRequest log download template save (update) request
type LogDownloadTemplateSaveRequest struct {
	TemplateID       int               `json:"template_id"`                  // Template ID
	TemplateName     string            `json:"template_name"`                // Template name
	GroupName        string            `json:"group_name"`                   // Group name
	GroupID          int               `json:"group_id"`                     // Group ID
	DataSource       string            `json:"data_source"`                  // Data source: ng, cc, waf
	Status           int               `json:"status"`                       // Status: 1-enabled, 0-disabled, default: 1
	DownloadFields   []string          `json:"download_fields"`              // Download fields
	SearchTerms      map[string]string `json:"search_terms"`                 // Search conditions (map format: key -> value)
	DomainSelectType int               `json:"domain_select_type,omitempty"` // Domain select type: 0-partial, 1-all
}

// LogDownloadTemplateSaveResponse log download template save response
type LogDownloadTemplateSaveResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTemplateDeleteRequest log download template delete request
type LogDownloadTemplateDeleteRequest struct {
	TemplateID int `json:"template_id"` // Template ID
}

// LogDownloadTemplateDeleteResponse log download template delete response
type LogDownloadTemplateDeleteResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTemplateBatchDeleteRequest log download template batch delete request
type LogDownloadTemplateBatchDeleteRequest struct {
	TemplateIDs []int `json:"template_ids"` // Template IDs
}

// LogDownloadTemplateBatchDeleteResponse log download template batch delete response
type LogDownloadTemplateBatchDeleteResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTemplateChangeStatusRequest log download template change status request
type LogDownloadTemplateChangeStatusRequest struct {
	TemplateID int `json:"template_id"` // Template ID
	Status     int `json:"status"`      // Status: 1-enabled, 0-disabled
}

// LogDownloadTemplateChangeStatusResponse log download template change status response
type LogDownloadTemplateChangeStatusResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTemplateBatchChangeStatusRequest log download template batch change status request
type LogDownloadTemplateBatchChangeStatusRequest struct {
	TemplateIDs []int `json:"template_ids"` // Template IDs
	Status      int   `json:"status"`       // Status: 1-enabled, 0-disabled
}

// LogDownloadTemplateBatchChangeStatusResponse log download template batch change status response
type LogDownloadTemplateBatchChangeStatusResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

// LogDownloadTemplateAllResponse log download template all response (for adding tasks)
type LogDownloadTemplateAllResponse struct {
	Status Status                              `json:"status"`
	Data   map[string]LogDownloadTemplateGroup `json:"data"`
}

// LogDownloadTemplateGroup template group information
type LogDownloadTemplateGroup struct {
	GroupID   int                         `json:"group_id"`
	GroupName string                      `json:"group_name"`
	Templates []LogDownloadTemplateSimple `json:"templates"`
}

// LogDownloadTemplateSimple simple template information
type LogDownloadTemplateSimple struct {
	TemplateID     int         `json:"template_id"`
	TemplateName   string      `json:"template_name"`
	DownloadFields []string    `json:"download_fields"`
	SearchTerms    interface{} `json:"search_terms"` // Can be map[string]string (response) or []LogDownloadSearchTerm (request)
}

// LogDownloadTemplateGroupAllResponse log download template group all response
type LogDownloadTemplateGroupAllResponse struct {
	Status Status                         `json:"status"`
	Data   []LogDownloadTemplateGroupInfo `json:"data"`
}

// LogDownloadTemplateGroupInfo template group information
type LogDownloadTemplateGroupInfo struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}
