# Example 3: Query SCDN Domain
# This example demonstrates how to query a single SCDN domain by domain name

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

variable "domain_name" {
  description = "Domain name to query"
  type        = string
}

# Query single SCDN domain
data "edgenext_scdn_domain" "example" {
  domain = var.domain_name
}

output "domain_info" {
  description = "Domain information"
  value = {
    id              = data.edgenext_scdn_domain.example.id
    domain          = data.edgenext_scdn_domain.example.domain
    remark          = data.edgenext_scdn_domain.example.remark
    access_progress = data.edgenext_scdn_domain.example.access_progress
    protect_status  = data.edgenext_scdn_domain.example.protect_status
    cname           = data.edgenext_scdn_domain.example.cname
  }
}
