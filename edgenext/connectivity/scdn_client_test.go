package connectivity

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewScdnClient(t *testing.T) {
	client := NewScdnClient("https://api.example.com", "test-key", "test-secret", 30*time.Second)

	assert.NotNil(t, client)
	assert.Equal(t, "https://api.example.com", client.baseURL)
	assert.Equal(t, "test-key", client.apiKey)
	assert.Equal(t, "test-secret", client.apiSecret)
	assert.Equal(t, 30*time.Second, client.timeout)
	assert.NotNil(t, client.sdk)
}

func TestNewScdnClientFromConfig(t *testing.T) {
	// Create a mock EdgeNextClient
	edgeClient := &EdgeNextClient{}

	client := NewScdnClientFromConfig(edgeClient)

	assert.NotNil(t, client)
	assert.Equal(t, "https://api.edgenextscdn.com", client.baseURL)
	assert.Equal(t, 30*time.Second, client.timeout)
}

func TestScdnClient_Get(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v5/domains", r.URL.Path)
		assert.Equal(t, "test-key", r.Header.Get("X-Auth-App-Id"))
		assert.NotEmpty(t, r.Header.Get("X-Auth-Sign"))
		assert.Equal(t, "en", r.Header.Get("X-Lang"))
		assert.Contains(t, r.Header.Get("Content-Type"), "application/json")

		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    1,
				Message: "Success",
			},
			Data: map[string]interface{}{
				"total": 1,
				"list": []map[string]interface{}{
					{
						"id":     1,
						"domain": "test.example.com",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	req := &ScdnRequest{
		Query: map[string]interface{}{
			"page":      1,
			"page_size": 10,
		},
	}

	resp, err := client.Get(context.Background(), "/api/v5/domains", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Status.Code)
	assert.Equal(t, "Success", resp.Status.Message)
	assert.NotNil(t, resp.Data)
}

func TestScdnClient_Post(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v5/domains", r.URL.Path)
		assert.Equal(t, "test-key", r.Header.Get("X-Auth-App-Id"))
		assert.NotEmpty(t, r.Header.Get("X-Auth-Sign"))
		assert.Equal(t, "en", r.Header.Get("X-Lang"))
		assert.Contains(t, r.Header.Get("Content-Type"), "application/json")

		// Read request body
		var requestData map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		assert.NoError(t, err)
		assert.Equal(t, "test.example.com", requestData["domain"])

		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    1,
				Message: "Domain created successfully",
			},
			Data: map[string]interface{}{
				"id":     123,
				"domain": "test.example.com",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	req := &ScdnRequest{
		Data: map[string]interface{}{
			"domain":   "test.example.com",
			"group_id": 123,
		},
	}

	resp, err := client.Post(context.Background(), "/api/v5/domains", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Status.Code)
	assert.Equal(t, "Domain created successfully", resp.Status.Message)
	assert.NotNil(t, resp.Data)
}

func TestScdnClient_Put(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/api/v5/domains", r.URL.Path)

		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    1,
				Message: "Domain updated successfully",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	req := &ScdnRequest{
		Data: map[string]interface{}{
			"domain_id": 123,
			"remark":    "Updated remark",
		},
	}

	resp, err := client.Put(context.Background(), "/api/v5/domains", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Status.Code)
	assert.Equal(t, "Domain updated successfully", resp.Status.Message)
}

func TestScdnClient_Delete(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/api/v5/domains", r.URL.Path)

		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    1,
				Message: "Domain deleted successfully",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	req := &ScdnRequest{
		Data: map[string]interface{}{
			"ids": []int{123, 456},
		},
	}

	resp, err := client.Delete(context.Background(), "/api/v5/domains", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Status.Code)
	assert.Equal(t, "Domain deleted successfully", resp.Status.Message)
}

func TestScdnClient_CustomHeaders(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "zh", r.Header.Get("X-Lang"))
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))

		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    1,
				Message: "Success",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	req := &ScdnRequest{
		Header: map[string]string{
			"X-Lang":          "zh",
			"X-Custom-Header": "custom-value",
		},
	}

	resp, err := client.Get(context.Background(), "/api/v5/domains", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Status.Code)
}

