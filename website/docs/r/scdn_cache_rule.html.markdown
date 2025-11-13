---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_rule"
sidebar_current: "docs-edgenext-resource-scdn_cache_rule"
description: |-
  Provides a resource to create and manage SCDN cache rules.
---

# edgenext_scdn_cache_rule

Provides a resource to create and manage SCDN cache rules.

## Example Usage

### Create cache rule with cache time

```hcl
resource "edgenext_scdn_cache_rule" "example" {
  business_id   = 12345
  business_type = "tpl"
  name          = "my-cache-rule"
  expr          = "uri.path == \"/static/*\""

  conf {
    nocache = false

    cache_rule {
      cachetime             = 3600
      ignore_cache_time     = false
      ignore_nocache_header = false
      action                = "cachetime"
    }
  }
}
```

### Create cache rule with no cache

```hcl
resource "edgenext_scdn_cache_rule" "example" {
  business_id   = 12345
  business_type = "tpl"
  name          = "no-cache-rule"
  expr          = "uri.path == \"/api/*\""

  conf {
    nocache = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID (template ID for 'tpl' type, domain ID for 'domain' type)
* `business_type` - (Required, String, ForceNew) Business type: 'tpl' (template) or 'domain'
* `conf` - (Required, List) Cache configuration
* `name` - (Required, String) Rule name
* `expr` - (Optional, String) Wirefilter rule. Empty string means 'allow all'. If not set (null), keeps existing value.
* `remark` - (Optional, String) Rule remark
* `rule_id` - (Optional, Int) Rule ID for updating existing rule. If provided, this will update the rule instead of creating a new one.

The `browser_cache_rule` object of `conf` supports the following:

* `cachetime` - (Required, Int) Cache time
* `ignore_cache_time` - (Required, Bool) Ignore source cache time (cache-control)
* `nocache` - (Required, Bool) Whether to cache

The `cache_errstatus` object of `conf` supports the following:

* `cachetime` - (Required, Int) Status code cache time
* `err_status` - (Required, List) Status code array

The `cache_rule` object of `conf` supports the following:

* `cachetime` - (Required, Int) Cache time
* `action` - (Optional, String) Cache action: 'default', 'nocache', 'cachetime', or 'force'
* `ignore_cache_time` - (Optional, Bool) Ignore source cache time
* `ignore_nocache_header` - (Optional, Bool) Ignore no-cache header
* `no_cache_control_op` - (Optional, String) No cache control operation

The `cache_share` object of `conf` supports the following:

* `scheme` - (Required, String) HTTP/HTTPS cache sharing method: 'http' or 'https'

The `cache_url_rewrite` object of `conf` supports the following:

* `ignore_case` - (Required, Bool) Ignore case
* `sort_args` - (Required, Bool) Parameter sorting
* `cookies` - (Optional, List) Cookie processing
* `queries` - (Optional, List) Query string processing

The `conf` object supports the following:

* `cache_share` - (Required, List) Cache sharing configuration
* `nocache` - (Required, Bool) Cache eligibility (true: bypass cache, false: cache)
* `browser_cache_rule` - (Optional, List) Browser cache configuration
* `cache_errstatus` - (Optional, List) Status code cache configuration
* `cache_rule` - (Optional, List) Edge TTL cache configuration
* `cache_url_rewrite` - (Optional, List) Custom cache key configuration

The `cookies` object of `cache_url_rewrite` supports the following:

* `args_method` - (Required, String) Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'
* `items` - (Required, List) Cookie keys

The `queries` object of `cache_url_rewrite` supports the following:

* `args_method` - (Required, String) Action: 'SAVE', 'DEL', 'IGNORE', or 'CUT'
* `items` - (Required, List) Parameter keys

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `id` - The rule ID
* `status` - Status (1: enabled, 2: disabled)
* `type` - Type: 'domain', 'tpl', or 'global'
* `weight` - Weight


## Import

SCDN cache rules can be imported using the rule ID:

```shell
terraform import edgenext_scdn_cache_rule.example 67890
```

