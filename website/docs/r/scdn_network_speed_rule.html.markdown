---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_network_speed_rule"
sidebar_current: "docs-edgenext-resource-scdn_network_speed_rule"
description: |-
  Provides a resource to create and manage SCDN network speed rules.
---

# edgenext_scdn_network_speed_rule

Provides a resource to create and manage SCDN network speed rules.

## Example Usage

### Create custom page rule

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

### Create upstream URI change rule

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

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID
* `business_type` - (Required, String, ForceNew) Business type: 'tpl' or 'global'
* `config_group` - (Required, String, ForceNew) Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'
* `custom_page` - (Optional, List) Custom page rule
* `customized_req_headers_rule` - (Optional, List) Customized request headers rule
* `resp_headers_rule` - (Optional, List) Response headers rule
* `rule_id` - (Optional, Int) Rule ID for updating existing rule. If provided, this will update the rule instead of creating a new one.
* `upstream_uri_change_rule` - (Optional, List) Upstream URI change rule

The `custom_page` object supports the following:

* `page_content` - (Required, String) Page content
* `page_type` - (Required, String) Page type
* `status_code` - (Required, Int) Status code

The `customized_req_headers_rule` object supports the following:

* `content` - (Required, String) Content
* `type` - (Required, String) Type
* `remark` - (Optional, String) Remark

The `resp_headers_rule` object supports the following:

* `content` - (Required, String) Content
* `type` - (Required, String) Type
* `remark` - (Optional, String) Remark

The `upstream_uri_change_rule` object supports the following:

* `action` - (Required, String) Action
* `match` - (Required, String) Match value
* `target` - (Required, String) Target value
* `typ` - (Required, String) Type

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `id` - The rule ID


## Import

SCDN network speed rules can be imported using the rule ID:

```shell
terraform import edgenext_scdn_network_speed_rule.example 67890
```

