package helper

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestWriteToFile(t *testing.T) {
	// Create a test schema.ResourceData
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"result_output_file": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}, map[string]interface{}{
		"result_output_file": "test_output.json",
	})

	// Test data
	testData := map[string]interface{}{
		"domain": "example.com",
		"status": "active",
		"config": map[string]interface{}{
			"cache_enabled": true,
		},
	}

	// Call WriteToFile
	err := WriteToFile(d, testData)
	if err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat("test_output.json"); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	// Read and verify file contents
	content, err := os.ReadFile("test_output.json")
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify content
	if result["domain"] != "example.com" {
		t.Errorf("Expected domain 'example.com', got '%v'", result["domain"])
	}

	if result["status"] != "active" {
		t.Errorf("Expected status 'active', got '%v'", result["status"])
	}

	// Clean up
	os.Remove("test_output.json")
}

func TestWriteToFileNoOutput(t *testing.T) {
	// Create a test schema.ResourceData without result_output_file
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"result_output_file": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}, map[string]interface{}{})

	// Test data
	testData := map[string]interface{}{
		"domain": "example.com",
	}

	// Call WriteToFile - should not create any file
	err := WriteToFile(d, testData)
	if err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}

	// Verify no file was created
	if _, err := os.Stat("test_output.json"); !os.IsNotExist(err) {
		t.Fatal("Output file should not have been created")
		os.Remove("test_output.json") // Clean up if somehow created
	}
}
