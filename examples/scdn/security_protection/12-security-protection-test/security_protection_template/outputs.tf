output "id" {
  description = "Terraform resource ID for the security protection template."
  value       = edgenext_scdn_security_protection_template.this.id
}

output "business_id" {
  description = "Security protection template business ID."
  value       = edgenext_scdn_security_protection_template.this.business_id
}

output "name" {
  description = "Security protection template name."
  value       = edgenext_scdn_security_protection_template.this.name
}

output "type" {
  description = "Security protection template type."
  value       = edgenext_scdn_security_protection_template.this.type
}

output "created_at" {
  description = "Security protection template creation time."
  value       = edgenext_scdn_security_protection_template.this.created_at
}

output "bind_domain_count" {
  description = "Number of domains bound to the security protection template."
  value       = edgenext_scdn_security_protection_template.this.bind_domain_count
}

output "fail_domains" {
  description = "Domains that failed during template creation."
  value       = edgenext_scdn_security_protection_template.this.fail_domains
}

output "batch_config_id" {
  description = "Terraform resource ID for the optional batch config."
  value       = try(edgenext_scdn_security_protection_template_batch_config.this[0].id, null)
}

output "fail_templates" {
  description = "Templates that failed during batch config."
  value       = try(edgenext_scdn_security_protection_template_batch_config.this[0].fail_templates, {})
}

output "effective_config" {
  description = "Effective desired security protection configuration managed by this module."
  value = {
    template = {
      name               = var.name
      remark             = var.remark
      template_source_id = var.template_source_id
      domain_ids         = var.domain_ids
      group_ids          = var.group_ids
      domains            = var.domains
      bind_all           = var.bind_all
    }
    ddos_config                   = var.ddos_config
    waf_rule_config               = var.waf_rule_config
    precise_access_control_config = var.precise_access_control_config
    batch_config_set              = local.apply_batch_config
  }
}
