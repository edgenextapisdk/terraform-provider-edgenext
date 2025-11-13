# Example 7: Sort Cache Rules
# This example demonstrates how to sort cache rules by priority

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
  description = "Business ID (template ID for 'tpl' type, domain ID for 'domain' type)"
  type        = number
}

variable "business_type" {
  description = "Business type: 'tpl' (template) or 'domain'"
  type        = string
  default     = "tpl"
}

variable "rule_ids_to_sort" {
  description = "List of rule IDs in the desired sorted order"
  type        = list(number)
}

# Sort cache rules
resource "edgenext_scdn_cache_rules_sort" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  ids           = var.rule_ids_to_sort
}

output "sorted_ids" {
  description = "The actual sorted rule IDs after API call"
  value       = edgenext_scdn_cache_rules_sort.example.sorted_ids
}

