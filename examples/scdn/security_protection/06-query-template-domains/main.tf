# Example 6: Query Template Domains
# This example demonstrates how to query domains bound to a template

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
  description = "Business ID (template ID)"
  type        = number
}

# Query template domains
data "edgenext_scdn_security_protection_template_domains" "example" {
  business_id = var.business_id
  page        = 1
  page_size   = 20
}

output "template_domains" {
  description = "Template domains"
  value = {
    total   = data.edgenext_scdn_security_protection_template_domains.example.total
    domains = data.edgenext_scdn_security_protection_template_domains.example.domains
  }
}

