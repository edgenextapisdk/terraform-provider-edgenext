Use this data source to query detailed information of CDN cache prefetch task.

Example Usage

Query CDN prefetch task by task ID

```hcl
data "edgenext_cdn_prefetch" "example" {
  task_id = "prefetch-task-123456"
  output_file = "prefetch_task.json"
}
```

Query CDN prefetch task by URL and time range

```hcl
data "edgenext_cdn_prefetch" "example" {
  url        = "https://example.com/static/old-file.jpg"
  start_time = "2024-01-01"
  end_time   = "2024-01-31"
  output_file = "prefetch_tasks.json"
}
```
