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

Import

SCDN cache rules can be imported using the rule ID:

```shell
terraform import edgenext_scdn_cache_rule.example 67890
```

