# Example 3: List Origin Groups
# This example demonstrates how to list origin groups

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

variable "name" {
  description = "Origin group name filter (optional)"
  type        = string
  default     = ""
}

variable "page" {
  description = "Page number"
  type        = number
  default     = 1
}

variable "page_size" {
  description = "Page size"
  type        = number
  default     = 20
}

# List origin groups
data "edgenext_scdn_origin_groups" "example" {
  page      = var.page
  page_size = var.page_size
  name      = var.name
}

output "origin_groups" {
  description = "Origin groups list"
  value = {
    total        = data.edgenext_scdn_origin_groups.example.total
    origin_groups = data.edgenext_scdn_origin_groups.example.origin_groups
  }
}

