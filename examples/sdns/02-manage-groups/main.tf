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

variable "group_name" {
  description = "The name of the DNS group"
  type        = string
}

variable "domain_name" {
  description = "Domain to bind to group"
  type        = string
}

resource "edgenext_sdns_domain" "example" {
  domain = var.domain_name
}

resource "edgenext_sdns_domain_group" "example" {
  group_name = var.group_name
  remark     = "Managed by Terraform"
  domain_ids = [edgenext_sdns_domain.example.id]
}

output "group_id" {
  value = edgenext_sdns_domain_group.example.id
}

data "edgenext_sdns_domain_groups" "matched" {
  group_name = edgenext_sdns_domain_group.example.group_name
}

output "matched_groups" {
  value = data.edgenext_sdns_domain_groups.matched.groups
}
