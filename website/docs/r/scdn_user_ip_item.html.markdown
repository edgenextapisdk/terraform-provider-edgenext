---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_user_ip_item"
sidebar_current: "docs-edgenext-resource-scdn_user_ip_item"
description: |-
  # edgenext_scdn_user_ip_item
---

# edgenext_scdn_user_ip_item

# edgenext_scdn_user_ip_item

Provides a resource to manage individual IP items within an SCDN User IP List.

## Example Usage

### Create a user IP item

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

* `ip` - (Required, String) The IP address or CIDR
* `user_ip_id` - (Required, Int) The ID of the IP list to which this item belongs
* `remark` - (Optional, String) Remark for the IP item

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `format_created_at` - Creation time
* `format_updated_at` - Last update time
* `id` - The ID (UUID) of the IP item


## Import

SCDN User IP Items can be imported using the combined `user_ip_id` and `id` separated by a colon, e.g.

```
$ terraform import edgenext_scdn_user_ip_item.example <user_ip_id>:<uuid-of-item>
# Example:
$ terraform import edgenext_scdn_user_ip_item.example 123:681c5d277ef73537f36fbdb6
```

