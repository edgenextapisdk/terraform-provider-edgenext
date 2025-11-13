---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_clean_tasks"
sidebar_current: "docs-edgenext-datasource-scdn_cache_clean_tasks"
description: |-
  Use this data source to query SCDN cache clean tasks.
---

# edgenext_scdn_cache_clean_tasks

Use this data source to query SCDN cache clean tasks.

## Example Usage

### Query all cache clean tasks

```hcl
data "edgenext_scdn_cache_clean_tasks" "all" {
  page     = 1
  per_page = 20
}

output "task_count" {
  value = data.edgenext_scdn_cache_clean_tasks.all.total
}

output "tasks" {
  value = data.edgenext_scdn_cache_clean_tasks.all.tasks
}
```

### Query with filters

```hcl
data "edgenext_scdn_cache_clean_tasks" "filtered" {
  page       = 1
  per_page   = 20
  start_time = "2024-01-01 00:00:00"
  end_time   = "2024-01-31 23:59:59"
  status     = "2" # completed
}
```

### Query and save to file

```hcl
data "edgenext_scdn_cache_clean_tasks" "all" {
  page               = 1
  per_page           = 20
  result_output_file = "cache_clean_tasks.json"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Optional, String) End time, format: YYYY-MM-DD HH:II:SS
* `page` - (Optional, Int) Page number
* `per_page` - (Optional, Int) Items per page
* `result_output_file` - (Optional, String) Used to save results to a file
* `start_time` - (Optional, String) Start time, format: YYYY-MM-DD HH:II:SS
* `status` - (Optional, String) Status: 1-executing, 2-completed

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tasks` - Task list
  * `created_at` - Creation time, ISO 8601 format
  * `failed` - Number of failed nodes
  * `ongoing` - Number of executing nodes
  * `operator_user_name` - Operator user name
  * `status` - Task status (can be null): Failed, Finished, etc.
  * `sub_type` - Task type: Directory, SubDomain, URL
  * `sub_user_id` - Sub user ID
  * `succeed` - Number of successful nodes
  * `task_id` - Task ID
  * `total` - Total number of nodes
  * `user_id` - User ID
* `total` - Total number of tasks


