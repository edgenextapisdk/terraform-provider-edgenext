Use this data source to query detailed information of CDN cache push task.

Example Usage

Query CDN push task by task ID

```hcl
data "edgenext_cdn_push" "example" {
  task_id = "push-task-123456"
  output_file = "push_task.json"
}
```

Query CDN push task by URL and time range

```hcl
data "edgenext_cdn_push" "example" {
  url        = "https://example.com/static/image.jpg"
  start_time = "2024-01-01"
  end_time   = "2024-01-31"
  output_file = "push_tasks.json"
}
```
