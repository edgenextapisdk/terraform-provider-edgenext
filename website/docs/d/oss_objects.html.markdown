---
subcategory: "Object Storage Service (OSS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_oss_objects"
sidebar_current: "docs-edgenext-datasource-oss_objects"
description: |-
  Use this data source to query a list of OSS objects in a bucket.
---

# edgenext_oss_objects

Use this data source to query a list of OSS objects in a bucket.

## Example Usage

### Query all objects in a bucket

```hcl
data "edgenext_oss_objects" "all" {
  bucket      = "my-bucket"
  output_file = "all_objects.json"
}

output "object_count" {
  value = length(data.edgenext_oss_objects.all.objects)
}

output "object_keys" {
  value = data.edgenext_oss_objects.all.keys
}
```

### Query objects with prefix filter

```hcl
data "edgenext_oss_objects" "logs" {
  bucket      = "my-bucket"
  prefix      = "logs/2024/"
  output_file = "log_files.json"
}

output "log_files" {
  value = data.edgenext_oss_objects.logs.objects[*].key
}
```

### Query objects with delimiter to list directories

```hcl
data "edgenext_oss_objects" "dirs" {
  bucket    = "my-bucket"
  prefix    = "data/"
  delimiter = "/"
}

output "subdirectories" {
  value = data.edgenext_oss_objects.dirs.common_prefixes
}
```

### Query objects with pagination

```hcl
data "edgenext_oss_objects" "limited" {
  bucket   = "my-bucket"
  prefix   = "archive/"
  max_keys = 100
}

output "first_100_objects" {
  value = data.edgenext_oss_objects.limited.objects[*].key
}
```

### Filter and process objects

```hcl
data "edgenext_oss_objects" "images" {
  bucket = "my-bucket"
  prefix = "images/"
}

locals {
  png_files = [
    for obj in data.edgenext_oss_objects.images.objects :
    obj.key if endswith(obj.key, ".png")
  ]

  total_size = sum([
    for obj in data.edgenext_oss_objects.images.objects : obj.size
  ])
}

output "png_count" {
  value = length(local.png_files)
}

output "total_image_size_mb" {
  value = local.total_size / 1024 / 1024
}
```

### List object details

```hcl
data "edgenext_oss_objects" "files" {
  bucket = "my-bucket"
  prefix = "data/"
}

output "file_details" {
  value = [
    for obj in data.edgenext_oss_objects.files.objects : {
      key           = obj.key
      size_kb       = obj.size / 1024
      etag          = obj.etag
      last_modified = obj.last_modified
      storage_class = obj.storage_class
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) The name of the bucket
* `delimiter` - (Optional, String) A delimiter is a character you use to group keys
* `max_keys` - (Optional, Int) Sets the maximum number of keys returned in the response
* `output_file` - (Optional, String) File name where to save data source results (after running `terraform plan`)
* `prefix` - (Optional, String) Limits the response to keys that begin with the specified prefix

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `common_prefixes` - A list of common prefixes (if delimiter is specified)
* `keys` - A list of object keys
* `objects` - A list of objects
  * `etag` - The ETag of the object
  * `key` - The object key
  * `last_modified` - The last modified date of the object
  * `size` - The size of the object in bytes
  * `storage_class` - The storage class of the object


