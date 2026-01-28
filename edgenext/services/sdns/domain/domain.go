package domain

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/domain/data"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/domain/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all domain-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_sdns_domain": resource.ResourceEdgenextDnsDomain(),
	}
}

// DataSources returns all domain-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_sdns_domains": data.DataSourceEdgenextDnsDomain(),
	}
}
