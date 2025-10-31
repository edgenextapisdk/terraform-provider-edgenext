---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domain_base_settings"
sidebar_current: "docs-edgenext-datasource-scdn_domain_base_settings"
description: |-
  Use this data source to query base settings of a specific SCDN domain.
---

# edgenext_scdn_domain_base_settings

Use this data source to query base settings of a specific SCDN domain.

## Example Usage

### Query domain base settings

```hcl
data "edgenext_scdn_domain_base_settings" "example" {
  domain_id = 12345
}

output "proxy_host" {
  value = data.edgenext_scdn_domain_base_settings.example.proxy_host
}

output "domain_redirect" {
  value = data.edgenext_scdn_domain_base_settings.example.domain_redirect
}
```

### Query and save to file

```hcl
data "edgenext_scdn_domain_base_settings" "example" {
  domain_id          = 12345
  result_output_file = "base_settings.json"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int) The ID of the domain to query base settings
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domain_redirect` - Domain redirect configuration
  * `jump_to` - Redirect target URL
  * `jump_type` - Redirect jump type
  * `status` - Redirect status
* `proxy_host` - Proxy host configuration
  * `proxy_host_type` - Proxy host type
  * `proxy_host` - Proxy host value
* `proxy_sni` - Proxy SNI configuration
  * `proxy_sni` - Proxy SNI value
  * `status` - Proxy SNI status


