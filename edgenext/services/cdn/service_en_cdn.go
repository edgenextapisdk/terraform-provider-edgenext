package cdn

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// FlexibleInt is a custom type that can unmarshal both int and string values as int
type FlexibleInt int

// UnmarshalJSON implements json.Unmarshaler interface to handle both int and string values
func (fi *FlexibleInt) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as int first
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*fi = FlexibleInt(intVal)
		return nil
	}

	// If that fails, try to unmarshal as string and convert to int
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if strVal == "" {
			*fi = FlexibleInt(0)
			return nil
		}
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return fmt.Errorf("cannot convert string '%s' to int: %v", strVal, err)
		}
		*fi = FlexibleInt(intVal)
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleInt", string(data))
}

// MarshalJSON implements json.Marshaler interface to always marshal as int
func (fi FlexibleInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(fi))
}

// Int returns the int value
func (fi FlexibleInt) Int() int {
	return int(fi)
}

// CdnService CDN domain service
type CdnService struct {
	client *connectivity.EdgeNextClient
}

// NewCdnService creates a new CDN domain service instance
func NewCdnService(client *connectivity.EdgeNextClient) *CdnService {
	return &CdnService{client: client}
}

// DomainCreateRequest domain creation request
type DomainCreateRequest struct {
	Domain string                 `json:"domain"`
	Area   string                 `json:"area,omitempty"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// DomainResponse domain response
type DomainResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

// DomainData domain data
type DomainData struct {
	ID         string `json:"id"`
	Domain     string `json:"domain"`
	Type       string `json:"type"`       //Domain type(page(web)|download(download)|video_demand(video on demand)|dynamic(dynamic)|video_live(live streaming))
	Status     string `json:"status"`     //Domain status(serving=serving, suspend=suspended, deploying=deploying)
	IcpStatus  string `json:"icp_status"` //ICP filing status(checking=checking, yes=filed, no=not filed)
	IcpNum     string `json:"icp_num"`    //ICP filing number
	Area       string `json:"area"`       //global(global),mainland_china(China),outside_mainland_china(overseas),rim
	Cname      string `json:"cname"`
	CreateTime string `json:"create_time,omitempty"`
	UpdateTime string `json:"update_time,omitempty"`
	Https      int    `json:"https,omitempty"` //Whether HTTPS is enabled(0=disabled, 1=enabled)
}

type GetDomainResponse struct {
	Code int          `json:"code"`
	Data []DomainData `json:"data"`
	Msg  string       `json:"msg,omitempty"`
}

type DomainListRequest struct {
	PageNumber   int    `json:"page_number"`
	PageSize     int    `json:"page_size"`
	DomainStatus string `json:"domain_status"`
}

// DomainListResponse domain list response
type DomainListResponse struct {
	Code int                    `json:"code"`
	Data DomainListResponseData `json:"data"`
	Msg  string                 `json:"msg,omitempty"`
}

type DomainListResponseData struct {
	List        []DomainData `json:"list"`
	TotalNumber string       `json:"total_number"`
	PageNumber  string       `json:"page_number"`
	PageSize    int          `json:"page_size"`
}

type DeleteDomainResponse struct {
	Code int                        `json:"code"`
	Msg  string                     `json:"message,omitempty"`
	Data []DeleteDomainResponseData `json:"data"`
}

type DeleteDomainResponseData struct {
	ID         string `json:"id"`
	Domain     string `json:"domain"`
	Status     string `json:"status"`
	DeleteTime string `json:"delete_time"`
}

const (
	DomainStatusServing   = "serving"
	DomainStatusSuspended = "suspend"
	DomainStatusDeleted   = "deleted"
)

// CreateDomain creates a CDN domain
func (c *CdnService) CreateDomain(req DomainCreateRequest) (*DomainResponse, error) {
	ctx := context.Background()

	var response DomainResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create CDN domain: %w", err)
	}
	err = apiClient.Post(ctx, "/v2/domain", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create CDN domain: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to create CDN domain: %s", response.Msg)
	}

	return &response, nil
}

// GetDomain queries domain details
func (c *CdnService) GetDomain(domains string) (*GetDomainResponse, error) {
	ctx := context.Background()

	query := map[string]string{
		"domains": domains,
	}

	var response GetDomainResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to query domain details: %w", err)
	}
	err = apiClient.GetWithQuery(ctx, "/v2/domain", query, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query domain details: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query domain details: %s", response.Msg)
	}

	return &response, nil
}

// ListDomains queries domain list
func (c *CdnService) ListDomains(req DomainListRequest) (*DomainListResponse, error) {
	ctx := context.Background()

	query := map[string]string{
		"page_number":   fmt.Sprintf("%d", req.PageNumber),
		"page_size":     fmt.Sprintf("%d", req.PageSize),
		"domain_status": req.DomainStatus,
	}

	var response DomainListResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to query domain list: %w", err)
	}
	err = apiClient.GetWithQuery(ctx, "/v2/domain/list", query, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query domain list: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query domain list: %s", response.Msg)
	}

	return &response, nil
}

// DeleteDomain deletes a domain
func (c *CdnService) DeleteDomain(domains string) error {
	ctx := context.Background()

	var response DeleteDomainResponse
	path := fmt.Sprintf("/v2/domain?domains=%s", domains)
	// Delete domain API
	apiClient, err := c.client.APIClient()
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}
	err = apiClient.Delete(ctx, path, &response)
	if err != nil {
		return fmt.Errorf("failed to delete domain: %w", err)
	}

	if response.Code != 0 {
		return fmt.Errorf("failed to delete domain: %s", response.Msg)
	}

	return nil
}

// Helper method: check domain status
func (d *DomainData) IsServing() bool {
	return d.Status == DomainStatusServing
}

func (d *DomainData) IsDeleted() bool {
	return d.Status == DomainStatusDeleted
}

func (d *DomainData) IsSuspended() bool {
	return d.Status == DomainStatusSuspended
}

// Domain configuration related structs and methods

// DomainConfigRequest domain configuration request
type DomainConfigRequest struct {
	Domains string                 `json:"domains"`
	Config  map[string]interface{} `json:"config"`
}

// DomainConfigResponse domain configuration response
type DomainConfigResponse struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg,omitempty"`
}

