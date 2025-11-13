package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Security Protection Service Methods
// ============================================================================

// ============================================================================
// DDoS Protection Methods
// ============================================================================

// GetDdosProtectionConfig gets DDoS protection configuration
func (s *ScdnService) GetDdosProtectionConfig(req DdosProtectionGetConfigRequest) (*DdosProtectionGetConfigResponse, error) {
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
	if len(req.Keys) > 0 {
		scdnReq.Query["keys"] = req.Keys
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointSecurityProtectionDdosConfigs, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get DDoS protection config: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &DdosProtectionGetConfigResponse{
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

		var configData DdosProtectionGetConfigData
		if err := json.Unmarshal(dataBytes, &configData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal DDoS protection config data: %w", err)
		}
		response.Data = configData
	}

	return response, nil
}

// UpdateDdosProtectionConfig updates DDoS protection configuration
func (s *ScdnService) UpdateDdosProtectionConfig(req DdosProtectionUpdateConfigRequest) (*DdosProtectionUpdateConfigResponse, error) {
	ctx := context.Background()

	var response DdosProtectionUpdateConfigResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointSecurityProtectionDdosConfigs, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update DDoS protection config: %w", err)
	}

	return &response, nil
}

// ============================================================================
// WAF Rule Config Methods
// ============================================================================

// GetWafRuleConfig gets WAF rule configuration
func (s *ScdnService) GetWafRuleConfig(req WafRuleConfigGetRequest) (*WafRuleConfigGetResponse, error) {
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
	if len(req.Keys) > 0 {
		scdnReq.Query["keys"] = req.Keys
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointSecurityProtectionWafRules, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get WAF rule config: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &WafRuleConfigGetResponse{
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

		var configData WafRuleConfigGetData
		if err := json.Unmarshal(dataBytes, &configData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal WAF rule config data: %w", err)
		}
		response.Data = configData
	}

	return response, nil
}

// UpdateWafRuleConfig updates WAF rule configuration
func (s *ScdnService) UpdateWafRuleConfig(req WafRuleConfigUpdateRequest) (*WafRuleConfigUpdateResponse, error) {
	ctx := context.Background()

	var response WafRuleConfigUpdateResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointSecurityProtectionWafRules, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update WAF rule config: %w", err)
	}

	return &response, nil
}

// ============================================================================
// Security Protection Template Methods
// ============================================================================

// GetMemberGlobalTemplate gets member global template
func (s *ScdnService) GetMemberGlobalTemplate() (*SecurityProtectionTemplateGetMemberGlobalResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Call SCDN API
	scdnReq := &connectivity.ScdnRequest{}
	scdnResp, err := scdnClient.Get(ctx, EndpointSecurityProtectionTemplateMemberGlobal, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get member global template: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	var response SecurityProtectionTemplateGetMemberGlobalResponse
	if scdnResp != nil {
		dataBytes, err := json.Marshal(scdnResp)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		if err := json.Unmarshal(dataBytes, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	}

	return &response, nil
}

// CreateSecurityProtectionTemplate creates a security protection template
func (s *ScdnService) CreateSecurityProtectionTemplate(req SecurityProtectionTemplateCreateRequest) (*SecurityProtectionTemplateCreateResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateCreateResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointSecurityProtectionTemplate, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create security protection template: %w", err)
	}

	return &response, nil
}

// CreateSecurityProtectionTemplateDomain creates a domain template
func (s *ScdnService) CreateSecurityProtectionTemplateDomain(req SecurityProtectionTemplateCreateDomainRequest) (*SecurityProtectionTemplateCreateDomainResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateCreateDomainResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointSecurityProtectionTemplateDomain, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain template: %w", err)
	}

	return &response, nil
}

// SearchSecurityProtectionTemplates searches template list
func (s *ScdnService) SearchSecurityProtectionTemplates(req SecurityProtectionTemplateSearchRequest) (*SecurityProtectionTemplateSearchResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateSearchResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointSecurityProtectionTemplateSearch, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to search security protection templates: %w", err)
	}

	return &response, nil
}

// GetSecurityProtectionTemplateBindDomains gets template bind domain list
func (s *ScdnService) GetSecurityProtectionTemplateBindDomains(req SecurityProtectionTemplateBindDomainSearchRequest) (*SecurityProtectionTemplateBindDomainSearchResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateBindDomainSearchResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointSecurityProtectionTemplateDomainBindSearch, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get template bind domains: %w", err)
	}

	return &response, nil
}

// BindSecurityProtectionTemplateDomain binds template domain
func (s *ScdnService) BindSecurityProtectionTemplateDomain(req SecurityProtectionTemplateBindDomainRequest) (*SecurityProtectionTemplateBindDomainResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateBindDomainResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointSecurityProtectionTemplateDomainBind, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to bind template domain: %w", err)
	}

	return &response, nil
}

// DeleteSecurityProtectionTemplate deletes a template
func (s *ScdnService) DeleteSecurityProtectionTemplate(req SecurityProtectionTemplateDeleteRequest) (*SecurityProtectionTemplateDeleteResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointSecurityProtectionTemplate, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete security protection template: %w", err)
	}

	return &response, nil
}

// BatchConfigSecurityProtectionTemplate batch configures templates
func (s *ScdnService) BatchConfigSecurityProtectionTemplate(req SecurityProtectionTemplateBatchConfigRequest) (*SecurityProtectionTemplateBatchConfigResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateBatchConfigResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointSecurityProtectionTemplateBatchConfig, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to batch config security protection template: %w", err)
	}

	return &response, nil
}

// GetSecurityProtectionTemplateUnboundDomains gets unbound template domain list
func (s *ScdnService) GetSecurityProtectionTemplateUnboundDomains(req SecurityProtectionTemplateUnboundDomainSearchRequest) (*SecurityProtectionTemplateUnboundDomainSearchResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateUnboundDomainSearchResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointSecurityProtectionTemplateDomainUnboundSearch, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get unbound template domains: %w", err)
	}

	return &response, nil
}

// EditSecurityProtectionTemplate edits a template
func (s *ScdnService) EditSecurityProtectionTemplate(req SecurityProtectionTemplateEditRequest) (*SecurityProtectionTemplateEditResponse, error) {
	ctx := context.Background()

	var response SecurityProtectionTemplateEditResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointSecurityProtectionTemplate, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to edit security protection template: %w", err)
	}

	return &response, nil
}

// ============================================================================
// Security Protection Iota Methods
// ============================================================================

// GetSecurityProtectionIota gets iota enum values
func (s *ScdnService) GetSecurityProtectionIota() (*SecurityProtectionIotaResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Call SCDN API
	scdnReq := &connectivity.ScdnRequest{}
	scdnResp, err := scdnClient.Get(ctx, EndpointSecurityProtectionIota, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get security protection iota: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	var response SecurityProtectionIotaResponse
	if scdnResp != nil {
		dataBytes, err := json.Marshal(scdnResp)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		if err := json.Unmarshal(dataBytes, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
		}
	}

	return &response, nil
}
