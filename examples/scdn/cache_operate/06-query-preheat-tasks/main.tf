# Example 6: Query SCDN Cache Preheat Tasks
# This example demonstrates how to query cache preheat task list

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

variable "page" {
  description = "Page number"
  type        = number
  default     = 1
}

variable "per_page" {
  description = "Items per page"
  type        = number
  default     = 20
}

variable "status" {
  description = "Status filter"
  type        = string
  default     = null
}

variable "url" {
  description = "URL filter"
  type        = string
  default     = null
}

# Query SCDN cache preheat tasks
data "edgenext_scdn_cache_preheat_tasks" "example" {
  page     = var.page
  per_page = var.per_page
  status   = var.status
  url      = var.url
}

output "total" {
  description = "Total number of tasks"
  value       = data.edgenext_scdn_cache_preheat_tasks.example.total
}

output "tasks" {
  description = "Task list"
  value       = data.edgenext_scdn_cache_preheat_tasks.example.tasks
}

