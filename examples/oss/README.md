# EdgeNext OSS Terraform Example

This example demonstrates how to use the EdgeNext Terraform Provider to manage OSS (Object Storage Service) resources.

## Features

This example demonstrates the following features:

### Resources

1. **Bucket Management**
   - Create private buckets for different purposes (data, backups)
   - Configure ACL permissions (private, public-read)
   - Enable force_destroy for development environments
   - Import existing buckets

2. **Object Management**
   - Upload JSON configuration files with dynamic content
   - Set object properties (Content-Type, Cache-Control, Content-Disposition, Expires)
   - Configure object ACL (private, public-read)
   - Add custom metadata
   - Batch create objects using for_each
   - Upload files from local directory

3. **Object Copy**
   - Copy objects within buckets
   - Replace metadata during copy
   - Change ACL permissions when copying

### Data Sources

1. **edgenext_oss_buckets** - List all buckets with optional prefix filter
2. **edgenext_oss_objects** - List objects in a bucket with prefix/delimiter support
3. **edgenext_oss_object** - Read content and metadata of a specific object

## Prerequisites

1. Install Terraform >= 0.13
2. Have an EdgeNext account and API credentials
3. Configure OSS service access permissions

## Usage Steps

### 1. Prepare Configuration Files

```bash
# Copy example configuration
cp terraform.tfvars.example terraform.tfvars

# Edit terraform.tfvars and fill in your credentials
vim terraform.tfvars
```

**Required variables in `terraform.tfvars`:**
```hcl
access_key  = "your-access-key"
secret_key  = "your-secret-key"
endpoint    = "your-edgenext-api-endpoint-here"  # Optional, has default
region      = "your-region-here"                          # Optional, has default
environment = "dev"                                        # Optional, default is "dev"
```

### 2. (Optional) Prepare Local Files Directory

If you want to test the batch upload feature, create the files directory:

```bash
mkdir -p files
# Place your files in the files directory
# The example will upload all files in this directory
```

### 3. Initialize Terraform

```bash
terraform init
```

### 4. Review Execution Plan

```bash
terraform plan
```

This will show the resources to be created:
- 2 Buckets (myapp-data-dev, myapp-backups-dev)
- 4+ Objects (app.json, redis.conf, nginx.conf, and any files in ./files/)
- 1 Object Copy (app-v2.json)
- 3 Data Sources (list buckets, list objects, read object)

### 5. Apply Configuration

```bash
terraform apply
```

After confirmation, enter `yes` to start creating resources.

### 6. View Outputs

```bash
terraform output
```

You will see:
- Created bucket details (names, ACL, force_destroy settings)

### 7. Verify Resources

The example creates data sources output to JSON files:
- `buckets.json` - List of all buckets with prefix "myapp"
- `configs.json` - List of objects in the config/ directory

### 8. Import Existing Resources

The example includes commented-out import resource blocks:

```hcl
# Import existing bucket
resource "edgenext_oss_bucket" "import_bucket" {
  # bucket = "test"
  # acl    = "public-read"
  # force_destroy = var.environment == "dev" ? true : false
}

# Import existing object
resource "edgenext_oss_object" "import_object" {
  # bucket        = "myapp-data-dev"
  # content_type  = "image/jpeg"
  # key           = "1739760288.jpg"
}
```

To import existing resources:
1. Uncomment and configure the resource block
2. Run: `terraform import edgenext_oss_bucket.import_bucket existing-bucket-name`
3. Run: `terraform import edgenext_oss_object.import_object bucket-name/object-key`

### 9. Clean Up Resources

```bash
terraform destroy
```

After confirmation, enter `yes` to delete all created resources.

## Directory Structure

```
examples/oss/
├── README.md                    # This file
├── main.tf                      # Main configuration file
├── terraform.tfvars.example     # Configuration example
├── terraform.tfvars             # Your actual configuration (not in git)
├── files/                       # (Optional) Local files directory
│   ├── logo.png                 # Example files
│   ├── buckets.json
│   └── configs.json
├── buckets.json                 # (Generated) List of buckets
├── configs.json                 # (Generated) List of config objects
└── terraform.tfstate            # (Auto-generated) Terraform state file
```

