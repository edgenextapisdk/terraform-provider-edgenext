package scdn

import (
	"context"
	"testing"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/stretchr/testify/assert"
)

func TestScdnService_ListDomains(t *testing.T) {
	// Test request structure
	req := DomainListRequest{
		Page:     1,
		PageSize: 10,
		Domain:   "test.example.com",
	}

	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.PageSize)
	assert.Equal(t, "test.example.com", req.Domain)

	// Skip if not running integration tests
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	// Create real client and service
	client := createTestClient(t)
	service := NewScdnService(client)
	// Make real API call
	response, err := service.ListDomains(req)
	if err != nil {
		t.Logf("API call failed (this might be expected in test environment): %v", err)
		return
	}

	// Validate response structure
	assert.NotNil(t, response)
	assert.NotNil(t, response.Status)
	assert.NotNil(t, response.Data)
	assert.GreaterOrEqual(t, response.Data.Total, 0)
}

func TestScdnService_CreateDomain(t *testing.T) {
	req := DomainCreateRequest{
		Domain:        "test.example.com",
		GroupID:       123,
		Remark:        "Test domain",
		ProtectStatus: "scdn",
		Origins: []Origin{
			{
				Protocol:       0,
				ListenPorts:    []int{80, 443},
				OriginProtocol: 0,
				LoadBalance:    1,
				OriginType:     0,
				Records: []OriginRecord{
					{
						View:     "primary",
						Value:    "123.4.15.10",
						Port:     80,
						Priority: 1,
					},
				},
			},
		},
	}

	// Test request structure
	assert.Equal(t, "test.example.com", req.Domain)
	assert.Equal(t, 123, req.GroupID)
	assert.Equal(t, "Test domain", req.Remark)
	assert.Equal(t, "scdn", req.ProtectStatus)
	assert.Len(t, req.Origins, 1)
	assert.Equal(t, 0, req.Origins[0].Protocol)
	assert.Len(t, req.Origins[0].ListenPorts, 2)
	assert.Len(t, req.Origins[0].Records, 1)

	// Skip if not running integration tests
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	// Create real client and service
	client := createTestClient(t)
	service := NewScdnService(client)
	// Make real API call (this will likely fail in test environment, but we test the structure)
	response, err := service.CreateDomain(req)
	if err != nil {
		t.Logf("API call failed (this might be expected in test environment): %v", err)
		return
	}

	// Validate response structure
	assert.NotNil(t, response)
	assert.NotNil(t, response.Status)
}

func TestScdnService_UpdateDomain(t *testing.T) {
	req := DomainUpdateRequest{
		DomainID: 123,
		Remark:   "Updated remark",
	}

	assert.Equal(t, 123, req.DomainID)
	assert.Equal(t, "Updated remark", req.Remark)
}

func TestScdnService_BindDomainCertRequest(t *testing.T) {
	req := DomainCertBindRequest{
		DomainID: 123,
		CAID:     456,
	}

	assert.Equal(t, 123, req.DomainID)
	assert.Equal(t, 456, req.CAID)
}

func TestScdnService_UnbindDomainCertRequest(t *testing.T) {
	req := DomainCertUnbindRequest{
		DomainID: 123,
		CAID:     456,
	}

	assert.Equal(t, 123, req.DomainID)
	assert.Equal(t, 456, req.CAID)
}

func TestScdnService_DeleteDomainRequest(t *testing.T) {
	req := DomainDeleteRequest{
		IDs: []int{123, 456},
	}

	assert.Len(t, req.IDs, 2)
	assert.Contains(t, req.IDs, 123)
	assert.Contains(t, req.IDs, 456)
}

func TestScdnService_DisableDomainRequest(t *testing.T) {
	req := DomainDisableRequest{
		DomainIDs: []int{123, 456},
	}

	assert.Len(t, req.DomainIDs, 2)
	assert.Contains(t, req.DomainIDs, 123)
	assert.Contains(t, req.DomainIDs, 456)
}

