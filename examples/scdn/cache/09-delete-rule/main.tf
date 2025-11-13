# Example 9: Delete Cache Rule
# This example demonstrates how to DELETE (not disable) a cache rule
# 
# IMPORTANT: This operation PERMANENTLY DELETES the cache rule from the system.
# This is different from disabling a rule (which uses edgenext_scdn_cache_rule_status resource).
# Once deleted, the rule cannot be recovered.

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

variable "rule_id" {
  description = "Cache rule ID to delete (permanently)"
  type        = number
}

# Delete cache rule permanently
# Note: The fields below are required by the schema but will be ignored during deletion.
# The delete operation only uses rule_id, business_id, and business_type.
resource "edgenext_scdn_cache_rule" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  rule_id       = var.rule_id
  
  # These fields are required by schema but not used for deletion
  name = ""  
  expr = ""  

  conf {
    nocache = false
    cache_share {
      scheme = "http"
    }
  }
}

# To delete this resource, run: terraform destroy
# This will permanently delete the cache rule from the EdgeNext SCDN system.

