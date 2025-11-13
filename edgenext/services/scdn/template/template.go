package template

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all rule template-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_rule_template":               ResourceEdgenextScdnRuleTemplate(),
		"edgenext_scdn_rule_template_domain_bind":   ResourceEdgenextScdnRuleTemplateDomainBind(),
		"edgenext_scdn_rule_template_domain_unbind": ResourceEdgenextScdnRuleTemplateDomainUnbind(),
	}
}

// DataSources returns all rule template-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_rule_template":         DataSourceEdgenextScdnRuleTemplate(),
		"edgenext_scdn_rule_templates":        DataSourceEdgenextScdnRuleTemplates(),
		"edgenext_scdn_rule_template_domains": DataSourceEdgenextScdnRuleTemplateDomains(),
	}
}
