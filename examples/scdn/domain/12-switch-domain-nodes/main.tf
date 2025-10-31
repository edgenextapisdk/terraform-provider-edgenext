# Example 12: Switch Domain Nodes
# This example demonstrates how to switch domain node type

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

variable "domain_id" {
  description = "Domain ID to switch nodes"
  type        = number
}

variable "protect_status" {
  description = "Edge node type. Valid values: back_source, scdn, exclusive"
  type        = string
  default     = "scdn"
}

variable "exclusive_resource_id" {
  description = "Exclusive resource ID (required if protect_status is exclusive)"
  type        = number
  default     = null
}

# Switch domain nodes
resource "edgenext_scdn_domain_node_switch" "example" {
  domain_id            = var.domain_id
  protect_status       = var.protect_status
  exclusive_resource_id = var.exclusive_resource_id
}

output "domain_id" {
  description = "Domain ID"
  value       = edgenext_scdn_domain_node_switch.example.domain_id
}

output "protect_status" {
  description = "Current protect status"
  value       = edgenext_scdn_domain_node_switch.example.protect_status
}

