package edgenext

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// mergeProviderResources registers all entries from src into dst. Panics if any key already exists.
func mergeProviderResources(dst map[string]*schema.Resource, source string, src map[string]*schema.Resource) {
	for k, v := range src {
		if _, exists := dst[k]; exists {
			panic(fmt.Sprintf("duplicate provider resource %q: already registered (conflict from %s)", k, source))
		}
		dst[k] = v
	}
}

// mergeProviderDataSources registers all entries from src into dst. Panics if any key already exists.
func mergeProviderDataSources(dst map[string]*schema.Resource, source string, src map[string]*schema.Resource) {
	for k, v := range src {
		if _, exists := dst[k]; exists {
			panic(fmt.Sprintf("duplicate provider data source %q: already registered (conflict from %s)", k, source))
		}
		dst[k] = v
	}
}
