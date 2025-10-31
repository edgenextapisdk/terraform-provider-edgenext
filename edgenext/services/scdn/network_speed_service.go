package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Network Speed Management Methods
// ============================================================================

// GetNetworkSpeedConfig gets template configuration
func (s *ScdnService) GetNetworkSpeedConfig(req NetworkSpeedGetConfigRequest) (*NetworkSpeedGetConfigResponse, error) {
	ctx := context.Background()

	var response NetworkSpeedGetConfigResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointNetworkSpeedGetConfig, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get network speed config: %w", err)
	}

	return &response, nil
}

// UpdateNetworkSpeedConfig updates template configuration
func (s *ScdnService) UpdateNetworkSpeedConfig(req NetworkSpeedUpdateConfigRequest) (*NetworkSpeedUpdateConfigResponse, error) {
	ctx := context.Background()

	var response NetworkSpeedUpdateConfigResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointNetworkSpeedUpdateConfig, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update network speed config: %w", err)
	}

	return &response, nil
}

// GetNetworkSpeedRules gets rules list
func (s *ScdnService) GetNetworkSpeedRules(req NetworkSpeedGetRulesRequest) (*NetworkSpeedGetRulesResponse, error) {
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
	scdnReq.Query["config_group"] = req.ConfigGroup

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointNetworkSpeedRules, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get network speed rules: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &NetworkSpeedGetRulesResponse{
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

		var rulesData NetworkSpeedGetRulesData
		if err := json.Unmarshal(dataBytes, &rulesData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal network speed rules data: %w", err)
		}
		response.Data = rulesData
	}

	return response, nil
}

// CreateNetworkSpeedRule creates a rule
func (s *ScdnService) CreateNetworkSpeedRule(req NetworkSpeedCreateRuleRequest) (*NetworkSpeedCreateRuleResponse, error) {
	ctx := context.Background()

	var response NetworkSpeedCreateRuleResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointNetworkSpeedRules, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create network speed rule: %w", err)
	}

	return &response, nil
}

// DeleteNetworkSpeedRule deletes rules
func (s *ScdnService) DeleteNetworkSpeedRule(req NetworkSpeedDeleteRuleRequest) (*NetworkSpeedDeleteRuleResponse, error) {
	ctx := context.Background()

	var response NetworkSpeedDeleteRuleResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointNetworkSpeedRule, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete network speed rule: %w", err)
	}

	return &response, nil
}

// SortNetworkSpeedRules sorts rules
func (s *ScdnService) SortNetworkSpeedRules(req NetworkSpeedSortRulesRequest) (*NetworkSpeedSortRulesResponse, error) {
	ctx := context.Background()

	var response NetworkSpeedSortRulesResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointNetworkSpeedRuleSort, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to sort network speed rules: %w", err)
	}

	return &response, nil
}

// UpdateNetworkSpeedRule updates a rule
func (s *ScdnService) UpdateNetworkSpeedRule(req NetworkSpeedUpdateRuleRequest) (*NetworkSpeedUpdateRuleResponse, error) {
	ctx := context.Background()

	var response NetworkSpeedUpdateRuleResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointNetworkSpeedRule, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update network speed rule: %w", err)
	}

	return &response, nil
}
