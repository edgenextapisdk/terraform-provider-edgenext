Provides a resource to manage SCDN network speed configuration.

Example Usage

Configure network speed for template

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

Argument Reference

The following arguments are supported:

* `business_id` - (Optional/Computed) Business ID (template ID for 'tpl' type, user ID for 'global' type). **Note**: This is required when creating a new resource but can be omitted during import as it will be parsed from the import ID.
* `business_type` - (Optional/Computed) Business type: 'tpl' (template) or 'global'. **Note**: This is required when creating a new resource but can be omitted during import.
* `domain_proxy_conf` - (Optional) Domain proxy configuration.
* `upstream_redirect` - (Optional) Upstream redirect configuration.
* `customized_req_headers` - (Optional) Customized request headers configuration.
* `resp_headers` - (Optional) Response headers configuration.
* `upstream_uri_change` - (Optional) Upstream URI change configuration.
* `source_site_protect` - (Optional) Source site protection configuration.
* `slice` - (Optional) Range request configuration.
* `https` - (Optional) HTTPS configuration.
* `page_gzip` - (Optional) Page Gzip configuration.
* `webp` - (Optional) WebP format configuration.
* `upload_file` - (Optional) Upload file configuration.
* `websocket` - (Optional) WebSocket configuration.
* `mobile_jump` - (Optional) Mobile jump configuration.
* `custom_page` - (Optional) Custom page configuration.
* `upstream_check` - (Optional) Upstream check configuration.

Import

SCDN network speed configuration can be imported using the resource ID in the format `<business_id>-<business_type>`. The provider will automatically parse the ID and populate the required fields.

```shell
terraform import edgenext_scdn_network_speed_config.example 12345-tpl
```

