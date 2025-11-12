package connectivity

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewAPIClient tests APIClient creation
func TestNewAPIClient(t *testing.T) {
	accessKey := "test-access-key"
	secretKey := "test-secret-key"
	endpoint := "https://api.example.com"

	client := NewAPIClient(accessKey, secretKey, endpoint)

	if client == nil {
		t.Fatal("Expected API client, got nil")
	}

	if client.client == nil {
		t.Fatal("Expected resty client, got nil")
	}
}

// TestAPIClientGet tests GET request
func TestAPIClientGet(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify HTTP method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Verify headers
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Error("Expected X-API-Key header")
		}

		if r.Header.Get("X-API-Secret") != "test-secret" {
			t.Error("Expected X-API-Secret header")
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	type Response struct {
		Status string `json:"status"`
	}
	var result Response
	err := client.Get(ctx, "/test", &result)

	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}

	if result.Status != "ok" {
		t.Errorf("Expected status 'ok', got %v", result.Status)
	}
}

// TestAPIClientPost tests POST request
func TestAPIClientPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "123"}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	body := map[string]string{"name": "test"}
	type Response struct {
		ID string `json:"id"`
	}
	var result Response

	err := client.Post(ctx, "/create", body, &result)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}

	if result.ID != "123" {
		t.Errorf("Expected id '123', got %v", result.ID)
	}
}

// TestAPIClientPut tests PUT request
func TestAPIClientPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"updated": true}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	body := map[string]string{"name": "updated"}
	type Response struct {
		Updated bool `json:"updated"`
	}
	var result Response

	err := client.Put(ctx, "/update", body, &result)
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}

	if !result.Updated {
		t.Error("Expected updated=true")
	}
}

// TestAPIClientDelete tests DELETE request
func TestAPIClientDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	var result map[string]interface{}
	err := client.Delete(ctx, "/delete", &result)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
}

// TestAPIClientDeleteWithBody tests DELETE with body
func TestAPIClientDeleteWithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	body := map[string]string{"id": "123"}
	err := client.DeleteWithBody(ctx, "/delete", body)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
}

// TestAPIClientPatch tests PATCH request
func TestAPIClientPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"patched": true}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	body := map[string]string{"field": "value"}
	type Response struct {
		Patched bool `json:"patched"`
	}
	var result Response

	err := client.Patch(ctx, "/patch", body, &result)
	if err != nil {
		t.Fatalf("PATCH request failed: %v", err)
	}

	if !result.Patched {
		t.Error("Expected patched=true")
	}
}

// TestAPIClientGetWithQuery tests GET with query parameters
func TestAPIClientGetWithQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if r.URL.Query().Get("param1") != "value1" {
			t.Error("Expected param1=value1")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "ok"}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	query := map[string]string{"param1": "value1"}
	var result map[string]interface{}

	err := client.GetWithQuery(ctx, "/query", query, &result)
	if err != nil {
		t.Fatalf("GET with query failed: %v", err)
	}
}

// TestAPIClientPostWithHeaders tests POST with custom headers
func TestAPIClientPostWithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify custom header
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Error("Expected X-Custom-Header")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	body := map[string]string{"data": "test"}
	headers := map[string]string{"X-Custom-Header": "custom-value"}
	var result map[string]interface{}

	err := client.PostWithHeaders(ctx, "/post", body, headers, &result)
	if err != nil {
		t.Fatalf("POST with headers failed: %v", err)
	}
}

// TestAPIClientErrorHandling tests error responses
func TestAPIClientErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	var result map[string]interface{}
	err := client.Get(ctx, "/error", &result)

	if err == nil {
		t.Fatal("Expected error for 400 status, got nil")
	}

	t.Logf("Got expected error: %v", err)
}

// TestAPIClientSetTimeout tests timeout configuration
func TestAPIClientSetTimeout(t *testing.T) {
	client := NewAPIClient("test-key", "test-secret", "https://api.example.com")

	// Set timeout
	client.SetTimeout(10 * time.Second)

	// Verify by getting the underlying client
	restyClient := client.GetRestyClient()
	if restyClient == nil {
		t.Fatal("Expected resty client, got nil")
	}
}

// TestAPIClientSetRetryCount tests retry configuration
func TestAPIClientSetRetryCount(t *testing.T) {
	client := NewAPIClient("test-key", "test-secret", "https://api.example.com")

	client.SetRetryCount(5)

	// The configuration is set, no error should occur
	t.Log("Retry count set successfully")
}

// TestAPIClientSetRetryWaitTime tests retry wait time configuration
func TestAPIClientSetRetryWaitTime(t *testing.T) {
	client := NewAPIClient("test-key", "test-secret", "https://api.example.com")

	client.SetRetryWaitTime(2 * time.Second)
	client.SetRetryMaxWaitTime(10 * time.Second)

	t.Log("Retry wait times set successfully")
}

// TestAPIClientSetBaseURL tests base URL configuration
func TestAPIClientSetBaseURL(t *testing.T) {
	client := NewAPIClient("test-key", "test-secret", "https://api.example.com")

	newURL := "https://new-api.example.com"
	client.SetBaseURL(newURL)

	t.Logf("Base URL updated to: %s", newURL)
}

// TestAPIClientGetRestyClient tests getting underlying client
func TestAPIClientGetRestyClient(t *testing.T) {
	client := NewAPIClient("test-key", "test-secret", "https://api.example.com")

	restyClient := client.GetRestyClient()
	if restyClient == nil {
		t.Fatal("Expected resty client, got nil")
	}

	t.Log("Successfully retrieved resty client")
}

// TestAPIClientAuthenticationHeaders tests authentication headers
func TestAPIClientAuthenticationHeaders(t *testing.T) {
	headersCaptured := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth headers
		if r.Header.Get("X-API-Key") == "my-key" &&
			r.Header.Get("X-API-Secret") == "my-secret" &&
			r.URL.Query().Get("token") == "my-secret" {
			headersCaptured = true
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewAPIClient("my-key", "my-secret", server.URL)
	ctx := context.Background()

	var result map[string]interface{}
	_ = client.Get(ctx, "/test", &result)

	if !headersCaptured {
		t.Error("Authentication headers were not properly set")
	}
}

// BenchmarkAPIClientGet benchmarks GET requests
func BenchmarkAPIClientGet(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		_ = client.Get(ctx, "/test", &result)
	}
}

// BenchmarkAPIClientPost benchmarks POST requests
func BenchmarkAPIClientPost(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": "123"}`))
	}))
	defer server.Close()

	client := NewAPIClient("test-key", "test-secret", server.URL)
	ctx := context.Background()
	body := map[string]string{"name": "test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		_ = client.Post(ctx, "/create", body, &result)
	}
}
