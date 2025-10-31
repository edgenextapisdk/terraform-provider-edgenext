# Example 4: List SCDN Domains
# This example demonstrates how to query a list of SCDN domains

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

variable "page" {
  description = "Page number"
  type        = number
  default     = 1
}

variable "page_size" {
  description = "Number of items per page"
  type        = number
  default     = 10
}

variable "group_id" {
  description = "Filter by group ID (optional)"
  type        = number
  default     = null
}

# Query list of SCDN domains
data "edgenext_scdn_domains" "example" {
  page      = var.page
  page_size = var.page_size
  group_id  = var.group_id
}

output "domains_list" {
  description = "List of SCDN domains"
  value = [
    for domain in data.edgenext_scdn_domains.example.domains : {
      id              = domain.id
      domain          = domain.domain
      remark          = domain.remark
      access_progress = domain.access_progress
      protect_status  = domain.protect_status
    }
  ]
}

output "total_count" {
  description = "Total number of domains"
  value       = length(data.edgenext_scdn_domains.example.domains)
}
