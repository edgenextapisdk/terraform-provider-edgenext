# Example 1: Create SCDN Domain
# This example demonstrates how to create a basic SCDN domain

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

variable "domain_name" {
  description = "SCDN domain name"
  type        = string
}

variable "group_id" {
  description = "Group ID for the domain"
  type        = number
  default     = 1
}

variable "remark" {
  description = "Remark for the domain"
  type        = string
  default     = "Created by Terraform"
}

variable "origin_ip" {
  description = "Origin server IP address"
  type        = string
  default     = "1.2.3.4"
}

variable "origin_port" {
  description = "Origin server port"
  type        = number
  default     = 80
}

# Create SCDN domain
resource "edgenext_scdn_domain" "example" {
  domain         = var.domain_name
  group_id       = var.group_id
  remark         = var.remark
  protect_status = "scdn"

  origins {
    protocol        = 0  # HTTP
    listen_ports    = [80, 443]
    origin_protocol = 0  # HTTP
    load_balance    = 1  # Round Robin
    origin_type     = 0  # IP

    records {
      view     = "primary"
      value    = var.origin_ip
      port     = var.origin_port
      priority = 10
    }
  }
}

output "domain_id" {
  description = "Created domain ID"
  value       = edgenext_scdn_domain.example.id
}

output "domain_name" {
  description = "Domain name"
  value       = edgenext_scdn_domain.example.domain
}

output "cname" {
  description = "CNAME record for the domain"
  value       = edgenext_scdn_domain.example.cname
}

output "access_progress" {
  description = "Domain access progress"
  value       = edgenext_scdn_domain.example.access_progress
}
