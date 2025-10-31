Use this data source to query details of a specific SCDN cache clean task.

Example Usage

Query cache clean task detail

```hcl
data "edgenext_scdn_cache_clean_task_detail" "example" {
  task_id = 12345
  page    = 1
  per_page = 20
}

output "task_details" {
  value = data.edgenext_scdn_cache_clean_task_detail.example.details
}
```

Query with result filter

```hcl
data "edgenext_scdn_cache_clean_task_detail" "example" {
  task_id = 12345
  result  = 1  # success
}
```

Query and save to file

```hcl
data "edgenext_scdn_cache_clean_task_detail" "example" {
  task_id            = 12345
  result_output_file = "task_detail.json"
}
```

