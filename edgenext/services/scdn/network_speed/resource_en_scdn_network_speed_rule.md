Provides a resource to create and manage SCDN network speed rules.

Example Usage

Create custom page rule

```hcl
resource "edgenext_scdn_network_speed_rule" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "custom_page"

  custom_page {
    status_code  = 404
    page_type    = "html"
    page_content = "<html><body>Not Found</body></html>"
  }
}
```

Create upstream URI change rule

```hcl
resource "edgenext_scdn_network_speed_rule" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "upstream_uri_change_rule"

  upstream_uri_change_rule {
    typ    = "prefix"
    action = "replace"
    match  = "/old"
    target = "/new"
  }
}
```

Import

SCDN network speed rules can be imported using the rule ID:

```shell
terraform import edgenext_scdn_network_speed_rule.example 67890
```

