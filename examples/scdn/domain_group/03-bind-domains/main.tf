terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "0.2.1"
    }
  }
}

provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
}

variable "access_key" {
  type        = string
  description = "EdgeNext access key"
  sensitive   = true
}

variable "secret_key" {
  type        = string
  description = "EdgeNext secret key"
  sensitive   = true
}

variable "endpoint" {
  type        = string
  description = "EdgeNext API endpoint"
  default     = "https://api.edgenextscdn.com"
}

variable "group_id" {
  type        = string
  description = "Existing domain group ID to bind domains to"
}

variable "group_name" {
  type        = string
  description = "The name of the existing group (required for Terraform state)"
}

variable "domains_to_bind" {
  type        = list(string)
  description = "List of domains to bind to the group"
  default     = []
}

# ============================================================
# IMPORTANT: Before running terraform apply, you must first
# import the existing group into Terraform state:
#
#   terraform import edgenext_scdn_domain_group.existing <group_id>
#
# Example:
#   terraform import edgenext_scdn_domain_group.existing 111
# ============================================================

# This resource manages an EXISTING domain group
# You MUST import it first before applying
resource "edgenext_scdn_domain_group" "existing" {
  group_name = var.group_name

  # Add domains to bind
  domains = var.domains_to_bind
}

# Query domains in the group after binding
data "edgenext_scdn_domain_group_domains" "group_domains" {
  group_id   = tonumber(var.group_id)
  depends_on = [edgenext_scdn_domain_group.existing]
}

output "group_id" {
  value       = edgenext_scdn_domain_group.existing.id
  description = "The domain group ID"
}

output "group_name" {
  value       = edgenext_scdn_domain_group.existing.group_name
  description = "The domain group name"
}

output "bound_domains" {
  value       = data.edgenext_scdn_domain_group_domains.group_domains.list
  description = "List of domains bound to the group"
}

output "total_domains" {
  value       = data.edgenext_scdn_domain_group_domains.group_domains.total
  description = "Total number of domains in the group"
}
