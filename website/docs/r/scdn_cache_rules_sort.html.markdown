---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_rules_sort"
sidebar_current: "docs-edgenext-resource-scdn_cache_rules_sort"
description: |-
  Provides a resource to sort SCDN cache rules.
---

# edgenext_scdn_cache_rules_sort

Provides a resource to sort SCDN cache rules.

## Example Usage

### Sort cache rules

```hcl
resource "edgenext_scdn_cache_rules_sort" "example" {
  business_id   = 12345
  business_type = "tpl"
  ids           = [3, 1, 2]
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID (template ID for 'tpl' type, domain ID for 'domain' type)
* `business_type` - (Required, String, ForceNew) Business type: 'tpl' (template) or 'domain'
* `ids` - (Required, List: [`Int`]) Sorted rule IDs array (order matters)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `sorted_ids` - Sorted rule IDs after sorting


## Import

SCDN cache rules sort can be imported using the business ID and business type:

```shell
terraform import edgenext_scdn_cache_rules_sort.example 12345-tpl
```

