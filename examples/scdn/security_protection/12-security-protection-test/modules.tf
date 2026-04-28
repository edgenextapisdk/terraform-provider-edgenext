locals {
  enable_security_protection          = var.security_template_name != null && var.security_template_name != ""
  enable_domain_template              = var.domain_template_id != null
}

# Multi-domain template (original)
module "security_protection" {
  for_each = local.enable_security_protection ? { "this" = true } : {}
  source   = "./security_protection_template"

  name               = var.security_template_name
  remark             = var.security_template_remark
  template_source_id = var.security_template_source_id
  domain_ids         = var.security_domain_ids
  group_ids          = var.security_group_ids
  domains            = var.security_domains
  bind_all           = var.security_bind_all

  ddos_config                   = var.security_ddos_config
  waf_rule_config               = var.security_waf_rule_config
  precise_access_control_config = var.security_precise_access_control_config
}

# Single domain template (new)
module "domain_template" {
  for_each = local.enable_domain_template ? { "this" = true } : {}
  source   = "./domain_template"

  domain_id          = var.domain_template_id
  template_source_id = var.domain_template_source_id

  ddos_config                   = var.domain_template_ddos_config
  waf_rule_config               = var.domain_template_waf_rule_config
  precise_access_control_config = var.domain_template_precise_access_control_config
}
