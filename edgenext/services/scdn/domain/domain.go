package domain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all domain-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_domain":               ResourceEdgenextScdnDomain(),
		"edgenext_scdn_origin":               ResourceEdgenextScdnOrigin(),
		"edgenext_scdn_cert_binding":         ResourceEdgenextScdnCertBinding(),
		"edgenext_scdn_domain_base_settings": ResourceEdgenextScdnDomainBaseSettings(),
		"edgenext_scdn_domain_status":        ResourceEdgenextScdnDomainStatus(),
		"edgenext_scdn_domain_node_switch":   ResourceEdgenextScdnDomainNodeSwitch(),
		"edgenext_scdn_domain_access_mode":   ResourceEdgenextScdnDomainAccessMode(),
	}
}

// DataSources returns all domain-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_domain":               DataSourceEdgenextScdnDomain(),
		"edgenext_scdn_domains":              DataSourceEdgenextScdnDomains(),
		"edgenext_scdn_origin":               DataSourceEdgenextScdnOrigin(),
		"edgenext_scdn_origins":              DataSourceEdgenextScdnOrigins(),
		"edgenext_scdn_domain_base_settings": DataSourceEdgenextScdnDomainBaseSettings(),
		"edgenext_scdn_access_progress":      DataSourceEdgenextScdnAccessProgress(),
		"edgenext_scdn_domain_templates":     DataSourceEdgenextScdnDomainTemplates(),
		"edgenext_scdn_brief_domains":        DataSourceEdgenextScdnBriefDomains(),
	}
}
