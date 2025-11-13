# Example 1: Create SCDN Certificate
# This example demonstrates how to create a basic SCDN certificate

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

variable "ca_name" {
  description = "Certificate name"
  type        = string
}

variable "ca_cert" {
  description = "Certificate public key (PEM format)"
  type        = string
  sensitive   = true
}

variable "ca_key" {
  description = "Certificate private key (PEM format)"
  type        = string
  sensitive   = true
}

# Create SCDN certificate
resource "edgenext_scdn_certificate" "example" {
  ca_name = var.ca_name
  ca_cert = var.ca_cert
  ca_key  = var.ca_key
}

output "certificate_id" {
  description = "Created certificate ID"
  value       = edgenext_scdn_certificate.example.id
}

output "certificate_name" {
  description = "Certificate name"
  value       = edgenext_scdn_certificate.example.ca_name
}

output "certificate_sn" {
  description = "Certificate serial number"
  value       = edgenext_scdn_certificate.example.ca_sn
}

output "issuer" {
  description = "Certificate issuer"
  value       = edgenext_scdn_certificate.example.issuer
}

output "issuer_expiry_time" {
  description = "Certificate expiry time"
  value       = edgenext_scdn_certificate.example.issuer_expiry_time
}

output "binded" {
  description = "Whether the certificate is bound"
  value       = edgenext_scdn_certificate.example.binded
}