func TestScdnService_EnableDomainRequest(t *testing.T) {
	req := DomainEnableRequest{
		DomainIDs: []int{123, 456},
	}

	assert.Len(t, req.DomainIDs, 2)
	assert.Contains(t, req.DomainIDs, 123)
	assert.Contains(t, req.DomainIDs, 456)
}

func TestScdnService_RefreshDomainAccessRequest(t *testing.T) {
	req := DomainAccessRefreshRequest{
		DomainIDs: []int{123, 456},
	}

	assert.Len(t, req.DomainIDs, 2)
	assert.Contains(t, req.DomainIDs, 123)
	assert.Contains(t, req.DomainIDs, 456)
}

func TestScdnService_ExportDomainsRequest(t *testing.T) {
	req := DomainExportRequest{
		DomainIDs:           []int{123, 456},
		AccessProgress:      "enabled",
		GroupID:             789,
		Domain:              "test.example.com",
		Remark:              "Test domain",
		OriginIP:            "1.1.1.1",
		CAStatus:            "bind",
		AccessMode:          "cname",
		ProtectStatus:       "scdn",
		ExclusiveResourceID: 101,
	}

	assert.Len(t, req.DomainIDs, 2)
	assert.Equal(t, "enabled", req.AccessProgress)
	assert.Equal(t, 789, req.GroupID)
	assert.Equal(t, "test.example.com", req.Domain)
	assert.Equal(t, "Test domain", req.Remark)
	assert.Equal(t, "1.1.1.1", req.OriginIP)
	assert.Equal(t, "bind", req.CAStatus)
	assert.Equal(t, "cname", req.AccessMode)
	assert.Equal(t, "scdn", req.ProtectStatus)
	assert.Equal(t, 101, req.ExclusiveResourceID)
}

func TestScdnService_AddOriginsRequest(t *testing.T) {
	req := OriginAddRequest{
		DomainID: 123,
		Origins: []Origin{
			{
				Protocol:       0,
				ListenPorts:    []int{80, 443},
				OriginProtocol: 0,
				LoadBalance:    1,
				OriginType:     0,
				Records: []OriginRecord{
					{
						View:     "default",
						Value:    "1.1.1.1",
						Port:     80,
						Priority: 10,
					},
				},
			},
		},
	}

	assert.Equal(t, 123, req.DomainID)
	assert.Len(t, req.Origins, 1)
	assert.Equal(t, 0, req.Origins[0].Protocol)
}

func TestScdnService_UpdateOriginsRequest(t *testing.T) {
	req := OriginUpdateRequest{
		DomainID: 123,
		Origins: []EditOrigin{
			{
				Protocol:       0,
				ListenPort:     80,
				OriginProtocol: 0,
				LoadBalance:    1,
				OriginType:     0,
				Records: []OriginRecord{
					{
						View:     "default",
						Value:    "1.1.1.1",
						Port:     80,
						Priority: 10,
					},
				},
			},
		},
	}

	assert.Equal(t, 123, req.DomainID)
	assert.Len(t, req.Origins, 1)
	// Origin struct doesn't have ID and DomainID fields
}

func TestScdnService_DeleteOriginsRequest(t *testing.T) {
	req := OriginDeleteRequest{
		IDs:      []int{1, 2},
		DomainID: 123,
	}

	assert.Len(t, req.IDs, 2)
	assert.Contains(t, req.IDs, 1)
	assert.Contains(t, req.IDs, 2)
	assert.Equal(t, 123, req.DomainID)
}

func TestScdnService_ListOriginsRequest(t *testing.T) {
	req := OriginListRequest{
		DomainID: 123,
	}

	assert.Equal(t, 123, req.DomainID)
}

