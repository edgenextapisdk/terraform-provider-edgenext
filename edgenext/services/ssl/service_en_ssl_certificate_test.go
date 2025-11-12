package ssl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// Mock HTTP server for testing
func createMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == "POST" && r.URL.Path == "/v2/domain/certificate":
			// Mock create/update SSL certificate
			var req SslCertificateRequest
			json.NewDecoder(r.Body).Decode(&req)

			if req.Certificate == "invalid-certificate" {
				response := SslCertificateResponse{
					Code: 1001,
					Msg:  "Invalid certificate format",
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := SslCertificateResponse{
				Code: 0,
				Data: SslCertificateData{
					CertID:         "12345",
					Name:           req.Name,
					Certificate:    req.Certificate,
					Key:            req.Key,
					BindDomains:    []string{"example.com"},
					CertStartTime:  "2024-01-01 00:00:00",
					CertExpireTime: "2025-01-01 00:00:00",
				},
			}
			json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && r.URL.Path == "/v2/domain/certificate":
			// Mock get SSL certificate or list certificates
			certID := r.URL.Query().Get("cert_id")
			pageNumber := r.URL.Query().Get("page_number")

			if certID != "" {
				// Single certificate query
				id, _ := strconv.Atoi(certID)
				if id == 99999 {
					// Certificate not found
					response := SslCertificateResponse{
						Code: 1002,
						Msg:  "Certificate not found",
					}
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(response)
					return
				}

				response := SslCertificateResponse{
					Code: 0,
					Data: SslCertificateData{
						CertID:         certID,
						Name:           "test-certificate",
						Certificate:    "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
						Key:            "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----",
						BindDomains:    []string{"example.com", "www.example.com"},
						CertStartTime:  "2024-01-01 00:00:00",
						CertExpireTime: "2025-01-01 00:00:00",
					},
				}
				json.NewEncoder(w).Encode(response)
			} else if pageNumber != "" {
				// Certificate list query
				page, _ := strconv.Atoi(pageNumber)
				pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

				var certificates []SslCertificateDataV2
				if page == 1 {
					certificates = []SslCertificateDataV2{
						{
							CertID:            "12345",
							Name:              "test-certificate-1",
							Certificate:       "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
							Key:               "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----",
							AssociatedDomains: []string{"example.com"},
							IncludeDomains:    []string{"*.example.com"},
							CertStartTime:     "2024-01-01 00:00:00",
							CertExpireTime:    "2025-01-01 00:00:00",
						},
						{
							CertID:            "12346",
							Name:              "test-certificate-2",
							Certificate:       "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
							Key:               "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----",
							AssociatedDomains: []string{"test.com"},
							IncludeDomains:    []string{"*.test.com"},
							CertStartTime:     "2024-01-01 00:00:00",
							CertExpireTime:    "2025-01-01 00:00:00",
						},
					}
				}

				response := SslCertificateListResponse{
					Code: 0,
					Data: SslCertificateListResponseData{
						List:        certificates,
						TotalNumber: len(certificates),
						PageNumber:  page,
						PageSize:    pageSize,
					},
				}
				json.NewEncoder(w).Encode(response)
			}

		case r.Method == "DELETE" && r.URL.Path == "/v2/domain/certificate":
			// Mock delete SSL certificate
			var req DeleteSslCertificateRequest
			json.NewDecoder(r.Body).Decode(&req)

			if req.CertID == 99999 {
				// Certificate not found
				response := SslCertificateDeleteResponse{
					Code: 1002,
					Data: "Certificate not found",
				}
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(response)
				return
			}

			response := SslCertificateDeleteResponse{
				Code: 0,
				Data: "Certificate deleted successfully",
			}
			json.NewEncoder(w).Encode(response)

		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"code": 404, "msg": "endpoint not found"}`)
		}
	}))
}

// Helper function to create test service with mock server
func createTestService(server *httptest.Server) *SslCertificateService {
	// Create EdgeNextClient with mock API client
	config := &connectivity.Config{
		AccessKey: "test-key",
		SecretKey: "test-secret",
		Endpoint:  server.URL,
	}

	client, _ := config.Client()

	return NewSslCertificateService(client)
}

func TestNewSslCertificateService(t *testing.T) {
	config := &connectivity.Config{
		AccessKey: "test-key",
		SecretKey: "test-secret",
		Endpoint:  "http://test.example.com",
	}
	client, _ := config.Client()
	service := NewSslCertificateService(client)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.client != client {
		t.Fatal("Expected service client to match input client")
	}
}

func TestCreateOrUpdateSslCertificate(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	t.Run("SuccessfulCreate", func(t *testing.T) {
		req := SslCertificateRequest{
			Name:        "test-certificate",
			Certificate: "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
			Key:         "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----",
		}

		response, err := service.CreateOrUpdateSslCertificate(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.Name != req.Name {
			t.Fatalf("Expected certificate name %s, got: %s", req.Name, response.Data.Name)
		}

		if response.Data.CertID == "" {
			t.Fatal("Expected certificate ID to be set")
		}
	})

	t.Run("SuccessfulUpdate", func(t *testing.T) {
		certID := 12345
		req := SslCertificateRequest{
			Name:        "updated-certificate",
			Certificate: "-----BEGIN CERTIFICATE-----\nUPDATED...\n-----END CERTIFICATE-----",
			Key:         "-----BEGIN PRIVATE KEY-----\nUPDATED...\n-----END PRIVATE KEY-----",
			CertID:      &certID,
		}

		response, err := service.CreateOrUpdateSslCertificate(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.Name != req.Name {
			t.Fatalf("Expected certificate name %s, got: %s", req.Name, response.Data.Name)
		}
	})

	t.Run("InvalidCertificate", func(t *testing.T) {
		req := SslCertificateRequest{
			Name:        "invalid-certificate",
			Certificate: "invalid-certificate",
			Key:         "invalid-key",
		}

		_, err := service.CreateOrUpdateSslCertificate(req)
		if err == nil {
			t.Fatal("Expected error for invalid certificate, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})

	t.Run("EmptyRequest", func(t *testing.T) {
		req := SslCertificateRequest{}

		response, err := service.CreateOrUpdateSslCertificate(req)
		if err != nil {
			t.Fatalf("Expected no error for empty request, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})
}

func TestGetSslCertificate(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	t.Run("SuccessfulGet", func(t *testing.T) {
		certID := 12345

		response, err := service.GetSslCertificate(certID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.CertID != strconv.Itoa(certID) {
			t.Fatalf("Expected certificate ID %d, got: %s", certID, response.Data.CertID)
		}

		if response.Data.Name != "test-certificate" {
			t.Fatalf("Expected certificate name 'test-certificate', got: %s", response.Data.Name)
		}

		if len(response.Data.BindDomains) == 0 {
			t.Fatal("Expected bind domains to be populated")
		}
	})

	t.Run("CertificateNotFound", func(t *testing.T) {
		certID := 99999

		_, err := service.GetSslCertificate(certID)
		if err == nil {
			t.Fatal("Expected error for non-existent certificate, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})

	t.Run("ValidCertificateID", func(t *testing.T) {
		certID := 226109

		response, err := service.GetSslCertificate(certID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.Certificate == "" {
			t.Fatal("Expected certificate content to be populated")
		}

		if response.Data.Key == "" {
			t.Fatal("Expected private key to be populated")
		}

		if response.Data.CertStartTime == "" {
			t.Fatal("Expected certificate start time to be populated")
		}

		if response.Data.CertExpireTime == "" {
			t.Fatal("Expected certificate expiration time to be populated")
		}
	})
}

func TestListSslCertificates(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	t.Run("SuccessfulList", func(t *testing.T) {
		pageNumber := 1
		pageSize := 10

		response, err := service.ListSslCertificates(pageNumber, pageSize)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.PageNumber != pageNumber {
			t.Fatalf("Expected page number %d, got: %d", pageNumber, response.Data.PageNumber)
		}

		if response.Data.PageSize != pageSize {
			t.Fatalf("Expected page size %d, got: %d", pageSize, response.Data.PageSize)
		}

		if len(response.Data.List) != response.Data.TotalNumber {
			t.Fatalf("Expected list length to match total number: %d vs %d", len(response.Data.List), response.Data.TotalNumber)
		}

		for _, cert := range response.Data.List {
			if cert.CertID == "" {
				t.Fatal("Expected certificate ID to be populated")
			}

			if cert.Name == "" {
				t.Fatal("Expected certificate name to be populated")
			}

			if len(cert.AssociatedDomains) == 0 {
				t.Fatal("Expected associated domains to be populated")
			}
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		pageNumber := 2
		pageSize := 10

		response, err := service.ListSslCertificates(pageNumber, pageSize)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) != 0 {
			t.Fatalf("Expected empty list, got: %d items", len(response.Data.List))
		}

		if response.Data.TotalNumber != 0 {
			t.Fatalf("Expected total number 0, got: %d", response.Data.TotalNumber)
		}
	})

	t.Run("DifferentPageSizes", func(t *testing.T) {
		testCases := []struct {
			pageNumber int
			pageSize   int
		}{
			{1, 5},
			{1, 20},
			{1, 50},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("Page%d_Size%d", tc.pageNumber, tc.pageSize), func(t *testing.T) {
				response, err := service.ListSslCertificates(tc.pageNumber, tc.pageSize)
				if err != nil {
					t.Fatalf("Expected no error, got: %v", err)
				}

				if response.Data.PageSize != tc.pageSize {
					t.Fatalf("Expected page size %d, got: %d", tc.pageSize, response.Data.PageSize)
				}
			})
		}
	})
}

func TestDeleteSslCertificate(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	t.Run("SuccessfulDelete", func(t *testing.T) {
		req := DeleteSslCertificateRequest{
			CertID: 12345,
		}

		err := service.DeleteSslCertificate(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
	})

	t.Run("CertificateNotFound", func(t *testing.T) {
		req := DeleteSslCertificateRequest{
			CertID: 99999,
		}

		err := service.DeleteSslCertificate(req)
		if err == nil {
			t.Fatal("Expected error for non-existent certificate, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})

	t.Run("ValidCertificateIDs", func(t *testing.T) {
		testCases := []int{226105, 12345, 67890}

		for _, certID := range testCases {
			t.Run(fmt.Sprintf("CertID_%d", certID), func(t *testing.T) {
				req := DeleteSslCertificateRequest{
					CertID: certID,
				}

				err := service.DeleteSslCertificate(req)
				if err != nil {
					t.Fatalf("Expected no error for certificate ID %d, got: %v", certID, err)
				}
			})
		}
	})
}

// Benchmark tests
func BenchmarkCreateOrUpdateSslCertificate(b *testing.B) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	req := SslCertificateRequest{
		Name:        "benchmark-certificate",
		Certificate: "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----",
		Key:         "-----BEGIN PRIVATE KEY-----\nMII...\n-----END PRIVATE KEY-----",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CreateOrUpdateSslCertificate(req)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGetSslCertificate(b *testing.B) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetSslCertificate(12345)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkListSslCertificates(b *testing.B) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ListSslCertificates(1, 10)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkDeleteSslCertificate(b *testing.B) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	req := DeleteSslCertificateRequest{
		CertID: 12345,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := service.DeleteSslCertificate(req)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Table-driven tests for error scenarios
func TestSslCertificateServiceErrorScenarios(t *testing.T) {
	server := createMockServer()
	defer server.Close()

	service := createTestService(server)

	errorTests := []struct {
		name        string
		testFunc    func() error
		expectError bool
		errorMsg    string
	}{
		{
			name: "CreateInvalidCertificate",
			testFunc: func() error {
				req := SslCertificateRequest{
					Certificate: "invalid-certificate",
				}
				_, err := service.CreateOrUpdateSslCertificate(req)
				return err
			},
			expectError: true,
			errorMsg:    "Invalid certificate format",
		},
		{
			name: "GetNonExistentCertificate",
			testFunc: func() error {
				_, err := service.GetSslCertificate(99999)
				return err
			},
			expectError: true,
			errorMsg:    "Certificate not found",
		},
		{
			name: "DeleteNonExistentCertificate",
			testFunc: func() error {
				req := DeleteSslCertificateRequest{CertID: 99999}
				return service.DeleteSslCertificate(req)
			},
			expectError: true,
			errorMsg:    "Certificate not found",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()

			if tt.expectError && err == nil {
				t.Fatalf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if tt.expectError && err != nil {
				if tt.errorMsg != "" && err.Error() != "" {
					t.Logf("Error message: %s", err.Error())
				}
			}
		})
	}
}
