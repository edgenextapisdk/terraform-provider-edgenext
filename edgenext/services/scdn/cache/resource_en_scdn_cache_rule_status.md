Provides a resource to manage the status (enable/disable) of SCDN cache rules.

Example Usage

Enable cache rules

```hcl
resource "edgenext_scdn_cache_rule_status" "example" {
  business_id   = 12345
  business_type = "tpl"
  rule_ids      = [1, 2, 3]
  status        = 1  # 1: enabled
}
```

Disable cache rules

```hcl
resource "edgenext_scdn_cache_rule_status" "example" {
  business_id   = 12345
  business_type = "tpl"
  rule_ids      = [1, 2, 3]
  status        = 2  # 2: disabled
}
```

Import

SCDN cache rule status can be imported using the business ID, business type, and rule IDs:

```shell
terraform import edgenext_scdn_cache_rule_status.example 12345-tpl-1,2,3
```