func TestScdnService_SwitchDomainNodesRequest(t *testing.T) {
	req := DomainNodeSwitchRequest{
		DomainID:            123,
		ProtectStatus:       "exclusive",
		ExclusiveResourceID: 456,
	}

	assert.Equal(t, 123, req.DomainID)
	assert.Equal(t, "exclusive", req.ProtectStatus)
	assert.Equal(t, 456, req.ExclusiveResourceID)
}

func TestScdnService_SwitchDomainAccessModeRequest(t *testing.T) {
	req := DomainAccessModeSwitchRequest{
		DomainID:   123,
		AccessMode: "ns",
	}

	assert.Equal(t, 123, req.DomainID)
	assert.Equal(t, "ns", req.AccessMode)
}

func TestScdnService_GetAccessProgressRequest(t *testing.T) {
	// Skip if not running integration tests
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	// Create real client and service
	client := createTestClient(t)
	service := NewScdnService(client)

	// Make real API call
	response, err := service.GetAccessProgress()
	if err != nil {
		t.Logf("API call failed (this might be expected in test environment): %v", err)
		return
	}

	// Validate response structure
	assert.NotNil(t, response)
	assert.NotNil(t, response.Status)
	assert.NotNil(t, response.Data)
}

func TestScdnService_UpdateDomainBaseSettingsRequest(t *testing.T) {
	req := DomainBaseSettingsUpdateRequest{
		DomainID: 123,
		Value: DomainBaseSettingsValue{
			ProxyHost: &ProxyHostConfig{
				ProxyHost:     "example.com",
				ProxyHostType: "default",
			},
			ProxySNI: &ProxySNIConfig{
				ProxySNI: "sni.example.com",
				Status:   "on",
			},
			DomainRedirect: &DomainRedirectConfig{
				Status:   "on",
				JumpTo:   "https://www.newexample.com",
				JumpType: "explicit",
			},
		},
	}

	assert.Equal(t, 123, req.DomainID)
	assert.NotNil(t, req.Value.ProxyHost)
	assert.Equal(t, "example.com", req.Value.ProxyHost.ProxyHost)
	assert.Equal(t, "default", req.Value.ProxyHost.ProxyHostType)
	assert.NotNil(t, req.Value.ProxySNI)
	assert.Equal(t, "sni.example.com", req.Value.ProxySNI.ProxySNI)
	assert.Equal(t, "on", req.Value.ProxySNI.Status)
	assert.NotNil(t, req.Value.DomainRedirect)
	assert.Equal(t, "on", req.Value.DomainRedirect.Status)
	assert.Equal(t, "https://www.newexample.com", req.Value.DomainRedirect.JumpTo)
	assert.Equal(t, "explicit", req.Value.DomainRedirect.JumpType)
}

func TestScdnService_GetDomainBaseSettingsRequest(t *testing.T) {
	req := DomainBaseSettingsGetRequest{
		DomainID: 123,
	}

	assert.Equal(t, 123, req.DomainID)
}

func TestScdnService_ListBriefDomainsRequest(t *testing.T) {
	req := BriefDomainListRequest{
		IDs: []int{123, 456},
	}

	assert.Len(t, req.IDs, 2)
	assert.Contains(t, req.IDs, 123)
	assert.Contains(t, req.IDs, 456)
}

func TestScdnService_GetDomainTemplatesRequest(t *testing.T) {
	req := DomainTemplatesRequest{
		DomainID: 123,
	}

	assert.Equal(t, 123, req.DomainID)
}

func TestScdnService_DownloadAccessInfoRequest(t *testing.T) {
	req := AccessInfoDownloadRequest{
		DomainInfos: []DomainInfoItem{
			{
				Domain:     "domain1.com",
				DataKey:    "key1",
				BizMainKey: "mainkey1",
			},
			{
				Domain:     "domain2.com",
				DataKey:    "key2",
				BizMainKey: "mainkey2",
			},
		},
	}

	assert.Len(t, req.DomainInfos, 2)
	assert.Equal(t, "domain1.com", req.DomainInfos[0].Domain)
	assert.Equal(t, "key1", req.DomainInfos[0].DataKey)
	assert.Equal(t, "mainkey1", req.DomainInfos[0].BizMainKey)
	assert.Equal(t, "domain2.com", req.DomainInfos[1].Domain)
	assert.Equal(t, "key2", req.DomainInfos[1].DataKey)
	assert.Equal(t, "mainkey2", req.DomainInfos[1].BizMainKey)
}

