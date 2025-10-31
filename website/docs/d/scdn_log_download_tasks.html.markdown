---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_log_download_tasks"
sidebar_current: "docs-edgenext-datasource-scdn_log_download_tasks"
description: |-
  Use this data source to query SCDN log download tasks.
---

# edgenext_scdn_log_download_tasks

Use this data source to query SCDN log download tasks.

## Example Usage

### Query all log download tasks

```hcl
data "edgenext_scdn_log_download_tasks" "all" {
  page     = 1
  per_page = 20
}

output "task_count" {
  value = data.edgenext_scdn_log_download_tasks.all.total
}

output "tasks" {
  value = data.edgenext_scdn_log_download_tasks.all.tasks
}
```

### Query with filters

```hcl
data "edgenext_scdn_log_download_tasks" "filtered" {
  page     = 1
  per_page = 20
  search_terms = [
    {
      key   = "domain"
      value = "example.com"
    }
  ]
}
```

### Query and save to file

```hcl
data "edgenext_scdn_log_download_tasks" "all" {
  page               = 1
  per_page           = 20
  result_output_file = "log_download_tasks.json"
}
```

## Argument Reference

The following arguments are supported:

* `data_source` - (Optional, String) Data source: ng, cc, waf
* `file_type` - (Optional, String) File type: xls, csv, json
* `page` - (Optional, Int) Page number
* `per_page` - (Optional, Int) Items per page
* `result_output_file` - (Optional, String) Used to save results to a file
* `status` - (Optional, Int) Task status: 0-not started, 1-in progress, 2-completed, 3-failed, 4-cancelled
* `task_name` - (Optional, String) Task name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tasks` - Task list
  * `created_at` - Creation timestamp
  * `data_source` - Data source
  * `download_fields` - Download fields
  * `download_url` - Download URL
  * `end_time` - End time
  * `file_type` - File type
  * `is_use_template` - Whether to use template
  * `member_id` - Member ID
  * `search_terms` - Search conditions
    * `key` - Search key
    * `value` - Search value
  * `start_time` - Start time
  * `status` - Task status
  * `task_id` - Task ID
  * `task_name` - Task name
  * `template_id` - Template ID
  * `updated_at` - Last update timestamp
* `total` - Total number of tasks


