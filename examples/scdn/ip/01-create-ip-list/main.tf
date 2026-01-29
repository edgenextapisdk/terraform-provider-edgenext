# Example 1: Create SCDN User IP List
# This example demonstrates how to create a User IP List

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
  description = "IP List Name"
  type        = string
  default     = "terraform-example-ip-list"
}

variable "remark" {
  description = "IP List Remark"
  type        = string
  default     = "Created via Terraform"
}

resource "edgenext_scdn_user_ip" "example" {
  name   = var.name
  remark = var.remark
}

output "ip_list_id" {
  value = edgenext_scdn_user_ip.example.id
}

output "ip_list_name" {
  value = edgenext_scdn_user_ip.example.name
}
