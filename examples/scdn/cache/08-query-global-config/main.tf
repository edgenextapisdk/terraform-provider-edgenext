# Example 8: Query Global Cache Config
# This example demonstrates how to query the default global cache configuration

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

# Query global cache configuration
data "edgenext_scdn_cache_global_config" "example" {
}

output "global_config" {
  description = "Global cache configuration"
  value = {
    id   = data.edgenext_scdn_cache_global_config.example.id
    name = data.edgenext_scdn_cache_global_config.example.name
    conf = data.edgenext_scdn_cache_global_config.example.conf
  }
}

