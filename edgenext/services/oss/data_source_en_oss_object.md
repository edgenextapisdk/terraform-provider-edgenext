Use this data source to query details of an OSS object.

Example Usage

Query object metadata

```hcl
data "edgenext_oss_object" "config" {
  bucket = "my-bucket"
  key    = "config/app.json"
}

output "object_info" {
  value = {
    etag          = data.edgenext_oss_object.config.etag
    size          = data.edgenext_oss_object.config.content_length
    content_type  = data.edgenext_oss_object.config.content_type
    last_modified = data.edgenext_oss_object.config.last_modified
  }
}
```

Query object with metadata

```hcl
data "edgenext_oss_object" "document" {
  bucket = "my-bucket"
  key    = "docs/report.pdf"
}

output "document_metadata" {
  value = data.edgenext_oss_object.document.metadata
}
```

Get object URL

```hcl
data "edgenext_oss_object" "image" {
  bucket = "my-bucket"
  key    = "images/logo.png"
}

output "image_url" {
  value = data.edgenext_oss_object.image.url
}
```

Query object HTTP headers

```hcl
data "edgenext_oss_object" "asset" {
  bucket = "my-bucket"
  key    = "static/app.js"
}

output "http_headers" {
  value = {
    cache_control       = data.edgenext_oss_object.asset.cache_control
    content_encoding    = data.edgenext_oss_object.asset.content_encoding
    content_disposition = data.edgenext_oss_object.asset.content_disposition
    expires             = data.edgenext_oss_object.asset.expires
  }
}
```