## Configuration Details

### Variables

| Variable | Description | Type | Default | Required |
|----------|-------------|------|---------|----------|
| access_key | EdgeNext Access Key | string | - | Yes |
| secret_key | EdgeNext Secret Key | string | - | Yes |
| endpoint | OSS service endpoint | string | - | No |
| region | OSS region | string | - | No |
| environment | Environment name (affects bucket naming and force_destroy) | string | dev | No |

### Resources Created

#### Buckets

1. **myapp-data-{environment}**
   - Purpose: Application data storage
   - ACL: private
   - force_destroy: true (dev only)

2. **myapp-backups-{environment}**
   - Purpose: Backup storage
   - ACL: private
   - force_destroy: true (dev only)

#### Objects

1. **config/app.json**
   - JSON configuration with dynamic content
   - Content-Type: application/json
   - Cache-Control: public, max-age=3600
   - Content-Disposition: attachment
   - Expires: 2025-11-05T16:00:00Z
   - ACL: private
   - Custom metadata: managed-by, environment, version

2. **config/redis.conf**
   - Redis configuration file
   - Content-Type: text/plain
   - Content-Disposition: attachment; filename=redis.conf

3. **config/nginx.conf**
   - Nginx configuration file
   - Content-Type: text/plain
   - Content-Disposition: attachment; filename=nginx.conf

