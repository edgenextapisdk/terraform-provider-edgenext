# Example 4: Update SCDN Certificate
# This example demonstrates how to update a SCDN certificate name

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

variable "certificate_id" {
  description = "Certificate ID to update"
  type        = string
}

variable "new_ca_name" {
  description = "New certificate name"
  type        = string
}

# Optional variables for certificate content update
# variable "ca_cert" {
#   description = "Certificate public key (PEM format) - required only if updating certificate content"
#   type        = string
#   sensitive   = true
#   default     = null
# }
#
# variable "ca_key" {
#   description = "Certificate private key (PEM format) - required only if updating certificate content"
#   type        = string
#   sensitive   = true
#   default     = null
# }

# Query existing certificate to get current values
data "edgenext_scdn_certificate" "example" {
  id = var.certificate_id
}

# Update certificate name only
# You can directly apply without importing by providing certificate_id
# Note: ca_cert and ca_key are optional for updates when only updating the name.
# If you want to update certificate content, uncomment and provide the actual values:
resource "edgenext_scdn_certificate" "example" {
  certificate_id = var.certificate_id  # Specify the certificate ID to update
  ca_name        = var.new_ca_name
  # ca_cert and ca_key are optional when only updating the name
  # If you want to update certificate content, uncomment and provide the actual values:
  # ca_cert = var.ca_cert  # Optional: provide if updating certificate content
  # ca_key  = var.ca_key   # Optional: provide if updating certificate content
}

output "certificate_id" {
  description = "Certificate ID"
  value       = edgenext_scdn_certificate.example.id
}

output "certificate_name" {
  description = "Updated certificate name"
  value       = edgenext_scdn_certificate.example.ca_name
}

output "certificate_sn" {
  description = "Certificate serial number"
  value       = edgenext_scdn_certificate.example.ca_sn
}

