---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_preheat_tasks"
sidebar_current: "docs-edgenext-datasource-scdn_cache_preheat_tasks"
description: |-
  Use this data source to query SCDN cache preheat tasks.
---

# edgenext_scdn_cache_preheat_tasks

Use this data source to query SCDN cache preheat tasks.

## Example Usage

### Query all cache preheat tasks

```hcl
data "edgenext_scdn_cache_preheat_tasks" "all" {
  page     = 1
  per_page = 20
}

output "task_count" {
  value = data.edgenext_scdn_cache_preheat_tasks.all.total
}

output "tasks" {
  value = data.edgenext_scdn_cache_preheat_tasks.all.tasks
}
```

### Query with filters

```hcl
data "edgenext_scdn_cache_preheat_tasks" "filtered" {
  page     = 1
  per_page = 20
  status   = "completed"
  url      = "https://example.com"
}
```

### Query and save to file

```hcl
data "edgenext_scdn_cache_preheat_tasks" "all" {
  page               = 1
  per_page           = 20
  result_output_file = "cache_preheat_tasks.json"
}
```

## Argument Reference

The following arguments are supported:

* `page` - (Optional, Int) Page number
* `per_page` - (Optional, Int) Items per page
* `result_output_file` - (Optional, String) Used to save results to a file
* `status` - (Optional, String) Status filter
* `url` - (Optional, String) URL filter

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tasks` - Task list
  * `domain_id` - Domain ID
  * `id` - Task ID
  * `operator_user_name` - Operator user name
  * `status` - Status: 1-Prefetch waiting, 2-Prefetch pending, 3-Prefetch successful, 4-Prefetch failed
  * `strategy_id` - Strategy ID
  * `strategy` - Strategy
  * `sub_user_id` - Sub user ID
  * `task_id` - Task ID
  * `time_create` - Creation time
  * `time_update` - Update time
  * `total` - Total
  * `url` - URL
  * `user_id` - User ID
  * `user_name` - User name
  * `weight` - Weight
* `total` - Total number of tasks


