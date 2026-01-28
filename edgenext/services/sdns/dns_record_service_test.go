package sdns

import (
	"testing"
)

func TestSdnsService_DnsRecordCRUD(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	client := createTestClient(t)
	service := NewSdnsService(client)

	// We need a domain for testing records
	domainName := "test-dns-record-domain.com"
	service.AddDnsDomain(domainName)
	defer func() {
		resp, _ := service.ListDnsDomains(DnsDomainListRequest{Domain: domainName})
		var ids []int
		for _, d := range resp.List {
			if d.Domain == domainName {
				ids = append(ids, d.ID)
			}
		}
		if len(ids) > 0 {
			service.DeleteDnsDomain(ids)
		}
	}()

	resp, _ := service.ListDnsDomains(DnsDomainListRequest{Domain: domainName})
	if len(resp.List) == 0 {
		t.Fatalf("Failed to prepare test domain")
	}
	domainID := resp.List[0].ID

	var recordID int

	// 1. Add
	t.Run("AddRecord", func(t *testing.T) {
		req := DnsRecordAddRequest{
			DomainID:    domainID,
			RecordName:  "www",
			RecordType:  "A",
			RecordView:  "any",
			RecordValue: "1.2.3.4",
			RecordTTL:   600,
		}
		id, err := service.AddDnsRecord(req)
		if err != nil {
			t.Fatalf("AddDnsRecord failed: %v", err)
		}
		recordID = id
	})

	// 2. List
	t.Run("ListRecords", func(t *testing.T) {
		resp, err := service.ListDnsRecords(DnsRecordListRequest{DomainID: domainID})
		if err != nil {
			t.Fatalf("ListDnsRecords failed: %v", err)
		}
		found := false
		for _, r := range resp.List {
			if r.ID == recordID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Record not found in list")
		}
	})

	// 3. Update
	t.Run("UpdateRecord", func(t *testing.T) {
		req := DnsRecordEditRequest{
			RecordID:    recordID,
			DomainID:    domainID,
			RecordName:  "www",
			RecordType:  "A",
			RecordView:  "any",
			RecordValue: "1.2.3.5",
			RecordTTL:   300,
		}
		err := service.UpdateDnsRecord(req)
		if err != nil {
			t.Fatalf("UpdateDnsRecord failed: %v", err)
		}
	})

	// 4. Delete
	t.Run("DeleteRecord", func(t *testing.T) {
		err := service.DeleteDnsRecord(recordID, domainID)
		if err != nil {
			t.Fatalf("DeleteDnsRecord failed: %v", err)
		}
	})
}
