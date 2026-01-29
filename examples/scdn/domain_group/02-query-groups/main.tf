terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "0.2.1"
    }
  }
}

provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
}

variable "access_key" {
  type        = string
  description = "EdgeNext access key"
  sensitive   = true
}

variable "secret_key" {
  type        = string
  description = "EdgeNext secret key"
  sensitive   = true
}

variable "endpoint" {
  type        = string
  description = "EdgeNext API endpoint"
  default     = "https://api.edgenextscdn.com"
}

variable "group_name" {
  type        = string
  description = "Filter by group name (optional)"
  default     = ""
}

variable "domain" {
  type        = string
  description = "Filter by domain (optional)"
  default     = ""
}

# Query all domain groups
data "edgenext_scdn_domain_groups" "all" {
  page     = 1
  per_page = 20
}

# Query domain groups by name
data "edgenext_scdn_domain_groups" "by_name" {
  count      = var.group_name != "" ? 1 : 0
  group_name = var.group_name
  page       = 1
  per_page   = 10
}

# Query domain groups by domain
data "edgenext_scdn_domain_groups" "by_domain" {
  count    = var.domain != "" ? 1 : 0
  domain   = var.domain
  page     = 1
  per_page = 10
}

output "all_groups" {
  value       = data.edgenext_scdn_domain_groups.all.list
  description = "All domain groups"
}

output "total_groups" {
  value       = data.edgenext_scdn_domain_groups.all.total
  description = "Total number of domain groups"
}

output "filtered_by_name" {
  value       = var.group_name != "" ? data.edgenext_scdn_domain_groups.by_name[0].list : []
  description = "Groups filtered by name"
}

output "filtered_by_domain" {
  value       = var.domain != "" ? data.edgenext_scdn_domain_groups.by_domain[0].list : []
  description = "Groups filtered by domain"
}
