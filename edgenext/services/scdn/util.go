package scdn

import (
	"encoding/json"
	"fmt"
)

// extractArrayFromData extracts array data from various formats (object, string, array)
// result must be a pointer to a slice type
func extractArrayFromData(data interface{}, result interface{}) error {
	if data == nil {
		// Set result to empty slice
		return json.Unmarshal([]byte("[]"), result)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Check if data is empty
	dataStr := string(dataBytes)
	if dataStr == "[]" || dataStr == "{}" || dataStr == "null" {
		return json.Unmarshal([]byte("[]"), result)
	}

	// Try to parse data - handle different formats
	var dataValue interface{}
	if err := json.Unmarshal(dataBytes, &dataValue); err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	// Handle different data formats
	var finalDataBytes []byte
	switch v := dataValue.(type) {
	case []interface{}:
		// If Data is already an array, use it directly
		finalDataBytes, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal array: %w", err)
		}
	case string:
		// If Data is a JSON string (like "{\"data\":[]}"), parse it
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(v), &jsonData); err != nil {
			return fmt.Errorf("failed to parse JSON string: %w", err)
		}
		if dataField, ok := jsonData["data"]; ok {
			return extractArrayFromData(dataField, result)
		}
		// Empty data field, return empty array
		return json.Unmarshal([]byte("[]"), result)
	case map[string]interface{}:
		// If Data is a map, check if it has a "data" field
		if dataField, ok := v["data"]; ok {
			return extractArrayFromData(dataField, result)
		}
		// If it's an empty object, return empty array
		return json.Unmarshal([]byte("[]"), result)
	default:
		// For other types, try to unmarshal directly
		finalDataBytes, err = json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal data value: %w", err)
		}
	}

	// Unmarshal into result
	if err := json.Unmarshal(finalDataBytes, result); err != nil {
		return fmt.Errorf("failed to unmarshal array data: %w", err)
	}

	return nil
}
