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
  description = "The ID of the domain group to import"
}

# This resource will be imported from an existing domain group
# Run: terraform import edgenext_scdn_domain_group.imported <group_id>
resource "edgenext_scdn_domain_group" "imported" {
  # These values will be populated from the import
  # You can modify them after import
  group_name = "imported-group-name"
  remark     = "Imported and managed by Terraform"

  # Optionally manage domains
  # domains = ["example.com"]
}

# Query the imported group's domains
data "edgenext_scdn_domain_group_domains" "imported_domains" {
  group_id   = tonumber(var.group_id)
  depends_on = [edgenext_scdn_domain_group.imported]
}

output "group_info" {
  value = {
    id         = edgenext_scdn_domain_group.imported.id
    name       = edgenext_scdn_domain_group.imported.group_name
    remark     = edgenext_scdn_domain_group.imported.remark
    created_at = edgenext_scdn_domain_group.imported.created_at
    updated_at = edgenext_scdn_domain_group.imported.updated_at
  }
  description = "Imported group information"
}

output "group_domains" {
  value       = data.edgenext_scdn_domain_group_domains.imported_domains.list
  description = "Domains in the imported group"
}
