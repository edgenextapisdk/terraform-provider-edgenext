Provides a resource to create and manage CDN cache purge tasks.

Example Usage

Basic CDN cache purge

```hcl
resource "edgenext_cdn_purge" "example" {
  urls = [
    "https://example.com/static/old-image.jpg",
    "https://example.com/static/old-style.css"
  ]
}
```

CDN cache purge with multiple URLs

```hcl
resource "edgenext_cdn_purge" "example" {
  urls = [
    "https://example.com/api/data.json",
    "https://example.com/static/images/photo1.jpg",
    "https://example.com/static/css/style.css"
  ]
}
```

Import

CDN purge tasks can be imported using the task ID:

```shell
terraform import edgenext_cdn_purge.example purge-task-123456
```
