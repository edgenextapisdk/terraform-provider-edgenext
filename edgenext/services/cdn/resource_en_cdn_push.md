Provides a resource to create and manage CDN cache push tasks.

Example Usage

Basic CDN cache push (URLs)

```hcl
resource "edgenext_cdn_push" "example" {
  type = "url"
  urls = [
    "https://example.com/static/image1.jpg",
    "https://example.com/static/image2.jpg"
  ]
}
```

CDN cache push (directories)

```hcl
resource "edgenext_cdn_push" "example" {
  type = "dir"
  urls = [
    "https://example.com/static/css/",
    "https://example.com/static/js/"
  ]
}
```

Import

CDN push tasks can be imported using the task ID:

```shell
terraform import edgenext_cdn_push.example push-task-123456
```
