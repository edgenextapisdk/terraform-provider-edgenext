# Example 14: List Brief Domains
# This example demonstrates how to query a brief list of domains

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

variable "domain_ids" {
  description = "List of domain IDs to query (optional, queries all if not specified)"
  type        = list(number)
  default     = []
}

# Query brief domains
data "edgenext_scdn_brief_domains" "example" {
  ids = var.domain_ids
}

output "brief_domains_list" {
  description = "List of brief domain information"
  value = [
    for domain in data.edgenext_scdn_brief_domains.example.list : {
      id     = domain.id
      domain = domain.domain
    }
  ]
}

output "total_count" {
  description = "Total number of domains"
  value       = data.edgenext_scdn_brief_domains.example.total
}