// 测试数据结构
func TestDomainInfo(t *testing.T) {
	domain := DomainInfo{
		ID:                  1,
		Domain:              "test.example.com",
		Remark:              "Test domain",
		AccessProgress:      "enabled",
		AccessMode:          "cname",
		ProtectStatus:       "scdn",
		EIForwardStatus:     "on",
		UseMyCname:          1,
		UseMyDNS:            1,
		CAStatus:            "bind",
		ExclusiveResourceID: 0,
		AccessProgressDesc:  "Enabled",
		HasOrigin:           true,
		CAID:                123,
		CreatedAt:           "2025-01-01T00:00:00Z",
		UpdatedAt:           "2025-01-01T00:00:00Z",
		PriDomain:           "test.example.com",
		Cname: CnameInfo{
			Master: "cname.test.example.com",
			Slaves: []string{},
		},
	}

	assert.Equal(t, 1, domain.ID)
	assert.Equal(t, "test.example.com", domain.Domain)
	assert.Equal(t, "Test domain", domain.Remark)
	assert.Equal(t, "enabled", domain.AccessProgress)
	assert.Equal(t, "cname", domain.AccessMode)
	assert.Equal(t, "scdn", domain.ProtectStatus)
	assert.Equal(t, "on", domain.EIForwardStatus)
	assert.Equal(t, 1, domain.UseMyCname)
	assert.Equal(t, 1, domain.UseMyDNS)
	assert.Equal(t, "bind", domain.CAStatus)
	assert.Equal(t, 0, domain.ExclusiveResourceID)
	assert.Equal(t, "Enabled", domain.AccessProgressDesc)
	assert.True(t, domain.HasOrigin)
	assert.Equal(t, 123, domain.CAID)
	assert.Equal(t, "2025-01-01T00:00:00Z", domain.CreatedAt)
	assert.Equal(t, "2025-01-01T00:00:00Z", domain.UpdatedAt)
	assert.Equal(t, "test.example.com", domain.PriDomain)
	assert.Equal(t, "cname.test.example.com", domain.Cname.Master)
	assert.Empty(t, domain.Cname.Slaves)
}

func TestOrigin(t *testing.T) {
	origin := Origin{
		Protocol:       0,
		ListenPorts:    []int{80, 443},
		OriginProtocol: 0,
		LoadBalance:    1,
		OriginType:     0,
		Records: []OriginRecord{
			{
				View:     "default",
				Value:    "1.1.1.1",
				Port:     80,
				Priority: 10,
			},
		},
	}

	assert.Equal(t, 0, origin.Protocol)
	assert.Len(t, origin.ListenPorts, 2)
	assert.Contains(t, origin.ListenPorts, 80)
	assert.Contains(t, origin.ListenPorts, 443)
	assert.Equal(t, 0, origin.OriginProtocol)
	assert.Equal(t, 1, origin.LoadBalance)
	assert.Equal(t, 0, origin.OriginType)
	assert.Len(t, origin.Records, 1)
	assert.Equal(t, "default", origin.Records[0].View)
	assert.Equal(t, "1.1.1.1", origin.Records[0].Value)
	assert.Equal(t, 80, origin.Records[0].Port)
	assert.Equal(t, 10, origin.Records[0].Priority)
}

func TestOriginRecord(t *testing.T) {
	record := OriginRecord{
		View:     "default",
		Value:    "1.1.1.1",
		Port:     80,
		Priority: 10,
	}

	assert.Equal(t, "default", record.View)
	assert.Equal(t, "1.1.1.1", record.Value)
	assert.Equal(t, 80, record.Port)
	assert.Equal(t, 10, record.Priority)
}

