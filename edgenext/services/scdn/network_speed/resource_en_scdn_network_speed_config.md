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

Import

SCDN network speed configuration can be imported using the business ID, business type, and config groups:

```shell
terraform import edgenext_scdn_network_speed_config.example 12345-tpl
```

