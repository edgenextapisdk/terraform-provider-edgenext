package scdn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// User IP Intelligence Methods
// ============================================================================

// ListUserIps lists user IP lists
func (s *ScdnService) ListUserIps(req UserIpListRequest) (*UserIpListResponse, error) {
	ctx := context.Background()

	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	if req.Page > 0 {
		scdnReq.Query["page"] = req.Page
	}
	if req.PerPage > 0 {
		scdnReq.Query["per_page"] = req.PerPage
	}

	scdnResp, err := scdnClient.Get(ctx, EndpointUserIpList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list user ips: %w", err)
	}

	response := &UserIpListResponse{
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

		if err := json.Unmarshal(dataBytes, &response.Data); err != nil {
			// Try to handle empty list case or different format if needed
			// But for now assume standard format
			return nil, fmt.Errorf("failed to unmarshal user ip list data: %w", err)
		}
	}

	return response, nil
}

// AddUserIp creates a new user IP list
func (s *ScdnService) AddUserIp(req UserIpAddRequest) (*UserIpAddResponse, error) {
	ctx := context.Background()

	var response UserIpAddResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointUserIpAdd, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to add user ip list: %w", err)
	}

	return &response, nil
}

// UpdateUserIp updates an existing user IP list
func (s *ScdnService) UpdateUserIp(req UserIpSaveRequest) (*UserIpSaveResponse, error) {
	ctx := context.Background()

	var response UserIpSaveResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointUserIpSave, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update user ip list: %w", err)
	}

	return &response, nil
}

// DeleteUserIp deletes one or more user IP lists
func (s *ScdnService) DeleteUserIp(req UserIpDelRequest) (*UserIpDelResponse, error) {
	ctx := context.Background()

	var response UserIpDelResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointUserIpDel, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user ip list: %w", err)
	}

	return &response, nil
}

// ListUserIpItems lists items in a user IP list
func (s *ScdnService) ListUserIpItems(req UserIpItemListRequest) (*UserIpItemListResponse, error) {
	ctx := context.Background()

	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	scdnReq.Query["user_ip_id"] = req.UserIpID
	if req.IP != "" {
		scdnReq.Query["ip"] = req.IP
	}
	if req.Page > 0 {
		scdnReq.Query["page"] = req.Page
	}
	if req.PerPage > 0 {
		scdnReq.Query["per_page"] = req.PerPage
	}

	scdnResp, err := scdnClient.Get(ctx, EndpointUserIpItemList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list user ip items: %w", err)
	}

	response := &UserIpItemListResponse{
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

		if err := json.Unmarshal(dataBytes, &response.Data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal user ip item list data: %w", err)
		}
	}

	return response, nil
}

// AddUserIpItem adds an item to a user IP list
func (s *ScdnService) AddUserIpItem(req UserIpItemAddRequest) (*UserIpItemAddResponse, error) {
	ctx := context.Background()

	var response UserIpItemAddResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointUserIpItemTextSave, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to add user ip item: %w", err)
	}

	return &response, nil
}

// UpdateUserIpItem updates an item in a user IP list
func (s *ScdnService) UpdateUserIpItem(req UserIpItemEditRequest) (*UserIpItemEditResponse, error) {
	ctx := context.Background()

	var response UserIpItemEditResponse
	err := s.callSCDNAPI(ctx, MethodPUT, EndpointUserIpItemEdit, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to update user ip item: %w", err)
	}

	return &response, nil
}

// DeleteUserIpItem deletes items from a user IP list
func (s *ScdnService) DeleteUserIpItem(req UserIpItemDelRequest) (*UserIpItemDelResponse, error) {
	ctx := context.Background()

	var response UserIpItemDelResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointUserIpItemDel, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user ip item: %w", err)
	}

	return &response, nil
}

// DeleteAllUserIpItems deletes all items from a user IP list
func (s *ScdnService) DeleteAllUserIpItems(req UserIpItemDelAllRequest) (*UserIpItemDelAllResponse, error) {
	ctx := context.Background()

	var response UserIpItemDelAllResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointUserIpItemAll, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete all user ip items: %w", err)
	}

	return &response, nil
}

// CopyUserIp copies a user IP list
func (s *ScdnService) CopyUserIp(req UserIpCopyRequest) (*UserIpCopyResponse, error) {
	ctx := context.Background()

	var response UserIpCopyResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointUserIpCopy, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to copy user ip list: %w", err)
	}

	return &response, nil
}

// UploadUserIpFile uploads a file to a user IP list
// Note: This requires multipart file upload which callSCDNAPI might not support directly if it expects JSON.
// We might need to construct the request manually using scdnClient.
func (s *ScdnService) UploadUserIpFile(userIpID string, filePath string, remark string) (*UserIpItemFileSaveResponse, error) {
	ctx := context.Background()

	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	params := map[string]string{
		"user_ip_id": userIpID,
	}
	if remark != "" {
		params["remark"] = remark
	}

	resp, err := scdnClient.Upload(ctx, EndpointUserIpItemFileSave, params, "file", filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload user ip file: %w", err)
	}

	response := &UserIpItemFileSaveResponse{
		Status: Status{
			Code:    resp.Status.Code,
			Message: resp.Status.Message,
		},
	}

	if resp.Data != nil {
		dataBytes, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		if err := json.Unmarshal(dataBytes, &response.Data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal user ip file save data: %w", err)
		}
	}

	return response, nil
}
