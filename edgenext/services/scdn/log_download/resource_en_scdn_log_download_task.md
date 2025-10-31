Provides a resource to create SCDN log download tasks.

Example Usage

Create log download task

```hcl
resource "edgenext_scdn_log_download_task" "example" {
  template_id = 12345
  start_time  = "2024-01-01 00:00:00"
  end_time    = "2024-01-01 23:59:59"
}
```

Create task with search terms

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

Import

SCDN log download tasks can be imported using the task ID:

```shell
terraform import edgenext_scdn_log_download_task.example 12345
```

