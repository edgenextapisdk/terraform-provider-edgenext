Provides a resource to create and manage OSS objects.

Example Usage

Upload object from local file

```hcl
resource "edgenext_oss_object" "file" {
  bucket       = "my-bucket"
  key          = "path/to/file.txt"
  source       = "./local/file.txt"
  content_type = "text/plain"
  acl          = "private"
}
```

Upload object from inline content

```hcl
resource "edgenext_oss_object" "config" {
  bucket  = "my-bucket"
  key     = "config.json"
  content = jsonencode({
    setting1 = "value1"
    setting2 = "value2"
  })
  content_type = "application/json"
}
```

Upload object with metadata and cache control

```hcl
resource "edgenext_oss_object" "asset" {
  bucket              = "my-bucket"
  key                 = "static/app.js"
  source              = "./dist/app.js"
  content_type        = "application/javascript"
  cache_control       = "public, max-age=31536000"
  content_disposition = "inline"
  content_encoding    = "gzip"
  acl                 = "public-read"

  metadata = {
    version     = "1.0.0"
    environment = "production"
  }
}
```

Upload object with expiration date

```hcl
resource "edgenext_oss_object" "temporary" {
  bucket       = "my-bucket"
  key          = "temp/data.json"
  source       = "./data.json"
  content_type = "application/json"
  expires      = "2025-12-31T23:59:59Z"
}
```

Upload multiple files using for_each

```hcl
locals {
  files = fileset("./static", "**")
}

resource "edgenext_oss_object" "static_files" {
  for_each = local.files

  bucket       = "my-bucket"
  key          = "static/${each.value}"
  source       = "./static/${each.value}"
  content_type = lookup(
    {
      "html" = "text/html"
      "css"  = "text/css"
      "js"   = "application/javascript"
      "json" = "application/json"
      "png"  = "image/png"
      "jpg"  = "image/jpeg"
    },
    element(split(".", each.value), length(split(".", each.value)) - 1),
    "application/octet-stream"
  )
}
```

Import

OSS objects can be imported using the bucket name and object key separated by a forward slash:

```shell
terraform import edgenext_oss_object.example my-bucket/path/to/object.txt
```
