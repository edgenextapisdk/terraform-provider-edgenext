---
subcategory: "Content Delivery Network (CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_prefetch"
sidebar_current: "docs-edgenext-datasource-cdn_prefetch"
description: |-
  Use this data source to query detailed information of CDN cache prefetch task.
---

# edgenext_cdn_prefetch

Use this data source to query detailed information of CDN cache prefetch task.

## Example Usage

### Query CDN prefetch task by task ID

```hcl
data "edgenext_cdn_prefetch" "example" {
  task_id     = "prefetch-task-123456"
  output_file = "prefetch_task.json"
}
```

### Query CDN prefetch task by URL and time range

```hcl
data "edgenext_cdn_prefetch" "example" {
  url         = "https://example.com/static/old-file.jpg"
  start_time  = "2024-01-01"
  end_time    = "2024-01-31"
  output_file = "prefetch_tasks.json"
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Task ID for querying the prefetch status of a specific task
* `output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of successfully submitted URLs
  * `complete_time` - Completion time
  * `create_time` - Creation time
  * `id` - URL ID
  * `status` - Status
  * `url` - URL
* `total` - Total number of records


