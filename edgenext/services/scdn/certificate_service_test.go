package scdn

import (
	"strings"
	"testing"
)

func TestScdnService_SaveCertificate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CATextSaveRequest
	}{
		{
			name: "Test SaveCertificate",
			req: CATextSaveRequest{
				CAName: "test",
				CACert: "-----BEGIN CERTIFICATE-----\nMIIFnDCCBISWR1xzhs4YY\n-----END CERTIFICATE-----",
				CAKey:  "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA8pe2L1RRe8JfdnsJiUUmOgtJA==\n-----END RSA PRIVATE KEY-----",
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
			got, err := service.SaveCertificate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SaveCertificate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ListCertificates(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CASelfListRequest
	}{
		{
			name: "Test ListCertificates",
			req: CASelfListRequest{
				Page:    1,
				PerPage: 10,
				Domain:  ".com",
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
			got, err := service.ListCertificates(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListCertificates() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetCertificateDetail(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CASelfDetailRequest
	}{
		{
			name: "Test GetCertificateDetail",
			req: CASelfDetailRequest{
				ID: 375, // Replace with actual certificate ID for testing
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
			got, err := service.GetCertificateDetail(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetCertificateDetail() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteCertificate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CASelfDeleteRequest
	}{
		{
			name: "Test DeleteCertificate",
			req: CASelfDeleteRequest{
				IDs: "1", // Replace with actual certificate IDs for testing
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
			got, err := service.DeleteCertificate(tt.req)
			if err != nil {
				if !strings.Contains(err.Error(), "code: 41000") {
					t.Errorf("ScdnService.DeleteCertificate() error = %v", err)
					return
				}
				t.Logf("err: %s", err.Error())
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_EditCertificateName(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CAEditNameRequest
	}{
		{
			name: "Test EditCertificateName",
			req: CAEditNameRequest{
				ID:     375, // Replace with actual certificate ID for testing
				CAName: "test-certificate-name",
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
			got, err := service.EditCertificateName(tt.req)
			if err != nil {
				t.Errorf("ScdnService.EditCertificateName() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ListCertificatesByDomains(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CABatchListRequest
	}{
		{
			name: "Test ListCertificatesByDomains",
			req: CABatchListRequest{
				Domains: []string{"test7.meipk.com"},
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
			got, err := service.ListCertificatesByDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListCertificatesByDomains() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ApplyCertificate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CAApplyAddRequest
	}{
		{
			name: "Test ApplyCertificate",
			req: CAApplyAddRequest{
				Domain: []string{"terraform.example.com"},
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
			got, err := service.ApplyCertificate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ApplyCertificate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ExportCertificate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  CASelfExportRequest
	}{
		{
			name: "Test ExportCertificate",
			req: CASelfExportRequest{
				ID: "372", // Replace with actual certificate ID for testing
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
			got, err := service.ExportCertificate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ExportCertificate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}
