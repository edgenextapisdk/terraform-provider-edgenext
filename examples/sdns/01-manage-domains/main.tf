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
  description = "EdgeNext API Endpoint"
  type        = string
  default     = "https://api.edgenextscdn.com"
}

variable "domain_name" {
  description = "The domain name to manage"
  type        = string
}

resource "edgenext_dns_domain" "example" {
  domain = var.domain_name
}

output "domain_id" {
  value = edgenext_dns_domain.example.id
}

output "domain_status" {
  value = edgenext_dns_domain.example.status
}

data "edgenext_dns_domain" "matched" {
  domain = edgenext_dns_domain.example.domain
}

output "matched_domains" {
  value = data.edgenext_dns_domain.matched.domains
}
