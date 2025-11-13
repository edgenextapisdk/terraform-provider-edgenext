# Example 3: List SCDN Certificates
# This example demonstrates how to list SCDN certificates with filters

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

variable "domain" {
  description = "Filter by domain name"
  type        = string
  default     = ""
}

variable "ca_name" {
  description = "Filter by certificate name"
  type        = string
  default     = ""
}

variable "binded" {
  description = "Filter by binding status: true-bound, false-unbound"
  type        = string
  default     = ""
}

variable "apply_status" {
  description = "Filter by application status: 1-applying, 2-issued, 3-review failed, 4-uploaded"
  type        = string
  default     = ""
}

variable "expiry_time" {
  description = "Filter by expiry status: true-expired, false-not expired, inno-about to expire"
  type        = string
  default     = ""
}

# List SCDN certificates
data "edgenext_scdn_certificates" "example" {
  page        = 1
  per_page    = 20
  domain      = var.domain != "" ? var.domain : null
  ca_name     = var.ca_name != "" ? var.ca_name : null
  binded      = var.binded != "" ? var.binded : null
  apply_status = var.apply_status != "" ? var.apply_status : null
  expiry_time = var.expiry_time != "" ? var.expiry_time : null
}

output "total" {
  description = "Total number of certificates"
  value       = data.edgenext_scdn_certificates.example.total
}

output "issuer_list" {
  description = "List of issuers"
  value       = data.edgenext_scdn_certificates.example.issuer_list
}

output "certificates" {
  description = "List of certificates"
  value       = data.edgenext_scdn_certificates.example.certificates
}

output "certificate_count" {
  description = "Number of certificates returned"
  value       = length(data.edgenext_scdn_certificates.example.certificates)
}

