package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetDataSourceStableID builds a deterministic ID from selected query arguments.
func SetDataSourceStableID(d *schema.ResourceData, keys ...string) {
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value, ok := d.GetOk(key)
		if !ok {
			parts = append(parts, fmt.Sprintf("%s=", key))
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", key, normalizeDataSourceIDValue(value)))
	}

	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	d.SetId(hex.EncodeToString(sum[:])[:32])
}

func normalizeDataSourceIDValue(value interface{}) string {
	switch v := value.(type) {
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}
		items := make([]string, 0, len(v))
		for _, item := range v {
			items = append(items, fmt.Sprintf("%v", item))
		}
		sort.Strings(items)
		return "[" + strings.Join(items, ",") + "]"
	default:
		return fmt.Sprintf("%v", value)
	}
}
