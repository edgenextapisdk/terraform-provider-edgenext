package scdn

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/stretchr/testify/assert"
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

	// Check if config file exists in current directory
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Try to find config file in test directory
		testDir := filepath.Dir(os.Args[0])
		configFile = filepath.Join(testDir, "test_config.json")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			// Try to find config file in parent directory
			configFile = filepath.Join("..", "test_config.json")
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				// No config file found, return default config
				return config
			}
		}
	}

	// Read and parse config file
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		// If we can't read the file, return default config
		return config
	}

	// Parse JSON config
	if err := json.Unmarshal(data, config); err != nil {
		// If parsing fails, return default config
		return config
	}

	return config
}

// createTestClient creates a real EdgeNextClient for integration tests
func createTestClient(t *testing.T) *connectivity.EdgeNextClient {
	config := getTestConfig()

	// Skip integration tests if not enabled or credentials are not provided
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

func TestNewScdnService(t *testing.T) {
	client := &connectivity.EdgeNextClient{}
	service := NewScdnService(client)

	assert.NotNil(t, service)
	assert.Equal(t, client, service.client)
}
