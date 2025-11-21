# Example 3: Query SCDN Domain
# This example demonstrates how to query a single SCDN domain by domain name or ID

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
  description = "Domain name to query (optional if id or domain_id is provided)"
  type        = string
  default     = ""
}

variable "domain_id" {
  description = "Domain ID to query (optional if domain_name or id is provided, deprecated - use id instead)"
  type        = string
  default     = ""
}

variable "id" {
  description = "Domain ID to query (optional if domain_name or domain_id is provided)"
  type        = string
  default     = ""
}

# Query single SCDN domain by name or ID
# You can use domain, id, or domain_id, but at least one must be provided
# Example 1: Query by domain name
#   domain = "terraform.example.com"
#
# Example 2: Query by domain ID (recommended)
#   id = "116038"
#
# Example 3: Query by domain_id (deprecated, use id instead)
#   domain_id = "116038"
#
# Example 4: Both provided (id takes priority)
#   domain = "terraform.example.com"
#   id = "116038"
data "edgenext_scdn_domain" "example" {
  domain    = var.domain_name
  id        = var.id
  domain_id = var.domain_id
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
