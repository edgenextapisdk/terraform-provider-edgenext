---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin_group_domain_bind"
sidebar_current: "docs-edgenext-resource-scdn_origin_group_domain_bind"
description: |-
  Provides a resource to bind domains to an SCDN origin group.
---

# edgenext_scdn_origin_group_domain_bind

Provides a resource to bind domains to an SCDN origin group.

## Example Usage

### Bind domains to origin group

```hcl
resource "edgenext_scdn_origin_group_domain_bind" "example" {
  origin_group_id = 12345
  domain_ids      = [67890, 11111]
}
```

### Bind by domain names

```hcl
resource "edgenext_scdn_origin_group_domain_bind" "example" {
  origin_group_id = 12345
  domains         = ["example.com", "www.example.com"]
}
```

### Bind by domain group IDs

```hcl
resource "edgenext_scdn_origin_group_domain_bind" "example" {
  origin_group_id  = 12345
  domain_group_ids = [1, 2]
}
```

## Argument Reference

The following arguments are supported:

* `origin_group_id` - (Required, Int, ForceNew) Origin group ID
* `domain_group_ids` - (Optional, List: [`Int`], ForceNew) Domain group ID array
* `domain_ids` - (Optional, List: [`Int`], ForceNew) Domain ID array
* `domains` - (Optional, List: [`String`], ForceNew) Domain array

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `job_id` - Batch job ID


## Import

SCDN origin group domain bindings can be imported using the origin group ID and domain IDs:

```shell
terraform import edgenext_scdn_origin_group_domain_bind.example 12345-67890,11111
```

