Use this data source to query a list of CDN cache purge tasks.

Example Usage

Query CDN purge tasks by time range

```hcl
data "edgenext_cdn_purges" "example" {
  start_time = "2024-01-01"
  end_time   = "2024-01-31"
  result_output_file = "purge_tasks.json"
}
```

Query CDN purge tasks with pagination and URL filter

```hcl
data "edgenext_cdn_purges" "example" {
  start_time  = "2024-01-01"
  end_time    = "2024-01-31"
  url         = "https://example.com/static/"
  page_number = "1"
  page_size   = "50"
  result_output_file = "purge_tasks.json"
}
```
