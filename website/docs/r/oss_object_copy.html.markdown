---
subcategory: "Object Storage Service (OSS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_oss_object_copy"
sidebar_current: "docs-edgenext-resource-oss_object_copy"
description: |-
  Provides a resource to copy OSS objects between buckets or within the same bucket.
---

# edgenext_oss_object_copy

Provides a resource to copy OSS objects between buckets or within the same bucket.

## Example Usage

### Copy object between buckets

```hcl
resource "edgenext_oss_object_copy" "backup" {
  source_bucket = "source-bucket"
  source_key    = "original/file.txt"
  bucket        = "backup-bucket"
  key           = "backup/file.txt"
}
```

### Copy object within the same bucket

```hcl
resource "edgenext_oss_object_copy" "duplicate" {
  source_bucket = "my-bucket"
  source_key    = "data/file.json"
  bucket        = "my-bucket"
  key           = "backup/file.json"
}
```

### Copy object with metadata replacement

```hcl
resource "edgenext_oss_object_copy" "with_metadata" {
  source_bucket = "my-bucket"
  source_key    = "original/file.txt"
  bucket        = "my-bucket"
  key           = "updated/file.txt"

  metadata_directive = "REPLACE"
  content_type       = "text/plain"
  cache_control      = "max-age=3600"

  metadata = {
    version     = "2.0"
    updated-by  = "terraform"
    update-date = "2024-01-01"
  }
}
```

### Copy object with ACL change

```hcl
resource "edgenext_oss_object_copy" "public_copy" {
  source_bucket = "private-bucket"
  source_key    = "data/report.pdf"
  bucket        = "public-bucket"
  key           = "reports/report.pdf"

  acl = "public-read"
}
```

### Copy object preserving original metadata

```hcl
resource "edgenext_oss_object_copy" "preserve" {
  source_bucket = "source-bucket"
  source_key    = "data/file.pdf"
  bucket        = "dest-bucket"
  key           = "archive/file.pdf"

  metadata_directive = "COPY"
}
```

### Copy object with all HTTP headers

```hcl
resource "edgenext_oss_object_copy" "full_headers" {
  source_bucket = "source-bucket"
  source_key    = "data/document.pdf"
  bucket        = "dest-bucket"
  key           = "docs/document.pdf"

  metadata_directive  = "REPLACE"
  content_type        = "application/pdf"
  content_encoding    = "identity"
  content_disposition = "attachment; filename=document.pdf"
  cache_control       = "no-cache"
  expires             = "2025-12-31T23:59:59Z"
  acl                 = "private"

  metadata = {
    document-type = "report"
    department    = "sales"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of the destination bucket
* `key` - (Required, String, ForceNew) The key of the destination object
* `source_bucket` - (Required, String, ForceNew) The name of the source bucket
* `source_key` - (Required, String, ForceNew) The key of the source object
* `acl` - (Optional, String) The canned ACL to apply to the object (private, public-read, public-read-write, authenticated-read)
* `cache_control` - (Optional, String) Specifies caching behavior along the request/reply chain, only used when metadata_directive is REPLACE
* `content_disposition` - (Optional, String) Specifies presentational information for the object, only used when metadata_directive is REPLACE
* `content_encoding` - (Optional, String) Specifies what content encodings have been applied to the object, only used when metadata_directive is REPLACE
* `content_type` - (Optional, String) A standard MIME type describing the format of the object data, only used when metadata_directive is REPLACE
* `expires` - (Optional, String) The date and time at which the object is no longer cacheable, only used when metadata_directive is REPLACE
* `metadata_directive` - (Optional, String) Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request (COPY or REPLACE)
* `metadata` - (Optional, Map) A map of metadata to store with the object, only used when metadata_directive is REPLACE

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `etag` - The ETag generated for the object
* `last_modified` - The last modified date of the object
* `size` - The size of the object in bytes