// DeleteDomainConfigRequest delete domain configuration request
type DeleteDomainConfigRequest struct {
	Domains string   `json:"domains"`
	Config  []string `json:"config"`
}

// DeleteDomainConfigResponse delete domain configuration response
type DeleteDomainConfigResponse struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg,omitempty"`
}

// GetDomainConfigResponse query domain configuration response
type GetDomainConfigResponse struct {
	Code int                `json:"code"`
	Data []DomainConfigItem `json:"data"`
	Msg  string             `json:"msg,omitempty"`
}

type DomainConfigItem struct {
	Domain   string     `json:"domain"`
	DomainID string     `json:"domain_id"`
	Status   string     `json:"status"`
	Config   ConfigItem `json:"config"`
	// Config   map[string]interface{} `json:"config"` // Only includes fields actually returned by API
}

type ConfigItem struct {
	Origin            *OriginItem              `json:"origin,omitempty"`
	OriginHost        *OriginHostItem          `json:"origin_host,omitempty"`
	CacheRule         []*CacheRuleItem         `json:"cache_rule,omitempty"`
	CacheRuleList     []*CacheRuleListItem     `json:"cache_rule_list,omitempty"`
	Referer           *RefererItem             `json:"referer,omitempty"`
	IPBlackList       *IPBlackListItem         `json:"ip_black_list,omitempty"`
	IPWhiteList       *IPWhiteListItem         `json:"ip_white_list,omitempty"`
	AddResponseHead   *AddResponseHeadItem     `json:"add_response_head,omitempty"`
	AddBackSourceHead []*AddBackSourceHeadItem `json:"add_back_source_head,omitempty"`
	HTTPS             *HTTPSItem               `json:"https,omitempty"`
	CompressResponse  *CompressResponseItem    `json:"compress_response,omitempty"`
	SpeedLimit        []*SpeedLimitItem        `json:"speed_limit,omitempty"`
	RateLimit         *RateLimitItem           `json:"rate_limit,omitempty"`
	CacheShare        *CacheShareItem          `json:"cache_share,omitempty"`
	HeadControl       *HeadControlItem         `json:"head_control,omitempty"`
	Timeout           *TimeoutItem             `json:"timeout,omitempty"`
	ConnectTimeout    *ConnectTimeoutItem      `json:"connect_timeout,omitempty"`
	DenyURL           *DenyURLItem             `json:"deny_url,omitempty"`
}

func (c *ConfigItem) UnmarshalJSON(data []byte) error {
	// If it's an empty array, ignore directly
	if string(data) == "[]" {
		return nil
	}
	// Otherwise parse as object
	type Alias ConfigItem
	return json.Unmarshal(data, (*Alias)(c))
}

