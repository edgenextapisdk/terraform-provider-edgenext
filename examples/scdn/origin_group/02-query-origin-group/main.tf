# Example 2: Query Origin Group
# This example demonstrates how to query an origin group

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

variable "origin_group_id" {
  description = "Origin group ID"
  type        = number
}

# Query origin group
data "edgenext_scdn_origin_group" "example" {
  origin_group_id = var.origin_group_id
}

output "origin_group" {
  description = "Origin group details"
  value = {
    origin_group_id = data.edgenext_scdn_origin_group.example.origin_group_id
    name            = data.edgenext_scdn_origin_group.example.name
    remark          = data.edgenext_scdn_origin_group.example.remark
    member_id       = data.edgenext_scdn_origin_group.example.member_id
    created_at      = data.edgenext_scdn_origin_group.example.created_at
    origins         = data.edgenext_scdn_origin_group.example.origins
  }
}

