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

Argument Reference

The following arguments are supported:

* `business_id` - (Optional, Computed, ForceNew) Business ID.
* `business_type` - (Optional, Computed, ForceNew) Business type: 'tpl' or 'global'.
* `config_group` - (Optional, Computed, ForceNew) Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'. If omitted, it will be inferred from the provided rule block.
* `rule_id` - (Optional, Computed) Rule ID for adopting an existing rule. If provided, Terraform will manage the existing rule instead of creating a new one.
* `custom_page` - (Optional) Custom page rule.
* `upstream_uri_change_rule` - (Optional) Upstream URI change rule.
* `resp_headers_rule` - (Optional) Response headers rule.
* `customized_req_headers_rule` - (Optional) Customized request headers rule.

Import

SCDN network speed rules can be imported using the composite ID: `{business_id}-{business_type}-{config_group}-{rule_id}`.

```shell
terraform import edgenext_scdn_network_speed_rule.example 12345-tpl-custom_page-67890
```