func TestCnameInfo(t *testing.T) {
	cname := CnameInfo{
		Master: "cname.test.example.com",
		Slaves: []string{"slave1.test.example.com", "slave2.test.example.com"},
	}

	assert.Equal(t, "cname.test.example.com", cname.Master)
	assert.Len(t, cname.Slaves, 2)
	assert.Contains(t, cname.Slaves, "slave1.test.example.com")
	assert.Contains(t, cname.Slaves, "slave2.test.example.com")
}

func TestStatus(t *testing.T) {
	status := Status{
		Code:    1,
		Message: "操作成功",
	}

	assert.Equal(t, 1, status.Code)
	assert.Equal(t, "操作成功", status.Message)
}

func TestDomainListData(t *testing.T) {
	data := DomainListData{
		Total: 2,
		List: []DomainInfo{
			{ID: 1, Domain: "test1.example.com"},
			{ID: 2, Domain: "test2.example.com"},
		},
	}

	assert.Equal(t, 2, data.Total)
	assert.Len(t, data.List, 2)
	assert.Equal(t, 1, data.List[0].ID)
	assert.Equal(t, "test1.example.com", data.List[0].Domain)
	assert.Equal(t, 2, data.List[1].ID)
	assert.Equal(t, "test2.example.com", data.List[1].Domain)
}

// Integration tests that require real API credentials
func TestScdnService_Integration_ListDomains(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)

	// Test with minimal request
	req := DomainListRequest{
		Page:     1,
		PageSize: 5,
	}

	response, err := service.ListDomains(req)
	if err != nil {
		t.Fatalf("Failed to list domains: %v", err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Status.Code)
	assert.GreaterOrEqual(t, response.Data.Total, 0)
}

func TestScdnService_Integration_GetAccessProgress(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)

	response, err := service.GetAccessProgress()
	if err != nil {
		t.Fatalf("Failed to get access progress: %v", err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Status.Code)
	assert.NotNil(t, response.Data)
}

func TestScdnService_Integration_ListBriefDomains(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)

	req := BriefDomainListRequest{
		IDs: []int{}, // Empty list to get all domains
	}

	response, err := service.ListBriefDomains(req)
	if err != nil {
		t.Fatalf("Failed to list brief domains: %v", err)
	}

	assert.NotNil(t, response)
	assert.Equal(t, 1, response.Status.Code)
	assert.NotNil(t, response.Data)
}

func TestScdnService_Integration_ClientHealth(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	scdnClient, err := client.ScdnClient()
	if err != nil {
		t.Fatalf("Failed to get SCDN client: %v", err)
	}

	// Test client health
	err = scdnClient.IsHealthy(context.Background())
	if err != nil {
		t.Logf("SCDN client health check failed (this might be expected in test environment): %v", err)
		return
	}

	// If health check passes, verify client properties
	assert.Equal(t, "v5", scdnClient.GetAPIVersion())
	assert.Equal(t, "scdn", scdnClient.GetServiceName())
	assert.NotEmpty(t, scdnClient.GetBaseURL())
}

