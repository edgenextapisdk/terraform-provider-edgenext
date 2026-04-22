package sdns

import (
	"context"
	"fmt"
)

// GetDnsDomainInfo Gets DNS domain info by listing and filtering
func (s *SdnsService) GetDnsDomainInfo(domainID int) (*DnsDomainInfo, error) {
	resp, err := s.ListDnsDomains(DnsDomainListRequest{
		Id:      domainID,
		PerPage: 1000,
	})

	if err != nil {
		return nil, err
	}
	for _, d := range resp.List {
		if d.ID == domainID {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("domain not found: %d", domainID)
}

// ListDnsDomains Lists DNS domains
func (s *SdnsService) ListDnsDomains(req DnsDomainListRequest) (*DnsDomainListData, error) {
	var resp DnsDomainListResponse
	err := s.callAPI(context.Background(), "GET", EndpointDnsDomainList, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AddDnsDomain Adds a new DNS domain
func (s *SdnsService) AddDnsDomain(domainName string) (*DnsDomainAddData, error) {
	req := DnsDomainAddRequest{Domain: domainName}
	var resp DnsDomainAddResponse
	err := s.callAPI(context.Background(), "POST", EndpointDnsDomainAdd, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteDnsDomain Deletes DNS domains
func (s *SdnsService) DeleteDnsDomain(domainIDs []int) error {
	req := DnsDomainDeleteRequest{DomainIDs: domainIDs}
	return s.callAPI(context.Background(), "DELETE", EndpointDnsDomainBatchDelete, req, nil)
}
