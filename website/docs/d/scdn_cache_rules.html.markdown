---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_rules"
sidebar_current: "docs-edgenext-datasource-scdn_cache_rules"
description: |-
  Use this data source to query SCDN cache rules.
---

# edgenext_scdn_cache_rules

Use this data source to query SCDN cache rules.

## Example Usage

### Query cache rules for template

```hcl
data "edgenext_scdn_cache_rules" "example" {
  business_id   = 12345
  business_type = "tpl"
}

output "rule_count" {
  value = data.edgenext_scdn_cache_rules.example.total
}

output "rules" {
  value = data.edgenext_scdn_cache_rules.example.list
}
```

### Query cache rules for domain

```hcl
data "edgenext_scdn_cache_rules" "example" {
  business_id   = 67890
  business_type = "domain"
  page          = 1
  page_size     = 50
}
```

### Query and save to file

```hcl
data "edgenext_scdn_cache_rules" "example" {
  business_id        = 12345
  business_type      = "tpl"
  result_output_file = "cache_rules.json"
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int) Business ID (template ID for 'tpl' type, domain ID for 'domain' type)
* `business_type` - (Required, String) Business type: 'tpl' (template) or 'domain'
* `page_size` - (Optional, Int) Page size
* `page` - (Optional, Int) Page number
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of cache rules
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
  * `expr` - Wirefilter rule
  * `id` - Rule ID
  * `name` - Rule name
  * `remark` - Remark
  * `status` - Status (1: enabled, 2: disabled)
  * `type` - Type: 'domain', 'tpl', or 'global'
  * `weight` - Weight
* `total` - Total number of rules


