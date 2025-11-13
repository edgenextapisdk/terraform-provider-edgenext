package log_download

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/log_download/data"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/log_download/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all log download-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_log_download_task":            resource.ResourceEdgenextScdnLogDownloadTask(),
		"edgenext_scdn_log_download_template":        resource.ResourceEdgenextScdnLogDownloadTemplate(),
		"edgenext_scdn_log_download_template_status": resource.ResourceEdgenextScdnLogDownloadTemplateStatus(),
	}
}

// DataSources returns all log download-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_log_download_tasks":     data.DataSourceEdgenextScdnLogDownloadTasks(),
		"edgenext_scdn_log_download_templates": data.DataSourceEdgenextScdnLogDownloadTemplates(),
		"edgenext_scdn_log_download_fields":    data.DataSourceEdgenextScdnLogDownloadFields(),
	}
}
