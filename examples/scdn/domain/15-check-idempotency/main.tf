# Example 15: Check Idempotency for SCDN Domain
# This example verifies that the edgenext_scdn_domain resource is idempotent.
# After the first successful apply, subsequent runs should show 'No changes'.

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
  description = "SCDN domain name for idempotency test"
  type        = string
  default     = "idempotency-test-example.com"
}

variable "exclusive_resource_id" {
  description = "ID of the exclusive resource package"
  type        = number
  default     = 1509
}

resource "edgenext_scdn_domain" "example" {
  domain                = var.domain_name
  exclusive_resource_id = var.exclusive_resource_id
  protect_status        = "exclusive"

  # Testing idempotency with multiple origins
  # The fix ensures these are sorted stably in the state
  origins {
    listen_ports    = [80]
    load_balance    = 1
    origin_protocol = 1
    origin_type     = 1
    protocol        = 0

    records {
      port     = 443
      priority = 1
      value    = "origin1-http.example.com"
      view     = "primary"
    }
  }

  origins {
    listen_ports    = [443]
    load_balance    = 1
    origin_protocol = 1
    origin_type     = 1
    protocol        = 1

    records {
      port     = 443
      priority = 1
      value    = "origin2-https.example.com"
      view     = "primary"
    }
  }
}

output "domain_id" {
  value = edgenext_scdn_domain.example.id
}

output "exclusive_resource_id" {
  value = edgenext_scdn_domain.example.exclusive_resource_id
}
