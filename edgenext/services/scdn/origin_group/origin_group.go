package origin_group

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/origin_group/data"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/origin_group/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all origin group-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_origin_group":             resource.ResourceEdgenextScdnOriginGroup(),
		"edgenext_scdn_origin_group_domain_bind": resource.ResourceEdgenextScdnOriginGroupDomainBind(),
		"edgenext_scdn_origin_group_domain_copy": resource.ResourceEdgenextScdnOriginGroupDomainCopy(),
	}
}

// DataSources returns all origin group-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_origin_group":      data.DataSourceEdgenextScdnOriginGroup(),
		"edgenext_scdn_origin_groups":     data.DataSourceEdgenextScdnOriginGroups(),
		"edgenext_scdn_origin_groups_all": data.DataSourceEdgenextScdnOriginGroupsAll(),
		// "edgenext_scdn_origin_group_bind_history": data.DataSourceEdgenextScdnOriginGroupBindHistory(), // Temporarily disabled due to API ambiguity
	}
}
