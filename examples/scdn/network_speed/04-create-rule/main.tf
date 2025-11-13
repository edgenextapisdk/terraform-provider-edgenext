# Example 4: Create Network Speed Rule
# This example demonstrates how to create a network speed rule

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

variable "customized_req_headers_rule" {
  description = "Customized request headers rule configuration"
  type = object({
    type    = string
    content = string
    remark  = string
  })
  default = {
    type    = "User-Agent"
    content = "test-content"
    remark  = "test-remark"
  }
}

# Example: Create a customized request headers rule
resource "edgenext_scdn_network_speed_rule" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  config_group  = var.config_group

  customized_req_headers_rule {
    type    = var.customized_req_headers_rule.type
    content = var.customized_req_headers_rule.content
    remark  = var.customized_req_headers_rule.remark
  }
}

output "rule_id" {
  description = "Created rule ID"
  value       = edgenext_scdn_network_speed_rule.example.id
}


