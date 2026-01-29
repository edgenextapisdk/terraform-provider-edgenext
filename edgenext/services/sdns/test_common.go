package sdns

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// Test configuration for integration tests
type TestConfig struct {
	AccessKey              string `json:"access_key"`
	SecretKey              string `json:"secret_key"`
	Endpoint               string `json:"endpoint"`
	TimeoutSeconds         int    `json:"timeout_seconds"`
	EnableIntegrationTests bool   `json:"enable_integration_tests"`
}

// getTestConfig returns test configuration from config file
func getTestConfig() *TestConfig {
	// Default configuration
	config := &TestConfig{
		AccessKey:              "",
		SecretKey:              "",
		Endpoint:               "https://api.edgenextscdn.com",
		TimeoutSeconds:         30,
		EnableIntegrationTests: false,
	}

	// Try to read from config file
	configFile := "test_config.json"

	// Check if config file exists in current directory or parents
	curr, _ := os.Getwd()
	for {
		path := filepath.Join(curr, configFile)
		if _, err := os.Stat(path); err == nil {
			data, err := ioutil.ReadFile(path)
			if err == nil {
				json.Unmarshal(data, config)
				return config
			}
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	return config
}

// createTestClient creates a real EdgeNextClient for integration tests
func createTestClient(t *testing.T) *connectivity.EdgeNextClient {
	config := getTestConfig()

	if !config.EnableIntegrationTests {
		t.Skip("Skipping integration test: integration tests are disabled in config file")
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		t.Skip("Skipping integration test: access_key and secret_key are required in config file")
	}

	connectivityConfig := &connectivity.Config{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
		Endpoint:  config.Endpoint,
	}

	client, err := connectivityConfig.Client()
	if err != nil {
		t.Fatalf("Failed to create EdgeNextClient: %v", err)
	}

	return client
}

// isIntegrationTest checks if we should run integration tests
func isIntegrationTest() bool {
	config := getTestConfig()
	return config.EnableIntegrationTests && config.AccessKey != "" && config.SecretKey != ""
}
