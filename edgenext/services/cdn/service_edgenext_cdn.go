package cdn

import (
	"context"
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// CdnService CDN domain service
type CdnService struct {
	client *connectivity.Client
}

// NewCdnService creates a new CDN domain service instance
func NewCdnService(client *connectivity.Client) *CdnService {
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
	err := c.client.Post(ctx, "/v2/domain", req, &response)
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
	err := c.client.GetWithQuery(ctx, "/v2/domain", query, &response)
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
	err := c.client.GetWithQuery(ctx, "/v2/domain/list", query, &response)
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
	err := c.client.Delete(ctx, path, &response)
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
	Domain   string                 `json:"domain"`
	DomainID string                 `json:"domain_id"`
	Status   string                 `json:"status"`
	Config   map[string]interface{} `json:"config"` // Only includes fields actually returned by API
}

// SetDomainConfig sets domain configuration
func (c *CdnService) SetDomainConfig(domains string, config map[string]interface{}) (*DomainConfigResponse, error) {
	ctx := context.Background()

	requestBody := DomainConfigRequest{
		Domains: domains,
		Config:  config,
	}
	// log.Printf("[DEBUG] Set domain config request body: %+v", requestBody)
	var response DomainConfigResponse
	err := c.client.Post(ctx, "/v2/domain/config", requestBody, &response)
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
	err := c.client.Get(ctx, path, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query domain configuration: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query domain configuration: %s", response.Msg)
	}
	log.Printf("DEBUG config:%s", config)
	log.Printf("DEBUG response:%+v", response)

	return &response, nil
}

// DeleteDomainConfig deletes domain configuration
func (c *CdnService) DeleteDomainConfig(req DeleteDomainConfigRequest) error {
	ctx := context.Background()

	var response DeleteDomainConfigResponse
	err := c.client.DeleteWithBodyAndResult(ctx, "/v2/domain/config", req, &response)
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
	err := c.client.Post(ctx, "/v2/cache/refresh", requestBody, &response)
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
	err := c.client.GetWithQuery(ctx, "/v2/cache/refresh", query, &response)
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
	err := c.client.Post(ctx, "/v2/cache/prefetch", requestBody, &response)
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
	err := c.client.GetWithQuery(ctx, "/v2/cache/prefetch", query, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query file purge status: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("failed to query file purge status: %s", response.Msg)
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
