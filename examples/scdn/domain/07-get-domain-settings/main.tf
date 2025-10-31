# Example 7: Query Domain Base Settings
# This example demonstrates how to query domain base settings

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

variable "domain_id" {
  description = "Domain ID to query base settings"
  type        = number
}

# Query domain base settings
data "edgenext_scdn_domain_base_settings" "example" {
  domain_id = var.domain_id
}

output "domain_settings" {
  description = "Domain base settings information"
  value = {
    domain_id       = var.domain_id
    proxy_host      = data.edgenext_scdn_domain_base_settings.example.proxy_host
    proxy_sni       = data.edgenext_scdn_domain_base_settings.example.proxy_sni
    domain_redirect = data.edgenext_scdn_domain_base_settings.example.domain_redirect
  }
}

