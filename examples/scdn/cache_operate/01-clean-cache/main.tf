# Example 1: Create SCDN Cache Clean Task
# This example demonstrates how to create a cache clean task

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

variable "wholesite" {
  description = "Whole site domains to clean"
  type        = list(string)
  default     = []
}

variable "specialurl" {
  description = "Special URLs to clean"
  type        = list(string)
  default     = []
}

variable "specialdir" {
  description = "Special directories to clean"
  type        = list(string)
  default     = []
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

# Create SCDN cache clean task
resource "edgenext_scdn_cache_clean_task" "example" {
  wholesite  = var.wholesite
  specialurl = var.specialurl
  specialdir = var.specialdir
  group_id   = var.group_id
  protocol   = var.protocol
  port       = var.port
}

output "task_id" {
  description = "Created cache clean task ID"
  value       = edgenext_scdn_cache_clean_task.example.id
}

