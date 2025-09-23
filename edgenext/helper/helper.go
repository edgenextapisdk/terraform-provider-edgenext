package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataResourceIdsHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return fmt.Sprintf("%d", HashString(buf.String()))
}

func HashString(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func MergeStringBoolMaps(maps ...map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// MergeMap merges maps with overwrite: existing fields in dst are overwritten by src
func MergeMap(dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = make(map[string]interface{})
	}

	for k, v := range src {
		if vMap, ok := v.(map[string]interface{}); ok {
			// If dst also contains a map, merge recursively
			if dv, ok := dst[k].(map[string]interface{}); ok {
				dst[k] = MergeMap(dv, vMap)
			} else {
				dst[k] = MergeMap(make(map[string]interface{}), vMap)
			}
		} else {
			// Direct overwrite
			dst[k] = v
		}
	}
	return dst
}

// WriteToFile writes data to the specified file in JSON format if result_output_file is provided
func WriteToFile(d *schema.ResourceData, data interface{}) error {
	outputFile := d.Get("result_output_file").(string)
	if outputFile == "" {
		return nil
	}

	// Convert data to JSON
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Log file write operation
	log.Printf("[INFO] Writing output to file: %s", outputFile)

	// Write to file
	err = os.WriteFile(outputFile, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outputFile, err)
	}

	log.Printf("[INFO] Successfully wrote output to file: %s", outputFile)
	return nil
}
