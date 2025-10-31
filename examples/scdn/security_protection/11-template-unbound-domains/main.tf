# Example 11: Query Template Unbound Domains
# This example demonstrates how to query domains that are not bound to any template

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

variable "domain" {
  description = "Domain filter (optional)"
  type        = string
  default     = ""
}

variable "page" {
  description = "Page number"
  type        = number
  default     = 1
}

variable "page_size" {
  description = "Page size"
  type        = number
  default     = 20
}

# Query template unbound domains
data "edgenext_scdn_security_protection_template_unbound_domains" "example" {
  domain    = var.domain
  page      = var.page
  page_size = var.page_size
}

output "template_unbound_domains" {
  description = "Template unbound domains"
  value = {
    total   = data.edgenext_scdn_security_protection_template_unbound_domains.example.total
    domains = data.edgenext_scdn_security_protection_template_unbound_domains.example.domains
  }
}

