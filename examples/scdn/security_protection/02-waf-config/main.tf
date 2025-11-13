# Example 2: WAF Rule Configuration
# This example demonstrates how to configure WAF rules

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
  description = "Business ID"
  type        = number
}

# Configure WAF rules
resource "edgenext_scdn_security_protection_waf_config" "example" {
  business_id = var.business_id

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

output "waf_config" {
  description = "WAF rule configuration"
  value = {
    business_id = edgenext_scdn_security_protection_waf_config.example.business_id
  }
}

