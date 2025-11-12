---
subcategory: "Object Storage Service (OSS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_oss_bucket"
sidebar_current: "docs-edgenext-resource-oss_bucket"
description: |-
  Provides a resource to create and manage OSS buckets with automatic CORS configuration.
---

# edgenext_oss_bucket

Provides a resource to create and manage OSS buckets with automatic CORS configuration.

## Example Usage

### Basic bucket creation

```hcl
resource "edgenext_oss_bucket" "example" {
  bucket = "my-example-bucket"
  acl    = "private"
}
```

### Bucket with public read access

```hcl
resource "edgenext_oss_bucket" "public" {
  bucket = "my-public-bucket"
  acl    = "public-read"
}
```

### Bucket with force destroy enabled

```hcl
resource "edgenext_oss_bucket" "temp" {
  bucket        = "my-temp-bucket"
  acl           = "private"
  force_destroy = true
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of the bucket (3-63 characters)
* `acl` - (Optional, String) The canned ACL to apply to the bucket (private, public-read, public-read-write, authenticated-read)
* `force_destroy` - (Optional, Bool) A boolean that indicates all objects should be deleted from the bucket so that the bucket can be destroyed without error

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

OSS buckets can be imported using the bucket name:

```shell
terraform import edgenext_oss_bucket.example my-bucket-name
```

