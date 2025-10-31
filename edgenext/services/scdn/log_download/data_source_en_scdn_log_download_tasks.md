Use this data source to query SCDN log download tasks.

Example Usage

Query all log download tasks

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

Query with filters

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

Query and save to file

```hcl
data "edgenext_scdn_log_download_tasks" "all" {
  page                = 1
  per_page            = 20
  result_output_file  = "log_download_tasks.json"
}
```

