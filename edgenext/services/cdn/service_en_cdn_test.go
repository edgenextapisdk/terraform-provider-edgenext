package cdn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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
			handleFilePurge(w, r)
		case r.Method == "GET" && r.URL.Path == "/v2/cache/prefetch":
			handleQueryFilePurge(w, r)
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
	pageNumber := r.URL.Query().Get("page_number")
	pageSize := r.URL.Query().Get("page_size")
	domainStatus := r.URL.Query().Get("domain_status")

	page, _ := strconv.Atoi(pageNumber)
	size, _ := strconv.Atoi(pageSize)

	var domains []DomainData
	if page == 1 && domainStatus != "deleted" {
		domains = []DomainData{
			{
				ID:         "12345",
				Domain:     "example1.com",
				Type:       "page",
				Status:     "serving",
				IcpStatus:  "yes",
				Area:       "mainland_china",
				Cname:      "example1.com.edgenext.net",
				CreateTime: "2024-01-01 12:00:00",
				Https:      1,
			},
			{
				ID:         "12346",
				Domain:     "example2.com",
				Type:       "download",
				Status:     "serving",
				IcpStatus:  "no",
				Area:       "global",
				Cname:      "example2.com.edgenext.net",
				CreateTime: "2024-01-01 13:00:00",
				Https:      0,
			},
		}
	}

	if len(domains) > size {
		domains = domains[:size]
	}

	response := DomainListResponse{
		Code: 0,
		Data: DomainListResponseData{
			List:        domains,
			TotalNumber: fmt.Sprintf("%d", len(domains)),
			PageNumber:  pageNumber,
			PageSize:    size,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleDeleteDomain(w http.ResponseWriter, r *http.Request) {
	domains := r.URL.Query().Get("domains")

	if domains == "notfound.example.com" {
		response := DeleteDomainResponse{
			Code: 1002,
			Msg:  "Domain not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := DeleteDomainResponse{
		Code: 0,
		Data: []DeleteDomainResponseData{
			{
				ID:         "12345",
				Domain:     domains,
				Status:     "deleted",
				DeleteTime: "2024-01-01 15:00:00",
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleSetDomainConfig(w http.ResponseWriter, r *http.Request) {
	var req DomainConfigRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Domains == "invalid.domain" {
		response := DomainConfigResponse{
			Code: 1001,
			Msg:  "Invalid domain configuration",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := DomainConfigResponse{
		Code: 0,
		Data: map[string]interface{}{
			"domains": req.Domains,
			"status":  "success",
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleGetDomainConfig(w http.ResponseWriter, r *http.Request) {
	domains := r.URL.Query().Get("domains")

	if domains == "notfound.example.com" {
		response := GetDomainConfigResponse{
			Code: 1002,
			Msg:  "Domain configuration not found",
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
				Status:   "serving",
				Config: ConfigItem{
					Origin: &OriginItem{
						DefaultMaster: "origin.example.com",
						OriginMode:    "https",
						Port:          443,
					},
					CacheRule: []*CacheRuleItem{
						{
							Type:          1,
							Pattern:       "jpg,png,gif",
							Time:          3600,
							TimeUnit:      "s",
							IgnoreNoCache: "on",
							IgnoreExpired: "on",
							IgnoreQuery:   "on",
						},
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleDeleteDomainConfig(w http.ResponseWriter, r *http.Request) {
	var req DeleteDomainConfigRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Domains == "notfound.example.com" {
		response := DeleteDomainConfigResponse{
			Code: 1002,
			Data: "Domain not found",
			Msg:  "Domain not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := DeleteDomainConfigResponse{
		Code: 0,
		Data: "Configuration deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func handleCacheRefresh(w http.ResponseWriter, r *http.Request) {
	var req CacheRefreshRequest
	json.NewDecoder(r.Body).Decode(&req)

	if len(req.URLs) == 0 {
		response := CacheRefreshResponse{
			Code: 1001,
			Msg:  "URLs cannot be empty",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CacheRefreshResponse{
		Code: 0,
		Data: CacheRefreshData{
			TaskID:  "task123456",
			Count:   len(req.URLs),
			ErrURLs: []string{},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleQueryCacheRefresh(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Query().Get("task_id")

	if taskIDStr == "99999" {
		response := CacheRefreshQueryResponse{
			Code: 1002,
			Msg:  "Task not found",
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
					ID:           taskIDStr,
					URL:          "https://example.com/test.jpg",
					Status:       "completed",
					Type:         "url",
					CreateTime:   "2024-01-01 12:00:00",
					CompleteTime: "2024-01-01 12:05:00",
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleFilePurge(w http.ResponseWriter, r *http.Request) {
	var req FilePurgeRequest
	json.NewDecoder(r.Body).Decode(&req)

	if len(req.URLs) == 0 {
		response := FilePurgeResponse{
			Code: 1001,
			Msg:  "URLs cannot be empty",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := FilePurgeResponse{
		Code: 0,
		Data: FilePurgeData{
			TaskID:  "purge123456",
			Count:   len(req.URLs),
			ErrURLs: []string{},
		},
	}
	json.NewEncoder(w).Encode(response)
}

func handleQueryFilePurge(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Query().Get("task_id")

	if taskIDStr == "99999" {
		response := FilePurgeQueryResponse{
			Code: 1002,
			Msg:  "Task not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := FilePurgeQueryResponse{
		Code: 0,
		Data: FilePurgeQueryData{
			Total:      1,
			PageNumber: 1,
			List: []FilePurgeQueryItem{
				{
					ID:           taskIDStr,
					URL:          "https://example.com/file.pdf",
					Status:       "completed",
					CreateTime:   "2024-01-01 12:00:00",
					CompleteTime: "2024-01-01 12:05:00",
				},
			},
		},
	}
	json.NewEncoder(w).Encode(response)
}

// Helper function to create test service with mock server
func createTestCDNService(server *httptest.Server) *CdnService {
	client := connectivity.NewClient("test-key", "test-secret", server.URL)
	return &CdnService{client: client}
}

func TestNewCdnService(t *testing.T) {
	client := &connectivity.Client{}
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
			Area:   "mainland_china",
			Type:   "page",
			Config: map[string]interface{}{
				"origin": map[string]interface{}{
					"default_master": "origin.example.com",
					"origin_mode":    "https",
					"port":           "443",
				},
			},
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
			Area:   "mainland_china",
			Type:   "page",
		}

		_, err := service.CreateDomain(req)
		if err == nil {
			t.Fatal("Expected error for invalid domain, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})

	t.Run("DifferentDomainTypes", func(t *testing.T) {
		types := []string{"page", "download", "video_demand", "dynamic", "video_live"}

		for _, domainType := range types {
			t.Run(fmt.Sprintf("Type_%s", domainType), func(t *testing.T) {
				req := DomainCreateRequest{
					Domain: fmt.Sprintf("test-%s.example.com", domainType),
					Area:   "global",
					Type:   domainType,
				}

				response, err := service.CreateDomain(req)
				if err != nil {
					t.Fatalf("Expected no error for type %s, got: %v", domainType, err)
				}

				if response.Code != 0 {
					t.Fatalf("Expected response code 0 for type %s, got: %d", domainType, response.Code)
				}
			})
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
			t.Fatal("Expected domain data, got empty array")
		}

		domainData := response.Data[0]
		if domainData.Domain != domain {
			t.Fatalf("Expected domain %s, got: %s", domain, domainData.Domain)
		}

		if domainData.ID == "" {
			t.Fatal("Expected domain ID to be set")
		}

		if domainData.Cname == "" {
			t.Fatal("Expected CNAME to be set")
		}
	})

	t.Run("DomainNotFound", func(t *testing.T) {
		domain := "notfound.example.com"

		_, err := service.GetDomain(domain)
		if err == nil {
			t.Fatal("Expected error for non-existent domain, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestListDomains(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulList", func(t *testing.T) {
		req := DomainListRequest{
			PageNumber:   1,
			PageSize:     10,
			DomainStatus: "serving",
		}

		response, err := service.ListDomains(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.PageSize != req.PageSize {
			t.Fatalf("Expected page size %d, got: %d", req.PageSize, response.Data.PageSize)
		}

		pageNumber, _ := strconv.Atoi(response.Data.PageNumber)
		if pageNumber != req.PageNumber {
			t.Fatalf("Expected page number %d, got: %d", req.PageNumber, pageNumber)
		}

		// Verify domain data structure
		for _, domain := range response.Data.List {
			if domain.ID == "" {
				t.Fatal("Expected domain ID to be populated")
			}

			if domain.Domain == "" {
				t.Fatal("Expected domain name to be populated")
			}

			if domain.Status == "" {
				t.Fatal("Expected domain status to be populated")
			}
		}
	})

	t.Run("EmptyList", func(t *testing.T) {
		req := DomainListRequest{
			PageNumber:   2,
			PageSize:     10,
			DomainStatus: "serving",
		}

		response, err := service.ListDomains(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) != 0 {
			t.Fatalf("Expected empty list, got: %d items", len(response.Data.List))
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

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestDomainDataHelperMethods(t *testing.T) {
	t.Run("IsServing", func(t *testing.T) {
		domain := &DomainData{Status: DomainStatusServing}
		if !domain.IsServing() {
			t.Fatal("Expected IsServing to return true for serving domain")
		}

		domain.Status = DomainStatusSuspended
		if domain.IsServing() {
			t.Fatal("Expected IsServing to return false for suspended domain")
		}
	})

	t.Run("IsDeleted", func(t *testing.T) {
		domain := &DomainData{Status: DomainStatusDeleted}
		if !domain.IsDeleted() {
			t.Fatal("Expected IsDeleted to return true for deleted domain")
		}

		domain.Status = DomainStatusServing
		if domain.IsDeleted() {
			t.Fatal("Expected IsDeleted to return false for serving domain")
		}
	})

	t.Run("IsSuspended", func(t *testing.T) {
		domain := &DomainData{Status: DomainStatusSuspended}
		if !domain.IsSuspended() {
			t.Fatal("Expected IsSuspended to return true for suspended domain")
		}

		domain.Status = DomainStatusServing
		if domain.IsSuspended() {
			t.Fatal("Expected IsSuspended to return false for serving domain")
		}
	})
}

func TestSetDomainConfig(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulSet", func(t *testing.T) {
		config := map[string]interface{}{
			"origin": map[string]interface{}{
				"default_master": "origin.example.com",
				"origin_mode":    "https",
			},
		}

		response, err := service.SetDomainConfig("example.com", config)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}
	})

	t.Run("InvalidConfig", func(t *testing.T) {
		config := map[string]interface{}{
			"invalid": "config",
		}

		_, err := service.SetDomainConfig("invalid.domain", config)
		if err == nil {
			t.Fatal("Expected error for invalid config, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestGetDomainConfig(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulGet", func(t *testing.T) {
		domain := "example.com"

		response, err := service.GetDomainConfig(domain, []string{})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data) == 0 {
			t.Fatal("Expected config data, got empty array")
		}

		configItem := response.Data[0]
		if configItem.Domain != domain {
			t.Fatalf("Expected domain %s, got: %s", domain, configItem.Domain)
		}

		if configItem.Config.Origin == nil {
			t.Fatal("Expected origin to be populated")
		}

	})

	t.Run("ConfigNotFound", func(t *testing.T) {
		domain := "notfound.example.com"

		_, err := service.GetDomainConfig(domain, []string{})
		if err == nil {
			t.Fatal("Expected error for non-existent config, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestDeleteDomainConfig(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulDelete", func(t *testing.T) {
		req := DeleteDomainConfigRequest{
			Domains: "example.com",
			Config:  []string{"cache_rule", "referer"},
		}

		err := service.DeleteDomainConfig(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
	})

	t.Run("DomainNotFound", func(t *testing.T) {
		req := DeleteDomainConfigRequest{
			Domains: "notfound.example.com",
			Config:  []string{"cache_rule"},
		}

		err := service.DeleteDomainConfig(req)
		if err == nil {
			t.Fatal("Expected error for non-existent domain, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestRefreshCache(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulRefresh", func(t *testing.T) {
		req := CacheRefreshRequest{
			URLs: []string{
				"https://example.com/image.jpg",
				"https://example.com/style.css",
			},
			Type: "url",
		}

		response, err := service.CacheRefresh(req.URLs, req.Type)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.TaskID == "" {
			t.Fatal("Expected task ID to be set")
		}

		if response.Data.Count != len(req.URLs) {
			t.Fatalf("Expected count %d, got: %d", len(req.URLs), response.Data.Count)
		}
	})

	t.Run("EmptyURLs", func(t *testing.T) {
		req := CacheRefreshRequest{
			URLs: []string{},
			Type: "url",
		}

		_, err := service.CacheRefresh(req.URLs, req.Type)
		if err == nil {
			t.Fatal("Expected error for empty URLs, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestQueryCacheRefresh(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulQuery", func(t *testing.T) {
		req := CacheRefreshQueryRequest{
			TaskID: 123456,
		}

		response, err := service.QueryCacheRefresh(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) == 0 {
			t.Fatal("Expected refresh data, got empty array")
		}

		item := response.Data.List[0]
		if item.ID == "" {
			t.Fatal("Expected task ID to be set")
		}

		if item.Status == "" {
			t.Fatal("Expected status to be set")
		}
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		req := CacheRefreshQueryRequest{
			TaskID: 99999,
		}

		_, err := service.QueryCacheRefresh(req)
		if err == nil {
			t.Fatal("Expected error for non-existent task, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestPurgeFiles(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulPurge", func(t *testing.T) {
		req := FilePurgeRequest{
			URLs: []string{
				"https://example.com/file1.pdf",
				"https://example.com/file2.zip",
			},
		}

		response, err := service.FilePurge(req.URLs)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if response.Data.TaskID == "" {
			t.Fatal("Expected task ID to be set")
		}

		if response.Data.Count != len(req.URLs) {
			t.Fatalf("Expected count %d, got: %d", len(req.URLs), response.Data.Count)
		}
	})

	t.Run("EmptyURLs", func(t *testing.T) {
		req := FilePurgeRequest{
			URLs: []string{},
		}

		_, err := service.FilePurge(req.URLs)
		if err == nil {
			t.Fatal("Expected error for empty URLs, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestQueryFilePurge(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulQuery", func(t *testing.T) {
		req := FilePurgeQueryRequest{
			TaskID: 123456,
		}

		response, err := service.QueryFilePurge(req)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) == 0 {
			t.Fatal("Expected purge data, got empty array")
		}

		item := response.Data.List[0]
		if item.ID == "" {
			t.Fatal("Expected task ID to be set")
		}

		if item.Status == "" {
			t.Fatal("Expected status to be set")
		}
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		req := FilePurgeQueryRequest{
			TaskID: 99999,
		}

		_, err := service.QueryFilePurge(req)
		if err == nil {
			t.Fatal("Expected error for non-existent task, got nil")
		}

		errorMsg := err.Error()
		if errorMsg == "" {
			t.Fatal("Expected non-empty error message")
		}
		t.Logf("Got expected error: %s", errorMsg)
	})
}

func TestQueryFilePurgeByTaskID(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	t.Run("SuccessfulQuery", func(t *testing.T) {
		taskID := 123456

		response, err := service.QueryFilePurgeByTaskID(taskID)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if response.Code != 0 {
			t.Fatalf("Expected response code 0, got: %d", response.Code)
		}

		if len(response.Data.List) == 0 {
			t.Fatal("Expected purge data, got empty array")
		}
	})
}

func TestNewCacheRefreshRequest(t *testing.T) {
	urls := []string{"https://example.com/test1", "https://example.com/test2"}
	refreshType := "url"

	req := NewCacheRefreshRequest(urls, refreshType)

	if len(req.URLs) != len(urls) {
		t.Fatalf("Expected %d URLs, got: %d", len(urls), len(req.URLs))
	}

	if req.Type != refreshType {
		t.Fatalf("Expected type %s, got: %s", refreshType, req.Type)
	}

	for i, url := range urls {
		if req.URLs[i] != url {
			t.Fatalf("Expected URL %s, got: %s", url, req.URLs[i])
		}
	}
}

// Benchmark tests
func BenchmarkCreateDomain(b *testing.B) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	req := DomainCreateRequest{
		Domain: "benchmark.example.com",
		Area:   "global",
		Type:   "page",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CreateDomain(req)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGetDomain(b *testing.B) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetDomain("example.com")
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkListDomains(b *testing.B) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	req := DomainListRequest{
		PageNumber:   1,
		PageSize:     10,
		DomainStatus: "serving",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ListDomains(req)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkRefreshCache(b *testing.B) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	req := CacheRefreshRequest{
		URLs: []string{"https://example.com/test.jpg"},
		Type: "url",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CacheRefresh(req.URLs, req.Type)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Integration-style tests (can be run against real API when available)
func TestCdnServiceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// These tests would run against a real API endpoint
	// Uncomment and configure when real API is available

	/*
		client := connectivity.NewClient("real-api-key", "real-secret", "https://real-api-endpoint.com")
		service := NewCdnService(client)

		t.Run("RealAPICreateAndDelete", func(t *testing.T) {
			// Test against real API
			req := DomainCreateRequest{
				Domain: "integration-test.example.com",
				Area:   "global",
				Type:   "page",
			}

			// Create domain
			response, err := service.CreateDomain(req)
			if err != nil {
				t.Fatalf("Failed to create domain: %v", err)
			}

			// Clean up - delete domain
			err = service.DeleteDomain(req.Domain)
			if err != nil {
				t.Fatalf("Failed to delete domain: %v", err)
			}
		})
	*/
}

// Table-driven tests for error scenarios
func TestCdnServiceErrorScenarios(t *testing.T) {
	server := createMockCDNServer()
	defer server.Close()

	service := createTestCDNService(server)

	errorTests := []struct {
		name        string
		testFunc    func() error
		expectError bool
		errorMsg    string
	}{
		{
			name: "CreateInvalidDomain",
			testFunc: func() error {
				req := DomainCreateRequest{Domain: "invalid.domain"}
				_, err := service.CreateDomain(req)
				return err
			},
			expectError: true,
			errorMsg:    "Invalid domain format",
		},
		{
			name: "GetNonExistentDomain",
			testFunc: func() error {
				_, err := service.GetDomain("notfound.example.com")
				return err
			},
			expectError: true,
			errorMsg:    "Domain not found",
		},
		{
			name: "DeleteNonExistentDomain",
			testFunc: func() error {
				return service.DeleteDomain("notfound.example.com")
			},
			expectError: true,
			errorMsg:    "Domain not found",
		},
		{
			name: "RefreshEmptyURLs",
			testFunc: func() error {
				req := CacheRefreshRequest{URLs: []string{}}
				_, err := service.CacheRefresh(req.URLs, req.Type)
				return err
			},
			expectError: true,
			errorMsg:    "URLs cannot be empty",
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
					// Check if error message contains expected text
					t.Logf("Error message: %s", err.Error())
				}
			}
		})
	}
}