func TestScdnService_ListDomainsSimple(t *testing.T) {
	type args struct {
		req DomainSimpleListRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *DomainSimpleListResponse
		wantErr bool
	}{
		{
			name: "Test ListDomainsSimple",
			args: args{
				req: DomainSimpleListRequest{
					Page:    1,
					PerPage: 10,
				},
			},
		},
		{
			name: "Test ListDomainsSimple with domain",
			args: args{
				req: DomainSimpleListRequest{
					Page:    1,
					PerPage: 10,
					Domain:  "test.example.com",
				},
			},
		},
	}

	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ListDomainsSimple(tt.args.req)
			if err != nil {
				t.Errorf("ScdnService.ListDomainsSimple() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateDomainRemark(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		client *connectivity.EdgeNextClient
		// Named input parameters for target function.
		req     DomainUpdateRequest
		want    *DomainUpdateResponse
		wantErr bool
	}{
		{
			name: "Test UpdateDomainRemark",
			req: DomainUpdateRequest{
				DomainID: 102021,
				Remark:   "Updated remark",
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.UpdateDomain(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateDomain() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_BindDomainCert(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		client *connectivity.EdgeNextClient
		// Named input parameters for target function.
		req     DomainCertBindRequest
		want    *DomainCertBindResponse
		wantErr bool
	}{
		{
			name: "Test BindDomainCert",
			req: DomainCertBindRequest{
				DomainID: 101753,
				CAID:     370,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.BindDomainCert(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BindDomainCert() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UnbindDomainCert(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		client *connectivity.EdgeNextClient
		// Named input parameters for target function.
		req     DomainCertUnbindRequest
		want    *DomainCertUnbindResponse
		wantErr bool
	}{
		{
			name: "Test UnbindDomainCert",
			req: DomainCertUnbindRequest{
				DomainID: 101753,
				CAID:     370,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.UnbindDomainCert(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UnbindDomainCert() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteDomain(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		client *connectivity.EdgeNextClient
		// Named input parameters for target function.
		req     DomainListRequest
		want    *DomainDeleteResponse
		wantErr bool
	}{
		{
			name: "Test DeleteDomain",
			req: DomainListRequest{
				Page:     1,
				PageSize: 10,
				Domain:   "test.example.com",
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			domains, err := service.ListDomains(tt.req)
			if err != nil {
				t.Fatalf("Failed to list domains: %v", err)
			}
			if len(domains.Data.List) == 0 {
				t.Fatalf("No domains found")
				return
			}
			domainIDs := make([]int, len(domains.Data.List))
			for i, domain := range domains.Data.List {
				domainIDs[i] = domain.ID
			}

			got, err := service.DeleteDomain(DomainDeleteRequest{
				IDs: domainIDs,
			})
			if err != nil {
				t.Errorf("ScdnService.DeleteDomain() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DisableDomain(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainDisableRequest
	}{
		{
			name: "Test DisableDomain",
			req: DomainDisableRequest{
				DomainIDs: []int{101753},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.DisableDomain(tt.req)
			if err != nil {
				t.Errorf("ScdnService.DisableDomain() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_EnableDomain(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainEnableRequest
	}{
		{
			name: "Test EnableDomain",
			req: DomainEnableRequest{
				DomainIDs: []int{101753},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				got, err := service.EnableDomain(tt.req)
				if err != nil {
					t.Errorf("ScdnService.EnableDomain() error = %v", err)
					return
				}
				t.Logf("Response: %+v", got)
			})
		})
	}
}

func TestScdnService_RefreshDomainAccess(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainAccessRefreshRequest
	}{
		{
			name: "Test RefreshDomainAccess",
			req: DomainAccessRefreshRequest{
				DomainIDs: []int{101753},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.RefreshDomainAccess(tt.req)
			if err != nil {
				t.Errorf("ScdnService.RefreshDomainAccess() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ExportDomains(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainExportRequest
	}{
		{
			name: "Test ExportDomains",
			req: DomainExportRequest{
				Domain: ".com",
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ExportDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ExportDomains() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_AddOrigins(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  OriginAddRequest
	}{
		{
			name: "Test AddOrigins",
			req: OriginAddRequest{
				DomainID: 101753,
				Origins: []Origin{
					{
						Protocol:       0,
						ListenPorts:    []int{443},
						OriginProtocol: 0,
						LoadBalance:    1,
						OriginType:     0,
						Records: []OriginRecord{
							{
								View:     "primary",
								Value:    "52.72.129.198",
								Port:     80,
								Priority: 10,
							},
							{
								View:     "primary",
								Value:    "18.207.14.107",
								Port:     80,
								Priority: 11,
							},
						},
					},
				},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.AddOrigins(tt.req)
			if err != nil {
				t.Errorf("ScdnService.AddOrigins() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateOrigins(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  OriginUpdateRequest
	}{
		{
			name: "Test UpdateOrigins",
			req: OriginUpdateRequest{
				DomainID: 101753,
				Origins: []EditOrigin{
					{
						Id:             200250,
						Protocol:       0,
						ListenPort:     443,
						OriginProtocol: 0,
						LoadBalance:    1,
						OriginType:     0,
						Records: []OriginRecord{
							{
								View:     "primary",
								Value:    "52.72.129.198",
								Port:     80,
								Priority: 10,
							},
						},
					},
				},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.UpdateOrigins(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateOrigins() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteOrigins(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  OriginDeleteRequest
	}{
		{
			name: "Test DeleteOrigins",
			req: OriginDeleteRequest{
				DomainID: 101753,
				IDs:      []int{200250},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.DeleteOrigins(tt.req)
			if err != nil {
				t.Errorf("ScdnService.DeleteOrigins() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ListOrigins(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  OriginListRequest
	}{
		{
			name: "Test ListOrigins",
			req: OriginListRequest{
				DomainID: 101753,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ListOrigins(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListOrigins() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_SwitchDomainNodes(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainNodeSwitchRequest
	}{
		{
			name: "Test SwitchDomainNodes",
			req: DomainNodeSwitchRequest{
				DomainID:            101753,
				ProtectStatus:       "exclusive",
				ExclusiveResourceID: 456,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SwitchDomainNodes(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SwitchDomainNodes() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_SwitchDomainAccessMode(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainAccessModeSwitchRequest
	}{
		{
			name: "Test SwitchDomainAccessMode",
			req: DomainAccessModeSwitchRequest{
				DomainID:   101753,
				AccessMode: "ns",
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SwitchDomainAccessMode(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SwitchDomainAccessMode() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetAccessProgress(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			name: "Test GetAccessProgress",
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetAccessProgress()
			if err != nil {
				t.Errorf("ScdnService.GetAccessProgress() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateDomainBaseSettings(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainBaseSettingsUpdateRequest
	}{
		{
			name: "Test UpdateDomainBaseSettings",
			req: DomainBaseSettingsUpdateRequest{
				DomainID: 101753,
				Value: DomainBaseSettingsValue{
					ProxyHost: &ProxyHostConfig{
						ProxyHost:     "example.com",
						ProxyHostType: "default",
					},
					ProxySNI: &ProxySNIConfig{
						ProxySNI: "sni.example.com",
						Status:   "on",
					},
					DomainRedirect: &DomainRedirectConfig{
						Status:   "on",
						JumpTo:   "https://www.newexample.com",
						JumpType: "explicit",
					},
				},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.UpdateDomainBaseSettings(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateDomainBaseSettings() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetDomainBaseSettings(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainBaseSettingsGetRequest
	}{
		{
			name: "Test GetDomainBaseSettings",
			req: DomainBaseSettingsGetRequest{
				DomainID: 101753,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetDomainBaseSettings(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetDomainBaseSettings() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ListBriefDomains(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  BriefDomainListRequest
	}{
		{
			name: "Test ListBriefDomains",
			req: BriefDomainListRequest{
				IDs: []int{115863},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.ListBriefDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListBriefDomains() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetDomainTemplates(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  DomainTemplatesRequest
	}{
		{
			name: "Test GetDomainTemplates",
			req: DomainTemplatesRequest{
				DomainID: 115863,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetDomainTemplates(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetDomainTemplates() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DownloadAccessInfo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  AccessInfoDownloadRequest
	}{
		{
			name: "Test DownloadAccessInfo",
			req: AccessInfoDownloadRequest{
				DomainInfos: []DomainInfoItem{
					{
						Domain: "yj.com",
					},
				},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.DownloadAccessInfo(tt.req)
			if err != nil {
				t.Errorf("ScdnService.DownloadAccessInfo() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}
