---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_user_ip_items"
sidebar_current: "docs-edgenext-datasource-scdn_user_ip_items"
description: |-
  # edgenext_scdn_user_ip_items
---

# edgenext_scdn_user_ip_items

# edgenext_scdn_user_ip_items

Use this data source to query items within a specific SCDN User IP List.

## Example Usage

### Query items in a user IP list

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

* `user_ip_id` - (Required, Int) User IP List ID
* `ip` - (Optional, String) Filter by IP
* `page` - (Optional, Int) Page number
* `per_page` - (Optional, Int) Items per page

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - List of user IP items
  * `format_created_at` - Created At
  * `format_updated_at` - Updated At
  * `id` - ID
  * `ip` - IP
  * `remark` - Remark
  * `user_ip_id` - User IP ID
* `total` - Total count


