terraform {
  required_providers {
    edgenext = {
      source = "edgenextapisdk/edgenext"
    }
  }
}

data "edgenext_scdn_security_protection_member_global_template" "global" {}

resource "edgenext_scdn_security_protection_domain_template" "this" {
  domain_id          = var.domain_id
  template_source_id = coalesce(var.template_source_id, data.edgenext_scdn_security_protection_member_global_template.global.template[0].id)
}

locals {
  apply_batch_config = var.ddos_config != null || var.waf_rule_config != null || var.precise_access_control_config != null
}

resource "edgenext_scdn_security_protection_template_batch_config" "this" {
  count = local.apply_batch_config ? 1 : 0

  template_ids = [edgenext_scdn_security_protection_domain_template.this.business_id]

  dynamic "ddos_config" {
    iterator = ddos
    for_each = var.ddos_config != null ? [var.ddos_config] : []
    content {
      dynamic "application_ddos_protection" {
        iterator = app_ddos
        for_each = try(ddos.value.application_ddos_protection, null) != null ? [ddos.value.application_ddos_protection] : []
        content {
          status                = try(app_ddos.value.status, null)
          ai_cc_status          = try(app_ddos.value.ai_cc_status, null)
          type                  = try(app_ddos.value.type, null)
          need_attack_detection = try(app_ddos.value.need_attack_detection, null)
          ai_status             = try(app_ddos.value.ai_status, null)
        }
      }

      dynamic "visitor_authentication" {
        iterator = visitor_auth
        for_each = try(ddos.value.visitor_authentication, null) != null ? [ddos.value.visitor_authentication] : []
        content {
          status           = try(visitor_auth.value.status, null)
          auth_token       = try(visitor_auth.value.auth_token, null)
          pass_still_check = try(visitor_auth.value.pass_still_check, null)
        }
      }
    }
  }

  dynamic "waf_rule_config" {
    iterator = waf
    for_each = var.waf_rule_config != null ? [var.waf_rule_config] : []
    content {
      dynamic "waf_rule_config" {
        iterator = waf_rule
        for_each = try(waf.value.waf_rule_config, null) != null ? [waf.value.waf_rule_config] : []
        content {
          status          = try(waf_rule.value.status, null)
          ai_status       = try(waf_rule.value.ai_status, null)
          waf_level       = try(waf_rule.value.waf_level, null)
          waf_mode        = try(waf_rule.value.waf_mode, null)
          waf_strategy_id = try(waf_rule.value.waf_strategy_id, null)
        }
      }

      dynamic "waf_intercept_page" {
        iterator = intercept_page
        for_each = try(waf.value.waf_intercept_page, null) != null ? [waf.value.waf_intercept_page] : []
        content {
          status  = try(intercept_page.value.status, null)
          type    = try(intercept_page.value.type, null)
          content = try(intercept_page.value.content, null)
        }
      }
    }
  }

  dynamic "precise_access_control_config" {
    iterator = precise
    for_each = var.precise_access_control_config != null ? [var.precise_access_control_config] : []
    content {
      action = precise.value.action

      dynamic "policies" {
        iterator = policy
        for_each = try(precise.value.policies, [])
        content {
          type        = try(policy.value.type, null)
          action      = try(policy.value.action, null)
          action_data = try(policy.value.action_data, null)
          rules       = try(policy.value.rules, null)
          from        = try(policy.value.from, null)
          status      = try(policy.value.status, null)
        }
      }
    }
  }
}