4. **{files from ./files/ directory}**
   - Batch uploaded files
   - Source: ./files/*

#### Object Copy

1. **config/app-v2.json**
   - Copied from config/app.json
   - Metadata directive: REPLACE
   - ACL changed to: public-read
   - Cache-Control: no-cache
   - New metadata: version=2.0.0, updated-by=terraform

### Data Sources

1. **edgenext_oss_buckets.all**
   - Lists buckets with prefix "myapp"
   - Max buckets: 100
   - Output to: buckets.json

2. **edgenext_oss_objects.configs**
   - Lists objects in config/ directory
   - Uses delimiter "/" for folder-like structure
   - Max keys: 100
   - Output to: configs.json

3. **edgenext_oss_object.current_config**
   - Reads config/app.json content and metadata

### Outputs

- **bucket_details** - Details of created buckets including name, ACL, and force_destroy settings

## Key Features Demonstrated

### 1. Dynamic Content Generation

The example uses `jsonencode()` to generate JSON configuration dynamically:

```hcl
content = jsonencode({
  version     = "1.0.0"
  environment = var.environment
  features = {
    auth_enabled  = true
    cache_enabled = true
    debug_mode    = var.environment == "dev"
  }
})
```

### 2. Conditional Configuration

Force destroy is enabled only in dev environment:

```hcl
force_destroy = var.environment == "dev" ? true : false
```

### 3. Object Metadata

Custom metadata using hyphens (not underscores) in keys:

```hcl
metadata = {
  managed-by  = "terraform"
  environment = var.environment
  version     = "1.0.0"
}
```

### 4. Batch Object Creation

Using `for_each` with local variables to create multiple objects:

```hcl
locals {
  config_files = {
    "redis.conf" = { ... }
    "nginx.conf" = { ... }
  }
}

resource "edgenext_oss_object" "configs" {
  for_each = local.config_files
  # ...
}
```

### 5. File Upload from Directory

Upload all files from a local directory:

```hcl
locals {
  files = fileset("./files", "*")
}

resource "edgenext_oss_object" "all_files" {
  for_each = { for f in local.files : f => "./files/${f}" }
  source   = each.value
  # ...
}
```

### 6. Object Copy with Metadata Replacement

Copy an object and replace its metadata:

```hcl
resource "edgenext_oss_object_copy" "config_with_metadata" {
  source_bucket      = edgenext_oss_bucket.app_data.id
  source_key         = "config/app.json"
  bucket             = edgenext_oss_bucket.app_data.id
  key                = "config/app-v2.json"
  metadata_directive = "REPLACE"
  acl                = "public-read"
  # ...
}
```

### 7. Data Source with Prefix and Delimiter

List objects in a folder-like structure:

```hcl
data "edgenext_oss_objects" "configs" {
  bucket    = edgenext_oss_bucket.app_data.id
  prefix    = "config/"
  delimiter = "/"  # Treats "/" as folder separator
  # ...
}
```

## FAQ

### Q: How to upload large files?

A: For large files (>5GB), it's recommended to use the AWS CLI or SDK's multipart upload feature instead of Terraform. Terraform is best for managing infrastructure configuration, not large data uploads.

### Q: Why use hyphens instead of underscores in metadata keys?

A: S3-compatible APIs typically use hyphens in HTTP headers and metadata keys. While underscores may work, hyphens are the standard convention (e.g., `Content-Type`, `Cache-Control`).

### Q: What if bucket deletion fails?

A: Ensure the bucket is empty, or set `force_destroy = true` to automatically delete all objects when destroying the bucket. **Warning**: Use force_destroy carefully in production!

### Q: How to update object content?

A: Simply modify the `content` or `source` parameter in your Terraform configuration and run `terraform apply`. Terraform will detect the change and update the object.

### Q: Can I use this with other S3-compatible services?

A: Yes! Just change the `endpoint` and `region` variables to match your S3-compatible service endpoint.

### Q: How do I access objects after creation?

A: For public objects, use: `https://{endpoint}/{bucket}/{key}`
   For private objects, you'll need to generate a pre-signed URL using the AWS SDK or CLI.

## Best Practices

1. **Use force_destroy Carefully**
   - Only enable in development/test environments
   - Always set to false in production to prevent accidental data loss

2. **Set Appropriate Cache-Control Headers**
   - Static resources: `public, max-age=31536000` (1 year)
   - Dynamic content: `no-cache` or short max-age
   - Private data: `private, no-cache`

3. **Use Metadata to Tag Resources**
   - Always add `managed-by = "terraform"` to track resource management
   - Include environment, version, and other useful information

4. **Organize Objects with Prefixes**
   - Use prefix patterns like `config/`, `assets/`, `backups/`
   - Makes it easier to list and manage objects
   - Enables better access control policies

5. **Use depends_on for Resource Dependencies**
   - Ensure buckets are created before objects
   - Ensure source objects exist before copying
   - Prevents race conditions during creation

6. **Sensitive Data Handling**
   - Mark sensitive variables with `sensitive = true`
   - Never commit `terraform.tfvars` with real credentials
   - Use environment variables or secure secret management

7. **Output Files for Verification**
   - Use `output_file` in data sources to save results
   - Helps verify resource creation
   - Useful for debugging

## Troubleshooting

### Error: Bucket already exists

If you get an error that the bucket already exists:
1. Change the bucket name in your configuration
2. Or import the existing bucket: `terraform import edgenext_oss_bucket.app_data existing-bucket-name`

### Error: Failed to read from directory

If the batch upload fails because `./files/` doesn't exist:
1. Create the directory: `mkdir -p files`
2. Or comment out the batch upload section in main.tf

### Error: Access Denied

Ensure your credentials have the necessary permissions:
- `s3:CreateBucket`, `s3:DeleteBucket` for bucket management
- `s3:PutObject`, `s3:GetObject`, `s3:DeleteObject` for object management
- `s3:ListBucket` for listing operations

## Related Resources

- [EdgeNext OSS Documentation](https://docs.edgenext.com/oss/)
- [Terraform Documentation](https://www.terraform.io/docs)
- [S3 API Reference](https://docs.aws.amazon.com/s3/)
- [AWS S3 Best Practices](https://docs.aws.amazon.com/AmazonS3/latest/userguide/best-practices.html)
