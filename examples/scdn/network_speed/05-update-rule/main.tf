# Example 5: Update Network Speed Rule
# This example demonstrates how to update an existing network speed rule

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

variable "business_id" {
  description = "Business ID (template ID for 'tpl' type)"
  type        = number
}

variable "business_type" {
  description = "Business type: 'tpl' (template) or 'global'"
  type        = string
  default     = "tpl"
}

variable "config_group" {
  description = "Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'"
  type        = string
}

variable "rule_id" {
  description = "Rule ID to update (import first: terraform import edgenext_scdn_network_speed_rule.example <business_id>-<business_type>-<config_group>-<rule_id>)"
  type        = number
}

variable "customized_req_headers_rule" {
  description = "Customized request headers rule configuration"
  type = object({
    type    = string
    content = string
    remark  = string
  })
  default = {
    type    = "User-Agent"
    content = "updated-content"
    remark  = "updated-remark"
  }
}

# Update network speed rule
# Note: You can either provide rule_id to update directly, or import the existing rule first
resource "edgenext_scdn_network_speed_rule" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  config_group  = var.config_group
  rule_id       = var.rule_id  # If provided, this will update the existing rule instead of creating a new one

  # Example: Update a customized request headers rule
  customized_req_headers_rule {
    type    = var.customized_req_headers_rule.type
    content = var.customized_req_headers_rule.content
    remark  = var.customized_req_headers_rule.remark
  }
}

output "rule_id" {
  description = "Updated rule ID"
  value       = edgenext_scdn_network_speed_rule.example.id
}

