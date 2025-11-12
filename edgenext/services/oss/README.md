# EdgeNext OSS Services

This package provides Terraform resources and data sources for managing EdgeNext OSS (Object Storage Service), which is S3-compatible object storage.

## Resources

### OSS Bucket
- **Resource**: `edgenext_oss_bucket` (`ResourceOSSBucket`)
- **File**: `resource_en_oss_bucket.go`
- **Description**: Manage OSS buckets with ACL control and automatic CORS configuration

### OSS Object
- **Resource**: `edgenext_oss_object` (`ResourceOSSObject`)
- **File**: `resource_en_oss_object.go`
- **Description**: Manage OSS objects with file/content upload, metadata, cache control, and HTTP headers

### OSS Object Copy
- **Resource**: `edgenext_oss_object_copy` (`ResourceOSSObjectCopy`)
- **File**: `resource_en_oss_object_copy.go`
- **Description**: Copy objects between buckets or within the same bucket with optional metadata updates

## Data Sources

### OSS Buckets List
- **Data Source**: `edgenext_oss_buckets` (`DataSourceOSSBuckets`)
- **File**: `data_source_en_oss_buckets.go`
- **Description**: Query a list of OSS buckets with optional filtering by prefix

### OSS Object
- **Data Source**: `edgenext_oss_object` (`DataSourceOSSObject`)
- **File**: `data_source_en_oss_object.go`
- **Description**: Query OSS object details including metadata, content type, ETag, and presigned URL

### OSS Objects List
- **Data Source**: `edgenext_oss_objects` (`DataSourceOSSObjects`)
- **File**: `data_source_en_oss_objects.go`
- **Description**: Query a list of OSS objects in a bucket with prefix/delimiter filtering

## File Structure

```
edgenext/services/oss/
├── README.md                           # This documentation
├── resource_en_oss_bucket.go           # OSS bucket resource implementation
├── resource_en_oss_bucket.md           # OSS bucket resource documentation
├── resource_en_oss_object.go           # OSS object resource implementation
├── resource_en_oss_object.md           # OSS object resource documentation
├── resource_en_oss_object_copy.go      # OSS object copy resource implementation
├── resource_en_oss_object_copy.md      # OSS object copy resource documentation
├── data_source_en_oss_buckets.go       # OSS buckets list data source implementation
├── data_source_en_oss_buckets.md       # OSS buckets list data source documentation
├── data_source_en_oss_object.go        # OSS object data source implementation
├── data_source_en_oss_object.md        # OSS object data source documentation
├── data_source_en_oss_objects.go       # OSS objects list data source implementation
└── data_source_en_oss_objects.md       # OSS objects list data source documentation
```

## Usage Examples

### Basic Bucket Creation

```hcl
resource "edgenext_oss_bucket" "example" {
  bucket = "my-example-bucket"
  acl    = "private"
}
```

### Bucket with Public Access

```hcl
resource "edgenext_oss_bucket" "public" {
  bucket = "my-public-bucket"
  acl    = "public-read"
}
```

### Bucket with Force Destroy

```hcl
resource "edgenext_oss_bucket" "temp" {
  bucket        = "my-temp-bucket"
  acl           = "private"
  force_destroy = true
}
```

### Upload Object from File

```hcl
resource "edgenext_oss_object" "file" {
  bucket       = edgenext_oss_bucket.example.bucket
  key          = "path/to/file.txt"
  source       = "./local/file.txt"
  content_type = "text/plain"
  acl          = "private"
}
```

### Upload Object from Content

```hcl
resource "edgenext_oss_object" "config" {
  bucket  = edgenext_oss_bucket.example.bucket
  key     = "config.json"
  content = jsonencode({
    setting1 = "value1"
    setting2 = "value2"
  })
  content_type = "application/json"
}
```

### Upload Object with Metadata and Cache Control

```hcl
resource "edgenext_oss_object" "asset" {
  bucket              = edgenext_oss_bucket.example.bucket
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

### Upload Multiple Objects from Directory

```hcl
locals {
  files = fileset("./static", "**")
}

resource "edgenext_oss_object" "static_files" {
  for_each = local.files

  bucket       = edgenext_oss_bucket.example.bucket
  key          = "static/${each.value}"
  source       = "./static/${each.value}"
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

### Copy Object Between Buckets

```hcl
resource "edgenext_oss_object_copy" "backup" {
  source_bucket = "source-bucket"
  source_key    = "original/file.txt"
  bucket        = "backup-bucket"
  key           = "backup/file.txt"
  acl           = "private"
}
```

### Copy Object with Metadata Replacement

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
  }
}
```

### Query Buckets

```hcl
data "edgenext_oss_buckets" "all" {
  bucket_prefix = "my-app-"
  max_buckets   = 100
  output_file   = "buckets.json"
}

