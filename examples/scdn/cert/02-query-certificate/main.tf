# Example 2: Query SCDN Certificate
# This example demonstrates how to query a SCDN certificate by ID

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
  description = "Certificate ID to query"
  type        = string
}

# Query SCDN certificate
data "edgenext_scdn_certificate" "example" {
  id = var.certificate_id
}

output "certificate_id" {
  description = "Certificate ID"
  value       = data.edgenext_scdn_certificate.example.id
}

output "certificate_name" {
  description = "Certificate name"
  value       = data.edgenext_scdn_certificate.example.ca_name
}

output "issuer" {
  description = "Certificate issuer"
  value       = data.edgenext_scdn_certificate.example.issuer
}

output "issuer_start_time" {
  description = "Certificate start time"
  value       = data.edgenext_scdn_certificate.example.issuer_start_time
}

output "issuer_expiry_time" {
  description = "Certificate expiry time"
  value       = data.edgenext_scdn_certificate.example.issuer_expiry_time
}

output "issuer_expiry_time_desc" {
  description = "Certificate expiry time description"
  value       = data.edgenext_scdn_certificate.example.issuer_expiry_time_desc
}

output "binded" {
  description = "Whether the certificate is bound"
  value       = data.edgenext_scdn_certificate.example.binded
}

output "ca_domain" {
  description = "Domains in the certificate"
  value       = data.edgenext_scdn_certificate.example.ca_domain
}

output "apply_status" {
  description = "Application status"
  value       = data.edgenext_scdn_certificate.example.apply_status
}

output "ca_type" {
  description = "Certificate type"
  value       = data.edgenext_scdn_certificate.example.ca_type
}

