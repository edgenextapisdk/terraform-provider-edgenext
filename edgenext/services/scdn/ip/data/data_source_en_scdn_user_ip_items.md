---
subcategory: "SCDN"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_user_ip_items"
sidebar_current: "docs-edgenext-datasource-scdn-user-ip-items"
description: |-
  Use this data source to query items in an SCDN User IP List.
---

# Data Source: edgenext_scdn_user_ip_items

Use this data source to query items within a specific SCDN User IP List.

## Example Usage

```hcl
data "edgenext_scdn_user_ip_items" "items" {
  user_ip_id = 123
  ip         = "192.168.1.1" # Optional filter
}

output "ip_items" {
  value = data.edgenext_scdn_user_ip_items.items
}
```

## Argument Reference

The following arguments are supported:

* `user_ip_id` - (Required) The ID of the User IP List to query.
* `page` - (Optional) The page number for pagination. Defaults to 1.
* `per_page` - (Optional) The number of items per page. Defaults to 10.
* `ip` - (Optional) Filter by IP address.

## Attributes Reference

The following attributes are exported:

* `total` - The total number of IP items matching the query.
* `items` - A list of IP items. Each item exports the following attributes:
  * `id` - The unique ID (UUID) of the IP item.
  * `ip` - The IP address.
  * `remark` - The remark for the item.
  * `user_ip_id` - The ID of the User IP List.
  * `format_created_at` - The formatted creation time.
  * `format_updated_at` - The formatted last update time.
