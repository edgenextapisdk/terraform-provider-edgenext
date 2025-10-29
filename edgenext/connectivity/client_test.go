package connectivity

import (
	"sync"
	"testing"
)

// TestConfigClient tests the Config.Client() method
func TestConfigClient(t *testing.T) {
	config := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Endpoint:  "https://api.example.com",
		Region:    "us-east-1",
	}

	client, err := config.Client()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client, got nil")
	}

	if client.config != config {
		t.Error("Client config should match input config")
	}
}

// TestAPIClientInitialization tests APIClient initialization
func TestAPIClientInitialization(t *testing.T) {
	config := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Endpoint:  "https://api.example.com",
		Region:    "us-east-1",
	}

	client, _ := config.Client()

	// First call should initialize
	apiClient1, err := client.APIClient()
	if err != nil {
		t.Fatalf("Failed to get API client: %v", err)
	}

	if apiClient1 == nil {
		t.Fatal("Expected API client, got nil")
	}

	// Second call should return the same instance
	apiClient2, err := client.APIClient()
	if err != nil {
		t.Fatalf("Failed to get API client on second call: %v", err)
	}

	if apiClient1 != apiClient2 {
		t.Error("API client should be a singleton")
	}
}

// TestOSSClientInitialization tests OSSClient initialization
func TestOSSClientInitialization(t *testing.T) {
	config := &Config{
		AccessKey: "8cafd7f44c0d499cb8324954b6202e47",
		SecretKey: "a46a0fcc3b4941148995bb0a5e51a40b",
		Endpoint:  "https://oss-as-central-5.edgenextcs.com",
		Region:    "oss-as-central-5",
	}

	client, _ := config.Client()

	// First call should initialize
	ossClient1, err := client.OSSClient()
	if err != nil {
		t.Fatalf("Failed to get OSS client: %v", err)
	}

	if ossClient1 == nil {
		t.Fatal("Expected OSS client, got nil")
	}

	// Second call should return the same instance
	ossClient2, err := client.OSSClient()
	if err != nil {
		t.Fatalf("Failed to get OSS client on second call: %v", err)
	}

	if ossClient1 != ossClient2 {
		t.Error("OSS client should be a singleton")
	}

	// Verify region
	if ossClient1.GetRegion() != config.Region {
		t.Errorf("Expected region %s, got %s", config.Region, ossClient1.GetRegion())
	}
}

// TestOSSClientWithoutRegion tests OSSClient initialization without region
func TestOSSClientWithoutRegion(t *testing.T) {
	config := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Endpoint:  "https://oss.example.com",
		Region:    "", // No region
	}

	client, _ := config.Client()

	ossClient, err := client.OSSClient()
	if err == nil {
		t.Fatal("Expected error for missing region, got nil")
	}

	if ossClient != nil {
		t.Error("Expected nil client when error occurs")
	}

	expectedErr := "region is required for OSS services"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got '%s'", expectedErr, err.Error())
	}
}

// TestAPIClientConcurrentAccess tests concurrent access to APIClient
func TestAPIClientConcurrentAccess(t *testing.T) {
	config := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Endpoint:  "https://api.example.com",
		Region:    "us-east-1",
	}

	client, _ := config.Client()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	clients := make([]*APIClient, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			apiClient, err := client.APIClient()
			if err != nil {
				t.Errorf("Failed to get API client: %v", err)
				return
			}
			clients[index] = apiClient
		}(i)
	}

	wg.Wait()

	// Verify all clients are the same instance
	firstClient := clients[0]
	if firstClient == nil {
		t.Fatal("First client is nil")
	}

	for i, c := range clients {
		if c == nil {
			t.Errorf("Client at index %d is nil", i)
			continue
		}
		if c != firstClient {
			t.Errorf("Client at index %d is different from first client", i)
		}
	}
}

// TestOSSClientConcurrentAccess tests concurrent access to OSSClient
func TestOSSClientConcurrentAccess(t *testing.T) {
	config := &Config{
		AccessKey: "8cafd7f44c0d499cb8324954b6202e47",
		SecretKey: "a46a0fcc3b4941148995bb0a5e51a40b",
		Endpoint:  "https://oss-as-central-5.edgenextcs.com",
		Region:    "oss-as-central-5",
	}

	client, _ := config.Client()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	clients := make([]*OSSClient, goroutines)
	errors := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			ossClient, err := client.OSSClient()
			clients[index] = ossClient
			errors[index] = err
		}(i)
	}

	wg.Wait()

	// Verify all calls returned the same result
	firstClient := clients[0]
	firstErr := errors[0]

	if firstErr != nil {
		t.Fatalf("Failed to get OSS client: %v", firstErr)
	}

	if firstClient == nil {
		t.Fatal("First OSS client is nil")
	}

	for i := range clients {
		if errors[i] != nil {
			t.Errorf("Error at index %d: %v", i, errors[i])
			continue
		}

		if clients[i] == nil {
			t.Errorf("Client at index %d is nil", i)
			continue
		}

		if clients[i] != firstClient {
			t.Errorf("Client at index %d is different from first client", i)
		}
	}

	// Verify region
	if firstClient.GetRegion() != config.Region {
		t.Errorf("Expected region %s, got %s", config.Region, firstClient.GetRegion())
	}
}

