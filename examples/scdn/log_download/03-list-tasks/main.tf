# Example 3: List SCDN Log Download Tasks
# This example demonstrates how to query log download tasks

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

variable "status" {
  description = "Task status: 0-not started, 1-in progress, 2-completed, 3-failed, 4-cancelled"
  type        = number
  default     = null
}

variable "task_name" {
  description = "Task name"
  type        = string
  default     = null
}

variable "file_type" {
  description = "File type: xls, csv, json"
  type        = string
  default     = null
}

variable "data_source" {
  description = "Data source: ng, cc, waf"
  type        = string
  default     = null
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

# Query SCDN log download tasks
data "edgenext_scdn_log_download_tasks" "example" {
  status     = var.status
  task_name  = var.task_name
  file_type  = var.file_type
  data_source = var.data_source
  page       = var.page
  per_page   = var.per_page
}

output "total" {
  description = "Total number of tasks"
  value       = data.edgenext_scdn_log_download_tasks.example.total
}

output "tasks" {
  description = "Task list"
  value       = data.edgenext_scdn_log_download_tasks.example.tasks
}

