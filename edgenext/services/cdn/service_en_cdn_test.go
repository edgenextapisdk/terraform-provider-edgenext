package cdn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// Mock HTTP server for CDN testing
func createMockCDNServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == "POST" && r.URL.Path == "/v2/domain":
			handleCreateDomain(w, r)
		case r.Method == "GET" && r.URL.Path == "/v2/domain":
			handleGetDomain(w, r)
		case r.Method == "GET" && r.URL.Path == "/v2/domain/list":
			handleListDomains(w, r)
		case r.Method == "DELETE" && r.URL.Path == "/v2/domain":
			handleDeleteDomain(w, r)
		case r.Method == "POST" && r.URL.Path == "/v2/domain/config":
			handleSetDomainConfig(w, r)
		case r.Method == "GET" && r.URL.Path == "/v2/domain/config":
			handleGetDomainConfig(w, r)
		case r.Method == "DELETE" && r.URL.Path == "/v2/domain/config":
			handleDeleteDomainConfig(w, r)
		case r.Method == "POST" && r.URL.Path == "/v2/cache/refresh":
			handleCacheRefresh(w, r)
		case r.Method == "GET" && r.URL.Path == "/v2/cache/refresh":
			handleQueryCacheRefresh(w, r)
		case r.Method == "POST" && r.URL.Path == "/v2/cache/prefetch":
			handleFilePrefetch(w, r)
		case r.Method == "GET" && r.URL.Path == "/v2/cache/prefetch":
			handleQueryFilePrefetch(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"code": 404, "msg": "endpoint not found"}`)
		}
	}))
}

func handleCreateDomain(w http.ResponseWriter, r *http.Request) {
	var req DomainCreateRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Domain == "invalid.domain" {
		response := DomainResponse{
			Code: 1001,
			Msg:  "Invalid domain format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := DomainResponse{Code: 0}
	json.NewEncoder(w).Encode(response)
}

func handleGetDomain(w http.ResponseWriter, r *http.Request) {
	domains := r.URL.Query().Get("domains")

	if domains == "notfound.example.com" {
		response := GetDomainResponse{
			Code: 1002,
			Msg:  "Domain not found",
			Data: []DomainData{},
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := GetDomainResponse{
		Code: 0,
		Data: []DomainData{
			{
				ID:         "12345",
				Domain:     domains,
				Type:       "page",
				Status:     "serving",
				IcpStatus:  "yes",
				IcpNum:     "ICP-12345678",
				Area:       "mainland_china",
				Cname:      "example.com.edgenext.net",
				CreateTime: "2024-01-01 12:00:00",
				UpdateTime: "2024-01-01 12:00:00",
				Https:      1,
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleListDomains(w http.ResponseWriter, r *http.Request) {
	response := DomainListResponse{
		Code: 0,
		Data: DomainListResponseData{
			List: []DomainData{
				{
					ID:         "12345",
					Domain:     "example1.com",
					Type:       "page",
					Status:     "serving",
					IcpStatus:  "yes",
					Area:       "mainland_china",
					Cname:      "example1.com.edgenext.net",
					CreateTime: "2024-01-01 12:00:00",
					UpdateTime: "2024-01-01 12:00:00",
					Https:      1,
				},
				{
					ID:         "12346",
					Domain:     "example2.com",
					Type:       "download",
					Status:     "serving",
					IcpStatus:  "yes",
					Area:       "global",
					Cname:      "example2.com.edgenext.net",
					CreateTime: "2024-01-02 12:00:00",
					UpdateTime: "2024-01-02 12:00:00",
					Https:      0,
				},
			},
			TotalNumber: "2",
			PageNumber:  "1",
			PageSize:    10,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleDeleteDomain(w http.ResponseWriter, r *http.Request) {
	domains := r.URL.Query().Get("domains")

	if domains == "notfound.example.com" {
		response := DomainResponse{
			Code: 1002,
			Msg:  "Domain not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := DomainResponse{Code: 0}
	json.NewEncoder(w).Encode(response)
}

func handleSetDomainConfig(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	response := DomainConfigResponse{
		Code: 0,
		Data: map[string]interface{}{"success": true},
	}
	json.NewEncoder(w).Encode(response)
}

func handleGetDomainConfig(w http.ResponseWriter, r *http.Request) {
	domains := r.URL.Query().Get("domains")

	if domains == "notfound.example.com" {
		response := GetDomainConfigResponse{
			Code: 1002,
			Msg:  "Domain not found",
			Data: []DomainConfigItem{},
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := GetDomainConfigResponse{
		Code: 0,
		Data: []DomainConfigItem{
			{
				Domain:   domains,
				DomainID: "12345",
				Config: ConfigItem{
					Origin: &OriginItem{
						DefaultMaster: "1.2.3.4",
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleDeleteDomainConfig(w http.ResponseWriter, r *http.Request) {
	response := DomainResponse{Code: 0}
	json.NewEncoder(w).Encode(response)
}

func handleCacheRefresh(w http.ResponseWriter, r *http.Request) {
	var req CacheRefreshRequest
	json.NewDecoder(r.Body).Decode(&req)

	if len(req.URLs) > 0 && req.URLs[0] == "http://invalid.domain/path" {
		response := CacheRefreshResponse{
			Code: 1001,
			Msg:  "Invalid URL",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CacheRefreshResponse{
		Code: 0,
		Data: CacheRefreshData{
			TaskID: "task-12345",
			Count:  len(req.URLs),
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleQueryCacheRefresh(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")

	if taskID == "999999" {
		response := CacheRefreshQueryResponse{
			Code: 1002,
			Msg:  "Task not found",
			Data: CacheRefreshQueryData{
				List: []CacheRefreshQueryItem{},
			},
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CacheRefreshQueryResponse{
		Code: 0,
		Data: CacheRefreshQueryData{
			Total:      1,
			PageNumber: 1,
			List: []CacheRefreshQueryItem{
				{
					ID:           "1",
					URL:          "http://example.com/path",
					Status:       "completed",
					CreateTime:   "2024-01-01 11:59:00",
					CompleteTime: "2024-01-01 12:00:00",
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleFilePrefetch(w http.ResponseWriter, r *http.Request) {
	var req FilePrefetchRequest
	json.NewDecoder(r.Body).Decode(&req)

	if len(req.URLs) > 0 && req.URLs[0] == "http://invalid.domain/file" {
		response := FilePrefetchResponse{
			Code: 1001,
			Msg:  "Invalid URL",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := FilePrefetchResponse{
		Code: 0,
		Data: FilePrefetchData{
			TaskID: "task-67890",
			Count:  len(req.URLs),
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleQueryFilePrefetch(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("task_id")

	if taskID == "999999" {
		response := FilePrefetchQueryResponse{
			Code: 1002,
			Msg:  "Task not found",
			Data: FilePrefetchQueryData{
				List: []FilePrefetchQueryItem{},
			},
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := FilePrefetchQueryResponse{
		Code: 0,
		Data: FilePrefetchQueryData{
			Total:      1,
			PageNumber: 1,
			List: []FilePrefetchQueryItem{
				{
					ID:           "1",
					URL:          "http://example.com/file.txt",
					Status:       "completed",
					CreateTime:   "2024-01-01 11:59:00",
					CompleteTime: "2024-01-01 12:00:00",
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func createTestCDNService(server *httptest.Server) *CdnService {
	config := &connectivity.Config{
		AccessKey: "test-key",
		SecretKey: "test-secret",
		Endpoint:  server.URL,
	}
	client, _ := config.Client()
	return NewCdnService(client)
}

func TestNewCdnService(t *testing.T) {
	config := &connectivity.Config{
		AccessKey: "test-key",
		SecretKey: "test-secret",
		Endpoint:  "http://test.example.com",
	}
	client, _ := config.Client()
	service := NewCdnService(client)

	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}

	if service.client != client {
		t.Fatal("Expected service client to match input client")
	}
}

func TestCreateDomain(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulCreate", func(t *testing.T) {
		req := DomainCreateRequest{
			Domain: "example.com",
			Area:   "global",
			Type:   "page",
			Config: map[string]interface{}{},
		}

		response, err := service.CreateDomain(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})

	t.Run("InvalidDomain", func(t *testing.T) {
		req := DomainCreateRequest{
			Domain: "invalid.domain",
			Type:   "page",
		}

		_, err := service.CreateDomain(req)
		if err == nil {
			t.Fatal("Expected error for invalid domain, got nil")
		}
	})
}

func TestGetDomain(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulGet", func(t *testing.T) {
		domain := "example.com"

		response, err := service.GetDomain(domain)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data) == 0 {
			t.Fatal("Expected domain data to be returned")
		}

		if response.Data[0].Domain != domain {
			t.Fatalf("Expected domain %s, got: %s", domain, response.Data[0].Domain)
		}
	})

	t.Run("DomainNotFound", func(t *testing.T) {
		domain := "notfound.example.com"

		_, err := service.GetDomain(domain)
		if err == nil {
			t.Fatal("Expected error for non-existent domain, got nil")
		}
	})
}

func TestListDomains(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulList", func(t *testing.T) {
		req := DomainListRequest{
			PageNumber: 1,
			PageSize:   10,
		}

		response, err := service.ListDomains(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) == 0 {
			t.Fatal("Expected domain list to be returned")
		}

		if response.Data.TotalNumber == "" {
			t.Fatal("Expected total number to be set")
		}
	})
}

func TestDeleteDomain(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulDelete", func(t *testing.T) {
		domain := "example.com"

		err := service.DeleteDomain(domain)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
	})

	t.Run("DomainNotFound", func(t *testing.T) {
		domain := "notfound.example.com"

		err := service.DeleteDomain(domain)
		if err == nil {
			t.Fatal("Expected error for non-existent domain, got nil")
		}
	})
}

func TestSetDomainConfig(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulSet", func(t *testing.T) {
		domain := "example.com"
		config := map[string]interface{}{
			"origin": map[string]interface{}{
				"default_master": "1.2.3.4",
			},
		}

		response, err := service.SetDomainConfig(domain, config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})

	t.Run("EmptyConfig", func(t *testing.T) {
		domain := "example.com"
		config := map[string]interface{}{}

		response, err := service.SetDomainConfig(domain, config)
		if err != nil {
			t.Fatalf("Expected no error for empty config, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})
}

func TestGetDomainConfig(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulGet", func(t *testing.T) {
		domain := "example.com"

		response, err := service.GetDomainConfig(domain, nil)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data) == 0 {
			t.Fatal("Expected config data to be returned")
		}
	})

	t.Run("WithConfigItems", func(t *testing.T) {
		domain := "example.com"
		configItems := []string{"origin", "cache_rule"}

		response, err := service.GetDomainConfig(domain, configItems)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})

	t.Run("DomainNotFound", func(t *testing.T) {
		domain := "notfound.example.com"

		_, err := service.GetDomainConfig(domain, nil)
		if err == nil {
			t.Fatal("Expected error for non-existent domain, got nil")
		}
	})
}

func TestDeleteDomainConfig(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulDelete", func(t *testing.T) {
		req := DeleteDomainConfigRequest{
			Domains: "example.com",
			Config:  []string{"origin"},
		}

		err := service.DeleteDomainConfig(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
	})
}

func TestCacheRefresh(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulRefresh", func(t *testing.T) {
		urls := []string{"http://example.com/path"}
		refreshType := "url"

		response, err := service.CacheRefresh(urls, refreshType)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.TaskID == "" {
			t.Fatal("Expected task ID to be set")
		}

		if response.Data.Count != len(urls) {
			t.Fatalf("Expected count %d, got: %d", len(urls), response.Data.Count)
		}
	})

	t.Run("InvalidURL", func(t *testing.T) {
		urls := []string{"http://invalid.domain/path"}
		refreshType := "url"

		_, err := service.CacheRefresh(urls, refreshType)
		if err == nil {
			t.Fatal("Expected error for invalid URL, got nil")
		}
	})
}

func TestQueryCacheRefresh(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulQueryByTaskID", func(t *testing.T) {
		req := CacheRefreshQueryRequest{
			TaskID: 12345,
		}

		response, err := service.QueryCacheRefresh(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) == 0 {
			t.Fatal("Expected task data to be returned")
		}
	})

	t.Run("QueryByTaskIDHelper", func(t *testing.T) {
		taskID := 12345

		response, err := service.QueryCacheRefreshByTaskID(taskID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		req := CacheRefreshQueryRequest{
			TaskID: 999999,
		}

		_, err := service.QueryCacheRefresh(req)
		if err == nil {
			t.Fatal("Expected error for non-existent task, got nil")
		}
	})
}

func TestFilePrefetch(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulPrefetch", func(t *testing.T) {
		urls := []string{"http://example.com/file.txt"}

		response, err := service.FilePrefetch(urls)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.TaskID == "" {
			t.Fatal("Expected task ID to be set")
		}

		if response.Data.Count != len(urls) {
			t.Fatalf("Expected count %d, got: %d", len(urls), response.Data.Count)
		}
	})

	t.Run("InvalidURL", func(t *testing.T) {
		urls := []string{"http://invalid.domain/file"}

		_, err := service.FilePrefetch(urls)
		if err == nil {
			t.Fatal("Expected error for invalid URL, got nil")
		}
	})
}

func TestQueryFilePrefetch(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulQueryByTaskID", func(t *testing.T) {
		req := FilePrefetchQueryRequest{
			TaskID: 67890,
		}

		response, err := service.QueryFilePrefetch(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) == 0 {
			t.Fatal("Expected task data to be returned")
		}
	})

	t.Run("QueryByTaskIDHelper", func(t *testing.T) {
		taskID := 67890

		response, err := service.QueryFilePrefetchByTaskID(taskID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		req := FilePrefetchQueryRequest{
			TaskID: 999999,
		}

		_, err := service.QueryFilePrefetch(req)
		if err == nil {
			t.Fatal("Expected error for non-existent task, got nil")
		}
	})
}

// Test FlexibleInt type
func TestFlexibleIntUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		jsonInput     string
		expectedValue int
		expectedError bool
	}{
		{
			name:          "Parse integer value",
			jsonInput:     "443",
			expectedValue: 443,
			expectedError: false,
		},
		{
			name:          "Parse string value",
			jsonInput:     "\"80\"",
			expectedValue: 80,
			expectedError: false,
		},
		{
			name:          "Parse zero integer",
			jsonInput:     "0",
			expectedValue: 0,
			expectedError: false,
		},
		{
			name:          "Parse empty string",
			jsonInput:     "\"\"",
			expectedValue: 0,
			expectedError: false,
		},
		{
			name:          "Parse invalid string",
			jsonInput:     "\"invalid\"",
			expectedValue: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi FlexibleInt
			err := json.Unmarshal([]byte(tt.jsonInput), &fi)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if fi.Int() != tt.expectedValue {
					t.Errorf("Expected %d, got %d", tt.expectedValue, fi.Int())
				}
			}
		})
	}
}

func TestFlexibleIntMarshal(t *testing.T) {
	fi := FlexibleInt(443)
	data, err := json.Marshal(fi)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "443"
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

// Benchmark tests
func BenchmarkCreateDomain(b *testing.B) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	req := DomainCreateRequest{
		Domain: "benchmark.com",
		Area:   "global",
		Type:   "page",
		Config: map[string]interface{}{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateDomain(req)
	}
}

func BenchmarkGetDomain(b *testing.B) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetDomain("example.com")
	}
}

func BenchmarkCacheRefresh(b *testing.B) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	urls := []string{"http://example.com/path"}
	refreshType := "url"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CacheRefresh(urls, refreshType)
	}
}

// Test error scenarios
func TestCDNServiceErrorScenarios(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	errorTests := []struct {
		name        string
		testFunc    func() error
		expectError bool
	}{
		{
			name: "CreateInvalidDomain",
			testFunc: func() error {
				req := DomainCreateRequest{Domain: "invalid.domain", Type: "page"}
				_, err := service.CreateDomain(req)
				return err
			},
			expectError: true,
		},
		{
			name: "GetNonExistentDomain",
			testFunc: func() error {
				_, err := service.GetDomain("notfound.example.com")
				return err
			},
			expectError: true,
		},
		{
			name: "DeleteNonExistentDomain",
			testFunc: func() error {
				return service.DeleteDomain("notfound.example.com")
			},
			expectError: true,
		},
		{
			name: "GetNonExistentDomainConfig",
			testFunc: func() error {
				_, err := service.GetDomainConfig("notfound.example.com", nil)
				return err
			},
			expectError: true,
		},
		{
			name: "CacheRefreshInvalidURL",
			testFunc: func() error {
				urls := []string{"http://invalid.domain/path"}
				_, err := service.CacheRefresh(urls, "url")
				return err
			},
			expectError: true,
		},
		{
			name: "QueryNonExistentTask",
			testFunc: func() error {
				req := CacheRefreshQueryRequest{TaskID: 999999}
				_, err := service.QueryCacheRefresh(req)
				return err
			},
			expectError: true,
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
				t.Logf("Got expected error: %s", err.Error())
			}
		})
	}
}
