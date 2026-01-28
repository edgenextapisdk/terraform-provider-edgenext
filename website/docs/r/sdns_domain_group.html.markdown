---
subcategory: "Security DNS (SDNS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_sdns_domain_group"
sidebar_current: "docs-edgenext-resource-sdns_domain_group"
description: |-
  Provides a resource to create and manage SDNS domain groups.
---

# edgenext_sdns_domain_group

Provides a resource to create and manage SDNS domain groups.

## Example Usage

### Create SDNS domain group

```hcl
resource "edgenext_sdns_domain_group" "example" {
  group_name = "my-domain-group"
  remark     = "This is a test group"
  domain_ids = ["123", "456"]
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String) The name of the DNS domain group
* `domain_ids` - (Optional, Set: [`String`]) List of domain IDs to bind to the group
* `remark` - (Optional, String) Remark for the domain group

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SDNS domain groups can be imported using the group ID:

```shell
terraform import edgenext_sdns_domain_group.example 67890
```

