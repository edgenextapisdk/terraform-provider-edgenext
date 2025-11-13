# Example 1: Create SCDN Log Download Task
# This example demonstrates how to create a log download task

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

variable "task_name" {
  description = "Task name"
  type        = string
}

variable "is_use_template" {
  description = "Whether to use template: 0-no, 1-yes"
  type        = number
  default     = 0
}

variable "template_id" {
  description = "Template ID (required when is_use_template is 1)"
  type        = number
  default     = null
}

variable "data_source" {
  description = "Data source: ng, cc, waf"
  type        = string
}

variable "download_fields" {
  description = "Download fields"
  type        = list(string)
}

variable "search_terms" {
  description = "Search conditions"
  type = list(object({
    key   = string
    value = string
  }))
  default = []
}

variable "file_type" {
  description = "File type: xls, csv, json"
  type        = string
}

variable "start_time" {
  description = "Start time (format: YYYY-MM-DD HH:MM:SS)"
  type        = string
}

variable "end_time" {
  description = "End time (format: YYYY-MM-DD HH:MM:SS)"
  type        = string
}

variable "lang" {
  description = "Language: zh_CN, en_US, default: zh_CN"
  type        = string
  default     = "zh_CN"
}

# Create SCDN log download task
resource "edgenext_scdn_log_download_task" "example" {
  task_name       = var.task_name
  is_use_template = var.is_use_template
  template_id     = var.template_id
  data_source     = var.data_source
  download_fields = var.download_fields
  
  dynamic "search_terms" {
    for_each = var.search_terms
    content {
      key   = search_terms.value.key
      value = search_terms.value.value
    }
  }
  
  file_type  = var.file_type
  start_time = var.start_time
  end_time   = var.end_time
  lang       = var.lang
}

output "task_id" {
  description = "Created log download task ID"
  value       = edgenext_scdn_log_download_task.example.task_id
}

output "status" {
  description = "Task status"
  value       = edgenext_scdn_log_download_task.example.status
}

output "download_url" {
  description = "Download URL (available when task is completed)"
  value       = edgenext_scdn_log_download_task.example.download_url
}

