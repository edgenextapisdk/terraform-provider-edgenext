package connectivity

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestNewECSClient(t *testing.T) {
	client := NewECSClient("test-access-key", "test-secret-key", "https://api.example.com")
	if client == nil {
		t.Fatal("Expected ECS client, got nil")
	}

	if client.client == nil {
		t.Fatal("Expected resty client, got nil")
	}
}

func TestECSClientGet(t *testing.T) {
	const (
		accessKey = "ENAK160fb5102e1cf78477a6"
		secretKey = "41544f054cccb05a64edb85405e3c7bd"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		if got := r.Header.Get("Authorization"); got != fmt.Sprintf("Bearer %s", accessKey) {
			t.Errorf("Unexpected Authorization header: %s", got)
		}

		timestamp := r.Header.Get("Edgenext-Timestamp")
		if timestamp == "" {
			t.Fatal("Expected Edgenext-Timestamp header")
		}

		expectedPayload := r.URL.Query().Encode()
		expectedSig := buildExpectedSignature(secretKey, timestamp, []byte(expectedPayload))
		if got := r.Header.Get("Signature"); got != expectedSig {
			t.Errorf("Unexpected Signature header, expected %s, got %s", expectedSig, got)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := NewECSClient(accessKey, secretKey, server.URL)
	query := map[string]interface{}{
		"start_time": "2023-10-23 12:00",
		"end_time":   "2023-10-23 12:05",
		"token":      "cdfee2bc43d63caeaa3b169ad31123123",
		"domains":    "cdn.api.baishan.com",
	}

	var result map[string]interface{}
	err := client.Get(context.Background(), "/cdn/v2/stat/request", query, &result)
	if err != nil {
		t.Fatalf("ECS GET request failed: %v", err)
	}
}

func TestECSClientPost(t *testing.T) {
	const (
		accessKey = "ENAK160fb5102e1cf78477a6"
		secretKey = "41544f054cccb05a64edb85405e3c7bd"
	)

	requestBody := map[string]interface{}{
		"startTime": "2023-10-10T02:00:00Z",
		"endTime":   "2023-10-10T03:05:00Z",
		"sn":        "725000408",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if got := r.Header.Get("Authorization"); got != fmt.Sprintf("Bearer %s", accessKey) {
			t.Errorf("Unexpected Authorization header: %s", got)
		}

		timestamp := r.Header.Get("Edgenext-Timestamp")
		if timestamp == "" {
			t.Fatal("Expected Edgenext-Timestamp header")
		}

		var gotBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		payload, err := json.Marshal(gotBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		expectedSig := buildExpectedSignature(secretKey, timestamp, payload)
		if got := r.Header.Get("Signature"); got != expectedSig {
			t.Errorf("Unexpected Signature header, expected %s, got %s", expectedSig, got)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	client := NewECSClient(accessKey, secretKey, server.URL)
	var result map[string]interface{}
	err := client.Post(context.Background(), "/bms/openapi/v1/getFlow", requestBody, &result)
	if err != nil {
		t.Fatalf("ECS POST request failed: %v", err)
	}
}

// TestECSClientPostGetEcsAggregateBandwidth matches:
//
//	curl -X POST 'http://<endpoint>/ecs/openapi/v1/ecs/getEcsAggregateBandwidth' \
//	  -H 'Content-Type: application/json' \
//	  -d '{"start_time":"2026-04-01 00:00:00","end_time":"2026-04-01 23:59:59"}'
//
// (Signature / Authorization / Edgenext-Timestamp are added by ECSClient.)
func TestECSClientPostGetEcsAggregateBandwidth(t *testing.T) {
	// Use placeholder credentials; real calls use TestECSClientPostGetEcsAggregateBandwidth_Integration + env vars.
	const (
		accessKey = "ENAKtest_access_key_for_mock"
		secretKey = "test_secret_key_for_mock_signature"
	)

	requestBody := map[string]interface{}{
		"start_time": "2026-04-01 00:00:00",
		"end_time":   "2026-04-01 23:59:59",
	}
	path := "/ecs/openapi/v1/ecs/getEcsAggregateBandwidth"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.URL.Path != path {
			t.Errorf("Expected path %q, got %q", path, r.URL.Path)
		}

		if got := r.Header.Get("Authorization"); got != fmt.Sprintf("Bearer %s", accessKey) {
			t.Errorf("Unexpected Authorization header: %s", got)
		}
		if r.Header.Get("Content-Type") == "" {
			t.Error("Expected Content-Type header")
		}

		timestamp := r.Header.Get("Edgenext-Timestamp")
		if timestamp == "" {
			t.Fatal("Expected Edgenext-Timestamp header")
		}

		var gotBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if gotBody["start_time"] != requestBody["start_time"] || gotBody["end_time"] != requestBody["end_time"] {
			t.Fatalf("Unexpected body: %#v", gotBody)
		}

		payload, err := json.Marshal(gotBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		expectedSig := buildExpectedSignature(secretKey, timestamp, payload)
		if got := r.Header.Get("Signature"); got != expectedSig {
			t.Errorf("Unexpected Signature header, expected %s, got %s", expectedSig, got)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{}}`))
	}))
	defer server.Close()

	client := NewECSClient(accessKey, secretKey, server.URL)
	var result map[string]interface{}
	err := client.Post(context.Background(), path, requestBody, &result)
	if err != nil {
		t.Fatalf("ECS POST getEcsAggregateBandwidth failed: %v", err)
	}
}

// TestECSClientPostGetEcsAggregateBandwidth_Integration calls the real OpenAPI (optional).
// Run with:
//
//	EDGENEXT_ECS_ACCESS_KEY=... EDGENEXT_ECS_SECRET_KEY=... EDGENEXT_ECS_ENDPOINT=edgenext-openapi.console.prxcdn.com \
//	  go test ./edgenext/connectivity -run TestECSClientPostGetEcsAggregateBandwidth_Integration -v
//
// Endpoint may be host only or full URL; http is assumed if no scheme is set.
func TestECSClientPostGetEcsAggregateBandwidth_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_ECS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_ECS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_ECS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_ECS_ACCESS_KEY, EDGENEXT_ECS_SECRET_KEY, EDGENEXT_ECS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewECSClient(accessKey, secretKey, endpoint)
	body := map[string]interface{}{
		"start_time": "2026-04-01 00:00:00",
		"end_time":   "2026-04-01 23:59:59",
	}
	path := "/ecs/openapi/v1/ecs/getEcsAggregateBandwidth"

	var result map[string]interface{}
	err := client.Post(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("ECS POST getEcsAggregateBandwidth integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestECSClientPostInstanceList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_ECS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_ECS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_ECS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_ECS_ACCESS_KEY, EDGENEXT_ECS_SECRET_KEY, EDGENEXT_ECS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewECSClient(accessKey, secretKey, endpoint)
	body := map[string]interface{}{
		"limit":  10,
		"name":   "dsd",
		"region": "lianyungang-a",
	}
	path := "/ecs/openapi/v2/instance/list"

	var result map[string]interface{}
	err := client.Post(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("ECS POST instance list integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("result: %#v", result)
}

func TestECSClientPostKeyPairCreate_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_ECS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_ECS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_ECS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_ECS_ACCESS_KEY, EDGENEXT_ECS_SECRET_KEY, EDGENEXT_ECS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewECSClient(accessKey, secretKey, endpoint)
	body := map[string]interface{}{
		"name":   "test-key-pair",
		"region": "lianyungang-a",
	}
	path := "/ecs/openapi/v2/keypair/create"

	var result map[string]interface{}
	err := client.Post(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("ECS POST key pair create integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("result: %#v", result)
}

func TestECSClientPostKeyPairList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_ECS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_ECS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_ECS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_ECS_ACCESS_KEY, EDGENEXT_ECS_SECRET_KEY, EDGENEXT_ECS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewECSClient(accessKey, secretKey, endpoint)
	body := map[string]interface{}{
		"region": "lianyungang-a",
	}
	path := "/ecs/openapi/v2/keypair/list"

	var result map[string]interface{}
	err := client.Post(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("ECS POST key pair list integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("result: %#v", result)
}

func TestECSClientPostImageList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_ECS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_ECS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_ECS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_ECS_ACCESS_KEY, EDGENEXT_ECS_SECRET_KEY, EDGENEXT_ECS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewECSClient(accessKey, secretKey, endpoint)
	body := map[string]interface{}{
		"region":     "lianyungang-a",
		"visibility": "public",
		"name":       "Debian 11.11 64-bit-v1.12",
		"page_num":   1,
		"page_size":  10,
	}
	path := "/ecs/openapi/v2/image/list"

	var result map[string]interface{}
	err := client.Post(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("ECS POST image list integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("result: %#v", result)
}

func TestECSClientPostTagList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_ECS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_ECS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_ECS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_ECS_ACCESS_KEY, EDGENEXT_ECS_SECRET_KEY, EDGENEXT_ECS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewECSClient(accessKey, secretKey, endpoint)
	body := map[string]interface{}{
		"pageNum":  1,
		"pageSize": 10,
		"tagKey":   "test-key",
		"tagValue": "test-value",
	}
	path := "/ecs/openapi/v2/tags/list"

	var result map[string]interface{}
	err := client.Get(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("ECS GET tag list integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("result: %#v", result)
}

func TestECSClientPostTagCreate_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_ECS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_ECS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_ECS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_ECS_ACCESS_KEY, EDGENEXT_ECS_SECRET_KEY, EDGENEXT_ECS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewECSClient(accessKey, secretKey, endpoint)
	body := map[string]interface{}{
		"tags": []map[string]interface{}{{
			"key":   "test-key",
			"value": "test-value",
		}},
	}
	path := "/ecs/openapi/v2/tags/create"

	var result map[string]interface{}
	err := client.Post(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("ECS POST tag create integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("result: %#v", result)
}

func TestSetQueryParamsFromValues(t *testing.T) {
	values := setQueryParamsFromValues(map[string]interface{}{
		"string_field": "abc",
		"number_field": 100,
		"bool_field":   true,
	})

	if values.Get("string_field") != "abc" {
		t.Fatalf("Expected string_field=abc, got %s", values.Get("string_field"))
	}
	if values.Get("number_field") != "100" {
		t.Fatalf("Expected number_field=100, got %s", values.Get("number_field"))
	}
	if values.Get("bool_field") != "true" {
		t.Fatalf("Expected bool_field=true, got %s", values.Get("bool_field"))
	}
}

func TestECSClientSign(t *testing.T) {
	client := NewECSClient("ak", "sk", "https://api.example.com")
	timestamp, signature := client.sign([]byte("a=1&b=2"))

	if timestamp == "" {
		t.Fatal("Expected non-empty timestamp")
	}
	if signature == "" {
		t.Fatal("Expected non-empty signature")
	}

	if _, err := strconv.ParseInt(timestamp, 10, 64); err != nil {
		t.Fatalf("Expected unix timestamp format, got %q", timestamp)
	}
}

func buildExpectedSignature(secretKey, timestamp string, payload []byte) string {
	key := fmt.Sprintf("%s-%s", secretKey, timestamp)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(payload)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
