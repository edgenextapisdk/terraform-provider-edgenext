# Example 2: Update Network Speed Config
# This example demonstrates how to update network speed configuration for a template

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

# Domain proxy configuration variables
variable "domain_proxy_conf" {
  description = "Domain proxy configuration"
  type = object({
    proxy_connect_timeout = number
    fails_timeout         = number
    keep_new_src_time     = number
    max_fails             = number
    proxy_keepalive       = number
  })
  default = {
    proxy_connect_timeout = 30
    fails_timeout         = 10
    keep_new_src_time     = 30
    max_fails             = 30
    proxy_keepalive       = 0
  }
}

# Upstream redirect configuration variables
variable "upstream_redirect" {
  description = "Upstream redirect configuration (301/302 follow)"
  type = object({
    status = string
  })
  default = {
    status = "off"
  }
}

# Customized request headers configuration variables
variable "customized_req_headers" {
  description = "Customized request headers configuration"
  type = object({
    status = string
  })
  default = {
    status = "off"
  }
}

# Response headers configuration variables
variable "resp_headers" {
  description = "Response headers configuration"
  type = object({
    status = string
  })
  default = {
    status = "off"
  }
}

# Upstream URI change configuration variables
variable "upstream_uri_change" {
  description = "Upstream URI change configuration"
  type = object({
    status = string
  })
  default = {
    status = "on"
  }
}

# Source site protection configuration variables
variable "source_site_protect" {
  description = "Source site protection configuration"
  type = object({
    status = string
    num    = number
    second = number
  })
  default = {
    status  = "off"
    num     = 1000
    second  = 10
  }
}

# Range request configuration variables
variable "slice" {
  description = "Range request configuration"
  type = object({
    status = string
  })
  default = {
    status = "off"
  }
}

# HTTPS configuration variables
variable "https_config" {
  description = "HTTPS configuration"
  type = object({
    status                   = string
    http2https               = string
    http2https_port          = number
    http2                    = string
    hsts                     = string
    ocsp_stapling            = string
    min_version              = string
    ciphers_preset           = string
    custom_encrypt_algorithm = list(string)
  })
  default = {
    status                   = "on"
    http2https               = "off"
    http2https_port          = 443
    http2                    = "off"
    hsts                     = "off"
    ocsp_stapling            = "off"
    min_version              = "TLSv1.2"
    ciphers_preset           = "default"
    custom_encrypt_algorithm = [
      "ECDHE-ECDSA",
      "ECDHE-RSA",
      "DHE-RSA",
      "EECDH+CHACHA20",
      "EECDH+AES128",
      "EECDH+AES256",
      "RSA+AES128",
      "RSA+AES256",
      "EECDH+3DES",
      "RSA+3DES"
    ]
  }
}

# Page Gzip configuration variables
variable "page_gzip" {
  description = "Page Gzip configuration"
  type = object({
    status = string
  })
  default = {
    status = "on"
  }
}

# WebP format configuration variables
variable "webp" {
  description = "WebP format configuration"
  type = object({
    status = string
  })
  default = {
    status = "off"
  }
}

# Upload file configuration variables
variable "upload_file" {
  description = "Upload file configuration"
  type = object({
    upload_size      = number
    upload_size_unit = string
  })
  default = {
    upload_size      = 100
    upload_size_unit = "MB"
  }
}

# WebSocket configuration variables
variable "websocket" {
  description = "WebSocket configuration"
  type = object({
    status = string
  })
  default = {
    status = "off"
  }
}

# Mobile jump configuration variables
variable "mobile_jump" {
  description = "Mobile jump configuration"
  type = object({
    status   = string
    jump_url = string
  })
  default = {
    status   = "off"
    jump_url = ""
  }
}

# Custom page configuration variables
variable "custom_page" {
  description = "Custom page configuration"
  type = object({
    status = string
  })
  default = {
    status = "off"
  }
}

# Upstream check configuration variables
variable "upstream_check" {
  description = "Upstream check configuration (health check)"
  type = object({
    status  = string
    fails   = number
    intval  = number
    rise    = number
    timeout = number
    type    = string
    op      = string  # Required when type is "http": "HEAD", "GET", or "AUTO"
    path    = string  # Required when type is "http", must start with "/"
  })
  default = {
    status  = "on"
    fails   = 3
    intval  = 10
    rise    = 2
    timeout = 5
    type    = "http"
    op      = "GET"
    path    = "/health"
  }
}