// OriginItem origin configuration
type OriginItem struct {
	DefaultMaster string      `json:"default_master,omitempty"`
	DefaultSlave  string      `json:"default_slave,omitempty"`
	OriginMode    string      `json:"origin_mode,omitempty"`
	OriHttps      string      `json:"ori_https,omitempty"`
	Port          FlexibleInt `json:"port,omitempty"`
}

// OriginHostItem origin HOST configuration
type OriginHostItem struct {
	Host string `json:"host,omitempty"`
}

// CacheRuleItem cache rule configuration
type CacheRuleItem struct {
	Type          int    `json:"type,omitempty"`
	Pattern       string `json:"pattern,omitempty"`
	Time          int    `json:"time,omitempty"`
	TimeUnit      string `json:"timeunit,omitempty"`
	IgnoreNoCache string `json:"ignore_no_cache,omitempty"`
	IgnoreExpired string `json:"ignore_expired,omitempty"`
	IgnoreQuery   string `json:"ignore_query,omitempty"`
}

// CacheRuleListItem new cache rule configuration
type CacheRuleListItem struct {
	MatchMethod          string `json:"match_method,omitempty"`
	Pattern              string `json:"pattern,omitempty"`
	CaseIgnore           string `json:"case_ignore,omitempty"`
	Expire               int    `json:"expire,omitempty"`
	ExpireUnit           string `json:"expire_unit,omitempty"`
	IgnoreNoCacheHeaders string `json:"ignore_no_cache_headers,omitempty"`
	FollowExpired        string `json:"follow_expired,omitempty"`
	QueryParamsOp        string `json:"query_params_op,omitempty"`
	Priority             int    `json:"priority,omitempty"`
	QueryParamsOpWay     string `json:"query_params_op_way,omitempty"`
	QueryParamsOpWhen    string `json:"query_params_op_when,omitempty"`
	Params               string `json:"params,omitempty"`
}

// RefererItem Referer blacklist/whitelist configuration
type RefererItem struct {
	Type       int      `json:"type,omitempty"`
	List       []string `json:"list,omitempty"`
	AllowEmpty bool     `json:"allow_empty,omitempty"`
}

// IPBlackListItem IP blacklist configuration
type IPBlackListItem struct {
	List []string `json:"list,omitempty"`
	Mode string   `json:"mode,omitempty"`
}

// IPWhiteListItem IP whitelist configuration
type IPWhiteListItem struct {
	List []string `json:"list,omitempty"`
	Mode string   `json:"mode,omitempty"`
}

// AddResponseHeadItem response header configuration
type AddResponseHeadItem struct {
	Type string                `json:"type,omitempty"`
	List []*ResponseHeaderItem `json:"list,omitempty"`
}

// ResponseHeaderItem response header item
type ResponseHeaderItem struct {
	Name    string `json:"name,omitempty"`
	Value   string `json:"value,omitempty"`
	Cover   string `json:"cover,omitempty"`
	OnlyHit string `json:"only_hit,omitempty"`
}

// AddBackSourceHeadItem add back source request header
type AddBackSourceHeadItem struct {
	Name            string `json:"head_name,omitempty"`
	Value           string `json:"head_value,omitempty"`
	WriteWhenExists string `json:"write_when_exists,omitempty"` // Whether to overwrite when the same request header exists, default value is yes when not filled.
}

// HTTPSItem HTTPS configuration
type HTTPSItem struct {
	CertID      int      `json:"cert_id,omitempty"`
	HTTP2       string   `json:"http2,omitempty"`
	ForceHTTPS  string   `json:"force_https,omitempty"`
	OCSP        string   `json:"ocsp,omitempty"`
	SSLProtocol []string `json:"ssl_protocol,omitempty"`
}

// CompressResponseItem compress response configuration
type CompressResponseItem struct {
	ContentType []string `json:"content_type,omitempty"`
	MinSize     int      `json:"min_size,omitempty"`
	MinSizeUnit string   `json:"min_size_unit,omitempty"`
}

// SpeedLimitItem single link speed limit configuration
type SpeedLimitItem struct {
	Type      string `json:"type,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Speed     int    `json:"speed,omitempty"`
	BeginTime string `json:"begin_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Priority  int    `json:"priority,omitempty"`
}

