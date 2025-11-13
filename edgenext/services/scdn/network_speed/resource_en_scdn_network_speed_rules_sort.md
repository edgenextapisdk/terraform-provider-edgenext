Provides a resource to sort SCDN network speed rules.

Example Usage

Sort network speed rules

```hcl
resource "edgenext_scdn_network_speed_rules_sort" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "custom_page"
  ids           = [3, 1, 2]
}
```

Import

SCDN network speed rules sort can be imported using the business ID, business type, and config group:

```shell
terraform import edgenext_scdn_network_speed_rules_sort.example 12345-tpl-custom_page
```

