# Example 1: Create SCDN Rule Template
# This example demonstrates how to create a basic SCDN rule template

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

variable "template_name" {
  description = "Rule template name"
  type        = string
}

variable "description" {
  description = "Rule template description"
  type        = string
  default     = ""
}

variable "app_type" {
  description = "Application type (e.g., 'network_speed')"
  type        = string
}

variable "domain_ids" {
  description = "List of domain IDs to bind to the template"
  type        = list(number)
  default     = []
}

# Create SCDN rule template
resource "edgenext_scdn_rule_template" "example" {
  name        = var.template_name
  description = var.description
  app_type    = var.app_type

  dynamic "bind_domain" {
    for_each = length(var.domain_ids) > 0 ? [1] : []
    content {
      domain_ids = var.domain_ids
      is_bind    = true
    }
  }
}

output "template_id" {
  description = "Created rule template ID"
  value       = edgenext_scdn_rule_template.example.id
}

output "template_name" {
  description = "Rule template name"
  value       = edgenext_scdn_rule_template.example.name
}

output "created_at" {
  description = "Template creation timestamp"
  value       = edgenext_scdn_rule_template.example.created_at
}

