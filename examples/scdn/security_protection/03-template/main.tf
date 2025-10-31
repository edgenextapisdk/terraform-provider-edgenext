# Example 3: Security Protection Template
# This example demonstrates how to create a security protection template

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

variable "template_name" {
  description = "Template name"
  type        = string
  default     = "test-security-template"
}

variable "remark" {
  description = "Template remark"
  type        = string
  default     = "Test security protection template"
}

# Create security protection template
resource "edgenext_scdn_security_protection_template" "example" {
  name   = var.template_name
  remark = var.remark
}

output "template" {
  description = "Security protection template"
  value = {
    business_id = edgenext_scdn_security_protection_template.example.business_id
    name        = edgenext_scdn_security_protection_template.example.name
    type        = edgenext_scdn_security_protection_template.example.type
  }
}

