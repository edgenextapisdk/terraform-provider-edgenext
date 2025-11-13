package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Certificate Management Methods
// ============================================================================

// SaveCertificate saves or updates a certificate using text format
func (s *ScdnService) SaveCertificate(req CATextSaveRequest) (*CATextSaveResponse, error) {
	ctx := context.Background()

	var response CATextSaveResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointCATextSave, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to save certificate: %w", err)
	}

	return &response, nil
}

// ListCertificates lists certificates with various filter options
func (s *ScdnService) ListCertificates(req CASelfListRequest) (*CASelfListResponse, error) {
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
	if req.Domain != "" {
		scdnReq.Query["domain"] = req.Domain
	}
	if req.ProductFlag != "" {
		scdnReq.Query["product_flag"] = req.ProductFlag
	}
	if req.CAName != "" {
		scdnReq.Query["ca_name"] = req.CAName
	}
	if req.Binded != "" {
		scdnReq.Query["binded"] = req.Binded
	}
	if req.ApplyStatus != "" {
		scdnReq.Query["apply_status"] = req.ApplyStatus
	}
	if req.Issuer != "" {
		scdnReq.Query["issuer"] = req.Issuer
	}
	if req.ExpiryTime != "" {
		scdnReq.Query["expiry_time"] = req.ExpiryTime
	}
	if req.IsExactSearch != "" {
		scdnReq.Query["is_exact_search"] = req.IsExactSearch
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointCASelfList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list certificates: %w", err)
	}

	// Convert response
	response := &CASelfListResponse{
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

		var certListData CASelfListData
		if err := json.Unmarshal(dataBytes, &certListData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal certificate list data: %w", err)
		}
		response.Data = certListData
	}

	return response, nil
}

// GetCertificateDetail gets certificate detail by ID
func (s *ScdnService) GetCertificateDetail(req CASelfDetailRequest) (*CASelfDetailResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: map[string]interface{}{
			"id": req.ID,
		},
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointCASelf, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate detail: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	var response CASelfDetailResponse
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

// DeleteCertificate deletes certificates
func (s *ScdnService) DeleteCertificate(req CASelfDeleteRequest) (*CASelfDeleteResponse, error) {
	ctx := context.Background()

	var response CASelfDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointCASelfDel, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete certificate: %w", err)
	}

	return &response, nil
}

// EditCertificateName edits certificate name
func (s *ScdnService) EditCertificateName(req CAEditNameRequest) (*CAEditNameResponse, error) {
	ctx := context.Background()

	var response CAEditNameResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointCAEditName, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to edit certificate name: %w", err)
	}

	return &response, nil
}

// ListCertificatesByDomains lists certificates by domain list
func (s *ScdnService) ListCertificatesByDomains(req CABatchListRequest) (*CABatchListResponse, error) {
	ctx := context.Background()

	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{
		Data: make(map[string]interface{}),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	if err := json.Unmarshal(reqBytes, &scdnReq.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Post(ctx, EndpointCABatchList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list certificates by domains: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CABatchListResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	// Extract array data from various formats
	if err := extractArrayFromData(scdnResp.Data, &response.Data); err != nil {
		return nil, fmt.Errorf("failed to extract array data: %w", err)
	}

	return response, nil
}

// ApplyCertificate applies for a certificate
func (s *ScdnService) ApplyCertificate(req CAApplyAddRequest) (*CAApplyAddResponse, error) {
	ctx := context.Background()

	var response CAApplyAddResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointCAApplyAdd, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to apply certificate: %w", err)
	}

	return &response, nil
}

// ExportCertificate exports certificates
func (s *ScdnService) ExportCertificate(req CASelfExportRequest) (*CASelfExportResponse, error) {
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
	if req.ProductFlag != "" {
		scdnReq.Query["product_flag"] = req.ProductFlag
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointCASelfExport, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to export certificate: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &CASelfExportResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}
	// Extract array data from various formats
	if err := extractArrayFromData(scdnResp.Data, &response.Data); err != nil {
		return nil, fmt.Errorf("failed to extract array data: %w", err)
	}

	return response, nil
}
