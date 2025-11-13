package scdn

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ============================================================================
// Cache Clean and Preheat Types
// ============================================================================

// FlexibleString is a type that can unmarshal from both string and number
type FlexibleString string

// UnmarshalJSON implements json.Unmarshaler interface
func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case string:
		*fs = FlexibleString(val)
	case float64:
		*fs = FlexibleString(strconv.FormatFloat(val, 'f', -1, 64))
	case int:
		*fs = FlexibleString(strconv.Itoa(val))
	case int64:
		*fs = FlexibleString(strconv.FormatInt(val, 10))
	default:
		*fs = FlexibleString(fmt.Sprintf("%v", val))
	}
	return nil
}

// String returns the string value
func (fs FlexibleString) String() string {
	return string(fs)
}

// CacheCleanGetConfigRequest get cache clean config request
type CacheCleanGetConfigRequest struct {
	// No parameters needed for this endpoint
}

// CacheCleanGetConfigResponse get cache clean config response
type CacheCleanGetConfigResponse struct {
	Status Status                  `json:"status"`
	Data   CacheCleanGetConfigData `json:"data"`
}

// CacheCleanGetConfigData cache clean config data
type CacheCleanGetConfigData struct {
	ID         string   `json:"id"`         // Config ID
	Wholesite  []string `json:"wholesite"`  // Whole site config
	Specialurl []string `json:"specialurl"` // Special URL config
	Specialdir []string `json:"specialdir"` // Special directory config
}

// CacheCleanSaveRequest save cache clean task request
type CacheCleanSaveRequest struct {
	GroupID    int      `json:"group_id,omitempty"`   // Group ID, can refresh cache by group
	Protocol   string   `json:"protocol,omitempty"`   // Protocol: http/https; only valid when refreshing by group
	Port       string   `json:"port,omitempty"`       // Website port, only needed for special ports; only valid when refreshing by group
	Wholesite  []string `json:"wholesite,omitempty"`  // Whole site
	Specialurl []string `json:"specialurl,omitempty"` // Special URL
	Specialdir []string `json:"specialdir,omitempty"` // Special directory
}

// CacheCleanSaveResponse save cache clean task response
type CacheCleanSaveResponse struct {
	Status Status             `json:"status"`
	Data   CacheCleanSaveData `json:"data"`
}

// CacheCleanSaveData cache clean save data
type CacheCleanSaveData struct {
	Wholesite  []string `json:"wholesite"`  // Whole site config
	Specialurl []string `json:"specialurl"` // Special URL config
	Specialdir []string `json:"specialdir"` // Special directory config
}

// CacheCleanTaskListRequest get cache clean task list request
type CacheCleanTaskListRequest struct {
	Page      int    `json:"page,omitempty"`       // Page number
	PerPage   int    `json:"per_page,omitempty"`   // Items per page
	StartTime string `json:"start_time,omitempty"` // Start time, format: YYYY-MM-DD HH:II:SS
	EndTime   string `json:"end_time,omitempty"`   // End time, format: YYYY-MM-DD HH:II:SS
	Status    string `json:"status,omitempty"`     // Status: 1-executing, 2-completed
}

// CacheCleanTaskListResponse get cache clean task list response
type CacheCleanTaskListResponse struct {
	Status Status                 `json:"status"`
	Data   CacheCleanTaskListData `json:"data"`
}

// CacheCleanTaskListData cache clean task list data
type CacheCleanTaskListData struct {
	Total FlexibleString       `json:"total"` // Total number of tasks
	List  []CacheCleanTaskInfo `json:"list"`  // Task list
}

// CacheCleanTaskInfo cache clean task information
type CacheCleanTaskInfo struct {
	UserID           int            `json:"user_id"`            // User ID
	SubUserID        int            `json:"sub_user_id"`        // Sub user ID
	Status           *string        `json:"status"`             // Task status (can be null): "Failed", "Finished", etc.
	TaskID           int            `json:"task_id"`            // Task ID
	SubType          string         `json:"sub_type"`           // Task type: "Directory", "SubDomain", "URL"
	Total            FlexibleString `json:"total"`              // Total number of nodes
	Succeed          FlexibleString `json:"succeed"`            // Number of successful nodes
	Failed           FlexibleString `json:"failed"`             // Number of failed nodes
	Ongoing          FlexibleString `json:"ongoing"`            // Number of executing nodes
	CreatedAt        string         `json:"created_at"`         // Creation time, ISO 8601 format: "2025-09-04T02:12:22Z"
	OperatorUserName string         `json:"operator_user_name"` // Operator user name
}

// CacheCleanTaskDetailRequest get cache clean task detail request
type CacheCleanTaskDetailRequest struct {
	TaskID  int `json:"task_id,omitempty"`  // Task ID
	Page    int `json:"page,omitempty"`     // Page number
	PerPage int `json:"per_page,omitempty"` // Items per page
	Result  int `json:"result,omitempty"`   // Result: 1-success, 2-failed, 3-executing
}

