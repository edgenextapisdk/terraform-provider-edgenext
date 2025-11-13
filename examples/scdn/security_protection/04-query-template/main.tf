# Example 4: Query Security Protection Template
# This example demonstrates how to query a security protection template

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

# Query security protection template
data "edgenext_scdn_security_protection_template" "example" {
  business_id = var.business_id
}

output "template" {
  description = "Security protection template details"
  value = {
    business_id       = data.edgenext_scdn_security_protection_template.example.business_id
    name              = data.edgenext_scdn_security_protection_template.example.name
    type              = data.edgenext_scdn_security_protection_template.example.type
    created_at        = data.edgenext_scdn_security_protection_template.example.created_at
    remark            = data.edgenext_scdn_security_protection_template.example.remark
    bind_domain_count = data.edgenext_scdn_security_protection_template.example.bind_domain_count
  }
}

