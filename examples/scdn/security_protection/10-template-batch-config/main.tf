# Example 10: Batch Config Security Protection Template
# This example demonstrates how to batch configure security protection templates

terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}

provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
}

variable "access_key" {
  description = "EdgeNext Access Key"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "EdgeNext Secret Key"
  type        = string
  sensitive   = true
}

variable "endpoint" {
  description = "EdgeNext SCDN API Endpoint"
  type        = string
  default     = "https://api.edgenextscdn.com"
}

variable "template_ids" {
  description = "Template ID list (only single domain templates allowed)"
  type        = list(number)
}

variable "domains" {
  description = "Domain list to apply configuration"
  type        = list(string)
  default     = []
}

# Batch configure security protection templates
resource "edgenext_scdn_security_protection_template_batch_config" "example" {
  template_ids = var.template_ids
  domains      = var.domains

  ddos_config {
    application_ddos_protection {
      status                = "on"
      ai_cc_status          = "on"
      type                  = "strict"
      need_attack_detection = 1
      ai_status             = "on"
    }
  }

  waf_rule_config {
    waf_rule_config {
      status    = "on"
      ai_status = "on"
      waf_level = "strict"
      waf_mode  = "block"
    }

    waf_intercept_page {
      status  = "on"
      type    = "default"
      content = ""
    }
  }

  # Precise access control configuration
  precise_access_control_config {
    action = "add"

    # Policy 1: Anti-CC protection for specific URL
    policies {
      action = "anticc"
      action_data = {
        level = "default"
      }
      type   = "plus"
      status = 1
      sort   = 1

      rules {
        rule_type = "url"
        logic     = "contains"
        data      = jsonencode(["/aaa"])
      }
    }

    # Policy 2: Deny access from specific referer
    policies {
      action = "deny"
      type   = "plus"
      from   = "aR"
      status = 1

      rules {
        rule_type = "referer_domain"
        logic     = "not_equals"
        data      = jsonencode(["home.console.prxcdn.com"])
      }

      rules {
        rule_type = "referer"
        logic     = "len_greater_than"
        data      = jsonencode({ len = 0 })
      }

      rules {
        rule_type = "postfix"
        logic     = "equals"
        data      = jsonencode(["css", "js", "txt", "img", "png", "jpg", "jpeg", "gif", "svg", "ico"])
      }
    }

    # Policy 3: Region-based access control
    policies {
      action = "deny"
      type   = "plus"
      from   = "zL"
      status = 1

      rules {
        rule_type = "region"
        logic     = "not_belongs"
        data      = jsonencode({ province = [], country = ["CN", "MN", "JP"] })
      }
    }
  }
}

output "template_batch_config" {
  description = "Template batch config result"
  value = {
    template_ids   = edgenext_scdn_security_protection_template_batch_config.example.template_ids
    fail_templates = edgenext_scdn_security_protection_template_batch_config.example.fail_templates
  }
}