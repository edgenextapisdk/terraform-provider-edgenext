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
  description = "Domain to create records for"
  type        = string
}

variable "record_name" {
  description = "Record name"
  type        = string
}

variable "record_value" {
  description = "Record value"
  type        = string
}

resource "edgenext_sdns_domain" "example" {
  domain = var.domain_name
}

resource "edgenext_sdns_record" "example" {
  domain_id = tonumber(edgenext_sdns_domain.example.id)
  name      = var.record_name
  type      = "A"
  view      = "any"
  value     = var.record_value
  ttl       = 600
  remark    = "Managed by ...."
}

output "record_id" {
  value = edgenext_sdns_record.example.id
}

data "edgenext_sdns_records" "matched" {
  domain_id = tonumber(edgenext_sdns_domain.example.id)
}

output "all_records" {
  value = data.edgenext_sdns_records.matched.records
}
