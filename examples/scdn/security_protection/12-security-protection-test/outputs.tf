output "security_protection_template_id" {
  description = "Security protection template business ID."
  value       = try(module.security_protection["this"].business_id, null)
}

output "security_protection_template_name" {
  description = "Security protection template name."
  value       = try(module.security_protection["this"].name, null)
}

output "security_protection_template_type" {
  description = "Security protection template type."
  value       = try(module.security_protection["this"].type, null)
}

output "security_protection_fail_domains" {
  description = "Domains that failed during template creation."
  value       = try(module.security_protection["this"].fail_domains, {})
}

output "security_protection_fail_templates" {
  description = "Templates that failed during batch config."
  value       = try(module.security_protection["this"].fail_templates, {})
}

output "security_protection_effective_config" {
  description = "Effective desired security protection config managed by this test root module."
  value       = try(module.security_protection["this"].effective_config, null)
}

# Domain template outputs (single domain)
output "domain_template_domain_id" {
  description = "Domain ID for the single domain template."
  value       = try(module.domain_template["this"].domain_id, null)
}

output "domain_template_business_id" {
  description = "Business ID (template ID) for the single domain template."
  value       = try(module.domain_template["this"].business_id, null)
}

output "domain_template_source_id" {
  description = "Source template ID for the single domain template."
  value       = try(module.domain_template["this"].template_source_id, null)
}

output "domain_template_global_template_id" {
  description = "Global template ID used for the single domain template."
  value       = try(module.domain_template["this"].global_template_id, null)
}

output "domain_template_fail_templates" {
  description = "Failed templates during batch config for single domain template."
  value       = try(module.domain_template["this"].fail_templates, {})
}

output "domain_template_effective_config" {
  description = "Effective desired config managed by the single domain template module."
  value       = try(module.domain_template["this"].effective_config, null)
}
