package record

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/record/data"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/record/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all record-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_sdns_record": resource.ResourceEdgenextDnsRecord(),
	}
}

// DataSources returns all record-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_sdns_records": data.DataSourceEdgenextDnsRecord(),
	}
}
