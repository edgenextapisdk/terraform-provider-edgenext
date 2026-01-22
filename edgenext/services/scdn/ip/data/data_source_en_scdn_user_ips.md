---
subcategory: "SCDN"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_user_ips"
sidebar_current: "docs-edgenext-datasource-scdn-user-ips"
description: |-
  Use this data source to query SCDN User IP Lists.
---

# Data Source: edgenext_scdn_user_ips

Use this data source to query SCDN User IP Lists.

## Example Usage

```hcl
data "edgenext_scdn_user_ips" "all" {
  page     = 1
  per_page = 10
}

output "ip_lists" {
  value = data.edgenext_scdn_user_ips.all.items
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional) The page number for pagination. Defaults to 1.
* `per_page` - (Optional) The number of items per page. Defaults to 10.

## Attributes Reference

The following attributes are exported:

* `total` - The total number of IP lists matching the query.
* `items` - A list of IP lists. Each item exports the following attributes:
  * `id` - The ID of the IP list.
  * `name` - The name of the IP list.
  * `remark` - The remark of the IP list.
  * `item_num` - The number of IPs in the list.
  * `created_at` - The creation time.
  * `updated_at` - The last update time.
