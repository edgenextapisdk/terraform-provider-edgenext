# Example 2: Create SCDN Log Download Template
# This example demonstrates how to create a log download template

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
  description = "Template name"
  type        = string
}

variable "group_name" {
  description = "Group name"
  type        = string
}

variable "group_id" {
  description = "Group ID"
  type        = number
}

variable "data_source" {
  description = "Data source: ng, cc, waf"
  type        = string
}

variable "status" {
  description = "Status: 1-enabled, 0-disabled, default: 1"
  type        = number
  default     = 1
}

variable "download_fields" {
  description = "Download fields"
  type        = list(string)
}

variable "search_terms" {
  description = "Search conditions"
  type = list(object({
    key   = string
    value = string
  }))
  default = []
}

variable "domain_select_type" {
  description = "Domain select type: 0-partial, 1-all, default: 0"
  type        = number
  default     = 0
}

# Create SCDN log download template
resource "edgenext_scdn_log_download_template" "example" {
  template_name     = var.template_name
  group_name        = var.group_name
  group_id          = var.group_id
  data_source       = var.data_source
  status            = var.status
  download_fields   = var.download_fields
  
  dynamic "search_terms" {
    for_each = var.search_terms
    content {
      key   = search_terms.value.key
      value = search_terms.value.value
    }
  }
  
  domain_select_type = var.domain_select_type
}

output "template_id" {
  description = "Created log download template ID"
  value       = edgenext_scdn_log_download_template.example.template_id
}

output "template_name" {
  description = "Template name"
  value       = edgenext_scdn_log_download_template.example.template_name
}

