# Example 13: Domain Template
# This example demonstrates how to create a domain-level security protection template for a single domain

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
  description = "Domain ID to create template for"
  type        = number
}

# Get the global template ID (required for template_source_id)
data "edgenext_scdn_security_protection_member_global_template" "global" {}

# Create a domain template for a single domain
# This creates a separate security protection template for the domain
resource "edgenext_scdn_security_protection_domain_template" "example" {
  domain_id          = var.domain_id
  template_source_id = data.edgenext_scdn_security_protection_member_global_template.global.template[0].id
}

output "domain_template" {
  description = "Domain template creation result"
  value = {
    domain_id          = edgenext_scdn_security_protection_domain_template.example.domain_id
    template_source_id = edgenext_scdn_security_protection_domain_template.example.template_source_id
    business_id        = edgenext_scdn_security_protection_domain_template.example.business_id
  }
}

output "global_template" {
  description = "Global template information"
  value = {
    template_id       = data.edgenext_scdn_security_protection_member_global_template.global.template[0].id
    template_name     = data.edgenext_scdn_security_protection_member_global_template.global.template[0].name
    bind_domain_count = data.edgenext_scdn_security_protection_member_global_template.global.bind_domain_count
  }
}