---
subcategory: "Content Delivery Network (CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_purge"
sidebar_current: "docs-edgenext-datasource-cdn_purge"
description: |-
  Use this data source to query detailed information of CDN cache purge task.
---

# edgenext_cdn_purge

Use this data source to query detailed information of CDN cache purge task.

## Example Usage

### Query CDN purge task by task ID

```hcl
data "edgenext_cdn_purge" "example" {
  task_id     = "purge-task-123456"
  output_file = "purge_task.json"
}
```

### Query CDN purge task by URL and time range

```hcl
data "edgenext_cdn_purge" "example" {
  url         = "https://example.com/static/image.jpg"
  start_time  = "2024-01-01"
  end_time    = "2024-01-31"
  output_file = "purge_tasks.json"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Optional, String) End time, format: YYYY-MM-DD, used together with start_time
* `output_file` - (Optional, String) Used to save results.
* `page_number` - (Optional, String) Page number to retrieve, default 1
* `page_size` - (Optional, String) Page size, default 50, range 1-500
* `start_time` - (Optional, String) Start time, format: YYYY-MM-DD, used together with end_time
* `task_id` - (Optional, String) Task ID for querying the purge status of a specific task
* `url` - (Optional, String) URL

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


