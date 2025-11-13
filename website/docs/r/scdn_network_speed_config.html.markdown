---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_network_speed_config"
sidebar_current: "docs-edgenext-resource-scdn_network_speed_config"
description: |-
  Provides a resource to manage SCDN network speed configuration.
---

# edgenext_scdn_network_speed_config

Provides a resource to manage SCDN network speed configuration.

## Example Usage

### Configure network speed for template

```hcl
resource "edgenext_scdn_network_speed_config" "example" {
  business_id   = 12345
  business_type = "tpl"

  domain_proxy_conf {
    proxy_connect_timeout = 30
    fails_timeout         = 10
    keep_new_src_time     = 60
    max_fails             = 3
    proxy_keepalive       = 1
  }

  upstream_redirect {
    status = "on"
  }
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID (template ID for 'tpl' type, user ID for 'global' type)
* `business_type` - (Required, String, ForceNew) Business type: 'tpl' (template) or 'global'
* `custom_page` - (Optional, List) Custom page configuration
* `customized_req_headers` - (Optional, List) Customized request headers configuration
* `domain_proxy_conf` - (Optional, List) Domain proxy configuration
* `https` - (Optional, List) HTTPS configuration
* `mobile_jump` - (Optional, List) Mobile jump configuration
* `page_gzip` - (Optional, List) Page Gzip configuration
* `resp_headers` - (Optional, List) Response headers configuration
* `slice` - (Optional, List) Range request configuration
* `source_site_protect` - (Optional, List) Source site protection configuration
* `upload_file` - (Optional, List) Upload file configuration
* `upstream_check` - (Optional, List) Upstream check configuration
* `upstream_redirect` - (Optional, List) Upstream redirect configuration
* `upstream_uri_change` - (Optional, List) Upstream URI change configuration
* `webp` - (Optional, List) WebP format configuration
* `websocket` - (Optional, List) WebSocket configuration

The `custom_page` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `customized_req_headers` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `domain_proxy_conf` object supports the following:

* `fails_timeout` - (Optional, Int) Failure timeout
* `keep_new_src_time` - (Optional, Int) Keep new source time
* `max_fails` - (Optional, Int) Max failures
* `proxy_connect_timeout` - (Optional, Int) Connection timeout
* `proxy_keepalive` - (Optional, Int) Keepalive (0 or 1)

The `https` object supports the following:

* `ciphers_preset` - (Optional, String) Ciphers preset: 'default', 'strong', or 'custom'
* `custom_encrypt_algorithm` - (Optional, List) Custom encryption algorithms
* `hsts` - (Optional, String) HSTS: 'on' or 'off'
* `http2` - (Optional, String) HTTP2: 'on' or 'off'
* `http2https_port` - (Optional, Int) Redirect port
* `http2https` - (Optional, String) HTTP to HTTPS redirect: 'off', 'all', or 'special'
* `min_version` - (Optional, String) Minimum TLS version
* `ocsp_stapling` - (Optional, String) OCSP Stapling: 'on' or 'off'
* `status` - (Optional, String) Status: 'on' or 'off'

The `mobile_jump` object supports the following:

* `jump_url` - (Optional, String) Jump URL
* `status` - (Optional, String) Status: 'on' or 'off'

The `page_gzip` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `resp_headers` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `slice` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `source_site_protect` object supports the following:

* `num` - (Optional, Int) Number of requests
* `second` - (Optional, Int) Time in seconds
* `status` - (Optional, String) Status: 'on' or 'off'

The `upload_file` object supports the following:

* `upload_size_unit` - (Optional, String) Unit (e.g., 'MB')
* `upload_size` - (Optional, Int) Upload size

The `upstream_check` object supports the following:

* `fails` - (Optional, Int) Consecutive unavailable times (1-10)
* `intval` - (Optional, Int) Check interval in seconds (3-300)
* `op` - (Optional, String) HTTP method: 'HEAD', 'GET', or 'AUTO' (required when type is 'http')
* `path` - (Optional, String) HTTP check path, must start with '/' (required when type is 'http')
* `rise` - (Optional, Int) Consecutive available times (1-10)
* `status` - (Optional, String) Status: 'on' or 'off'
* `timeout` - (Optional, Int) TCP connection timeout in seconds (1-10)
* `type` - (Optional, String) Check type: 'tcp' or 'http'

The `upstream_redirect` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `upstream_uri_change` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `webp` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

The `websocket` object supports the following:

* `status` - (Optional, String) Status: 'on' or 'off'

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SCDN network speed configuration can be imported using the business ID, business type, and config groups:

```shell
terraform import edgenext_scdn_network_speed_config.example 12345-tpl
```

