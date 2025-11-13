Provides a resource to sort SCDN cache rules.

Example Usage

Sort cache rules

```hcl
resource "edgenext_scdn_cache_rules_sort" "example" {
  business_id   = 12345
  business_type = "tpl"
  ids           = [3, 1, 2]
}
```

Import

SCDN cache rules sort can be imported using the business ID and business type:

```shell
terraform import edgenext_scdn_cache_rules_sort.example 12345-tpl
```

