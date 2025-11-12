---
subcategory: "Object Storage Service (OSS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_oss_object"
sidebar_current: "docs-edgenext-resource-oss_object"
description: |-
  Provides a resource to create and manage OSS objects.
---

# edgenext_oss_object

Provides a resource to create and manage OSS objects.

## Example Usage

### Upload object from local file

```hcl
resource "edgenext_oss_object" "file" {
  bucket       = "my-bucket"
  key          = "path/to/file.txt"
  source       = "./local/file.txt"
  content_type = "text/plain"
  acl          = "private"
}
```

### Upload object from inline content

```hcl
resource "edgenext_oss_object" "config" {
  bucket = "my-bucket"
  key    = "config.json"
  content = jsonencode({
    setting1 = "value1"
    setting2 = "value2"
  })
  content_type = "application/json"
}
```

### Upload object with metadata and cache control

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

### Upload object with expiration date

```hcl
resource "edgenext_oss_object" "temporary" {
  bucket       = "my-bucket"
  key          = "temp/data.json"
  source       = "./data.json"
  content_type = "application/json"
  expires      = "2025-12-31T23:59:59Z"
}
```

### Upload multiple files using for_each

```hcl
locals {
  files = fileset("./static", "**")
}

resource "edgenext_oss_object" "static_files" {
  for_each = local.files

  bucket = "my-bucket"
  key    = "static/${each.value}"
  source = "./static/${each.value}"
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

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of the bucket to put the object in
* `key` - (Required, String, ForceNew) The name of the object once it is in the bucket
* `acl` - (Optional, String) The canned ACL to apply to the object (private, public-read, public-read-write, authenticated-read)
* `cache_control` - (Optional, String) Specifies caching behavior along the request/reply chain
* `content_disposition` - (Optional, String) Specifies presentational information for the object
* `content_encoding` - (Optional, String) Specifies what content encodings have been applied to the object
* `content_type` - (Optional, String) A standard MIME type describing the format of the object data
* `content` - (Optional, String) Literal string value to use as the object content, conflicts with source
* `expires` - (Optional, String) The date and time at which the object is no longer cacheable
* `metadata` - (Optional, Map) A map of metadata to store with the object
* `source` - (Optional, String) The path to a file that will be read and uploaded as raw bytes for the object content, conflicts with content

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `etag` - The ETag generated for the object
* `last_modified` - The last modified date of the object
* `size` - The size of the object in bytes


## Import

OSS objects can be imported using the bucket name and object key separated by a forward slash:

```shell
terraform import edgenext_oss_object.example my-bucket/path/to/object.txt
```

