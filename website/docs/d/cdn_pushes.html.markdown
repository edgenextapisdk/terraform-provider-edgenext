---
subcategory: "Content Delivery Network(CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_pushes"
sidebar_current: "docs-edgenext-datasource-cdn_pushes"
description: |-
  Use this data source to query a list of CDN cache push tasks.
---

# edgenext_cdn_pushes

Use this data source to query a list of CDN cache push tasks.

## Example Usage

### Query CDN push tasks by time range

```hcl
data "edgenext_cdn_pushes" "example" {
  start_time         = "2024-01-01"
  end_time           = "2024-01-31"
  result_output_file = "push_tasks.json"
}
```

### Query CDN push tasks with pagination and URL filter

```hcl
data "edgenext_cdn_pushes" "example" {
  start_time         = "2024-01-01"
  end_time           = "2024-01-31"
  url                = "https://example.com/static/"
  page_number        = "1"
  page_size          = "50"
  result_output_file = "push_tasks.json"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time, format: YYYY-MM-DD
* `start_time` - (Required, String) Start time, format: YYYY-MM-DD
* `page_number` - (Optional, String) Page number to retrieve, default 1
* `page_size` - (Optional, String) Page size, default 50, range 1-500
* `result_output_file` - (Optional, String) Used to save results.
* `url` - (Optional, String) URL to query refresh status for

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of successfully submitted URLs
  * `complete_time` - Completion time
  * `create_time` - Creation time
  * `id` - URL ID
  * `status` - Status
  * `type` - URL type
  * `url` - URL/Directory
* `total` - Total number of records


