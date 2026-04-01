Provides a resource to create and manage SCDN cache rules.

Example Usage

Create cache rule with cache time

```hcl
resource "edgenext_scdn_cache_rule" "example" {
  business_id   = 12345
  business_type = "tpl"
  name          = "my-cache-rule"
  expr          = "uri.path == \"/static/*\""

  conf {
    nocache = false

    cache_rule {
      cachetime              = 3600
      ignore_cache_time      = false
      ignore_nocache_header  = false
      action                 = "cachetime"
    }
  }
}
```

Create cache rule with no cache

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

Create cache rule with minimal config (server provides defaults)

```hcl
# When only nocache = false is set, the server will populate default values for
# cache_rule and other sub-blocks. Terraform accepts these server-side defaults
# without showing drift on subsequent plans.
resource "edgenext_scdn_cache_rule" "example" {
  business_id   = 12345
  business_type = "tpl"
  name          = "minimal-cache-rule"
  remark        = "minimal config"

  conf {
    nocache = false
  }
}
```

Argument Reference

The following arguments are supported:

* `business_id` - (Optional, Computed, ForceNew) Business ID (template ID for 'tpl' type, domain ID for 'domain' type).
* `business_type` - (Optional, Computed, ForceNew) Business type: 'tpl' (template) or 'domain'.
* `name` - (Required) Rule name.
* `expr` - (Optional) Wirefilter rule. Empty string means 'allow all'.
* `remark` - (Optional) Rule remark.
* `conf` - (Required) Cache configuration. The `conf` block supports:
  * `nocache` - (Required) Cache eligibility (true: bypass cache, false: cache).
  * `cache_rule` - (Optional, Computed) Edge TTL cache configuration. When not specified, the server may return default values which Terraform stores in state without causing drift. The `cache_rule` block supports:
    * `cachetime` - (Required) Cache time in seconds.
    * `ignore_cache_time` - (Optional) Ignore source cache time.
    * `ignore_nocache_header` - (Optional) Ignore no-cache header.
    * `no_cache_control_op` - (Optional) No cache control operation.
    * `action` - (Optional) Cache action: 'default', 'nocache', 'cachetime', or 'force'.
  * `browser_cache_rule` - (Optional, Computed) Browser cache configuration. The `browser_cache_rule` block supports:
    * `cachetime` - (Required) Cache time in seconds.
    * `ignore_cache_time` - (Required) Ignore source cache time (cache-control).
    * `nocache` - (Required) Whether to cache.
  * `cache_errstatus` - (Optional, Computed) Status code cache configuration. Multiple blocks are supported. Each block supports:
    * `cachetime` - (Required) Status code cache time in seconds.
    * `err_status` - (Required) List of HTTP status codes.
  * `cache_url_rewrite` - (Optional, Computed) Custom cache key configuration. The `cache_url_rewrite` block supports:
    * `sort_args` - (Required) Parameter sorting.
    * `ignore_case` - (Required) Ignore case.
    * `queries` - (Optional) Query string processing.
    * `cookies` - (Optional) Cookie processing.
  * `cache_share` - (Optional, Computed) Cache sharing configuration. The `cache_share` block supports:
    * `scheme` - (Optional, Computed) HTTP/HTTPS cache sharing method: '', 'http' or 'https'.
* `rule_id` - (Optional, Computed) Rule ID for adopting an existing rule. If provided, Terraform will manage the existing rule instead of creating a new one.

> **Note:** Sub-blocks under `conf` (`cache_rule`, `browser_cache_rule`, `cache_errstatus`, `cache_url_rewrite`, `cache_share`) are all `Optional + Computed`. This means you can configure them explicitly in HCL, or omit them — in which case the server may return default values that Terraform stores in state without treating them as drift on subsequent plans.

Import

SCDN cache rules can be imported using the composite ID: `{business_id}-{business_type}-{rule_id}`.

```shell
terraform import edgenext_scdn_cache_rule.example 12345-tpl-67890
```

