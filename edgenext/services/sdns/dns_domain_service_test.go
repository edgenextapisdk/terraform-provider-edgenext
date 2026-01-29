package sdns

import (
	"testing"
)

func TestSdnsService_DnsDomainCRUD(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	client := createTestClient(t)
	service := NewSdnsService(client)

	domainName := "test-dns-domain.com"

	// 1. Add
	t.Run("AddDomain", func(t *testing.T) {
		resp, err := service.AddDnsDomain(domainName)
		if err != nil {
			t.Fatalf("AddDnsDomain failed: %v", err)
		}
		t.Logf("Added domain ID: %d", resp.ID)
	})

	// 2. List
	var domainID int
	t.Run("ListDomains", func(t *testing.T) {
		resp, err := service.ListDnsDomains(DnsDomainListRequest{Domain: domainName})
		if err != nil {
			t.Fatalf("ListDnsDomains failed: %v", err)
		}
		if len(resp.List) == 0 {
			t.Fatalf("Domain not found in list")
		}
		for _, d := range resp.List {
			if d.Domain == domainName {
				domainID = d.ID
				break
			}
		}
	})

	// 3. Info
	if domainID > 0 {
		t.Run("GetDomainInfo", func(t *testing.T) {
			info, err := service.GetDnsDomainInfo(domainID)
			if err != nil {
				t.Fatalf("GetDnsDomainInfo failed: %v", err)
			}
			if info.Domain != domainName {
				t.Errorf("Domain name mismatch: expected %s, got %s", domainName, info.Domain)
			}
		})
	}

	// 4. Delete
	t.Run("DeleteDomain", func(t *testing.T) {
		resp, _ := service.ListDnsDomains(DnsDomainListRequest{Domain: domainName})
		var ids []int
		for _, d := range resp.List {
			if d.Domain == domainName {
				ids = append(ids, d.ID)
			}
		}
		if len(ids) > 0 {
			err := service.DeleteDnsDomain(ids)
			if err != nil {
				t.Fatalf("DeleteDnsDomain failed: %v", err)
			}
		}
	})
}
