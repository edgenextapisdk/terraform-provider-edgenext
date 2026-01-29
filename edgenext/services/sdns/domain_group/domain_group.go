package domain_group

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/domain_group/data"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns/domain_group/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all domain group-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_sdns_domain_group": resource.ResourceEdgenextDnsGroup(),
	}
}

// DataSources returns all domain group-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_sdns_domain_groups": data.DataSourceEdgenextDnsGroup(),
	}
}
