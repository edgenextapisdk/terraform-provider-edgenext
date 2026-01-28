package sdns

import (
	"context"
	"fmt"
)

// ListDnsGroups Lists DNS domain groups
func (s *SdnsService) ListDnsGroups(req DnsGroupListRequest) (*DnsGroupListData, error) {
	var resp DnsGroupListResponse
	err := s.callAPI(context.Background(), "GET", EndpointDnsGroupList, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AddDnsGroup Adds a new DNS domain group
func (s *SdnsService) AddDnsGroup(req DnsGroupAddRequest) (*DnsGroupAddResponse, error) {
	var resp DnsGroupAddResponse
	err := s.callAPI(context.Background(), "POST", EndpointDnsGroupAdd, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDnsGroup Updates a DNS domain group
func (s *SdnsService) UpdateDnsGroup(req DnsGroupSaveRequest) error {
	return s.callAPI(context.Background(), "PUT", EndpointDnsGroupUpdate, req, nil)
}

// DeleteDnsGroup Deletes a DNS domain group
func (s *SdnsService) DeleteDnsGroup(groupID int) error {
	req := DnsGroupDelRequest{GroupID: groupID}
	return s.callAPI(context.Background(), "DELETE", EndpointDnsGroupDelete, req, nil)
}

// GetDnsGroupInfo Gets DNS domain group info by listing and filtering
func (s *SdnsService) GetDnsGroupInfo(groupID int) (*DnsGroup, error) {
	resp, err := s.ListDnsGroups(DnsGroupListRequest{
		Id:      groupID,
		PerPage: 1000,
	})

	if err != nil {
		return nil, err
	}
	for _, g := range resp.List {
		if g.ID == groupID {
			return &g, nil
		}
	}
	return nil, fmt.Errorf("group not found: %d", groupID)
}

// BindDomainsToGroup Binds or unbinds domains to a group
func (s *SdnsService) BindDomainsToGroup(req DnsGroupDomainSaveRequest) error {
	return s.callAPI(context.Background(), "POST", EndpointDnsGroupRecordRelation, req, nil)
}
