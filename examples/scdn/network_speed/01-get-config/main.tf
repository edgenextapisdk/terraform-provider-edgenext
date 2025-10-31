# Example 1: Get Network Speed Config
# This example demonstrates how to query network speed configuration for a template

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

variable "config_groups" {
  description = "Configuration groups to retrieve (optional, defaults to all)"
  type        = list(string)
  default     = []
}

# Query network speed configuration
data "edgenext_scdn_network_speed_config" "example" {
  business_id   = var.business_id
  business_type = var.business_type
  config_groups = length(var.config_groups) > 0 ? var.config_groups : null
}

output "domain_proxy_conf" {
  description = "Domain proxy configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.domain_proxy_conf[0], null)
}

output "upstream_redirect" {
  description = "Upstream redirect configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.upstream_redirect[0], null)
}

output "customized_req_headers" {
  description = "Customized request headers configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.customized_req_headers[0], null)
}

output "resp_headers" {
  description = "Response headers configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.resp_headers[0], null)
}

output "upstream_uri_change" {
  description = "Upstream URI change configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.upstream_uri_change[0], null)
}

output "source_site_protect" {
  description = "Source site protection configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.source_site_protect[0], null)
}

output "slice" {
  description = "Range request configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.slice[0], null)
}

output "https_config" {
  description = "HTTPS configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.https[0], null)
}

output "page_gzip" {
  description = "Page Gzip configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.page_gzip[0], null)
}

output "webp" {
  description = "WebP format configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.webp[0], null)
}

output "upload_file" {
  description = "Upload file configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.upload_file[0], null)
}

output "websocket" {
  description = "WebSocket configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.websocket[0], null)
}

output "mobile_jump" {
  description = "Mobile jump configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.mobile_jump[0], null)
}

output "custom_page" {
  description = "Custom page configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.custom_page[0], null)
}

output "upstream_check" {
  description = "Upstream check configuration"
  value       = try(data.edgenext_scdn_network_speed_config.example.upstream_check[0], null)
}

# Output all configurations as a single object
output "all_configs" {
  description = "All network speed configurations"
  value = {
    domain_proxy_conf      = try(data.edgenext_scdn_network_speed_config.example.domain_proxy_conf[0], null)
    upstream_redirect      = try(data.edgenext_scdn_network_speed_config.example.upstream_redirect[0], null)
    customized_req_headers  = try(data.edgenext_scdn_network_speed_config.example.customized_req_headers[0], null)
    resp_headers            = try(data.edgenext_scdn_network_speed_config.example.resp_headers[0], null)
    upstream_uri_change      = try(data.edgenext_scdn_network_speed_config.example.upstream_uri_change[0], null)
    source_site_protect     = try(data.edgenext_scdn_network_speed_config.example.source_site_protect[0], null)
    slice                   = try(data.edgenext_scdn_network_speed_config.example.slice[0], null)
    https                   = try(data.edgenext_scdn_network_speed_config.example.https[0], null)
    page_gzip               = try(data.edgenext_scdn_network_speed_config.example.page_gzip[0], null)
    webp                    = try(data.edgenext_scdn_network_speed_config.example.webp[0], null)
    upload_file             = try(data.edgenext_scdn_network_speed_config.example.upload_file[0], null)
    websocket               = try(data.edgenext_scdn_network_speed_config.example.websocket[0], null)
    mobile_jump             = try(data.edgenext_scdn_network_speed_config.example.mobile_jump[0], null)
    custom_page             = try(data.edgenext_scdn_network_speed_config.example.custom_page[0], null)
    upstream_check          = try(data.edgenext_scdn_network_speed_config.example.upstream_check[0], null)
  }
}


