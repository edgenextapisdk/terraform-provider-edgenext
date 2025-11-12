Provides a resource to create and manage OSS buckets with automatic CORS configuration.

Example Usage

Basic bucket creation

```hcl
resource "edgenext_oss_bucket" "example" {
  bucket = "my-example-bucket"
  acl    = "private"
}
```

Bucket with public read access

```hcl
resource "edgenext_oss_bucket" "public" {
  bucket = "my-public-bucket"
  acl    = "public-read"
}
```

Bucket with force destroy enabled

```hcl
resource "edgenext_oss_bucket" "temp" {
  bucket        = "my-temp-bucket"
  acl           = "private"
  force_destroy = true
}
```

Import

OSS buckets can be imported using the bucket name:

```shell
terraform import edgenext_oss_bucket.example my-bucket-name
```
