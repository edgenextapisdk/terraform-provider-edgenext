Provides a resource to copy OSS objects between buckets or within the same bucket.

Example Usage

Copy object between buckets

```hcl
resource "edgenext_oss_object_copy" "backup" {
  source_bucket = "source-bucket"
  source_key    = "original/file.txt"
  bucket        = "backup-bucket"
  key           = "backup/file.txt"
}
```

Copy object within the same bucket

```hcl
resource "edgenext_oss_object_copy" "duplicate" {
  source_bucket = "my-bucket"
  source_key    = "data/file.json"
  bucket        = "my-bucket"
  key           = "backup/file.json"
}
```

Copy object with metadata replacement

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

Copy object with ACL change

```hcl
resource "edgenext_oss_object_copy" "public_copy" {
  source_bucket = "private-bucket"
  source_key    = "data/report.pdf"
  bucket        = "public-bucket"
  key           = "reports/report.pdf"
  
  acl = "public-read"
}
```

Copy object preserving original metadata

```hcl
resource "edgenext_oss_object_copy" "preserve" {
  source_bucket = "source-bucket"
  source_key    = "data/file.pdf"
  bucket        = "dest-bucket"
  key           = "archive/file.pdf"
  
  metadata_directive = "COPY"
}
```

Copy object with all HTTP headers

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
