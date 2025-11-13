---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_rule_status"
sidebar_current: "docs-edgenext-resource-scdn_cache_rule_status"
description: |-
  Provides a resource to manage the status (enable/disable) of SCDN cache rules.
---

# edgenext_scdn_cache_rule_status

Provides a resource to manage the status (enable/disable) of SCDN cache rules.

## Example Usage

### Enable cache rules

```hcl
resource "edgenext_scdn_cache_rule_status" "example" {
  business_id   = 12345
  business_type = "tpl"
  rule_ids      = [1, 2, 3]
  status        = 1 # 1: enabled
}
```

### Disable cache rules

```hcl
resource "edgenext_scdn_cache_rule_status" "example" {
  business_id   = 12345
  business_type = "tpl"
  rule_ids      = [1, 2, 3]
  status        = 2 # 2: disabled
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID (template ID for 'tpl' type, domain ID for 'domain' type)
* `business_type` - (Required, String, ForceNew) Business type: 'tpl' (template) or 'domain'
* `rule_ids` - (Required, List: [`Int`]) Rule IDs array to update status
* `status` - (Required, Int) Status: 1 (enabled) or 2 (disabled)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `updated_ids` - Rule IDs that were updated


## Import

SCDN cache rule status can be imported using the business ID, business type, and rule IDs:

```shell
terraform import edgenext_scdn_cache_rule_status.example 12345-tpl-1,2,3
```

