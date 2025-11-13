---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domain_node_switch"
sidebar_current: "docs-edgenext-resource-scdn_domain_node_switch"
description: |-
  Provides a resource to switch the node type of an SCDN domain.
---

# edgenext_scdn_domain_node_switch

Provides a resource to switch the node type of an SCDN domain.

## Example Usage

### Switch to SCDN node

```hcl
resource "edgenext_scdn_domain_node_switch" "example" {
  domain_id      = 12345
  protect_status = "scdn"
}
```

### Switch to exclusive node

```hcl
resource "edgenext_scdn_domain_node_switch" "example" {
  domain_id             = 12345
  protect_status        = "exclusive"
  exclusive_resource_id = 999
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int, ForceNew) The ID of the domain to switch nodes
* `protect_status` - (Required, String) The edge node type. Valid values: back_source, scdn, exclusive
* `exclusive_resource_id` - (Optional, Int) The ID of the exclusive resource package (required if protect_status is exclusive)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SCDN domain node switch can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain_node_switch.example 12345
```

