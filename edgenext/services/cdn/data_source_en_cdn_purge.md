Use this data source to query detailed information of CDN cache purge task.

Example Usage

Query CDN purge task by task ID

```hcl
data "edgenext_cdn_purge" "example" {
  task_id = "purge-task-123456"
  output_file = "purge_task.json"
}
```

Query CDN purge task by URL and time range

```hcl
data "edgenext_cdn_purge" "example" {
  url        = "https://example.com/static/image.jpg"
  start_time = "2024-01-01"
  end_time   = "2024-01-31"
  output_file = "purge_tasks.json"
}
```