func TestScdnClient_QueryParameters(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "10", r.URL.Query().Get("page_size"))
		assert.Equal(t, "test", r.URL.Query().Get("domain"))

		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    1,
				Message: "Success",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	req := &ScdnRequest{
		Query: map[string]interface{}{
			"page":      1,
			"page_size": 10,
			"domain":    "test",
		},
	}

	resp, err := client.Get(context.Background(), "/api/v5/domains", req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Status.Code)
}

func TestScdnClient_APIError(t *testing.T) {
	// Create mock server that returns API error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    0,
				Message: "API Error: Invalid parameters",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	req := &ScdnRequest{}
	resp, err := client.Get(context.Background(), "/api/v5/domains", req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "API error: API Error: Invalid parameters")
}

func TestScdnClient_SetTimeout(t *testing.T) {
	client := NewScdnClient("https://api.example.com", "test-key", "test-secret", 30*time.Second)

	client.SetTimeout(60 * time.Second)
	assert.Equal(t, 60*time.Second, client.GetTimeout())
}

func TestScdnClient_SetBaseURL(t *testing.T) {
	client := NewScdnClient("https://api.example.com", "test-key", "test-secret", 30*time.Second)

	client.SetBaseURL("https://new-api.example.com")
	assert.Equal(t, "https://new-api.example.com", client.GetBaseURL())
}

func TestScdnClient_SetCredentials(t *testing.T) {
	client := NewScdnClient("https://api.example.com", "test-key", "test-secret", 30*time.Second)

	client.SetCredentials("new-key", "new-secret")
	apiKey, apiSecret := client.GetCredentials()
	assert.Equal(t, "new-key", apiKey)
	assert.Equal(t, "new-secret", apiSecret)
}

func TestScdnClient_IsHealthy(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ScdnResponse{
			Status: ScdnStatus{
				Code:    1,
				Message: "Success",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	err := client.IsHealthy(context.Background())
	assert.NoError(t, err)
}

func TestScdnClient_IsHealthy_Error(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewScdnClient(server.URL, "test-key", "test-secret", 30*time.Second)

	err := client.IsHealthy(context.Background())
	assert.Error(t, err)
}

func TestScdnClient_GetAPIVersion(t *testing.T) {
	client := NewScdnClient("https://api.example.com", "test-key", "test-secret", 30*time.Second)

	assert.Equal(t, "v5", client.GetAPIVersion())
}

func TestScdnClient_GetServiceName(t *testing.T) {
	client := NewScdnClient("https://api.example.com", "test-key", "test-secret", 30*time.Second)

	assert.Equal(t, "scdn", client.GetServiceName())
}

func TestScdnRequest(t *testing.T) {
	req := &ScdnRequest{
		Data: map[string]interface{}{
			"domain": "test.example.com",
		},
		Query: map[string]interface{}{
			"page": 1,
		},
		Header: map[string]string{
			"X-Lang": "en",
		},
	}

	assert.NotNil(t, req.Data)
	assert.NotNil(t, req.Query)
	assert.NotNil(t, req.Header)
	assert.Equal(t, "test.example.com", req.Data["domain"])
	assert.Equal(t, 1, req.Query["page"])
	assert.Equal(t, "en", req.Header["X-Lang"])
}

func TestScdnResponse(t *testing.T) {
	resp := &ScdnResponse{
		Status: ScdnStatus{
			Code:    1,
			Message: "Success",
		},
		Data: map[string]interface{}{
			"id": 123,
		},
	}

	assert.Equal(t, 1, resp.Status.Code)
	assert.Equal(t, "Success", resp.Status.Message)
	assert.NotNil(t, resp.Data)
}

func TestScdnStatus(t *testing.T) {
	status := ScdnStatus{
		Code:    1,
		Message: "Success",
	}

	assert.Equal(t, 1, status.Code)
	assert.Equal(t, "Success", status.Message)
}
