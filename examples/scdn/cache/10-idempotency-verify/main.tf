# Example 10: Idempotency Verification for Cache Rule
#
# This example verifies that creating a cache rule with minimal configuration
# (only conf.nocache) produces a stable plan on the second apply.
#
# Previously, the second `terraform apply` would attempt to remove the
# server-returned default `cache_rule` block (action, cachetime, etc.),
# causing a spurious in-place update. This has been fixed by marking
# optional sub-blocks as Computed so Terraform accepts API defaults.
#
# Usage:
#   1. Copy terraform.tfvars.example to terraform.tfvars and fill in values
#   2. terraform init
#   3. terraform apply          # First apply: creates the resource
#   4. terraform plan           # Second plan: should show "No changes"
#   5. terraform apply          # Second apply: should be a no-op
#   6. terraform destroy        # Cleanup

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

# ============================================================================
# Minimal cache rule: only sets conf.nocache
# The API will return default values for cache_rule, browser_cache_rule, etc.
# With the fix, these server-generated defaults should NOT cause drift.
# ============================================================================
resource "edgenext_scdn_cache_rule" "minimal" {
  business_id   = var.business_id
  business_type = var.business_type
  name          = "idempotency-test-minimal"
  remark        = "Idempotency verification - minimal"

  conf {
    nocache = false
    cache_share {
      scheme = "http"
    }
  }
}

output "minimal_rule_id" {
  description = "Minimal cache rule ID"
  value       = edgenext_scdn_cache_rule.minimal.id
}

output "minimal_conf" {
  description = "Minimal rule conf (shows server-populated defaults)"
  value       = edgenext_scdn_cache_rule.minimal.conf
}
