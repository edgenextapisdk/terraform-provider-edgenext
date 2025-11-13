# Example 8: Query Member Global Template
# This example demonstrates how to query the member global security protection template

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

# Query member global template
data "edgenext_scdn_security_protection_member_global_template" "example" {
}

output "member_global_template" {
  description = "Member global template details"
  value = {
    bind_domain_count = data.edgenext_scdn_security_protection_member_global_template.example.bind_domain_count
    template          = data.edgenext_scdn_security_protection_member_global_template.example.template
  }
}

