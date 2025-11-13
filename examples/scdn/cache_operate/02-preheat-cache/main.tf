# Example 2: Create SCDN Cache Preheat Task
# This example demonstrates how to create a cache preheat task

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

variable "preheat_url" {
  description = "Preheat URLs"
  type        = list(string)
}

variable "group_id" {
  description = "Group ID, can refresh cache by group"
  type        = number
  default     = null
}

variable "protocol" {
  description = "Protocol: http/https; only valid when refreshing by group"
  type        = string
  default     = null
}

variable "port" {
  description = "Website port, only needed for special ports; only valid when refreshing by group"
  type        = string
  default     = null
}

# Create SCDN cache preheat task
resource "edgenext_scdn_cache_preheat_task" "example" {
  preheat_url = var.preheat_url
  group_id    = var.group_id
  protocol    = var.protocol
  port        = var.port
}

output "task_id" {
  description = "Created cache preheat task ID"
  value       = edgenext_scdn_cache_preheat_task.example.id
}

output "error_url" {
  description = "List of URLs with preheat errors"
  value       = edgenext_scdn_cache_preheat_task.example.error_url
}

