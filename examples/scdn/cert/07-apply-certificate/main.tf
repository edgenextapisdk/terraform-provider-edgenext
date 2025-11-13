# Example 7: Apply for SCDN Certificate
# This example demonstrates how to apply for a certificate (Let's Encrypt, etc.)

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

variable "domains" {
  description = "List of domains to apply for certificate"
  type        = list(string)
}

# Apply for certificate
resource "edgenext_scdn_certificate_apply" "example" {
  domain = var.domains
}

output "certificate_application_id" {
  description = "Certificate application resource ID"
  value       = edgenext_scdn_certificate_apply.example.id
}

output "ca_id_domains" {
  description = "Mapping of domain_id to domain"
  value       = edgenext_scdn_certificate_apply.example.ca_id_domains
}

output "ca_id_names" {
  description = "Mapping of ca_id to ca_name"
  value       = edgenext_scdn_certificate_apply.example.ca_id_names
}

output "domain_count" {
  description = "Number of domains applied"
  value       = length(var.domains)
}

output "certificate_ids" {
  description = "List of certificate IDs created"
  value = [
    for ca_id, _ in edgenext_scdn_certificate_apply.example.ca_id_names : ca_id
  ]
}

output "certificate_names" {
  description = "List of certificate names created"
  value = [
    for _, ca_name in edgenext_scdn_certificate_apply.example.ca_id_names : ca_name
  ]
}