// TestOSSClientConcurrentAccessWithoutRegion tests concurrent access when region is missing
func TestOSSClientConcurrentAccessWithoutRegion(t *testing.T) {
	config := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Endpoint:  "https://oss.example.com",
		Region:    "", // No region
	}

	client, _ := config.Client()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	errors := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			_, err := client.OSSClient()
			errors[index] = err
		}(i)
	}

	wg.Wait()

	// All errors should be the same
	firstErr := errors[0]
	if firstErr == nil {
		t.Fatal("Expected error for missing region, got nil")
	}

	expectedErrMsg := "region is required for OSS services"
	if firstErr.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, firstErr.Error())
	}

	for i, err := range errors {
		if err == nil {
			t.Errorf("Error at index %d is nil, expected error", i)
			continue
		}
		if err.Error() != firstErr.Error() {
			t.Errorf("Error at index %d differs: got %v, want %v", i, err, firstErr)
		}
	}
}

// TestBothClientsConcurrentAccess tests concurrent access to both clients
func TestBothClientsConcurrentAccess(t *testing.T) {
	config := &Config{
		AccessKey: "8cafd7f44c0d499cb8324954b6202e47",
		SecretKey: "a46a0fcc3b4941148995bb0a5e51a40b",
		Endpoint:  "https://oss-as-central-5.edgenextcs.com",
		Region:    "oss-as-central-5",
	}

	client, _ := config.Client()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	apiClients := make([]*APIClient, goroutines)
	ossClients := make([]*OSSClient, goroutines)

	// Concurrently call both APIClient() and OSSClient()
	for i := 0; i < goroutines; i++ {
		// API Client goroutines
		go func(index int) {
			defer wg.Done()
			apiClient, err := client.APIClient()
			if err != nil {
				t.Errorf("Failed to get API client: %v", err)
				return
			}
			apiClients[index] = apiClient
		}(i)

		// OSS Client goroutines
		go func(index int) {
			defer wg.Done()
			ossClient, err := client.OSSClient()
			if err != nil {
				t.Errorf("Failed to get OSS client: %v", err)
				return
			}
			ossClients[index] = ossClient
		}(i)
	}

	wg.Wait()

	// Verify API clients
	firstAPIClient := apiClients[0]
	if firstAPIClient == nil {
		t.Fatal("First API client is nil")
	}

	for i, c := range apiClients {
		if c == nil {
			t.Errorf("API client at index %d is nil", i)
			continue
		}
		if c != firstAPIClient {
			t.Errorf("API client at index %d is different", i)
		}
	}

	// Verify OSS clients
	firstOSSClient := ossClients[0]
	if firstOSSClient == nil {
		t.Fatal("First OSS client is nil")
	}

	for i, c := range ossClients {
		if c == nil {
			t.Errorf("OSS client at index %d is nil", i)
			continue
		}
		if c != firstOSSClient {
			t.Errorf("OSS client at index %d is different", i)
		}
	}
}

// BenchmarkAPIClientAccess benchmarks APIClient access
func BenchmarkAPIClientAccess(b *testing.B) {
	config := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Endpoint:  "https://api.example.com",
		Region:    "us-east-1",
	}

	client, _ := config.Client()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.APIClient()
	}
}

// BenchmarkOSSClientAccess benchmarks OSSClient access
func BenchmarkOSSClientAccess(b *testing.B) {
	config := &Config{
		AccessKey: "8cafd7f44c0d499cb8324954b6202e47",
		SecretKey: "a46a0fcc3b4941148995bb0a5e51a40b",
		Endpoint:  "https://oss-as-central-5.edgenextcs.com",
		Region:    "oss-as-central-5",
	}

	client, _ := config.Client()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.OSSClient()
	}
}

// BenchmarkAPIClientParallel benchmarks parallel APIClient access
func BenchmarkAPIClientParallel(b *testing.B) {
	config := &Config{
		AccessKey: "test-access-key",
		SecretKey: "test-secret-key",
		Endpoint:  "https://api.example.com",
		Region:    "us-east-1",
	}

	client, _ := config.Client()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.APIClient()
		}
	})
}

// BenchmarkOSSClientParallel benchmarks parallel OSSClient access
func BenchmarkOSSClientParallel(b *testing.B) {
	config := &Config{
		AccessKey: "8cafd7f44c0d499cb8324954b6202e47",
		SecretKey: "a46a0fcc3b4941148995bb0a5e51a40b",
		Endpoint:  "https://oss-as-central-5.edgenextcs.com",
		Region:    "oss-as-central-5",
	}

	client, _ := config.Client()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = client.OSSClient()
		}
	})
}
