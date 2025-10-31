# Example 8: Update Domain Base Settings
# This example demonstrates how to update domain base settings

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

variable "domain_id" {
  description = "Domain ID to update base settings"
  type        = number
}

variable "enable_proxy_host" {
  description = "Whether to configure proxy host"
  type        = bool
  default     = false
}

variable "proxy_host" {
  description = "Proxy host value"
  type        = string
  default     = ""
}

variable "proxy_host_type" {
  description = "Proxy host type"
  type        = string
  default     = ""
}

variable "enable_proxy_sni" {
  description = "Whether to configure proxy SNI"
  type        = bool
  default     = false
}

variable "proxy_sni" {
  description = "Proxy SNI value"
  type        = string
  default     = ""
}

variable "proxy_sni_status" {
  description = "Proxy SNI status (on/off)"
  type        = string
  default     = "off"
}

variable "enable_redirect" {
  description = "Whether to configure domain redirect"
  type        = bool
  default     = false
}

variable "redirect_status" {
  description = "Redirect status (on/off)"
  type        = string
  default     = "off"
}

variable "redirect_jump_to" {
  description = "Redirect target URL"
  type        = string
  default     = ""
}

variable "redirect_jump_type" {
  description = "Redirect jump type"
  type        = string
  default     = ""
}

# Update domain base settings
resource "edgenext_scdn_domain_base_settings" "example" {
  domain_id = var.domain_id

  # Proxy host configuration (optional)
  dynamic "proxy_host" {
    for_each = var.enable_proxy_host ? [1] : []
    content {
      proxy_host     = var.proxy_host
      proxy_host_type = var.proxy_host_type
    }
  }

  # Proxy SNI configuration (optional)
  dynamic "proxy_sni" {
    for_each = var.enable_proxy_sni ? [1] : []
    content {
      proxy_sni = var.proxy_sni
      status    = var.proxy_sni_status
    }
  }

  # Domain redirect configuration (optional)
  dynamic "domain_redirect" {
    for_each = var.enable_redirect ? [1] : []
    content {
      status   = var.redirect_status
      jump_to  = var.redirect_jump_to
      jump_type = var.redirect_jump_type
    }
  }
}

output "domain_id" {
  description = "Domain ID"
  value       = edgenext_scdn_domain_base_settings.example.domain_id
}

