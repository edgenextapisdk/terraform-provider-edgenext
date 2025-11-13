package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Domain Management Methods
// ============================================================================

// ListDomains lists domains with various filter options
func (s *ScdnService) ListDomains(req DomainListRequest) (*DomainListResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	if req.Page > 0 {
		scdnReq.Query["page"] = req.Page
	}
	if req.PageSize > 0 {
		scdnReq.Query["page_size"] = req.PageSize
	}
	if req.ID > 0 {
		scdnReq.Query["id"] = req.ID
	}
	if req.AccessProgress != "" {
		scdnReq.Query["access_progress"] = req.AccessProgress
	}
	if req.GroupID > 0 {
		scdnReq.Query["group_id"] = req.GroupID
	}
	if req.Domain != "" {
		scdnReq.Query["domain"] = req.Domain
	}
	if req.Remark != "" {
		scdnReq.Query["remark"] = req.Remark
	}
	if req.OriginIP != "" {
		scdnReq.Query["origin_ip"] = req.OriginIP
	}
	if req.CAStatus != "" {
		scdnReq.Query["ca_status"] = req.CAStatus
	}
	if req.AccessMode != "" {
		scdnReq.Query["access_mode"] = req.AccessMode
	}
	if req.ProtectStatus != "" {
		scdnReq.Query["protect_status"] = req.ProtectStatus
	}
	if req.ExclusiveResourceID > 0 {
		scdnReq.Query["exclusive_resource_id"] = req.ExclusiveResourceID
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointDomains, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list domains: %w", err)
	}

	// Convert response to DomainListResponse
	response := &DomainListResponse{
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

		var domainListData DomainListData
		if err := json.Unmarshal(dataBytes, &domainListData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal domain list data: %w", err)
		}
		response.Data = domainListData
	}

	return response, nil
}

// ListDomainsSimple lists simple domains
func (s *ScdnService) ListDomainsSimple(req DomainSimpleListRequest) (*DomainSimpleListResponse, error) {
	ctx := context.Background()

	query := make(map[string]string)
	if req.Domain != "" {
		query["domain"] = req.Domain
	}
	if req.Page > 0 {
		query["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.PerPage > 0 {
		query["per_page"] = fmt.Sprintf("%d", req.PerPage)
	}

	var response DomainSimpleListResponse
	err := s.callSCDNAPI(ctx, MethodGET, EndpointDomainsSimple, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list simple domains: %w", err)
	}

	if response.Status.Code != 1 {
		return nil, fmt.Errorf("failed to list simple domains: %s", response.Status.Message)
	}

	return &response, nil
}

// CreateDomain creates a new domain
func (s *ScdnService) CreateDomain(req DomainCreateRequest) (*DomainCreateResponse, error) {
	ctx := context.Background()

	var response DomainCreateResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomains, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain: %w", err)
	}

	return &response, nil
}

// UpdateDomain updates an existing domain
func (s *ScdnService) UpdateDomain(req DomainUpdateRequest) (*DomainUpdateResponse, error) {
	ctx := context.Background()

	var response DomainUpdateResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointDomains, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update domain: %w", err)
	}

	return &response, nil
}

// BindDomainCert binds a certificate to a domain
func (s *ScdnService) BindDomainCert(req DomainCertBindRequest) (*DomainCertBindResponse, error) {
	ctx := context.Background()

	var response DomainCertBindResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsBindCert, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to bind domain certificate: %w", err)
	}

	return &response, nil
}

// UnbindDomainCert unbinds a certificate from a domain
func (s *ScdnService) UnbindDomainCert(req DomainCertUnbindRequest) (*DomainCertUnbindResponse, error) {
	ctx := context.Background()

	var response DomainCertUnbindResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsUnbindCert, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unbind domain certificate: %w", err)
	}

	return &response, nil
}

// DeleteDomain deletes domains
func (s *ScdnService) DeleteDomain(req DomainDeleteRequest) (*DomainDeleteResponse, error) {
	ctx := context.Background()

	var response DomainDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointDomains, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete domain: %w", err)
	}

	return &response, nil
}

// DisableDomain disables domains
func (s *ScdnService) DisableDomain(req DomainDisableRequest) (*DomainDisableResponse, error) {
	ctx := context.Background()

	var response DomainDisableResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsDisable, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to disable domain: %w", err)
	}

	return &response, nil
}

// EnableDomain enables domains
func (s *ScdnService) EnableDomain(req DomainEnableRequest) (*DomainEnableResponse, error) {
	ctx := context.Background()

	var response DomainEnableResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsEnable, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to enable domain: %w", err)
	}

	return &response, nil
}

