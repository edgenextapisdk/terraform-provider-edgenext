terraform {
  required_providers {
    edgenext = {
      source = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}

provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
  region     = var.region
}

variable "access_key" {
  description = "EdgeNext Access Key"
  type        = string
  sensitive   = true
}

variable "secret_key" {
  description = "EdgeNext Secret Key"
  type        = string
  sensitive   = true
}

variable "endpoint" {
  description = "EdgeNext OSS Endpoint"
  type        = string
}

variable "region" {
  description = "EdgeNext OSS Region"
  type        = string
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
}

# ========================================
# Bucket Resources
# ========================================


# Private bucket for application data
resource "edgenext_oss_bucket" "app_data" {
  bucket        = "myapp-data-${var.environment}"
  acl           = "private"
  force_destroy = var.environment == "dev" ? true : false
}

resource "edgenext_oss_bucket" "app_backup" {
  bucket        = "myapp-backup-${var.environment}"
  acl           = "private"
  force_destroy = var.environment == "dev" ? true : false
}

# # ========================================
# # Object Resources - Configuration Files
# # ========================================


# # JSON configuration file
resource "edgenext_oss_object" "app_config" {
  bucket = edgenext_oss_bucket.app_data.id
  key    = "config/app.json"
  content = jsonencode({
    version     = "1.0.0"
    environment = var.environment
    features = {
      auth_enabled  = true
      cache_enabled = true
      debug_mode    = var.environment == "dev"
    }
  })
  content_type        = "application/json"
  cache_control       = "public, max-age=3600"
  content_disposition = "attachment"
  expires             = "2025-11-05T16:00:00Z"
  acl                 = "private"

  metadata = {
    managed-by  = "terraform" # Use hyphens, not underscores!
    environment = var.environment
    version     = "1.0.0"
  }

  depends_on = [edgenext_oss_bucket.app_data]
}

# Image file
resource "edgenext_oss_object" "app_image" {
  bucket = edgenext_oss_bucket.app_data.id
  key    = "images/logo.png"
  source = "${path.module}/images/logo.png"
  content_type        = "image/png"
  cache_control       = "public, max-age=3600"
  content_disposition = "attachment; filename=logo.png"
  expires             = "2025-11-05T16:00:00Z"
  acl                 = "public-read"

  metadata = {
    managed-by  = "terraform" # Use hyphens, not underscores!
    environment = var.environment
    version     = "1.0.0"
  }

  depends_on = [edgenext_oss_bucket.app_data]
}

# ========================================
# Object Copy Resources
# ========================================

# Copy object with new metadata and public access
resource "edgenext_oss_object_copy" "config_with_metadata" {
  source_bucket = edgenext_oss_bucket.app_data.id
  source_key    = "config/app.json"
  bucket        = edgenext_oss_bucket.app_backup.id
  key           = "config/app-v2.json"

  acl = "public-read"

  metadata_directive = "REPLACE"
  content_type       = "application/json"
  cache_control      = "no-cache"
  metadata = {
    version    = "2.0.0"
    updated-by = "terraform"
  }

  depends_on = [edgenext_oss_bucket.app_data, edgenext_oss_object.app_config]
}


# ========================================
# Data Sources
# ========================================

# List all buckets
data "edgenext_oss_buckets" "all" {
  bucket_prefix = "myapp"
  max_buckets   = 100
  output_file   = "outputs/buckets.json"
  depends_on = [
    edgenext_oss_bucket.app_data,
    edgenext_oss_bucket.app_backup
  ]
}

# List objects in config directory with delimiter (folder-like structure)
data "edgenext_oss_objects" "configs" {
  bucket      = edgenext_oss_bucket.app_data.id
  prefix      = "config/"
  delimiter   = "/"
  max_keys    = 100
  output_file = "outputs/configs.json"

  depends_on = [edgenext_oss_object.app_config]
}

# # Read specific configuration object
data "edgenext_oss_object" "current_config" {
  bucket = edgenext_oss_bucket.app_data.id
  key    = edgenext_oss_object.app_config.key

  depends_on = [edgenext_oss_object.app_config]
}

data "edgenext_oss_object" "backup_config" {
  bucket = edgenext_oss_bucket.app_backup.id
  key    = edgenext_oss_object_copy.config_with_metadata.key

  depends_on = [edgenext_oss_object_copy.config_with_metadata]
}

# # ========================================
# # Outputs
# # ========================================

output "bucket_details" {
  description = "Details of created buckets"
  value = {
    app_data = {
      name          = edgenext_oss_bucket.app_data.bucket
      acl           = edgenext_oss_bucket.app_data.acl
      force_destroy = edgenext_oss_bucket.app_data.force_destroy
    }
    app_backup = {
      name          = edgenext_oss_bucket.app_backup.bucket
      acl           = edgenext_oss_bucket.app_backup.acl
      force_destroy = edgenext_oss_bucket.app_backup.force_destroy
    }
  }
}


# ========================================
# Example: Multiple objects with for_each
# ========================================

locals {
  config_files = {
    "redis.conf" = {
      content             = "maxmemory 256mb\nmaxmemory-policy allkeys-lru"
      content_type        = "text/plain"
      content_disposition = "attachment; filename=redis.conf"
    }
    "nginx.conf" = {
      content             = <<-EOT
        server {
            listen 80;
            server_name example.com;
            root /var/www/html;
        }
      EOT
      content_type        = "text/plain"
      content_disposition = "attachment; filename=nginx.conf"
    }
  }
}

resource "edgenext_oss_object" "configs" {
  for_each = local.config_files

  bucket              = edgenext_oss_bucket.app_data.id
  key                 = "config/${each.key}"
  content             = each.value.content
  content_type        = each.value.content_type
  content_disposition = each.value.content_disposition

  metadata = {
    config-name = each.key
    managed-by  = "terraform"
  }

  depends_on = [edgenext_oss_bucket.app_data]
}

# ========================================
# Example: Upload all files in a directory
# ========================================

locals {
  files = fileset("./files", "*")  # Get all files in directory
}

resource "edgenext_oss_object" "all_files" {
  for_each = { for f in local.files : f => "./files/${f}" }

  bucket = edgenext_oss_bucket.app_data.id
  key    = "config/${each.key}"
  source = each.value

  depends_on = [edgenext_oss_bucket.app_data]
}

# ========================================
# Import existing resources
# ========================================

# Import existing bucket
# resource "edgenext_oss_bucket" "import_bucket" {}

# Import existing object
# resource "edgenext_oss_object" "import_object" {}