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


# Query the Items in the list using a variable ID
variable "user_ip_id" {
  description = "The ID of the User IP List to query"
  type        = number
}

data "edgenext_scdn_user_ip_items" "items" {
  user_ip_id = var.user_ip_id
  page       = 1
  per_page   = 10
}


output "queried_items" {
  value = data.edgenext_scdn_user_ip_items.items.items
}

output "total_count" {
  value = data.edgenext_scdn_user_ip_items.items.total
}
