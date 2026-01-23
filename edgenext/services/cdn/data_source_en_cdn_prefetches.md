Use this data source to query a list of CDN cache prefetch tasks.

Example Usage

Query CDN prefetch tasks by time range

```hcl
data "edgenext_cdn_prefetches" "example" {
  start_time = "2024-01-01"
  end_time   = "2024-01-31"
  output_file = "prefetch_tasks.json"
}
```

Query CDN prefetch tasks with pagination and URL filter

```hcl
data "edgenext_cdn_prefetches" "example" {
  start_time  = "2024-01-01"
  end_time    = "2024-01-31"
  url         = "https://example.com/static/"
  page_number = "1"
  page_size   = "50"
  output_file = "prefetch_tasks.json"
}
```
