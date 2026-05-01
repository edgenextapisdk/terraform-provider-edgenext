package connectivity

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestNewRDSClient(t *testing.T) {
	client := NewRDSClient("test-access-key", "test-secret-key", "https://api.example.com", "lianyungang-a")
	if client == nil {
		t.Fatal("Expected RDS client, got nil")
	}
	if client.client == nil {
		t.Fatal("Expected resty client, got nil")
	}
}

func TestRDSClientPostRDSInstanceList_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	accessKey := os.Getenv("EDGENEXT_RDS_ACCESS_KEY")
	secretKey := os.Getenv("EDGENEXT_RDS_SECRET_KEY")
	endpoint := strings.TrimSpace(os.Getenv("EDGENEXT_RDS_ENDPOINT"))
	if accessKey == "" || secretKey == "" || endpoint == "" {
		t.Skip("set EDGENEXT_RDS_ACCESS_KEY, EDGENEXT_RDS_SECRET_KEY, EDGENEXT_RDS_ENDPOINT to run this test")
	}
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := NewRDSClient(accessKey, secretKey, endpoint, "lianyungang-a")
	body := map[string]interface{}{
		"page_number":    1,
		"page_size":      10,
		"datastore_type": "mysql",
	}
	path := "/rds/openapi/v2/instances/list"

	var result map[string]interface{}
	err := client.Post(context.Background(), path, body, &result)
	if err != nil {
		t.Fatalf("RDS POST rds instance list integration failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("result: %#v", result)
}
