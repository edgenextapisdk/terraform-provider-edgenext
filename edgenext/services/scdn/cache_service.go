package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Cache Rule Management Methods
// ============================================================================

// GetCacheRules gets cache rules list
func (s *ScdnService) GetCacheRules(req CacheRuleGetRulesRequest) (*CacheRuleGetRulesResponse, error) {
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

	scdnReq.Query["business_id"] = req.BusinessID
	scdnReq.Query["business_type"] = req.BusinessType
	if req.Page > 0 {
		scdnReq.Query["page"] = req.Page
	}
	if req.PageSize > 0 {
		scdnReq.Query["page_size"] = req.PageSize
	}
	// Convert rule_id to id query parameter for API
	if req.ID > 0 {
		scdnReq.Query["id"] = req.ID
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointCacheRules, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache rules: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CacheRuleGetRulesResponse{
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

		var rulesData CacheRuleGetRulesData
		if err := json.Unmarshal(dataBytes, &rulesData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cache rules data: %w", err)
		}
		response.Data = rulesData
	}

	return response, nil
}

// CreateCacheRule creates a cache rule
func (s *ScdnService) CreateCacheRule(req CacheRuleCreateRequest) (*CacheRuleCreateResponse, error) {
	ctx := context.Background()

	var response CacheRuleCreateResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointCacheRules, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache rule: %w", err)
	}

	return &response, nil
}

// UpdateCacheRule updates cache rule name/remark
func (s *ScdnService) UpdateCacheRule(req CacheRuleUpdateRequest) (*CacheRuleUpdateResponse, error) {
	ctx := context.Background()

	var response CacheRuleUpdateResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointCacheRule, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update cache rule: %w", err)
	}

	return &response, nil
}

// UpdateCacheRuleConfig updates cache rule configuration
func (s *ScdnService) UpdateCacheRuleConfig(req CacheRuleUpdateConfigRequest) (*CacheRuleUpdateConfigResponse, error) {
	ctx := context.Background()

	var response CacheRuleUpdateConfigResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointCacheRuleConf, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update cache rule config: %w", err)
	}

	return &response, nil
}

// UpdateCacheRuleStatus updates cache rule status (enable/disable)
func (s *ScdnService) UpdateCacheRuleStatus(req CacheRuleUpdateStatusRequest) (*CacheRuleUpdateStatusResponse, error) {
	ctx := context.Background()

	var response CacheRuleUpdateStatusResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointCacheRuleStatus, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update cache rule status: %w", err)
	}

	return &response, nil
}

// SortCacheRules sorts cache rules
func (s *ScdnService) SortCacheRules(req CacheRuleSortRequest) (*CacheRuleSortResponse, error) {
	ctx := context.Background()

	var response CacheRuleSortResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointCacheRuleSort, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to sort cache rules: %w", err)
	}

	return &response, nil
}

// DeleteCacheRule deletes cache rules
func (s *ScdnService) DeleteCacheRule(req CacheRuleDeleteRequest) (*CacheRuleDeleteResponse, error) {
	ctx := context.Background()

	var response CacheRuleDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointCacheRule, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete cache rule: %w", err)
	}

	return &response, nil
}

// GetCacheGlobalConfig gets global cache configuration
func (s *ScdnService) GetCacheGlobalConfig() (*CacheGlobalConfigGetResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Call SCDN API (GET request with no body)
	scdnReq := &connectivity.ScdnRequest{}
	scdnResp, err := scdnClient.Get(ctx, EndpointCacheGlobalConf, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache global config: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CacheGlobalConfigGetResponse{
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

		var configData CacheGlobalConfigGetData
		if err := json.Unmarshal(dataBytes, &configData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cache global config data: %w", err)
		}
		response.Data = configData
	}

	return response, nil
}
