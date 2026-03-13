# Example: Import SCDN Network Speed Config
# This example demonstrates how to import an existing network speed configuration into Terraform

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

# Managed resource for import
resource "edgenext_scdn_network_speed_config" "example" {
  # business_id   = var.business_id
  # business_type = var.business_type

  # You can now import without specifying these in advance!
  # terraform import edgenext_scdn_network_speed_config.example <business_id>-<business_type>
  # based on the imported state to avoid unwanted changes during the next apply.


  domain_proxy_conf {
    proxy_connect_timeout = 60
    fails_timeout         = 10
    keep_new_src_time     = 10
    max_fails             = 3
    proxy_keepalive       = 1
  }
}

output "imported_resource_id" {
  description = "The ID of the imported resource"
  value       = edgenext_scdn_network_speed_config.example.id
}