// RateLimitItem rate limit configuration
type RateLimitItem struct {
	MaxRateCount     int    `json:"max_rate_count"`
	LeadingFlowCount int    `json:"leading_flow_count"`
	LeadingFlowUnit  string `json:"leading_flow_unit"`
	MaxRateUnit      string `json:"max_rate_unit"`
}

// CacheShareItem shared cache configuration
type CacheShareItem struct {
	ShareWay string `json:"share_way"`
	Domain   string `json:"domain"`
}

// HeadControlItem HTTP header control configuration
type HeadControlItem struct {
	List []*HeadControlItemData `json:"list,omitempty"`
}

type HeadControlItemData struct {
	Regex         string `json:"regex,omitempty"`
	HeadOp        string `json:"head_op,omitempty"`
	HeadDirection string `json:"head_direction,omitempty"`
	Value         string `json:"value,omitempty"`
	FunName       string `json:"fun_name,omitempty"`
	Key           string `json:"key,omitempty"`
}

// TimeoutItem origin read timeout configuration
type TimeoutItem struct {
	Time string `json:"time,omitempty"`
}

// ConnectTimeoutItem origin connection timeout configuration
type ConnectTimeoutItem struct {
	OriginConnectTimeout string `json:"origin_connect_timeout,omitempty"`
}

// DenyURLItem block illegal URL configuration
type DenyURLItem struct {
	URLs []string `json:"urls"`
}

// SetDomainConfig sets domain configuration
func (c *CdnService) SetDomainConfig(domains string, config map[string]interface{}) (*DomainConfigResponse, error) {
	ctx := context.Background()

	requestBody := DomainConfigRequest{
		Domains: domains,
		Config:  config,
	}

	var response DomainConfigResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to set domain config: %w", err)
	}
	err = apiClient.Post(ctx, "/v2/domain/config", requestBody, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to set domain config: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to set domain config: %s", response.Msg)
	}

	return &response, nil
}

// GetDomainConfig queries domain configuration
func (c *CdnService) GetDomainConfig(domains string, config []string) (*GetDomainConfigResponse, error) {
	ctx := context.Background()

	// Build query path
	path := "/v2/domain/config?domains=" + domains

	// Add config[] query parameter for each configuration item
	for _, configItem := range config {
		path += "&config[]=" + configItem
	}

	var response GetDomainConfigResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to query domain configuration: %w", err)
	}
	err = apiClient.Get(ctx, path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query domain configuration: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query domain configuration: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query domain configuration: %s", response.Msg)
	}

	return &response, nil
}

// DeleteDomainConfig deletes domain configuration
func (c *CdnService) DeleteDomainConfig(req DeleteDomainConfigRequest) error {
	ctx := context.Background()

	var response DeleteDomainConfigResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return fmt.Errorf("failed to delete domain configuration: %w", err)
	}
	err = apiClient.DeleteWithBodyAndResult(ctx, "/v2/domain/config", req, &response)
	if err != nil {
		return fmt.Errorf("failed to delete domain configuration: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete domain configuration: %w", err)
	}

	if response.Code != 0 {
		return fmt.Errorf("failed to delete domain configuration: %s", response.Msg)
	}

	return nil
}

// Cache refresh related structs and methods

// CacheRefreshRequest cache refresh request
type CacheRefreshRequest struct {
	URLs []string `json:"urls"`
	Type string   `json:"type"`
}

// CacheRefreshResponse cache refresh response
type CacheRefreshResponse struct {
	Code int              `json:"code"`
	Data CacheRefreshData `json:"data"`
	Msg  string           `json:"msg,omitempty"`
}

// CacheRefreshData cache refresh data
type CacheRefreshData struct {
	TaskID  string   `json:"task_id"`
	Count   int      `json:"count"`
	ErrURLs []string `json:"err_urls,omitempty"`
}

// CacheRefreshQueryRequest cache refresh query request
type CacheRefreshQueryRequest struct {
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	URL        string `json:"url,omitempty"`
	PageNumber string `json:"page_number,omitempty"`
	PageSize   string `json:"page_size,omitempty"`
	TaskID     int    `json:"task_id,omitempty"`
}

// CacheRefreshQueryResponse cache refresh query response
type CacheRefreshQueryResponse struct {
	Code int                   `json:"code"`
	Msg  string                `json:"message,omitempty"`
	Data CacheRefreshQueryData `json:"data"`
}

