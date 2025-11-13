# Example 8: Copy Origin Group to Domain
# This example demonstrates how to copy an origin group to a domain

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

variable "origin_group_id" {
  description = "Origin group ID"
  type        = number
}

variable "domain_id" {
  description = "Domain ID"
  type        = number
}

# Copy origin group to domain
resource "edgenext_scdn_origin_group_domain_copy" "example" {
  origin_group_id = var.origin_group_id
  domain_id       = var.domain_id
}

output "copy_result" {
  description = "Copy result"
  value = {
    origin_group_id = edgenext_scdn_origin_group_domain_copy.example.origin_group_id
    domain_id       = edgenext_scdn_origin_group_domain_copy.example.domain_id
  }
}

