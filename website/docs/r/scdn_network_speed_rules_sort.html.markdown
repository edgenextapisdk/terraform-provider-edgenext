---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_network_speed_rules_sort"
sidebar_current: "docs-edgenext-resource-scdn_network_speed_rules_sort"
description: |-
  Provides a resource to sort SCDN network speed rules.
---

# edgenext_scdn_network_speed_rules_sort

Provides a resource to sort SCDN network speed rules.

## Example Usage

### Sort network speed rules

```hcl
resource "edgenext_scdn_network_speed_rules_sort" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "custom_page"
  ids           = [3, 1, 2]
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID (template ID for 'tpl' type, user ID for 'global' type)
* `business_type` - (Required, String, ForceNew) Business type: 'tpl' (template) or 'global'
* `config_group` - (Required, String, ForceNew) Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'
* `ids` - (Required, List: [`Int`]) Sorted rule IDs array (order matters)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `sorted_ids` - Sorted rule IDs after sorting


## Import

SCDN network speed rules sort can be imported using the business ID, business type, and config group:

```shell
terraform import edgenext_scdn_network_speed_rules_sort.example 12345-tpl-custom_page
```

