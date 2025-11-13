# Example 6: Bind SSL Certificate to SCDN Domain
# This example demonstrates how to bind an SSL certificate to an existing SCDN domain

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

variable "domain_id" {
  description = "ID of the SCDN domain to bind the certificate to"
  type        = number
}

variable "ca_id" {
  description = "ID of the SSL certificate to bind"
  type        = number
}

variable "enabled" {
  description = "Whether to bind the certificate (true) or unbind it (false)"
  type        = bool
  default     = true
}

# Bind or unbind SSL certificate to domain
# Set enabled = true to bind, enabled = false to unbind
resource "edgenext_scdn_cert_binding" "example" {
  count = var.enabled ? 1 : 0

  domain_id = var.domain_id
  ca_id     = var.ca_id
}

output "binding_id" {
  description = "Certificate binding ID (format: domain_id-ca_id)"
  value       = var.enabled ? edgenext_scdn_cert_binding.example[0].id : null
}

output "domain_id" {
  description = "Domain ID"
  value       = var.domain_id
}

output "ca_id" {
  description = "Certificate ID"
  value       = var.ca_id
}

output "binding_status" {
  description = "Current binding status"
  value       = var.enabled ? "bound" : "unbound"
}