output "bucket_names" {
  value = data.edgenext_oss_buckets.all.buckets[*].name
}
```

### Query Object Details

```hcl
data "edgenext_oss_object" "config" {
  bucket = "my-bucket"
  key    = "config.json"
}

output "object_info" {
  value = {
    etag          = data.edgenext_oss_object.config.etag
    size          = data.edgenext_oss_object.config.content_length
    content_type  = data.edgenext_oss_object.config.content_type
    last_modified = data.edgenext_oss_object.config.last_modified
    url           = data.edgenext_oss_object.config.url
  }
}
```

### Query Objects in Bucket

```hcl
data "edgenext_oss_objects" "logs" {
  bucket = "my-bucket"
  prefix = "logs/2024/"
  output_file = "log_files.json"
}

output "log_files" {
  value = data.edgenext_oss_objects.logs.objects[*].key
}

output "object_keys" {
  value = data.edgenext_oss_objects.logs.keys
}
```

## Key Features

### Bucket Features
- **ACL Management**: Control bucket access with canned ACLs (private, public-read, public-read-write, authenticated-read)
- **Automatic CORS**: Automatically configures CORS rules allowing all origins and methods
- **Force Destroy**: Option to automatically delete all objects when destroying a bucket

### Object Features
- **Multiple Upload Methods**: 
  - Upload from local file (`source`)
  - Upload from inline content (`content`)
- **HTTP Headers**: Full support for cache-control, content-type, content-encoding, content-disposition, expires
- **Metadata Management**: Set custom metadata key-value pairs
- **ACL Control**: Per-object access control
- **ETag Verification**: Automatic ETag generation for integrity verification
- **Import Support**: Import existing objects using bucket/key format

### Object Copy Features
- **Cross-Bucket Copy**: Copy objects between different buckets
- **Same-Bucket Copy**: Duplicate objects within the same bucket
- **Metadata Directive**: Choose to COPY (preserve) or REPLACE (update) metadata
- **Metadata Override**: Update metadata, HTTP headers, and ACL during copy
- **ACL Control**: Set ACL on the copied object

### Data Source Features
- **Bucket Listing**: List all buckets with optional prefix filtering
- **Bucket ACL**: Query bucket ACL information
- **Object Details**: Query complete object metadata and HTTP headers
- **Presigned URL**: Get presigned URL for object access
- **Object Listing**: List objects with prefix and delimiter support
- **Common Prefixes**: Discover directory structure using delimiter

## Best Practices

1. **Bucket Naming**
   - Use lowercase letters, numbers, and hyphens
   - Keep names between 3-63 characters
   - Make names globally unique
   - Use meaningful prefixes for organization

2. **Security**
   - Always use `private` ACL unless public access is required
   - CORS is automatically configured (all origins, all methods)
   - Regularly audit bucket permissions

3. **Object Management**
   - Use `force_destroy = true` for temporary/development buckets
   - Set appropriate content types for better browser handling
   - Use metadata for version tracking and management
   - Use ETags for cache validation

4. **Performance**
   - Use key prefixes to organize objects logically
   - Consider using `for_each` for batch operations
   - Set cache-control headers for static content
   - Use delimiter in queries to navigate directory structures

5. **Cost Optimization**
   - Delete unnecessary objects promptly
   - Use `force_destroy` carefully in production
   - Monitor bucket usage regularly

## Import Examples

### Import Bucket
```bash
terraform import edgenext_oss_bucket.example my-bucket-name
```

### Import Object
```bash
terraform import edgenext_oss_object.example my-bucket-name/path/to/object.txt
```

## Related Services

- [CDN Services](../cdn/README.md) - Content Delivery Network for OSS objects
- [SSL Services](../ssl/README.md) - SSL certificates for custom domains

## Notes

- OSS service is S3-compatible and uses AWS SDK v2 for Go
- All operations require proper authentication credentials configured in the provider
- Bucket names must be globally unique across the EdgeNext platform
- CORS configuration is automatically applied to all buckets with permissive rules
- Objects support standard HTTP headers for cache control and content negotiation