# Update network speed configuration
resource "edgenext_scdn_network_speed_config" "example" {
  business_id   = var.business_id
  business_type = var.business_type

  # Update domain proxy configuration
  domain_proxy_conf {
    proxy_connect_timeout = var.domain_proxy_conf.proxy_connect_timeout
    fails_timeout         = var.domain_proxy_conf.fails_timeout
    keep_new_src_time     = var.domain_proxy_conf.keep_new_src_time
    max_fails             = var.domain_proxy_conf.max_fails
    proxy_keepalive       = var.domain_proxy_conf.proxy_keepalive
  }

  # Update upstream redirect configuration
  upstream_redirect {
    status = var.upstream_redirect.status
  }

  # Update customized request headers configuration
  customized_req_headers {
    status = var.customized_req_headers.status
  }

  # Update response headers configuration
  resp_headers {
    status = var.resp_headers.status
  }

  # Update upstream URI change configuration
  upstream_uri_change {
    status = var.upstream_uri_change.status
  }

  # Update source site protection configuration
  source_site_protect {
    status  = var.source_site_protect.status
    num     = var.source_site_protect.num
    second  = var.source_site_protect.second
  }

  # Update range request configuration
  slice {
    status = var.slice.status
  }

  # Update HTTPS configuration
  https {
    status                   = var.https_config.status
    http2https               = var.https_config.http2https
    http2https_port          = var.https_config.http2https_port
    http2                    = var.https_config.http2
    hsts                     = var.https_config.hsts
    ocsp_stapling            = var.https_config.ocsp_stapling
    min_version              = var.https_config.min_version
    ciphers_preset           = var.https_config.ciphers_preset
    custom_encrypt_algorithm = var.https_config.custom_encrypt_algorithm
  }

  # Update page Gzip configuration
  page_gzip {
    status = var.page_gzip.status
  }

  # Update WebP format configuration
  webp {
    status = var.webp.status
  }

  # Update upload file configuration
  upload_file {
    upload_size      = var.upload_file.upload_size
    upload_size_unit = var.upload_file.upload_size_unit
  }

  # Update WebSocket configuration
  websocket {
    status = var.websocket.status
  }

  # Update mobile jump configuration
  mobile_jump {
    status   = var.mobile_jump.status
    jump_url = var.mobile_jump.jump_url
  }

  # Update custom page configuration
  custom_page {
    status = var.custom_page.status
  }

  # Update upstream check configuration
  upstream_check {
    status  = var.upstream_check.status
    fails   = var.upstream_check.fails
    intval  = var.upstream_check.intval
    rise    = var.upstream_check.rise
    timeout = var.upstream_check.timeout
    type    = var.upstream_check.type
    op      = var.upstream_check.op
    path    = var.upstream_check.path
  }
}

output "config_id" {
  description = "Network speed configuration resource ID"
  value       = edgenext_scdn_network_speed_config.example.id
}

output "domain_proxy_conf" {
  description = "Updated domain proxy configuration"
  value       = try(edgenext_scdn_network_speed_config.example.domain_proxy_conf[0], null)
}

output "upstream_redirect" {
  description = "Updated upstream redirect configuration"
  value       = try(edgenext_scdn_network_speed_config.example.upstream_redirect[0], null)
}

output "customized_req_headers" {
  description = "Updated customized request headers configuration"
  value       = try(edgenext_scdn_network_speed_config.example.customized_req_headers[0], null)
}

output "resp_headers" {
  description = "Updated response headers configuration"
  value       = try(edgenext_scdn_network_speed_config.example.resp_headers[0], null)
}

output "upstream_uri_change" {
  description = "Updated upstream URI change configuration"
  value       = try(edgenext_scdn_network_speed_config.example.upstream_uri_change[0], null)
}

output "source_site_protect" {
  description = "Updated source site protection configuration"
  value       = try(edgenext_scdn_network_speed_config.example.source_site_protect[0], null)
}

output "slice" {
  description = "Updated range request configuration"
  value       = try(edgenext_scdn_network_speed_config.example.slice[0], null)
}

output "https_config" {
  description = "Updated HTTPS configuration"
  value       = try(edgenext_scdn_network_speed_config.example.https[0], null)
}

output "page_gzip" {
  description = "Updated page Gzip configuration"
  value       = try(edgenext_scdn_network_speed_config.example.page_gzip[0], null)
}

output "webp" {
  description = "Updated WebP format configuration"
  value       = try(edgenext_scdn_network_speed_config.example.webp[0], null)
}

output "upload_file" {
  description = "Updated upload file configuration"
  value       = try(edgenext_scdn_network_speed_config.example.upload_file[0], null)
}

output "websocket" {
  description = "Updated WebSocket configuration"
  value       = try(edgenext_scdn_network_speed_config.example.websocket[0], null)
}

output "mobile_jump" {
  description = "Updated mobile jump configuration"
  value       = try(edgenext_scdn_network_speed_config.example.mobile_jump[0], null)
}

output "custom_page" {
  description = "Updated custom page configuration"
  value       = try(edgenext_scdn_network_speed_config.example.custom_page[0], null)
}

output "upstream_check" {
  description = "Updated upstream check configuration"
  value       = try(edgenext_scdn_network_speed_config.example.upstream_check[0], null)
}

