---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domain_status"
sidebar_current: "docs-edgenext-resource-scdn_domain_status"
description: |-
  Provides a resource to manage the status (enable/disable) of an SCDN domain.
---

# edgenext_scdn_domain_status

Provides a resource to manage the status (enable/disable) of an SCDN domain.

## Example Usage

### Enable domain

```hcl
resource "edgenext_scdn_domain_status" "example" {
  domain_id = 12345
  enabled   = true
}
```

### Disable domain

```hcl
resource "edgenext_scdn_domain_status" "example" {
  domain_id = 12345
  enabled   = false
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int, ForceNew) The ID of the domain to manage status
* `enabled` - (Required, Bool) Whether the domain is enabled (true) or disabled (false)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SCDN domain status can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain_status.example 12345
```

