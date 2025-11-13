---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_network_speed_config"
sidebar_current: "docs-edgenext-datasource-scdn_network_speed_config"
description: |-
  Use this data source to query SCDN network speed configuration.
---

# edgenext_scdn_network_speed_config

Use this data source to query SCDN network speed configuration.

## Example Usage

### Query network speed config for template

```hcl
data "edgenext_scdn_network_speed_config" "example" {
  business_id   = 12345
  business_type = "tpl"
}

output "domain_proxy_conf" {
  value = data.edgenext_scdn_network_speed_config.example.domain_proxy_conf
}
```

### Query specific config groups

```hcl
data "edgenext_scdn_network_speed_config" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_groups = ["domain_proxy_conf", "upstream_redirect"]
}
```

### Query and save to file

```hcl
data "edgenext_scdn_network_speed_config" "example" {
  business_id        = 12345
  business_type      = "tpl"
  result_output_file = "network_speed_config.json"
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int) Business ID (template ID for 'tpl' type, user ID for 'global' type)
* `business_type` - (Required, String) Business type: 'tpl' (template) or 'global'
* `config_groups` - (Optional, List: [`String`]) Configuration groups to retrieve
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `custom_page` - Custom page configuration
  * `status` - Status: 'on' or 'off'
* `customized_req_headers` - Customized request headers configuration
  * `status` - Status: 'on' or 'off'
* `domain_proxy_conf` - Domain proxy configuration
  * `fails_timeout` - Failure timeout
  * `keep_new_src_time` - Keep new source time
  * `max_fails` - Max failures
  * `proxy_connect_timeout` - Connection timeout
  * `proxy_keepalive` - Keepalive (0 or 1)
* `https` - HTTPS configuration
  * `ciphers_preset` - Ciphers preset: 'default', 'strong', or 'custom'
  * `custom_encrypt_algorithm` - Custom encryption algorithms
  * `hsts` - HSTS: 'on' or 'off'
  * `http2` - HTTP2: 'on' or 'off'
  * `http2https_port` - Redirect port
  * `http2https` - HTTP to HTTPS redirect: 'off', 'all', or 'special'
  * `min_version` - Minimum TLS version
  * `ocsp_stapling` - OCSP Stapling: 'on' or 'off'
  * `status` - Status: 'on' or 'off'
* `mobile_jump` - Mobile jump configuration
  * `jump_url` - Jump URL
  * `status` - Status: 'on' or 'off'
* `page_gzip` - Page Gzip configuration
  * `status` - Status: 'on' or 'off'
* `resp_headers` - Response headers configuration
  * `status` - Status: 'on' or 'off'
* `slice` - Range request configuration
  * `status` - Status: 'on' or 'off'
* `source_site_protect` - Source site protection configuration
  * `num` - Number of requests
  * `second` - Time in seconds
  * `status` - Status: 'on' or 'off'
* `upload_file` - Upload file configuration
  * `upload_size_unit` - Unit (e.g., 'MB')
  * `upload_size` - Upload size
* `upstream_check` - Upstream check configuration
  * `fails` - Consecutive unavailable times (1-10)
  * `intval` - Check interval in seconds (3-300)
  * `op` - HTTP method: 'HEAD', 'GET', or 'AUTO' (when type is 'http')
  * `path` - HTTP check path (when type is 'http')
  * `rise` - Consecutive available times (1-10)
  * `status` - Status: 'on' or 'off'
  * `timeout` - TCP connection timeout in seconds (1-10)
  * `type` - Check type: 'tcp' or 'http'
* `upstream_redirect` - Upstream redirect configuration
  * `status` - Status: 'on' or 'off'
* `upstream_uri_change` - Upstream URI change configuration
  * `status` - Status: 'on' or 'off'
* `webp` - WebP format configuration
  * `status` - Status: 'on' or 'off'
* `websocket` - WebSocket configuration
  * `status` - Status: 'on' or 'off'


