---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domain_base_settings"
sidebar_current: "docs-edgenext-resource-scdn_domain_base_settings"
description: |-
  Provides a resource to manage base settings of an SCDN domain.
---

# edgenext_scdn_domain_base_settings

Provides a resource to manage base settings of an SCDN domain.

## Example Usage

### Update domain base settings

```hcl
resource "edgenext_scdn_domain_base_settings" "example" {
  domain_id = 12345

  proxy_host {
    proxy_host      = "example.com"
    proxy_host_type = "custom"
  }

  proxy_sni {
    proxy_sni = "example.com"
    status    = "on"
  }

  domain_redirect {
    status    = "on"
    jump_to   = "https://www.example.com"
    jump_type = "301"
  }
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int, ForceNew) The ID of the domain to update base settings
* `domain_redirect` - (Optional, List) Domain redirect configuration
* `proxy_host` - (Optional, List) Proxy host configuration
* `proxy_sni` - (Optional, List) Proxy SNI configuration

The `domain_redirect` object supports the following:

* `jump_to` - (Optional, String) Redirect target URL
* `jump_type` - (Optional, String) Redirect jump type
* `status` - (Optional, String) Redirect status (on/off)

The `proxy_host` object supports the following:

* `proxy_host_type` - (Optional, String) Proxy host type
* `proxy_host` - (Optional, String) Proxy host value

The `proxy_sni` object supports the following:

* `proxy_sni` - (Optional, String) Proxy SNI value
* `status` - (Optional, String) Proxy SNI status

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SCDN domain base settings can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain_base_settings.example 12345
```

