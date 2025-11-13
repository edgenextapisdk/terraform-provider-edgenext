---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin_group_domain_copy"
sidebar_current: "docs-edgenext-resource-scdn_origin_group_domain_copy"
description: |-
  Provides a resource to copy an SCDN origin group configuration to a domain.
---

# edgenext_scdn_origin_group_domain_copy

Provides a resource to copy an SCDN origin group configuration to a domain.

## Example Usage

### Copy origin group to domain

```hcl
resource "edgenext_scdn_origin_group_domain_copy" "example" {
  origin_group_id = 12345
  domain_id       = 67890
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int, ForceNew) Domain ID
* `origin_group_id` - (Required, Int, ForceNew) Origin group ID

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SCDN origin group domain copies can be imported using the origin group ID and domain ID:

```shell
terraform import edgenext_scdn_origin_group_domain_copy.example 12345-67890
```

