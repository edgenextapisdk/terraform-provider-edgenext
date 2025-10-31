package security_protect

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/security_protect/data"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/security_protect/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Note: Import paths need to be adjusted based on actual package structure

// Resources returns all security protection-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_security_protection_ddos_config":           resource.ResourceEdgenextScdnSecurityProtectionDdosConfig(),
		"edgenext_scdn_security_protection_waf_config":            resource.ResourceEdgenextScdnSecurityProtectionWafConfig(),
		"edgenext_scdn_security_protection_template":              resource.ResourceEdgenextScdnSecurityProtectionTemplate(),
		"edgenext_scdn_security_protection_template_domain_bind":  resource.ResourceEdgenextScdnSecurityProtectionTemplateDomainBind(),
		"edgenext_scdn_security_protection_template_batch_config": resource.ResourceEdgenextScdnSecurityProtectionTemplateBatchConfig(),
	}
}

// DataSources returns all security protection-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_security_protection_ddos_config":              data.DataSourceEdgenextScdnSecurityProtectionDdosConfig(),
		"edgenext_scdn_security_protection_waf_config":               data.DataSourceEdgenextScdnSecurityProtectionWafConfig(),
		"edgenext_scdn_security_protection_template":                 data.DataSourceEdgenextScdnSecurityProtectionTemplate(),
		"edgenext_scdn_security_protection_templates":                data.DataSourceEdgenextScdnSecurityProtectionTemplates(),
		"edgenext_scdn_security_protection_template_domains":         data.DataSourceEdgenextScdnSecurityProtectionTemplateDomains(),
		"edgenext_scdn_security_protection_template_unbound_domains": data.DataSourceEdgenextScdnSecurityProtectionTemplateUnboundDomains(),
		"edgenext_scdn_security_protection_member_global_template":   data.DataSourceEdgenextScdnSecurityProtectionMemberGlobalTemplate(),
		"edgenext_scdn_security_protection_iota":                     data.DataSourceEdgenextScdnSecurityProtectionIota(),
	}
}
