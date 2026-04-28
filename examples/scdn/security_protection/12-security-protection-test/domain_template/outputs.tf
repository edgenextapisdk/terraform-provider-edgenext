output "domain_id" {
  description = "Domain ID."
  value       = edgenext_scdn_security_protection_domain_template.this.domain_id
}

output "business_id" {
  description = "Business ID (template ID)."
  value       = edgenext_scdn_security_protection_domain_template.this.business_id
}

output "template_source_id" {
  description = "Source template ID."
  value       = edgenext_scdn_security_protection_domain_template.this.template_source_id
}

output "global_template_id" {
  description = "Global template ID."
  value       = data.edgenext_scdn_security_protection_member_global_template.global.template[0].id
}

output "fail_templates" {
  description = "Failed templates during batch config."
  value       = try(edgenext_scdn_security_protection_template_batch_config.this[0].fail_templates, {})
}

output "effective_config" {
  description = "Effective desired config managed by this module."
  value = {
    domain_id                   = var.domain_id
    template_source_id          = coalesce(var.template_source_id, data.edgenext_scdn_security_protection_member_global_template.global.template[0].id)
    batch_config_set            = var.ddos_config != null || var.waf_rule_config != null || var.precise_access_control_config != null
    ddos_config                 = var.ddos_config
    waf_rule_config             = var.waf_rule_config
    precise_access_control_config = var.precise_access_control_config
  }
}
