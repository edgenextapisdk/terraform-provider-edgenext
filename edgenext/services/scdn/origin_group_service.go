package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Origin Group Management Methods
// ============================================================================

// ListOriginGroups lists origin groups
func (s *ScdnService) ListOriginGroups(req OriginGroupListRequest) (*OriginGroupListResponse, error) {
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
	if req.PageSize > 0 {
		scdnReq.Query["page_size"] = req.PageSize
	}
	if req.Name != "" {
		scdnReq.Query["name"] = req.Name
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointOriginGroups, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list origin groups: %w", err)
	}

	// Check business status code
	// Code -27 means no data, which is acceptable for list operations
	if scdnResp.Status.Code != 1 && scdnResp.Status.Code != -27 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &OriginGroupListResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	// Parse data field
	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var listData OriginGroupListData
		if err := json.Unmarshal(dataBytes, &listData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal origin group list data: %w", err)
		}
		response.Data = listData
	}

	return response, nil
}

// GetOriginGroupDetail gets origin group detail
func (s *ScdnService) GetOriginGroupDetail(req OriginGroupDetailRequest) (*OriginGroupDetailResponse, error) {
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

	scdnReq.Query["id"] = req.ID

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointOriginGroupsDetail, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get origin group detail: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &OriginGroupDetailResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	// Parse data field
	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var detailData OriginGroupDetailData
		if err := json.Unmarshal(dataBytes, &detailData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal origin group detail data: %w", err)
		}
		response.Data = detailData
	}

	return response, nil
}

// CreateOriginGroup creates an origin group
func (s *ScdnService) CreateOriginGroup(req OriginGroupCreateRequest) (*OriginGroupCreateResponse, error) {
	ctx := context.Background()

	var response OriginGroupCreateResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointOriginGroups, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create origin group: %w", err)
	}

	return &response, nil
}

// UpdateOriginGroup updates an origin group
func (s *ScdnService) UpdateOriginGroup(req OriginGroupUpdateRequest) (*OriginGroupUpdateResponse, error) {
	ctx := context.Background()

	var response OriginGroupUpdateResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointOriginGroups, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update origin group: %w", err)
	}

	return &response, nil
}

// DeleteOriginGroups deletes origin groups
func (s *ScdnService) DeleteOriginGroups(req OriginGroupDeleteRequest) (*OriginGroupDeleteResponse, error) {
	ctx := context.Background()

	var response OriginGroupDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointOriginGroups, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete origin groups: %w", err)
	}

	return &response, nil
}

// BindOriginGroupToDomains binds origin group to domains
func (s *ScdnService) BindOriginGroupToDomains(req OriginGroupBindDomainsRequest) (*OriginGroupBindDomainsResponse, error) {
	ctx := context.Background()

	var response OriginGroupBindDomainsResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointOriginGroupsBindDomains, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to bind origin group to domains: %w", err)
	}

	return &response, nil
}

// GetAllOriginGroups gets all origin groups
func (s *ScdnService) GetAllOriginGroups(req OriginGroupAllRequest) (*OriginGroupAllResponse, error) {
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

	scdnReq.Query["protect_status"] = req.ProtectStatus

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointOriginGroupsAll, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get all origin groups: %w", err)
	}

	// Check business status code
	// Code -27 means no data, which is acceptable for list operations
	if scdnResp.Status.Code != 1 && scdnResp.Status.Code != -27 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &OriginGroupAllResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	// Parse data field
	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var listData OriginGroupListData
		if err := json.Unmarshal(dataBytes, &listData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal origin group list data: %w", err)
		}
		response.Data = listData
	} else if scdnResp.Status.Code == -27 {
		// If code is -27, return empty list
		response.Data = OriginGroupListData{
			Total: 0,
			List:  []OriginGroupInfo{},
		}
	}

	return response, nil
}

// CopyOriginGroupToDomain copies origin group to domain
func (s *ScdnService) CopyOriginGroupToDomain(req OriginGroupCopyRequest) (*OriginGroupCopyResponse, error) {
	ctx := context.Background()

	var response OriginGroupCopyResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointOriginGroupsCopy, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to copy origin group to domain: %w", err)
	}

	return &response, nil
}

// GetOriginGroupBindHistory gets latest bind history
func (s *ScdnService) GetOriginGroupBindHistory(req OriginGroupBindHistoryRequest) (*OriginGroupBindHistoryResponse, error) {
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

	scdnReq.Query["origin_group_id"] = req.OriginGroupID

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointOriginGroupsBindHistory, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get origin group bind history: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &OriginGroupBindHistoryResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	// Parse data field
	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var historyData OriginGroupBindHistoryData
		if err := json.Unmarshal(dataBytes, &historyData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal origin group bind history data: %w", err)
		}
		response.Data = historyData
	}

	return response, nil
}
