package edgenext

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestMergeProviderResourcesDetectsDuplicate(t *testing.T) {
	dst := map[string]*schema.Resource{
		"edgenext_example": {},
	}
	src := map[string]*schema.Resource{
		"edgenext_example": {},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate resource key")
		}
	}()
	mergeProviderResources(dst, "test/module", src)
}

func TestMergeProviderDataSourcesDetectsDuplicate(t *testing.T) {
	dst := map[string]*schema.Resource{
		"edgenext_example": {},
	}
	src := map[string]*schema.Resource{
		"edgenext_example": {},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on duplicate data source key")
		}
	}()
	mergeProviderDataSources(dst, "test/module", src)
}

func TestProviderRegistrationNoDuplicateKeys(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Provider() panicked during registration: %v", r)
		}
	}()
	p := Provider()
	if len(p.ResourcesMap) == 0 || len(p.DataSourcesMap) == 0 {
		t.Fatal("expected non-empty ResourcesMap and DataSourcesMap")
	}
}

func TestMergeProviderResourcesNoConflict(t *testing.T) {
	dst := map[string]*schema.Resource{}
	mergeProviderResources(dst, "a", map[string]*schema.Resource{"edgenext_a": {}})
	mergeProviderResources(dst, "b", map[string]*schema.Resource{"edgenext_b": {}})
	if len(dst) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(dst))
	}
}
