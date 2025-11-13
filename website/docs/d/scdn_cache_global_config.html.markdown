---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_global_config"
sidebar_current: "docs-edgenext-datasource-scdn_cache_global_config"
description: |-
  Use this data source to query SCDN cache global configuration.
---

# edgenext_scdn_cache_global_config

Use this data source to query SCDN cache global configuration.

## Example Usage

### Query cache global config

```hcl
data "edgenext_scdn_cache_global_config" "example" {
}

output "cache_config" {
  value = data.edgenext_scdn_cache_global_config.example.conf
}
```

### Query and save to file

```hcl
data "edgenext_scdn_cache_global_config" "example" {
  result_output_file = "cache_global_config.json"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `conf` - Cache configuration
  * `browser_cache_rule` - Browser cache configuration
    * `cachetime` - Cache time
    * `ignore_cache_time` - Ignore source cache time (cache-control)
    * `nocache` - Whether to cache
  * `cache_errstatus` - Status code cache configuration
    * `cachetime` - Status code cache time
    * `err_status` - Status code array
  * `cache_rule` - Edge TTL cache configuration
    * `action` - Cache action: 'default', 'nocache', 'cachetime', or 'force'
    * `cachetime` - Cache time
    * `ignore_cache_time` - Ignore source cache time
    * `ignore_nocache_header` - Ignore no-cache header
    * `no_cache_control_op` - No cache control operation
  * `cache_share` - Cache sharing configuration
    * `scheme` - HTTP/HTTPS cache sharing method: 'http' or 'https'
  * `cache_url_rewrite` - Custom cache key configuration
    * `cookies` - Cookie processing
      * `args_method` - Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'
      * `items` - Cookie keys
    * `ignore_case` - Ignore case
    * `queries` - Query string processing
      * `args_method` - Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'
      * `items` - Parameter keys
    * `sort_args` - Parameter sorting
  * `nocache` - Cache eligibility (true: bypass cache, false: cache)
* `id` - Rule ID (as string to match Terraform's resource ID format)
* `name` - Rule name


