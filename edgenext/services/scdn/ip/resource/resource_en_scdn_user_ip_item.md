---
subcategory: "SCDN"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_user_ip_item"
sidebar_current: "docs-edgenext-resource-scdn-user-ip-item"
description: |-
  Provides a resource to manage SCDN User IP Items.
---

# Resource: edgenext_scdn_user_ip_item

Provides a resource to manage individual IP items within an SCDN User IP List.

## Example Usage

```hcl
resource "edgenext_scdn_user_ip" "list" {
  name = "example-list"
}

resource "edgenext_scdn_user_ip_item" "item" {
  user_ip_id = edgenext_scdn_user_ip.list.id
  ip         = "192.168.1.1"
  remark     = "Office IP"
}
```

## Argument Reference

The following arguments are supported:

* `user_ip_id` - (Required) The ID of the User IP List to which this item belongs.
* `ip` - (Required) The IP address or CIDR block.
* `remark` - (Optional) The remark or description for the IP item.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID (UUID) of the IP item.
* `format_created_at` - The formatted creation time.
* `format_updated_at` - The formatted last update time.

## Import

SCDN User IP Items can be imported using the combined `user_ip_id` and `id` separated by a colon, e.g.

```
$ terraform import edgenext_scdn_user_ip_item.example <user_ip_id>:<uuid-of-item>
# Example:
$ terraform import edgenext_scdn_user_ip_item.example 123:681c5d277ef73537f36fbdb6
```
