Provides a resource to create and manage CDN cache prefetch tasks.

Example Usage

Basic CDN cache prefetch

```hcl
resource "edgenext_cdn_prefetch" "example" {
  urls = [
    "https://example.com/static/old-image.jpg",
    "https://example.com/static/old-style.css"
  ]
}
```

CDN cache prefetch with multiple URLs

```hcl
resource "edgenext_cdn_prefetch" "example" {
  urls = [
    "https://example.com/api/data.json",
    "https://example.com/static/images/photo1.jpg",
    "https://example.com/static/css/style.css"
  ]
}
```

Import

CDN prefetch tasks can be imported using the task ID:

```shell
terraform import edgenext_cdn_prefetch.example prefetch-task-123456
```