// RefreshDomainAccess refreshes domain access status
func (s *ScdnService) RefreshDomainAccess(req DomainAccessRefreshRequest) (*DomainAccessRefreshResponse, error) {
	ctx := context.Background()

	var response DomainAccessRefreshResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsAccessRefresh, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh domain access: %w", err)
	}

	return &response, nil
}

// ExportDomains exports domains
func (s *ScdnService) ExportDomains(req DomainExportRequest) (*DomainExportResponse, error) {
	ctx := context.Background()

	var response DomainExportResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsExport, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to export domains: %w", err)
	}

	return &response, nil
}

// AddOrigins adds origins to a domain
func (s *ScdnService) AddOrigins(req OriginAddRequest) (*OriginAddResponse, error) {
	ctx := context.Background()

	var response OriginAddResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsOrigins, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to add origins: %w", err)
	}

	return &response, nil
}

// UpdateOrigins updates origins for a domain
func (s *ScdnService) UpdateOrigins(req OriginUpdateRequest) (*OriginUpdateResponse, error) {
	ctx := context.Background()

	var response OriginUpdateResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointDomainsOrigins, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update origins: %w", err)
	}

	return &response, nil
}

// DeleteOrigins deletes origins from a domain
func (s *ScdnService) DeleteOrigins(req OriginDeleteRequest) (*OriginDeleteResponse, error) {
	ctx := context.Background()

	var response OriginDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointDomainsOrigins, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete origins: %w", err)
	}

	return &response, nil
}

// ListOrigins lists origins for a domain
func (s *ScdnService) ListOrigins(req OriginListRequest) (*OriginListResponse, error) {
	ctx := context.Background()

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	if req.DomainID > 0 {
		scdnReq.Query["domain_id"] = req.DomainID
	}

	var response OriginListResponse
	err := s.requestSCDNAPI(ctx, MethodGET, EndpointDomainsOrigins, scdnReq, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list origins: %w", err)
	}

	return &response, nil
}

// SwitchDomainNodes switches domain nodes
func (s *ScdnService) SwitchDomainNodes(req DomainNodeSwitchRequest) (*DomainNodeSwitchResponse, error) {
	ctx := context.Background()

	var response DomainNodeSwitchResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsNodesSwitch, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to switch domain nodes: %w", err)
	}

	return &response, nil
}

// SwitchDomainAccessMode switches domain access mode
func (s *ScdnService) SwitchDomainAccessMode(req DomainAccessModeSwitchRequest) (*DomainAccessModeSwitchResponse, error) {
	ctx := context.Background()

	var response DomainAccessModeSwitchResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsAccessSwitch, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to switch domain access mode: %w", err)
	}

	return &response, nil
}

// GetAccessProgress gets access progress status list
func (s *ScdnService) GetAccessProgress() (*AccessProgressResponse, error) {
	ctx := context.Background()

	var response AccessProgressResponse
	err := s.callSCDNAPI(ctx, MethodGET, EndpointDomainsAccessProgress, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get access progress: %w", err)
	}

	return &response, nil
}

// UpdateDomainBaseSettings updates domain base settings
func (s *ScdnService) UpdateDomainBaseSettings(req DomainBaseSettingsUpdateRequest) (*DomainBaseSettingsUpdateResponse, error) {
	ctx := context.Background()

	var response DomainBaseSettingsUpdateResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointDomainsBaseSettings, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update domain base settings: %w", err)
	}

	return &response, nil
}

// GetDomainBaseSettings gets domain base settings
func (s *ScdnService) GetDomainBaseSettings(req DomainBaseSettingsGetRequest) (*DomainBaseSettingsGetResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: map[string]interface{}{
			"domain_id": req.DomainID,
		},
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointDomainsBaseSettings, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain base settings: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	var response DomainBaseSettingsGetResponse
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

// ListBriefDomains lists brief domains
func (s *ScdnService) ListBriefDomains(req BriefDomainListRequest) (*BriefDomainListResponse, error) {
	ctx := context.Background()

	var response BriefDomainListResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointBriefDomains, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list brief domains: %w", err)
	}

	return &response, nil
}

// GetDomainTemplates gets domain templates
func (s *ScdnService) GetDomainTemplates(req DomainTemplatesRequest) (*DomainTemplatesResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: map[string]interface{}{
			"domain_id": req.DomainID,
		},
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointDomainsTemplates, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain templates: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	var response DomainTemplatesResponse
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

// DownloadAccessInfo downloads access information
func (s *ScdnService) DownloadAccessInfo(req AccessInfoDownloadRequest) (*AccessInfoDownloadResponse, error) {
	ctx := context.Background()

	var response AccessInfoDownloadResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointDomainsAccessInfoDownload, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to download access info: %w", err)
	}

	return &response, nil
}
