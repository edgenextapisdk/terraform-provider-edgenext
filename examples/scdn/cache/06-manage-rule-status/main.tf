# Example 6: Manage Cache Rule Status
# This example demonstrates how to enable or disable cache rules

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

variable "rule_ids" {
  description = "List of cache rule IDs to update status"
  type        = list(number)
}

variable "status" {
  description = "Status: 1 (enabled) or 2 (disabled)"
  type        = number
}

# Manage cache rule status
resource "edgenext_scdn_cache_rule_status" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  rule_ids      = var.rule_ids
  status        = var.status
}

output "updated_ids" {
  description = "Rule IDs that were updated"
  value       = edgenext_scdn_cache_rule_status.example.updated_ids
}

