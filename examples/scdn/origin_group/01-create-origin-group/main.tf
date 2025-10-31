# Example 1: Create Origin Group
# This example demonstrates how to create an origin group

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

variable "name" {
  description = "Origin group name"
  type        = string
  default     = "test-origin-group"
}

variable "remark" {
  description = "Remark"
  type        = string
  default     = "Test origin group"
}

# Create origin group
resource "edgenext_scdn_origin_group" "example" {
  name   = var.name
  remark = var.remark

  origins {
    origin_type = 0 # IP

    records {
      value    = "54.85.23.59"
      port     = 80
      priority = 10
      view     = "primary"
      host     = "example.com"
    }

    protocol_ports {
      protocol     = 0 # http
      listen_ports = [80, 8080]
    }

    origin_protocol = 0 # http
    load_balance    = 1  # round_robin
  }
}

output "origin_group" {
  description = "Origin group details"
  value = {
    origin_group_id = edgenext_scdn_origin_group.example.origin_group_id
    name            = edgenext_scdn_origin_group.example.name
    created_at      = edgenext_scdn_origin_group.example.created_at
  }
}

