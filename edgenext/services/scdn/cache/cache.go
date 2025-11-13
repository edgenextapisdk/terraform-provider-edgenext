package cache

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all cache-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_cache_rule":        ResourceEdgenextScdnCacheRule(),
		"edgenext_scdn_cache_rule_status": ResourceEdgenextScdnCacheRuleStatus(),
		"edgenext_scdn_cache_rules_sort":  ResourceEdgenextScdnCacheRulesSort(),
	}
}

// DataSources returns all cache-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_cache_rules":         DataSourceEdgenextScdnCacheRules(),
		"edgenext_scdn_cache_global_config": DataSourceEdgenextScdnCacheGlobalConfig(),
	}
}
