---
subcategory: "Object Storage Service (OSS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_oss_object"
sidebar_current: "docs-edgenext-datasource-oss_object"
description: |-
  Use this data source to query details of an OSS object.
---

# edgenext_oss_object

Use this data source to query details of an OSS object.

## Example Usage

### Query object metadata

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

### Query object with metadata

```hcl
data "edgenext_oss_object" "document" {
  bucket = "my-bucket"
  key    = "docs/report.pdf"
}

output "document_metadata" {
  value = data.edgenext_oss_object.document.metadata
}
```

### Get object URL

```hcl
data "edgenext_oss_object" "image" {
  bucket = "my-bucket"
  key    = "images/logo.png"
}

output "image_url" {
  value = data.edgenext_oss_object.image.url
}
```

### Query object HTTP headers

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

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) The name of the bucket
* `key` - (Required, String) The object key

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `acl` - ACL of the object
* `cache_control` - Specifies caching behavior
* `content_disposition` - Presentational information for the object
* `content_encoding` - What content encodings have been applied to the object
* `content_language` - Language the content is in
* `content_length` - Size of the body in bytes
* `content_type` - A standard MIME type describing the format of the object data
* `etag` - ETag generated for the object
* `expires` - Date and time at which the object is no longer cacheable
* `last_modified` - Last modified date of the object
* `metadata` - A map of metadata stored with the object
* `storage_class` - Storage class of the object
* `url` - Presigned URL of the object


