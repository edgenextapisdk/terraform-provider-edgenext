# Example: Import Existing SCDN Resources
# This example demonstrates how to import existing EdgeNext SCDN resources into Terraform

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

# Resource IDs to import (set these in terraform.tfvars)
variable "domain_id_to_import" {
  description = "Domain ID to import"
  type        = number
  default     = null
}

variable "origin_group_id_to_import" {
  description = "Origin Group ID to import"
  type        = number
  default     = null
}

variable "certificate_id_to_import" {
  description = "Certificate ID to import"
  type        = number
  default     = null
}

variable "origin_id_to_import" {
  description = "Origin ID to import"
  type        = number
  default     = null
}

variable "domain_id_for_cert_binding" {
  description = "Domain ID for certificate binding"
  type        = number
  default     = null
}

variable "certificate_id_for_binding" {
  description = "Certificate ID for certificate binding"
  type        = number
  default     = null
}

# ============================================================================
# Domain Resource
# ============================================================================
# To import: terraform import edgenext_scdn_domain.example <domain_id>
# Example: terraform import edgenext_scdn_domain.example 102008
# Note: The origins block below is a placeholder. After import, run
#       'terraform show' to see the actual configuration and update this file.
resource "edgenext_scdn_domain" "example" {
  domain = "example.com" # This will be updated after import

  # Placeholder origins block (required for validation)
  # This will be replaced with actual values after import
  origins {
    protocol        = 0
    listen_ports    = [80]
    origin_protocol = 0
    load_balance    = 1
    origin_type     = 0

    records {
      view     = "primary"
      value    = "1.1.1.1"
      port     = 80
      priority = 10
    }
  }
}

# ============================================================================
# Origin Group Resource
# ============================================================================
# To import: terraform import edgenext_scdn_origin_group.example <origin_group_id>
# Example: terraform import edgenext_scdn_origin_group.example 85
# Note: The origins block below is a placeholder. After import, run
#       'terraform show' to see the actual configuration and update this file.
resource "edgenext_scdn_origin_group" "example" {
  name = "imported-origin-group" # This will be updated after import

  # Placeholder origins block (required for validation)
  # This will be replaced with actual values after import
  origins {
    origin_type = 0

    records {
      value    = "1.1.1.1"
      port     = 80
      priority = 10
      view     = "primary"
    }

    protocol_ports {
      protocol     = 0
      listen_ports = [80]
    }

    origin_protocol = 0
    load_balance    = 1
  }
}

# ============================================================================
# Certificate Resource
# ============================================================================
# To import: terraform import edgenext_scdn_certificate.example <certificate_id>
# Example: terraform import edgenext_scdn_certificate.example 456
# Note: The ca_name below is a placeholder. After import, run
#       'terraform show' to see the actual configuration and update this file.
resource "edgenext_scdn_certificate" "example" {
  ca_name = "placeholder-cert-name" # This will be updated after import
  # Other fields will be populated automatically after import
}

# ============================================================================
# Certificate Binding Resource
# ============================================================================
# To import: terraform import edgenext_scdn_cert_binding.example "<domain_id>:<certificate_id>"
# Example: terraform import edgenext_scdn_cert_binding.example "123:456"
# Note: Only created if both domain_id_for_cert_binding and certificate_id_for_binding are provided
resource "edgenext_scdn_cert_binding" "example" {
  count = var.domain_id_for_cert_binding != null && var.certificate_id_for_binding != null ? 1 : 0

  domain_id = var.domain_id_for_cert_binding
  ca_id     = var.certificate_id_for_binding
}

# ============================================================================
# Origin Resource
# ============================================================================
# To import: terraform import edgenext_scdn_origin.example <origin_id>
# Note: The fields below are placeholders. After import, run
#       'terraform show' to see the actual configuration and update this file.
# Note: Only created if both origin_id_to_import and domain_id_to_import are provided
resource "edgenext_scdn_origin" "example" {
  count = var.origin_id_to_import != null && var.domain_id_to_import != null ? 1 : 0

  domain_id = var.domain_id_to_import # Set this after import if needed

  # Placeholder values (required for validation)
  # These will be replaced with actual values after import
  protocol        = 0
  listen_ports    = [80]
  origin_protocol = 0
  load_balance    = 1
  origin_type     = 0

  records {
    view     = "primary"
    value    = "1.1.1.1"
    port     = 80
    priority = 10
  }
}

# ============================================================================
# Outputs
# ============================================================================
output "imported_domain_id" {
  description = "Imported domain ID"
  value       = edgenext_scdn_domain.example.id
}

output "imported_domain_name" {
  description = "Imported domain name"
  value       = edgenext_scdn_domain.example.domain
}

output "imported_origin_group_id" {
  description = "Imported origin group ID"
  value       = edgenext_scdn_origin_group.example.origin_group_id
}

output "imported_origin_group_name" {
  description = "Imported origin group name"
  value       = edgenext_scdn_origin_group.example.name
}

output "imported_certificate_id" {
  description = "Imported certificate ID"
  value       = edgenext_scdn_certificate.example.id
}

output "imported_cert_binding" {
  description = "Imported certificate binding (if imported)"
  value = length(edgenext_scdn_cert_binding.example) > 0 ? {
    domain_id = edgenext_scdn_cert_binding.example[0].domain_id
    ca_id     = edgenext_scdn_cert_binding.example[0].ca_id
  } : null
}

output "imported_origin" {
  description = "Imported origin (if imported)"
  value = length(edgenext_scdn_origin.example) > 0 ? {
    id        = edgenext_scdn_origin.example[0].id
    domain_id = edgenext_scdn_origin.example[0].domain_id
  } : null
}

output "import_instructions" {
  description = "Instructions for importing resources"
  value       = <<-EOT
    To import resources, use the following commands:
    
    1. Domain:
       terraform import edgenext_scdn_domain.example ${var.domain_id_to_import != null ? var.domain_id_to_import : "<domain_id>"}
    
    2. Origin Group:
       terraform import edgenext_scdn_origin_group.example ${var.origin_group_id_to_import != null ? var.origin_group_id_to_import : "<origin_group_id>"}
    
    3. Certificate:
       terraform import edgenext_scdn_certificate.example ${var.certificate_id_to_import != null ? var.certificate_id_to_import : "<certificate_id>"}
    
    4. Certificate Binding (if both domain_id and certificate_id are set):
       terraform import edgenext_scdn_cert_binding.example[0] "${var.domain_id_for_cert_binding != null ? var.domain_id_for_cert_binding : "<domain_id>"}:${var.certificate_id_for_binding != null ? var.certificate_id_for_binding : "<certificate_id>"}"
    
    5. Origin (if both origin_id and domain_id are set):
       terraform import edgenext_scdn_origin.example[0] ${var.origin_id_to_import != null ? var.origin_id_to_import : "<origin_id>"}
    
    After importing, run:
    - terraform show (to view imported state)
    - terraform plan (to verify configuration consistency)
  EOT
}

