# Example: Import SCDN Network Speed Config
# This example demonstrates how to import an existing network speed configuration into Terraform

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

resource "edgenext_scdn_cache_rule" "example" {
  name = "no-cache-rule1"
  expr = "(http.request.uri.path eq \"/api\")"

  conf {
    nocache = true
  }
}


resource "edgenext_scdn_network_speed_rule" "example" {
  customized_req_headers_rule {
    type    = "add"
    content = "X-Test-Header: ads"
    remark  = "testing drift fix"
  }
}
