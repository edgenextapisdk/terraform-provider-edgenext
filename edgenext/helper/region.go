package helper

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NormalizeRegion normalizes a region input value for API requests.
// It trims spaces and converts it to lower-case to match backend expectations.
func NormalizeRegion(region string) string {
	return strings.ToLower(strings.TrimSpace(region))
}

// NormalizeRegionStateFunc is used as schema.SchemaStateFunc.
// Terraform passes raw interface{} values, so we normalize safely.
func NormalizeRegionStateFunc(v interface{}) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return NormalizeRegion(s)
}

// RegionDiffSuppressFunc ignores case-only region changes.
func RegionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

// RegionDataSchema returns the standard ECS data source region schema.
func RegionDataSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		StateFunc:        NormalizeRegionStateFunc,
		DiffSuppressFunc: RegionDiffSuppressFunc,
		Description:      description,
	}
}

// RegionResourceSchema returns the standard ECS resource region schema.
func RegionResourceSchema(description string) *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		StateFunc:        NormalizeRegionStateFunc,
		ForceNew:         true,
		DiffSuppressFunc: RegionDiffSuppressFunc,
		Description:      description,
	}
}
