Use this data source to query a list of OSS buckets.

Example Usage

Query all buckets

```hcl
data "edgenext_oss_buckets" "all" {
  output_file = "buckets.json"
}

output "bucket_names" {
  value = data.edgenext_oss_buckets.all.buckets[*].name
}
```

Query buckets with prefix filter

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

Query buckets and display details

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
