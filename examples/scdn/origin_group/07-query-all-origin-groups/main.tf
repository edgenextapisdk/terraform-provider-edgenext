# Example 7: Query All Origin Groups
# This example demonstrates how to query all origin groups for domain configuration

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

variable "protect_status" {
  description = "Protection status: scdn-shared nodes, exclusive-dedicated nodes"
  type        = string
  default     = "scdn"
}

# Query all origin groups
data "edgenext_scdn_origin_groups_all" "example" {
  protect_status = var.protect_status
}

output "origin_groups" {
  description = "All origin groups"
  value = {
    total        = data.edgenext_scdn_origin_groups_all.example.total
    origin_groups = data.edgenext_scdn_origin_groups_all.example.origin_groups
  }
}

