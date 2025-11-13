# Example 5: Query Domains Bound to Rule Template
# This example demonstrates how to query domains bound to a specific rule template

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
  description = "Rule template ID"
  type        = number
}

variable "app_type" {
  description = "Application type (e.g., 'network_speed')"
  type        = string
}

variable "page" {
  description = "Page number for pagination"
  type        = number
  default     = 1
}

variable "page_size" {
  description = "Items per page"
  type        = number
  default     = 10
}

variable "domain" {
  description = "Filter by domain name"
  type        = string
  default     = ""
}

# Query domains bound to rule template
data "edgenext_scdn_rule_template_domains" "example" {
  id        = var.template_id
  app_type  = var.app_type
  page      = var.page
  page_size = var.page_size
  domain    = var.domain != "" ? var.domain : null
}

output "total" {
  description = "Total number of domains bound to template"
  value       = data.edgenext_scdn_rule_template_domains.example.total
}

output "domains" {
  description = "List of domains bound to the template"
  value       = data.edgenext_scdn_rule_template_domains.example.list
}

