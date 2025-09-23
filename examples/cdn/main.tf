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
# resource "edgenext_cdn_domain" "web_domain" {
#   domain = "example2.com"
#   area   = "mainland_china"
#   type   = "page"

#   config {
#     origin {
#       default_master = "1.2.3.5"
#       default_slave = "2.3.4.5"
#       origin_mode    = "custom"
#       ori_https = "no"
#       port = 800
#     }

#     origin_host {
#       host = "www.example.com"
#     }

#     # cache_rule {
#     #   type     = 1
#     #   pattern  = "jpg,png,gif"
#     #   time     = 360
#     #   timeunit = "s"
#     #   ignore_no_cache = "on"
#     #   ignore_expired = "on"
#     #   ignore_query = "on"
#     # }

#     cache_rule_list {
#       match_method = "ext"
#       pattern = "jpg,png,gif"
#       case_ignore = "yes"
#       expire = 3600
#       expire_unit = "s"
#       ignore_no_cache_headers = "no"
#       follow_expired = "no"
#       query_params_op = "customer"
#       priority = 1
#       cache_or_not = "yes"
#       query_params_op_way = "keep"
#       query_params_op_when = "cache_back_source"
#       params = "index"
#     }

#     # Test single Map configuration item (MaxItems: 1)
#     referer {
#       type = 2 # 2 is whitelist, 1 is blacklist
#       list = ["trusted.com", "another.com"]
#       allow_empty = true
#     }

#     # Test IP whitelist configuration (update list)
#     # ip_white_list {
#     #   list = ["192.168.1.0/24", "10.0.0.0/8"]
#     # }

#     ip_black_list {
#       list = ["192.168.11.0/24"]
#     }

#     # https {
#     #   cert_id = 226116
#     #   http2 = "off"
#     #   ocsp = "off"
#     #   ssl_protocol = ["TLSv1.2"]
#     #   force_https = "301"
#     # }

#     add_response_head {
#       list {
#         name = "X-Test"
#         value = "test"
#       }

#     }

#     add_back_source_head {
#       head_name = "X-Test"
#       head_value = "test"
#       write_when_exists = "no"
#     }

#     compress_response {
#       content_type = ["text/html", "text/css", "text/javascript"]
#       min_size = 1024
#       min_size_unit = "KB"
#     }
#     speed_limit {
#       speed = 180
#       type = "ext"
#       pattern = "jpg"
#       begin_time = "08:00"
#       end_time = "18:00"
#     }
#     head_control {
#       list {
#         regex = ".*"
#         head_op = "ADD"
#         head_direction = "CLI_REQ"
#         head = "test_header"
#         value = "X-test"
#         order = 1
#       }

#     }

#     timeout {
#       time = "30"
#     }
   
#     cache_share {
#       domain = "example2.com"
#       share_way = "inner_share"
#     }
    
#     rate_limit {
#       max_rate_count = 1024
#       # leading_flow_count = 1000
#       # leading_flow_unit = "MB"
#       # max_rate_unit = "MB"
#     }

#     # deny_url {
#     #   urls = ["http://www.test.com/5.txt"]
#     # }

#     # connect_timeout {
#     #   origin_connect_timeout = "30"
#     # }
#   }
# }

# 3. Query domain and configuration
# data "edgenext_cdn_domain" "web_config_info" {
#   domain = edgenext_cdn_domain.web_domain.domain

#   # Specify configuration items to query
#   # config_item = [
#   #   "cache_rule",
#   #   "referer", 
#   #   "ip_white_list",
#   #   "add_response_head"
#   # ]
  
#   depends_on = [edgenext_cdn_domain.web_domain]
# }

# data "edgenext_cdn_domains" "web_domains" {
#   # page_number = 1
#   # page_size = 100
#   domain_status = "serving,suspend"
# }

# =============================================================================
# CDN Push Task Examples
# =============================================================================

# 1. Push specific URLs
# resource "edgenext_cdn_push" "url_refresh" {
#   type = "url"
#   urls = [
#     "https://${edgenext_cdn_domain.web_domain.domain}/config.mp4",
#     "https://${edgenext_cdn_domain.web_domain.domain}/data/abc.jpg"
#   ]

#   depends_on = [edgenext_cdn_domain.web_domain]
# }

# 2. Query push task status
# data "edgenext_cdn_push" "static_push_status" {
#   task_id = edgenext_cdn_push.url_refresh.task_id

#   depends_on = [edgenext_cdn_push.url_refresh]
# }

# 3. Query push tasks within a time period
# data "edgenext_cdn_pushes" "recent_pushes" {
#   start_time = "2025-01-01"
#   end_time   = "2025-12-31"

#   depends_on = [edgenext_cdn_push.url_refresh]
# }
