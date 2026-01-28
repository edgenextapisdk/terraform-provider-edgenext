package sdns

import (
	"context"
)

// ListDnsRecords Lists DNS domain records
func (s *SdnsService) ListDnsRecords(req DnsRecordListRequest) (*DnsRecordListData, error) {
	var resp DnsRecordListResponse
	err := s.callAPI(context.Background(), "GET", EndpointDnsRecordList, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AddDnsRecord Adds a new DNS domain record
func (s *SdnsService) AddDnsRecord(req DnsRecordAddRequest) (int, error) {
	var resp DnsRecordResponse
	err := s.callAPI(context.Background(), "POST", EndpointDnsRecordAdd, req, &resp)
	if err != nil {
		return 0, err
	}
	return resp.Data.ID, nil
}

// UpdateDnsRecord Updates a DNS domain record
func (s *SdnsService) UpdateDnsRecord(req DnsRecordEditRequest) error {
	return s.callAPI(context.Background(), "PUT", EndpointDnsRecordEdit, req, nil)
}

// DeleteDnsRecord Deletes a DNS domain record
func (s *SdnsService) DeleteDnsRecord(recordID, domainID int) error {
	req := DnsRecordDeleteRequest{RecordID: recordID, DomainID: domainID}
	return s.callAPI(context.Background(), "DELETE", EndpointDnsRecordDelete, req, nil)
}
