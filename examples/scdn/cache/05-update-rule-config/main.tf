# Example 5: Update Cache Rule Configuration
# This example demonstrates how to update a cache rule's configuration

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

variable "business_id" {
  description = "Business ID (template ID for 'tpl' type, domain ID for 'domain' type)"
  type        = number
}

variable "business_type" {
  description = "Business type: 'tpl' (template) or 'domain'"
  type        = string
  default     = "tpl"
}

variable "rule_id" {
  description = "Cache rule ID to update"
  type        = number
}

variable "rule_name" {
  description = "Cache rule name"
  type        = string
}

variable "expr" {
  description = "Wirefilter rule expression"
  type        = string
}

variable "remark" {
  description = "Rule remark (optional)"
  type        = string
  default     = ""
}

variable "conf" {
  description = "Updated cache rule configuration"
  type = object({
    nocache = bool
    cache_rule = optional(object({
      cachetime             = number
      ignore_cache_time     = optional(bool)
      ignore_nocache_header = optional(bool)
      no_cache_control_op  = optional(string)
      action                = optional(string)
    }))
    browser_cache_rule = optional(object({
      cachetime        = number
      ignore_cache_time = bool
      nocache          = bool
    }))
    cache_errstatus = optional(list(object({
      cachetime  = number
      err_status = list(number)
    })))
    cache_url_rewrite = optional(object({
      sort_args   = bool
      ignore_case = bool
      queries = optional(object({
        args_method = string
        items       = list(string)
      }))
      cookies = optional(object({
        args_method = string
        items       = list(string)
      }))
    }))
    cache_share = object({
      scheme = string
    })
  })
}

# Update cache rule configuration
resource "edgenext_scdn_cache_rule" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  rule_id       = var.rule_id
  name          = var.rule_name
  expr          = var.expr
  remark        = var.remark

  conf {
    nocache = var.conf.nocache

    dynamic "cache_rule" {
      for_each = var.conf.cache_rule != null ? [var.conf.cache_rule] : []
      content {
        cachetime             = cache_rule.value.cachetime
        ignore_cache_time     = cache_rule.value.ignore_cache_time
        ignore_nocache_header = cache_rule.value.ignore_nocache_header
        no_cache_control_op  = cache_rule.value.no_cache_control_op
        action                = cache_rule.value.action
      }
    }

    dynamic "browser_cache_rule" {
      for_each = var.conf.browser_cache_rule != null ? [var.conf.browser_cache_rule] : []
      content {
        cachetime        = browser_cache_rule.value.cachetime
        ignore_cache_time = browser_cache_rule.value.ignore_cache_time
        nocache          = browser_cache_rule.value.nocache
      }
    }

    dynamic "cache_errstatus" {
      for_each = var.conf.cache_errstatus != null ? var.conf.cache_errstatus : []
      content {
        cachetime  = cache_errstatus.value.cachetime
        err_status = cache_errstatus.value.err_status
      }
    }

    dynamic "cache_url_rewrite" {
      for_each = var.conf.cache_url_rewrite != null ? [var.conf.cache_url_rewrite] : []
      content {
        sort_args   = cache_url_rewrite.value.sort_args
        ignore_case = cache_url_rewrite.value.ignore_case

        dynamic "queries" {
          for_each = cache_url_rewrite.value.queries != null ? [cache_url_rewrite.value.queries] : []
          content {
            args_method = queries.value.args_method
            items       = queries.value.items
          }
        }

        dynamic "cookies" {
          for_each = cache_url_rewrite.value.cookies != null ? [cache_url_rewrite.value.cookies] : []
          content {
            args_method = cookies.value.args_method
            items       = cookies.value.items
          }
        }
      }
    }

    cache_share {
      scheme = var.conf.cache_share.scheme
    }
  }
}

output "rule_id" {
  description = "Updated cache rule ID"
  value       = edgenext_scdn_cache_rule.example.id
}

