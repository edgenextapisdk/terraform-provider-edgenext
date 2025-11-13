---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin_groups"
sidebar_current: "docs-edgenext-datasource-scdn_origin_groups"
description: |-
  Use this data source to query a list of SCDN origin groups with optional filters.
---

# edgenext_scdn_origin_groups

Use this data source to query a list of SCDN origin groups with optional filters.

## Example Usage

### Query all origin groups

```hcl
data "edgenext_scdn_origin_groups" "all" {
  page      = 1
  page_size = 20
}

output "origin_group_count" {
  value = data.edgenext_scdn_origin_groups.all.total
}

output "origin_groups" {
  value = data.edgenext_scdn_origin_groups.all.origin_groups
}
```

### Query with name filter

```hcl
data "edgenext_scdn_origin_groups" "filtered" {
  name      = "my-group"
  page      = 1
  page_size = 20
}
```

### Query and save to file

```hcl
data "edgenext_scdn_origin_groups" "all" {
  page               = 1
  page_size          = 20
  result_output_file = "origin_groups.json"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Origin group name filter
* `page_size` - (Optional, Int) Page size
* `page` - (Optional, Int) Page number
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origin_groups` - Origin group list
  * `created_at` - Creation time
  * `id` - Origin group ID
  * `member_id` - Member ID
  * `name` - Origin group name
  * `remark` - Remark
  * `updated_at` - Update time
  * `username` - Username
* `total` - Total number of origin groups


