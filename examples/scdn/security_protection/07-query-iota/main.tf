# Example 7: Query Security Protection Iota
# This example demonstrates how to query security protection enum values

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

# Query security protection iota
data "edgenext_scdn_security_protection_iota" "example" {
}

output "iota" {
  description = "Security protection enum values"
  value       = data.edgenext_scdn_security_protection_iota.example.iota
}

