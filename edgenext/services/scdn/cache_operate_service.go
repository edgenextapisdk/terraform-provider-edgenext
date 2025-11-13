package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Cache Clean and Preheat Management Methods
// ============================================================================

// GetCacheCleanConfig gets cache clean configuration list
func (s *ScdnService) GetCacheCleanConfig(req CacheCleanGetConfigRequest) (*CacheCleanGetConfigResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointCacheCleanGetConfig, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache clean config: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CacheCleanGetConfigResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var configData CacheCleanGetConfigData
		if err := json.Unmarshal(dataBytes, &configData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cache clean config data: %w", err)
		}
		response.Data = configData
	}

	return response, nil
}

// SaveCacheCleanTask submits a cache clean task
func (s *ScdnService) SaveCacheCleanTask(req CacheCleanSaveRequest) (*CacheCleanSaveResponse, error) {
	ctx := context.Background()

	var response CacheCleanSaveResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointCacheCleanSave, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to save cache clean task: %w", err)
	}

	return &response, nil
}

// GetCacheCleanTaskList gets cache clean task list
func (s *ScdnService) GetCacheCleanTaskList(req CacheCleanTaskListRequest) (*CacheCleanTaskListResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	if req.Page > 0 {
		scdnReq.Query["page"] = req.Page
	}
	if req.PerPage > 0 {
		scdnReq.Query["per_page"] = req.PerPage
	}
	if req.StartTime != "" {
		scdnReq.Query["start_time"] = req.StartTime
	}
	if req.EndTime != "" {
		scdnReq.Query["end_time"] = req.EndTime
	}
	if req.Status != "" {
		scdnReq.Query["status"] = req.Status
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointCacheCleanTaskList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache clean task list: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CacheCleanTaskListResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var taskListData CacheCleanTaskListData
		if err := json.Unmarshal(dataBytes, &taskListData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cache clean task list data: %w", err)
		}
		response.Data = taskListData
	}

	return response, nil
}

// GetCacheCleanTaskDetail gets cache clean task detail
func (s *ScdnService) GetCacheCleanTaskDetail(req CacheCleanTaskDetailRequest) (*CacheCleanTaskDetailResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	if req.TaskID > 0 {
		scdnReq.Query["task_id"] = req.TaskID
	}
	if req.Page > 0 {
		scdnReq.Query["page"] = req.Page
	}
	if req.PerPage > 0 {
		scdnReq.Query["per_page"] = req.PerPage
	}
	if req.Result > 0 {
		scdnReq.Query["result"] = req.Result
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointCacheCleanTaskDetail, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache clean task detail: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CacheCleanTaskDetailResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var taskDetailData CacheCleanTaskDetailData
		if err := json.Unmarshal(dataBytes, &taskDetailData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cache clean task detail data: %w", err)
		}
		response.Data = taskDetailData
	}

	return response, nil
}

// GetCachePreheatTaskList gets preheat task list
func (s *ScdnService) GetCachePreheatTaskList(req CachePreheatTaskListRequest) (*CachePreheatTaskListResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Prepare request body (POST request with JSON body)
	requestBody := make(map[string]interface{})
	if req.Page > 0 {
		requestBody["page"] = req.Page
	}
	if req.PerPage > 0 {
		requestBody["per_page"] = req.PerPage
	}
	if req.PageSize > 0 {
		requestBody["page_size"] = req.PageSize
	}
	if req.Pagesize > 0 {
		requestBody["pagesize"] = req.Pagesize
	}
	if req.Size > 0 {
		requestBody["size"] = req.Size
	}
	if req.Status != "" {
		requestBody["status"] = req.Status
	}
	if req.URL != "" {
		requestBody["url"] = req.URL
	}

	scdnReq := &connectivity.ScdnRequest{
		Data: requestBody,
	}

	// Call SCDN API (POST request)
	scdnResp, err := scdnClient.Post(ctx, EndpointCachePreheatTaskList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get preheat task list: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CachePreheatTaskListResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var taskListData CachePreheatTaskListData
		if err := json.Unmarshal(dataBytes, &taskListData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal preheat task list data: %w", err)
		}
		response.Data = taskListData
	}

	return response, nil
}

// SaveCachePreheatTask submits a preheat task
func (s *ScdnService) SaveCachePreheatTask(req CachePreheatSaveRequest) (*CachePreheatSaveResponse, error) {
	ctx := context.Background()

	var response CachePreheatSaveResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointCachePreheatSave, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to save preheat task: %w", err)
	}

	return &response, nil
}
