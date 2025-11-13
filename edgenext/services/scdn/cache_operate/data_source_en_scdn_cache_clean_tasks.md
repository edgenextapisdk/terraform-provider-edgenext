Use this data source to query SCDN cache clean tasks.

Example Usage

Query all cache clean tasks

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

Query with filters

```hcl
data "edgenext_scdn_cache_clean_tasks" "filtered" {
  page      = 1
  per_page  = 20
  start_time = "2024-01-01 00:00:00"
  end_time   = "2024-01-31 23:59:59"
  status     = "2"  # completed
}
```

Query and save to file

```hcl
data "edgenext_scdn_cache_clean_tasks" "all" {
  page                = 1
  per_page            = 20
  result_output_file  = "cache_clean_tasks.json"
}
```