// CacheRefreshQueryData cache refresh query data
type CacheRefreshQueryData struct {
	Total      int                     `json:"total"`
	PageNumber int                     `json:"page_number"`
	List       []CacheRefreshQueryItem `json:"list"`
}

// CacheRefreshQueryItem cache refresh query item
type CacheRefreshQueryItem struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
	CompleteTime string `json:"complete_time,omitempty"`
}

// RefreshStatus refresh status constants
const (
	RefreshStatusCompleted  = "completed"  // Completed
	RefreshStatusWaiting    = "waiting"    // Waiting
	RefreshStatusProcessing = "processing" // Processing
	RefreshStatusFailed     = "failed"     // Processing failed
)

// RefreshType refresh type constants
const (
	RefreshTypeURL = "url" // URL refresh
	RefreshTypeDir = "dir" // Directory refresh
)

// CacheRefresh cache refresh
func (c *CdnService) CacheRefresh(urls []string, refreshType string) (*CacheRefreshResponse, error) {
	ctx := context.Background()

	requestBody := CacheRefreshRequest{
		URLs: urls,
		Type: refreshType,
	}

	var response CacheRefreshResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("cache refresh failed: %w", err)
	}
	err = apiClient.Post(ctx, "/v2/cache/refresh", requestBody, &response)
	if err != nil {
		return nil, fmt.Errorf("cache refresh failed: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("cache refresh failed: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("cache refresh failed: %s", response.Msg)
	}

	return &response, nil
}

// QueryCacheRefresh queries cache refresh status
func (c *CdnService) QueryCacheRefresh(req CacheRefreshQueryRequest) (*CacheRefreshQueryResponse, error) {
	ctx := context.Background()

	query := make(map[string]string)

	// Build query parameters based on query method
	if req.TaskID > 0 {
		// Query by task ID
		query["task_id"] = fmt.Sprintf("%d", req.TaskID)
	} else {
		// Query by time range
		if req.StartTime != "" {
			query["start_time"] = req.StartTime
		}
		if req.EndTime != "" {
			query["end_time"] = req.EndTime
		}
		if req.URL != "" {
			query["url"] = req.URL
		}
		if req.PageNumber != "" {
			query["page_number"] = req.PageNumber
		}
		if req.PageSize != "" {
			query["page_size"] = req.PageSize
		}
	}

	var response CacheRefreshQueryResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to query cache refresh status: %w", err)
	}
	err = apiClient.GetWithQuery(ctx, "/v2/cache/refresh", query, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query cache refresh status: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query cache refresh status: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query cache refresh status: %s", response.Msg)
	}

	return &response, nil
}

// QueryCacheRefreshByTaskID queries cache refresh status by task ID
func (c *CdnService) QueryCacheRefreshByTaskID(taskID int) (*CacheRefreshQueryResponse, error) {
	req := CacheRefreshQueryRequest{
		TaskID: taskID,
	}
	return c.QueryCacheRefresh(req)
}

// QueryCacheRefreshByTimeRange queries cache refresh status by time range
func (c *CdnService) QueryCacheRefreshByTimeRange(startTime, endTime, url, pageNumber, pageSize string) (*CacheRefreshQueryResponse, error) {
	req := CacheRefreshQueryRequest{
		StartTime:  startTime,
		EndTime:    endTime,
		URL:        url,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
	return c.QueryCacheRefresh(req)
}

// Helper method: check refresh status
func (item *CacheRefreshQueryItem) IsCompleted() bool {
	return item.Status == RefreshStatusCompleted
}

func (item *CacheRefreshQueryItem) IsWaiting() bool {
	return item.Status == RefreshStatusWaiting
}

func (item *CacheRefreshQueryItem) IsProcessing() bool {
	return item.Status == RefreshStatusProcessing
}

func (item *CacheRefreshQueryItem) IsFailed() bool {
	return item.Status == RefreshStatusFailed
}

// Helper method: build cache refresh request
func NewCacheRefreshRequest(urls []string, refreshType string) CacheRefreshRequest {
	return CacheRefreshRequest{
		URLs: urls,
		Type: refreshType,
	}
}

// File purge related structs and methods

// FilePurgeRequest file purge request
type FilePurgeRequest struct {
	URLs []string `json:"urls"`
}

// FilePurgeResponse file purge response
type FilePurgeResponse struct {
	Code int           `json:"code"`
	Data FilePurgeData `json:"data"`
	Msg  string        `json:"message,omitempty"`
}

// FilePurgeData file purge data
type FilePurgeData struct {
	TaskID  string   `json:"task_id"`
	Count   int      `json:"count"`
	ErrURLs []string `json:"err_urls,omitempty"`
}

// FilePurgeQueryRequest file purge query request
type FilePurgeQueryRequest struct {
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
	URL        string `json:"url,omitempty"`
	PageNumber string `json:"page_number,omitempty"`
	PageSize   string `json:"page_size,omitempty"`
	TaskID     int    `json:"task_id,omitempty"`
}

// FilePurgeQueryResponse file purge query response
type FilePurgeQueryResponse struct {
	Code int                `json:"code"`
	Msg  string             `json:"message,omitempty"`
	Data FilePurgeQueryData `json:"data"`
}

// FilePurgeQueryData file purge query data
type FilePurgeQueryData struct {
	Total      int                  `json:"total"`
	PageNumber int                  `json:"page_number"`
	List       []FilePurgeQueryItem `json:"list"`
}

// FilePurgeQueryItem file purge query item
type FilePurgeQueryItem struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	Status       string `json:"status"`
	CreateTime   string `json:"create_time"`
	CompleteTime string `json:"complete_time,omitempty"`
}

