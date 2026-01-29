package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Rule Template Management Methods
// ============================================================================

// CreateRuleTemplate creates a new rule template
func (s *ScdnService) CreateRuleTemplate(req RuleTemplateCreateRequest) (*RuleTemplateCreateResponse, error) {
	ctx := context.Background()

	var response RuleTemplateCreateResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointRuleTemplates, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule template: %w", err)
	}

	return &response, nil
}

// UpdateRuleTemplate updates an existing rule template
func (s *ScdnService) UpdateRuleTemplate(req RuleTemplateUpdateRequest) (*RuleTemplateUpdateResponse, error) {
	ctx := context.Background()

	var response RuleTemplateUpdateResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointRuleTemplates, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update rule template: %w", err)
	}

	return &response, nil
}

// DeleteRuleTemplate deletes a rule template
func (s *ScdnService) DeleteRuleTemplate(req RuleTemplateDeleteRequest) (*RuleTemplateDeleteResponse, error) {
	ctx := context.Background()

	var response RuleTemplateDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointRuleTemplates, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete rule template: %w", err)
	}

	return &response, nil
}

// ListRuleTemplates lists rule templates with various filter options
func (s *ScdnService) ListRuleTemplates(req RuleTemplateListRequest) (*RuleTemplateListResponse, error) {
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
	if req.Domain != "" {
		scdnReq.Query["domain"] = req.Domain
	}
	if req.AppType != "" {
		scdnReq.Query["app_type"] = req.AppType
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointRuleTemplates, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list rule templates: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &RuleTemplateListResponse{
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

		var templateListData RuleTemplateListData
		if err := json.Unmarshal(dataBytes, &templateListData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rule template list data: %w", err)
		}
		response.Data = templateListData
	}

	return response, nil
}

// BindRuleTemplateDomains binds domains to a rule template
func (s *ScdnService) BindRuleTemplateDomains(req RuleTemplateBindDomainRequest) (*RuleTemplateBindDomainResponse, error) {
	ctx := context.Background()

	var response RuleTemplateBindDomainResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointRuleTemplatesBind, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to bind domains to rule template: %w", err)
	}

	return &response, nil
}

// UnbindRuleTemplateDomains unbinds domains from a rule template
func (s *ScdnService) UnbindRuleTemplateDomains(req RuleTemplateUnbindDomainRequest) (*RuleTemplateUnbindDomainResponse, error) {
	ctx := context.Background()

	var response RuleTemplateUnbindDomainResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointRuleTemplatesUnbind, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unbind domains from rule template: %w", err)
	}

	return &response, nil
}

// ListRuleTemplateDomains lists domains bound to a specific rule template
func (s *ScdnService) ListRuleTemplateDomains(req RuleTemplateListDomainsRequest) (*RuleTemplateListDomainsResponse, error) {
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
	scdnReq.Query["app_type"] = req.AppType

	if req.Page > 0 {
		scdnReq.Query["page"] = req.Page
	}
	if req.PageSize > 0 {
		scdnReq.Query["page_size"] = req.PageSize
	}
	if req.Domain != "" {
		scdnReq.Query["domain"] = req.Domain
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointRuleTemplatesDomains, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list rule template domains: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &RuleTemplateListDomainsResponse{
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

		var templateDomainsData RuleTemplateListDomainsData
		if err := json.Unmarshal(dataBytes, &templateDomainsData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal rule template domains data: %w", err)
		}
		response.Data = templateDomainsData
	}

	return response, nil
}

// SwitchDomainTemplate switches domains from their current template to a new template
func (s *ScdnService) SwitchDomainTemplate(req RuleTemplateSwitchDomainRequest) (*RuleTemplateSwitchDomainResponse, error) {
	ctx := context.Background()

	var response RuleTemplateSwitchDomainResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointRuleTemplatesSwitchDomain, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to switch domain template: %w", err)
	}

	return &response, nil
}
