package cache_operate

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cache_operate/data"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/cache_operate/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all cache operate-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_cache_clean_task":   resource.ResourceEdgenextScdnCacheCleanTask(),
		"edgenext_scdn_cache_preheat_task": resource.ResourceEdgenextScdnCachePreheatTask(),
	}
}

// DataSources returns all cache operate-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_cache_clean_config":      data.DataSourceEdgenextScdnCacheCleanConfig(),
		"edgenext_scdn_cache_clean_tasks":       data.DataSourceEdgenextScdnCacheCleanTasks(),
		"edgenext_scdn_cache_clean_task_detail": data.DataSourceEdgenextScdnCacheCleanTaskDetail(),
		"edgenext_scdn_cache_preheat_tasks":     data.DataSourceEdgenextScdnCachePreheatTasks(),
	}
}
