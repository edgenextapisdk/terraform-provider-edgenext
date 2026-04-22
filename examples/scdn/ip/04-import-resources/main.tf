# Example 4: Import Resources
# This example demonstrates how to import existing resources into Terraform state.

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

# 1. User IP List Import
# Create a resource block that matches the configuration of the existing resource.
resource "edgenext_scdn_user_ip" "imported_list" {
}

# 2. User IP Item Import
# Create a resource block for the item.
resource "edgenext_scdn_user_ip_item" "imported_item" {
}
