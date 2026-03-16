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

Argument Reference

The following arguments are supported:

* `business_id` - (Optional, Computed, ForceNew) Business ID (template ID for 'tpl' type, domain ID for 'domain' type).
* `business_type` - (Optional, Computed, ForceNew) Business type: 'tpl' (template) or 'domain'.
* `name` - (Required) Rule name.
* `expr` - (Optional) Wirefilter rule. Empty string means 'allow all'.
* `remark` - (Optional) Rule remark.
* `conf` - (Required) Cache configuration. See [Conf Block](#conf-block) below. The `conf` block supports:
  * `nocache` - (Required) Cache eligibility (true: bypass cache, false: cache).
  * `cache_rule` - (Optional) Edge TTL cache configuration.
  * `browser_cache_rule` - (Optional) Browser cache configuration.
  * `cache_errstatus` - (Optional) Status code cache configuration.
  * `cache_url_rewrite` - (Optional) Custom cache key configuration.
  * `cache_share` - (Optional) Cache sharing configuration.
* `rule_id` - (Optional, Computed) Rule ID for adopting an existing rule. If provided, Terraform will manage the existing rule instead of creating a new one.

Import

SCDN cache rules can be imported using the composite ID: `{business_id}-{business_type}-{rule_id}`.

```shell
terraform import edgenext_scdn_cache_rule.example 12345-tpl-67890
```

