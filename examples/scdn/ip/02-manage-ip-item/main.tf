# Example 2: Manage SCDN User IP Item
# This example demonstrates how to create an IP List and add an Item to it

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

resource "edgenext_scdn_user_ip" "list" {
  name   = "terraform-ip-list-with-items"
  remark = "List for items demo"
}

resource "edgenext_scdn_user_ip_item" "item1" {
  user_ip_id = edgenext_scdn_user_ip.list.id
  ip         = "1.1.1.1"
  remark     = "First IP"
}

resource "edgenext_scdn_user_ip_item" "item2" {
  user_ip_id = edgenext_scdn_user_ip.list.id
  ip         = "2.2.2.2"
  remark     = "Second IP"
}

output "list_id" {
  value = edgenext_scdn_user_ip.list.id
}

output "item1_id" {
  value = edgenext_scdn_user_ip_item.item1.id
}
