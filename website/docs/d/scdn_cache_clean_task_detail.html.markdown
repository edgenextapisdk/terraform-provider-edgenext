---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_clean_task_detail"
sidebar_current: "docs-edgenext-datasource-scdn_cache_clean_task_detail"
description: |-
  Use this data source to query details of a specific SCDN cache clean task.
---

# edgenext_scdn_cache_clean_task_detail

Use this data source to query details of a specific SCDN cache clean task.

## Example Usage

### Query cache clean task detail

```hcl
data "edgenext_scdn_cache_clean_task_detail" "example" {
  task_id  = 12345
  page     = 1
  per_page = 20
}

output "task_details" {
  value = data.edgenext_scdn_cache_clean_task_detail.example.details
}
```

### Query with result filter

```hcl
data "edgenext_scdn_cache_clean_task_detail" "example" {
  task_id = 12345
  result  = 1 # success
}
```

### Query and save to file

```hcl
data "edgenext_scdn_cache_clean_task_detail" "example" {
  task_id            = 12345
  result_output_file = "task_detail.json"
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, Int) Task ID
* `page` - (Optional, Int) Page number
* `per_page` - (Optional, Int) Items per page
* `result_output_file` - (Optional, String) Used to save results to a file
* `result` - (Optional, Int) Result filter: 1-success, 2-failed, 3-executing

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `details` - Task detail list
  * `created_at` - Creation time
  * `directory` - Directory (present when this task type)
  * `message` - Execution message
  * `result` - Execution result
  * `subdomain` - Subdomain (present when this task type)
  * `updated_at` - Update time
  * `url` - URL (present when this task type)
* `total` - Total number of tasks


