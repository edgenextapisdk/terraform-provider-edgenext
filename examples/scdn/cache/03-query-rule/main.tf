# Example 3: Query Cache Rule
# This example demonstrates how to query a specific cache rule by id parameter

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
  description = "Cache rule ID to query"
  type        = number
}

# Query cache rule by rule_id (converted to id query parameter in API)
resource "edgenext_scdn_cache_rule" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  rule_id       = var.rule_id
  name          = ""  # Will be read from API
  expr          = ""  # Will be read from API

  conf {
    nocache = false  # Will be read from API
    cache_share {
      scheme = "http"  # Will be read from API
    }
  }
}

output "rule" {
  description = "Cache rule details"
  value = {
    id     = edgenext_scdn_cache_rule.example.id
    name   = edgenext_scdn_cache_rule.example.name
    remark = edgenext_scdn_cache_rule.example.remark
    expr   = edgenext_scdn_cache_rule.example.expr
    status = edgenext_scdn_cache_rule.example.status
    weight = edgenext_scdn_cache_rule.example.weight
    type   = edgenext_scdn_cache_rule.example.type
    conf = length(edgenext_scdn_cache_rule.example.conf) > 0 ? {
      nocache            = edgenext_scdn_cache_rule.example.conf[0].nocache
      cache_rule         = length(edgenext_scdn_cache_rule.example.conf[0].cache_rule) > 0 ? edgenext_scdn_cache_rule.example.conf[0].cache_rule[0] : null
      browser_cache_rule = length(edgenext_scdn_cache_rule.example.conf[0].browser_cache_rule) > 0 ? edgenext_scdn_cache_rule.example.conf[0].browser_cache_rule[0] : null
      cache_errstatus    = edgenext_scdn_cache_rule.example.conf[0].cache_errstatus
      cache_url_rewrite   = length(edgenext_scdn_cache_rule.example.conf[0].cache_url_rewrite) > 0 ? edgenext_scdn_cache_rule.example.conf[0].cache_url_rewrite[0] : null
      cache_share         = length(edgenext_scdn_cache_rule.example.conf[0].cache_share) > 0 ? edgenext_scdn_cache_rule.example.conf[0].cache_share[0] : null
    } : null
  }
}
