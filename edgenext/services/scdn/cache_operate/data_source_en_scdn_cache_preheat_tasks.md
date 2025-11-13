Use this data source to query SCDN cache preheat tasks.

Example Usage

Query all cache preheat tasks

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

Query with filters

```hcl
data "edgenext_scdn_cache_preheat_tasks" "filtered" {
  page     = 1
  per_page = 20
  status   = "completed"
  url      = "https://example.com"
}
```

Query and save to file

```hcl
data "edgenext_scdn_cache_preheat_tasks" "all" {
  page                = 1
  per_page            = 20
  result_output_file  = "cache_preheat_tasks.json"
}
```

