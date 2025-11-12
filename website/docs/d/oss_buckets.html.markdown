---
subcategory: "Object Storage Service (OSS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_oss_buckets"
sidebar_current: "docs-edgenext-datasource-oss_buckets"
description: |-
  Use this data source to query a list of OSS buckets.
---

# edgenext_oss_buckets

Use this data source to query a list of OSS buckets.

## Example Usage

### Query all buckets

```hcl
data "edgenext_oss_buckets" "all" {
  output_file = "buckets.json"
}

output "bucket_names" {
  value = data.edgenext_oss_buckets.all.buckets[*].name
}
```

### Query buckets with prefix filter

```hcl
data "edgenext_oss_buckets" "app_buckets" {
  bucket_prefix = "myapp-"
  max_buckets   = 100
  output_file   = "app_buckets.json"
}

output "app_bucket_count" {
  value = length(data.edgenext_oss_buckets.app_buckets.buckets)
}
```

### Query buckets and display details

```hcl
data "edgenext_oss_buckets" "prod" {
  bucket_prefix = "prod-"
}

output "bucket_details" {
  value = [
    for bucket in data.edgenext_oss_buckets.prod.buckets : {
      name          = bucket.name
      creation_date = bucket.creation_date
      acl           = bucket.acl
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `bucket_prefix` - (Optional, String) A prefix string to filter results by bucket name
* `max_buckets` - (Optional, Int) The maximum number of buckets to return.
* `output_file` - (Optional, String) File name where to save data source results (after running `terraform plan`)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `buckets` - A list of buckets
  * `acl` - The access control list (ACL) of the bucket
  * `creation_date` - The creation date of the bucket
  * `name` - The name of the bucket


