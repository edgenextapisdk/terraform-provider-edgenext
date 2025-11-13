# Example 2: Query SCDN Rule Template
# This example demonstrates how to query a specific SCDN rule template

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

variable "template_id" {
  description = "Rule template ID to query"
  type        = string
}

variable "app_type" {
  description = "Application type (e.g., 'network_speed')"
  type        = string
}

# Query SCDN rule template
data "edgenext_scdn_rule_template" "example" {
  id       = var.template_id
  app_type = var.app_type
}

output "template_id" {
  description = "Rule template ID"
  value       = data.edgenext_scdn_rule_template.example.id
}

output "template_name" {
  description = "Rule template name"
  value       = data.edgenext_scdn_rule_template.example.name
}

output "description" {
  description = "Rule template description"
  value       = data.edgenext_scdn_rule_template.example.description
}

output "created_at" {
  description = "Template creation timestamp"
  value       = data.edgenext_scdn_rule_template.example.created_at
}

output "bind_domains" {
  description = "List of domains bound to this template"
  value       = data.edgenext_scdn_rule_template.example.bind_domains
}