// PurgeStatus purge status constants
const (
	PurgeStatusCompleted  = "completed"  // Completed
	PurgeStatusWaiting    = "waiting"    // Waiting
	PurgeStatusProcessing = "processing" // Processing
	PurgeStatusFailed     = "failed"     // Processing failed
)

// FilePurge file purge
func (c *CdnService) FilePurge(urls []string) (*FilePurgeResponse, error) {
	ctx := context.Background()

	requestBody := FilePurgeRequest{
		URLs: urls,
	}

	var response FilePurgeResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("file purge failed: %w", err)
	}
	err = apiClient.Post(ctx, "/v2/cache/prefetch", requestBody, &response)
	if err != nil {
		return nil, fmt.Errorf("file purge failed: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("file purge failed: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("file purge failed: %s", response.Msg)
	}

	return &response, nil
}

// QueryFilePurge queries file purge status
func (c *CdnService) QueryFilePurge(req FilePurgeQueryRequest) (*FilePurgeQueryResponse, error) {
	ctx := context.Background()

	query := make(map[string]string)

	// Build query parameters based on query method
	if req.TaskID > 0 {
		// Query by task ID
		query["task_id"] = fmt.Sprintf("%d", req.TaskID)
	} else {
		// Query by time range
		if req.StartTime != "" {
			query["start_time"] = req.StartTime
		}
		if req.EndTime != "" {
			query["end_time"] = req.EndTime
		}
		if req.URL != "" {
			query["url"] = req.URL
		}
		if req.PageNumber != "" {
			query["page_number"] = req.PageNumber
		}
		if req.PageSize != "" {
			query["page_size"] = req.PageSize
		}
	}

	var response FilePurgeQueryResponse
	apiClient, err := c.client.APIClient()
	if err != nil {
		return nil, fmt.Errorf("failed to query file purge: %w", err)
	}
	err = apiClient.GetWithQuery(ctx, "/v2/cache/prefetch", query, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query file purge: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query file purge: %s", response.Msg)
	}

	return &response, nil
}

// QueryFilePurgeByTaskID queries file purge status by task ID
func (c *CdnService) QueryFilePurgeByTaskID(taskID int) (*FilePurgeQueryResponse, error) {
	req := FilePurgeQueryRequest{
		TaskID: taskID,
	}
	return c.QueryFilePurge(req)
}

// QueryFilePurgeByTimeRange queries file purge status by time range
func (c *CdnService) QueryFilePurgeByTimeRange(startTime, endTime, url, pageNumber, pageSize string) (*FilePurgeQueryResponse, error) {
	req := FilePurgeQueryRequest{
		StartTime:  startTime,
		EndTime:    endTime,
		URL:        url,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
	return c.QueryFilePurge(req)
}

// Helper method: check purge status
func (item *FilePurgeQueryItem) IsCompleted() bool {
	return item.Status == PurgeStatusCompleted
}

func (item *FilePurgeQueryItem) IsWaiting() bool {
	return item.Status == PurgeStatusWaiting
}

func (item *FilePurgeQueryItem) IsProcessing() bool {
	return item.Status == PurgeStatusProcessing
}

func (item *FilePurgeQueryItem) IsFailed() bool {
	return item.Status == PurgeStatusFailed
}
