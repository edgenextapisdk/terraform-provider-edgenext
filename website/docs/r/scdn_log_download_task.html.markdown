---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_log_download_task"
sidebar_current: "docs-edgenext-resource-scdn_log_download_task"
description: |-
  Provides a resource to create SCDN log download tasks.
---

# edgenext_scdn_log_download_task

Provides a resource to create SCDN log download tasks.

## Example Usage

### Create log download task

```hcl
resource "edgenext_scdn_log_download_task" "example" {
  template_id = 12345
  start_time  = "2024-01-01 00:00:00"
  end_time    = "2024-01-01 23:59:59"
}
```

### Create task with search terms

```hcl
resource "edgenext_scdn_log_download_task" "example" {
  template_id = 12345
  start_time  = "2024-01-01 00:00:00"
  end_time    = "2024-01-01 23:59:59"
  search_terms = [
    {
      key   = "domain"
      value = "example.com"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `data_source` - (Required, String) Data source: ng, cc, waf
* `download_fields` - (Required, List: [`String`]) Download fields
* `end_time` - (Required, String) End time (format: YYYY-MM-DD HH:MM:SS)
* `file_type` - (Required, String) File type: xls, csv, json
* `is_use_template` - (Required, Int) Whether to use template: 0-no, 1-yes
* `start_time` - (Required, String) Start time (format: YYYY-MM-DD HH:MM:SS)
* `task_name` - (Required, String) Task name
* `lang` - (Optional, String) Language: zh_CN, en_US, default: zh_CN
* `search_terms` - (Optional, List) Search conditions
* `template_id` - (Optional, Int) Template ID (required when is_use_template is 1)

The `search_terms` object supports the following:

* `key` - (Required, String) Search key
* `value` - (Required, String) Search value

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Creation timestamp
* `download_url` - Download URL (available when task is completed)
* `id` - The ID of the log download task
* `status` - Task status: 0-not started, 1-in progress, 2-completed, 3-failed, 4-cancelled
* `task_id` - The task ID
* `updated_at` - Last update timestamp


## Import

SCDN log download tasks can be imported using the task ID:

```shell
terraform import edgenext_scdn_log_download_task.example 12345
```

