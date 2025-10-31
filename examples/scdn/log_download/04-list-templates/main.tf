# Example 4: List SCDN Log Download Templates
# This example demonstrates how to query log download templates

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

variable "status" {
  description = "Status: 1-enabled, 0-disabled"
  type        = number
  default     = null
}

variable "group_id" {
  description = "Group ID"
  type        = number
  default     = null
}

variable "template_name" {
  description = "Template name"
  type        = string
  default     = null
}

variable "data_source" {
  description = "Data source: ng, cc, waf"
  type        = string
  default     = null
}

variable "page" {
  description = "Page number"
  type        = number
  default     = 1
}

variable "per_page" {
  description = "Items per page"
  type        = number
  default     = 20
}

# Query SCDN log download templates
data "edgenext_scdn_log_download_templates" "example" {
  status       = var.status
  group_id     = var.group_id
  template_name = var.template_name
  data_source  = var.data_source
  page         = var.page
  per_page     = var.per_page
}

output "total" {
  description = "Total number of templates"
  value       = data.edgenext_scdn_log_download_templates.example.total
}

output "templates" {
  description = "Template list"
  value       = data.edgenext_scdn_log_download_templates.example.templates
}

