# Example 9: Query Access Progress Status
# This example demonstrates how to query access progress status options

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

# Query access progress status options
data "edgenext_scdn_access_progress" "example" {
}

output "access_progress_options" {
  description = "Available access progress status options"
  value       = data.edgenext_scdn_access_progress.example.progress
}

output "progress_count" {
  description = "Number of available progress status options"
  value       = length(data.edgenext_scdn_access_progress.example.progress)
}