// CacheCleanTaskDetailResponse get cache clean task detail response
type CacheCleanTaskDetailResponse struct {
	Status Status                   `json:"status"`
	Data   CacheCleanTaskDetailData `json:"data"`
}

// CacheCleanTaskDetailData cache clean task detail data
type CacheCleanTaskDetailData struct {
	Total FlexibleString             `json:"total"` // Total number of tasks
	List  []CacheCleanTaskDetailInfo `json:"list"`  // Task list
}

// CacheCleanTaskDetailInfo cache clean task detail information
type CacheCleanTaskDetailInfo struct {
	Result    string `json:"result"`              // Execution result: 成功/失败/执行中
	Message   string `json:"message"`             // Execution message
	CreatedAt string `json:"created_at"`          // Creation time, format: YYYY-MM-DD HH:II:SS
	UpdatedAt string `json:"updated_at"`          // Update time, format: YYYY-MM-DD HH:II:SS
	Directory string `json:"directory,omitempty"` // Directory, present when this task type
	Subdomain string `json:"subdomain,omitempty"` // Subdomain, present when this task type
	URL       string `json:"url,omitempty"`       // URL, present when this task type
}

// CachePreheatTaskListRequest get preheat task list request
type CachePreheatTaskListRequest struct {
	Page     int    `json:"page,omitempty"`      // Page number
	PerPage  int    `json:"per_page,omitempty"`  // Items per page
	PageSize int    `json:"page_size,omitempty"` // Page size (alternative to per_page)
	Pagesize int    `json:"pagesize,omitempty"`  // Page size (alternative to per_page)
	Size     int    `json:"size,omitempty"`      // Size (alternative to per_page)
	Status   string `json:"status,omitempty"`    // Status filter
	URL      string `json:"url,omitempty"`       // URL filter
}

// CachePreheatTaskListResponse get preheat task list response
type CachePreheatTaskListResponse struct {
	Status Status                   `json:"status"`
	Data   CachePreheatTaskListData `json:"data"`
}

// CachePreheatTaskListData preheat task list data
type CachePreheatTaskListData struct {
	Total     FlexibleString         `json:"total"`      // Total number of tasks
	StatusMap map[string]string      `json:"status_map"` // Status map: "1": "Prefetch waiting", etc.
	List      []CachePreheatTaskInfo `json:"list"`       // Task list
}

// CachePreheatTaskInfo preheat task information
type CachePreheatTaskInfo struct {
	ID               int    `json:"id"`                 // Task ID
	UserID           int    `json:"user_id"`            // User ID
	TimeCreate       string `json:"time_create"`        // Creation time: "2025-09-18 19:02:42"
	TimeUpdate       string `json:"time_update"`        // Update time: "2025-09-18 19:02:42"
	TaskID           int    `json:"task_id"`            // Task ID
	DomainID         int    `json:"domain_id"`          // Domain ID
	URL              string `json:"url"`                // URL
	Status           int    `json:"status"`             // Status: 1-Prefetch waiting, 2-Prefetch pending, 3-Prefetch successful, 4-Prefetch failed
	Total            int    `json:"total"`              // Total
	Weight           int    `json:"weight"`             // Weight
	SubUserID        int    `json:"sub_user_id"`        // Sub user ID
	UserName         string `json:"user_name"`          // User name
	StrategyID       int    `json:"strategy_id"`        // Strategy ID
	Strategy         string `json:"strategy"`           // Strategy
	OperatorUserName string `json:"operator_user_name"` // Operator user name
}

// CachePreheatSaveRequest save preheat task request
type CachePreheatSaveRequest struct {
	GroupID    int      `json:"group_id,omitempty"`    // Group ID, can refresh cache by group
	Protocol   string   `json:"protocol,omitempty"`    // Protocol: http/https; only valid when refreshing by group
	Port       string   `json:"port,omitempty"`        // Website port, only needed for special ports; only valid when refreshing by group
	PreheatURL []string `json:"preheat_url,omitempty"` // Preheat URLs
}

// CachePreheatSaveResponse save preheat task response
type CachePreheatSaveResponse struct {
	Status Status               `json:"status"`
	Data   CachePreheatSaveData `json:"data"`
}

// CachePreheatSaveData preheat save data
type CachePreheatSaveData struct {
	ErrorURL interface{} `json:"error_url"` // List of URLs with preheat errors (can be array, object, or null)
}

// GetErrorURLs extracts error URLs from the ErrorURL field
func (d *CachePreheatSaveData) GetErrorURLs() []string {
	if d.ErrorURL == nil {
		return []string{}
	}

	// Try to convert to []string
	switch v := d.ErrorURL.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	default:
		// If it's an object or other type, return empty array
		return []string{}
	}
}
