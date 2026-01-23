package domain_group

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
)

type DomainGroupService struct {
	client *connectivity.EdgeNextClient
}

func NewDomainGroupService(client *connectivity.EdgeNextClient) *DomainGroupService {
	return &DomainGroupService{client: client}
}

// Helper to call API
func (s *DomainGroupService) callAPI(ctx context.Context, method, api string, reqData interface{}, respData interface{}) error {
	client, err := s.client.ScdnClient()
	if err != nil {
		return err
	}

	// Prepare Request
	scdnReq := &connectivity.ScdnRequest{
		Data:  make(map[string]interface{}),
		Query: make(map[string]interface{}),
	}

	if reqData != nil {
		// Marshal to JSON then Unmarshal to map
		dataBytes, err := json.Marshal(reqData)
		if err != nil {
			return fmt.Errorf("failed to marshal request data: %w", err)
		}

		var wrapper map[string]interface{}
		if err := json.Unmarshal(dataBytes, &wrapper); err != nil {
			return fmt.Errorf("failed to unmarshal request data to map: %w", err)
		}

		// For GET requests, parameters are usually queries
		if method == "GET" {
			scdnReq.Query = wrapper
		} else {
			scdnReq.Data = wrapper
		}
	}

	var scdnResp *connectivity.ScdnResponse
	switch method {
	case "GET":
		scdnResp, err = client.Get(ctx, api, scdnReq)
	case "POST":
		scdnResp, err = client.Post(ctx, api, scdnReq)
	case "DELETE":
		scdnResp, err = client.Delete(ctx, api, scdnReq)
	default:
		return fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return err
	}

	// Map Status
	// We need to set status on respData which is a pointer to a struct containing Status field
	// Use JSON roundtrip to populate respData including Status and Data
	// Construct a temporary structure matching the JSON of respData

	// Create a map representing the full response structure
	fullResp := map[string]interface{}{
		"status": map[string]interface{}{
			"code":    scdnResp.Status.Code,
			"message": scdnResp.Status.Message,
		},
		"data": scdnResp.Data,
	}

	fullRespBytes, err := json.Marshal(fullResp)
	if err != nil {
		return fmt.Errorf("failed to marshal full response: %w", err)
	}

	if err := json.Unmarshal(fullRespBytes, respData); err != nil {
		return fmt.Errorf("failed to unmarshal into response struct: %w", err)
	}

	return nil
}

// AddDomainGroup adds a new domain group
func (s *DomainGroupService) AddDomainGroup(req DomainGroupSaveRequest) (*DomainGroupSaveResponse, error) {
	var resp DomainGroupSaveResponse
	err := s.callAPI(context.Background(), "POST", scdn.EndpointDomainGroupAdd, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDomainGroup updates a domain group
func (s *DomainGroupService) UpdateDomainGroup(req DomainGroupSaveRequest) (*DomainGroupSaveResponse, error) {
	var resp DomainGroupSaveResponse
	err := s.callAPI(context.Background(), "POST", scdn.EndpointDomainGroupSave, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDomainGroup deletes a domain group
func (s *DomainGroupService) DeleteDomainGroup(req DomainGroupDelRequest) (*DomainGroupDelResponse, error) {
	var resp DomainGroupDelResponse
	err := s.callAPI(context.Background(), "POST", scdn.EndpointDomainGroupDel, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListDomainGroups lists domain groups
func (s *DomainGroupService) ListDomainGroups(req DomainGroupListRequest) (*DomainGroupListResponse, error) {
	var resp DomainGroupListResponse
	err := s.callAPI(context.Background(), "GET", scdn.EndpointDomainGroupList, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDomainGroupInfo gets domain group details
func (s *DomainGroupService) GetDomainGroupInfo(groupID int) (*DomainGroup, error) {
	req := map[string]interface{}{
		"group_id": groupID,
	}

	type DomainGroupInfoResponse struct {
		Status Status      `json:"status"`
		Data   DomainGroup `json:"data"`
	}

	var resp DomainGroupInfoResponse
	err := s.callAPI(context.Background(), "GET", scdn.EndpointDomainGroupInfo, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListGroupDomains lists domains in a group
func (s *DomainGroupService) ListGroupDomains(req DomainGroupDomainListRequest) (*DomainGroupDomainListResponse, error) {
	var resp DomainGroupDomainListResponse
	err := s.callAPI(context.Background(), "GET", scdn.EndpointDomainGroupDomainList, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// BindDomainsToGroup binds or unbinds domains to a group
func (s *DomainGroupService) BindDomainsToGroup(req DomainGroupDomainSaveRequest) (*DomainGroupDomainSaveResponse, error) {
	var resp DomainGroupDomainSaveResponse
	err := s.callAPI(context.Background(), "POST", scdn.EndpointDomainGroupDomainSave, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListUndistributedDomains lists domains not bound to a group
func (s *DomainGroupService) ListUndistributedDomains(req DomainGroupDomainListRequest) (*DomainGroupDomainListResponse, error) {
	var resp DomainGroupDomainListResponse
	err := s.callAPI(context.Background(), "GET", scdn.EndpointDomainGroupUndistributedDomainList, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// MoveDomains moves domains between groups
func (s *DomainGroupService) MoveDomains(req DomainGroupMoveDomainRequest) (*DomainGroupMoveDomainResponse, error) {
	var resp DomainGroupMoveDomainResponse
	err := s.callAPI(context.Background(), "POST", scdn.EndpointDomainGroupMoveDomain, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Helper to find a group by ID string (terraform uses string IDs)
func (s *DomainGroupService) GetDomainGroupByID(id string) (*DomainGroup, error) {
	groupID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID: %s", id)
	}
	return s.GetDomainGroupInfo(groupID)
}
