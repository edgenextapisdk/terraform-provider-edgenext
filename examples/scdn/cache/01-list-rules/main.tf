# Example 1: List Cache Rules
# This example demonstrates how to list cache rules for a template or domain

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
  description = "Business ID (template ID for 'tpl' type, domain ID for 'domain' type)"
  type        = number
}

variable "business_type" {
  description = "Business type: 'tpl' (template) or 'domain'"
  type        = string
  default     = "tpl"
}

variable "page" {
  description = "Page number (optional)"
  type        = number
  default     = null
}

variable "page_size" {
  description = "Page size (optional)"
  type        = number
  default     = null
}

# Query cache rules
data "edgenext_scdn_cache_rules" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  page          = var.page
  page_size     = var.page_size
}

output "total" {
  description = "Total number of cache rules"
  value       = data.edgenext_scdn_cache_rules.example.total
}

output "rules" {
  description = "List of cache rules"
  value       = data.edgenext_scdn_cache_rules.example.list
}

