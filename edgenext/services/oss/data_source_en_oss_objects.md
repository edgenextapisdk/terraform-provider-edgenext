Use this data source to query a list of OSS objects in a bucket.

Example Usage

Query all objects in a bucket

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

Query objects with prefix filter

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

Query objects with delimiter to list directories

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

Query objects with pagination

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

Filter and process objects

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

List object details

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
