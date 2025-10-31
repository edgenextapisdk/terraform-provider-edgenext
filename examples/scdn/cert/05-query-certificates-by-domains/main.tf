# Example 5: Query SCDN Certificates by Domains
# This example demonstrates how to query certificates by a list of domains

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
  description = "List of domains to query certificates for"
  type        = list(string)
}

# Query certificates by domains
data "edgenext_scdn_certificates_by_domains" "example" {
  domains = var.domains
}

output "certificates" {
  description = "List of certificates for the specified domains"
  value       = data.edgenext_scdn_certificates_by_domains.example.certificates
}

output "certificate_count" {
  description = "Number of certificates found"
  value       = length(data.edgenext_scdn_certificates_by_domains.example.certificates)
}

output "certificate_details" {
  description = "Detailed information about certificates"
  value = {
    for idx, cert in data.edgenext_scdn_certificates_by_domains.example.certificates : idx => {
      id                = cert.id
      ca_name           = cert.ca_name
      issuer            = cert.issuer
      ca_domain         = cert.ca_domain
      issuer_expiry_time = cert.issuer_expiry_time
      binded            = cert.binded
      apply_status      = cert.apply_status
    }
  }
}

