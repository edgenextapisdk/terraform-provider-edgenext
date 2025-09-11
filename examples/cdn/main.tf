# EdgeNext CDN Management Example
# Terraform configuration for EdgeNext CDN API

# =============================================================================
# Provider Configuration
# =============================================================================

terraform {
  required_providers {
    edgenext = {
      source = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}

provider "edgenext" {
  api_key  = "your-edgenext-api-key-here"
  secret   = "your-edgenext-api-secret-here"
  endpoint = "your-edgenext-api-endpoint-here"
  timeout  = 300
}

# =============================================================================
# CDN Domain and Configuration Resource Examples
# =============================================================================

# 1. Import domain
# resource "edgenext_cdn_domain" "import_domain" {
#   // Empty fields for import
# }

# 2. Create web acceleration domain
resource "edgenext_cdn_domain_config" "web_domain" {
  domain = "example.com"
  area   = "mainland_china"
  type   = "page"

  config {
    origin {
      default_master = "1.2.3.4"
      origin_mode    = "http"
      ori_https = "no"
      port = "80"
    }

    cache_rule {
      type     = 1
      pattern  = "jpg,png,gif"
      time     = 3600
      timeunit = "s"
    }

    # Test single Map configuration item (MaxItems: 1)
    referer {
      type = 2 # 2 is whitelist, 1 is blacklist
      list = ["trusted.com", "another.com"]
    }

    # Test IP whitelist configuration (update list)
    ip_white_list {
      list = ["192.168.1.0/24", "10.0.0.0/8", "172.16.0.0/12"]
    }

    # origin_host {
    #   host = "www.example.com"
    # }
    # Test HTTPS configuration (temporarily removed)
    # https {
    #   cert_id = 12345
    #   force_https = "yes"
    # }
  }
}

# 3. Query domain and configuration
data "edgenext_cdn_domain_config" "web_config_info" {
  domain = edgenext_cdn_domain_config.web_domain.domain

  # Specify configuration items to query
  config_item = [
    "cache_rule",
    "referer", 
    "ip_white_list",
    "add_response_head"
  ]
  
  depends_on = [edgenext_cdn_domain_config.web_domain]
}

# =============================================================================
# CDN Push Task Examples
# =============================================================================

# 1. Push specific URLs
resource "edgenext_cdn_push" "url_refresh" {
  type = "url"
  urls = [
    "https://${edgenext_cdn_domain_config.web_domain.domain}/config.mp4",
    "https://${edgenext_cdn_domain_config.web_domain.domain}/data/abc.jpg"
  ]

  depends_on = [edgenext_cdn_domain_config.web_domain]
}

# 2. Query push task status
data "edgenext_cdn_push" "static_push_status" {
  task_id = edgenext_cdn_push.url_refresh.task_id

  depends_on = [edgenext_cdn_push.url_refresh]
}

# 3. Query push tasks within a time period
data "edgenext_cdn_pushes" "recent_pushes" {
  start_time = "2025-01-01"
  end_time   = "2025-12-31"

  depends_on = [edgenext_cdn_push.url_refresh]
}
