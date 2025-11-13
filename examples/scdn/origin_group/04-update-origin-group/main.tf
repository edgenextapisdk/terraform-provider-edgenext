# Example 4: Update Origin Group
# This example demonstrates how to update an origin group

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

variable "origin_group_id" {
  description = "Origin group ID"
  type        = number
}

variable "name" {
  description = "Origin group name"
  type        = string
  default     = "updated-origin-group"
}

variable "remark" {
  description = "Remark"
  type        = string
  default     = "Updated remark"
}

# Update origin group
resource "edgenext_scdn_origin_group" "example" {
  origin_group_id = var.origin_group_id
  name            = var.name
  remark          = var.remark

  origins {
    origin_type = 0 # IP

    records {
      value    = "2.2.2.2"
      port     = 80
      priority = 20
      view     = "primary"
    }

    protocol_ports {
      protocol     = 0 # http
      listen_ports = [80]
    }

    origin_protocol = 0 # http
    load_balance    = 1  # round_robin
  }
}

output "origin_group" {
  description = "Updated origin group details"
  value = {
    origin_group_id = edgenext_scdn_origin_group.example.origin_group_id
    name            = edgenext_scdn_origin_group.example.name
    updated_at      = edgenext_scdn_origin_group.example.updated_at
  }
}

