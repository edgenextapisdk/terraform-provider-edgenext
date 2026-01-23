Provides a resource to create and manage CDN cache purge tasks.

Example Usage

Basic CDN cache purge (URLs)

```hcl
resource "edgenext_cdn_purge" "example" {
  type = "url"
  urls = [
    "https://example.com/static/image1.jpg",
    "https://example.com/static/image2.jpg"
  ]
}
```

CDN cache purge (directories)

```hcl
resource "edgenext_cdn_purge" "example" {
  type = "dir"
  urls = [
    "https://example.com/static/css/",
    "https://example.com/static/js/"
  ]
}
```

Import

CDN purge tasks can be imported using the task ID:

```shell
terraform import edgenext_cdn_purge.example purge-task-123456
```
