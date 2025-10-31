# Example 5: Query SCDN Log Download Fields
# This example demonstrates how to query available log download fields

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

# Query SCDN log download fields
data "edgenext_scdn_log_download_fields" "example" {
  # Optional: filter by data source (ng, cc, waf)
  # data_source = "ng"
}

output "configs" {
  description = "Field configurations by data source"
  value       = data.edgenext_scdn_log_download_fields.example.configs
}

# Example: Get download fields for ng data source
output "ng_download_fields" {
  description = "Download fields for ng data source"
  value = [
    for config in data.edgenext_scdn_log_download_fields.example.configs : config.download_fields
    if config.data_source == "ng"
  ]
}

