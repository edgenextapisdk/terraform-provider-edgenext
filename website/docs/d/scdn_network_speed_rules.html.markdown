---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_network_speed_rules"
sidebar_current: "docs-edgenext-datasource-scdn_network_speed_rules"
description: |-
  Use this data source to query SCDN network speed rules.
---

# edgenext_scdn_network_speed_rules

Use this data source to query SCDN network speed rules.

## Example Usage

### Query custom page rules

```hcl
data "edgenext_scdn_network_speed_rules" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "custom_page"
}

output "rule_count" {
  value = data.edgenext_scdn_network_speed_rules.example.total
}

output "rules" {
  value = data.edgenext_scdn_network_speed_rules.example.list
}
```

### Query upstream URI change rules

```hcl
data "edgenext_scdn_network_speed_rules" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "upstream_uri_change_rule"
}
```

### Query and save to file

```hcl
data "edgenext_scdn_network_speed_rules" "example" {
  business_id        = 12345
  business_type      = "tpl"
  config_group       = "custom_page"
  result_output_file = "network_speed_rules.json"
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int) Business ID
* `business_type` - (Required, String) Business type: 'tpl' or 'global'
* `config_group` - (Required, String) Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of rules
  * `business_id` - Business ID
  * `business_type` - Business type
  * `config_group` - Config group
  * `custom_page` - Custom page rule
    * `page_content` - Page content
    * `page_type` - Page type
    * `status_code` - Status code
  * `customized_req_headers_rule` - Customized request headers rule
    * `content` - Content
    * `remark` - Remark
    * `type` - Type
  * `id` - Rule ID
  * `resp_headers_rule` - Response headers rule
    * `content` - Content
    * `remark` - Remark
    * `type` - Type
  * `upstream_uri_change_rule` - Upstream URI change rule
    * `action` - Action
    * `match` - Match value
    * `target` - Target value
    * `typ` - Type
* `total` - Total number of rules


