---
layout: "edgenext"
page_title: "Provider: EdgeNext"
sidebar_current: "docs-edgenext-index"
description: |-
  The EdgeNext provider is used to interact with EdgeNext services.
---

# EdgeNext Provider

The EdgeNext Provider can be used to configure infrastructure in [EdgeNext](https://www.edgenext.com) using the EdgeNext Resource Manager API's. Documentation regarding the Data Sources and Resources supported by the EdgeNext Provider can be found in the navigation to the left.

-> **Note:** This provider requires EdgeNext API credentials (access key and secret key).

## Example Usage

```hcl
terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}

# Configure the EdgeNext Provider
provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
  region     = var.region
}
```

## Authentication

The EdgeNext provider offers a flexible means of providing credentials for authentication. The following methods are supported, in order of precedence:

1. **Static credentials** in the provider configuration block
2. **Environment variables**

### Static Credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

Static credentials can be provided by adding `access_key` and `secret_key` in-line in the EdgeNext provider block:

```hcl
provider "edgenext" {
  access_key = "your-access-key"
  secret_key = "your-secret-key"
  endpoint   = "https://api.edgenext.com"
  region     = "us-east-1"
}
```

### Environment Variables

You can provide your credentials via the `EDGENEXT_ACCESS_KEY`, `EDGENEXT_SECRET_KEY`, `EDGENEXT_ENDPOINT` and `EDGENEXT_REGION` environment variables:

```hcl
provider "edgenext" {}
```

Usage:

```bash
export EDGENEXT_ACCESS_KEY="your-access-key"
export EDGENEXT_SECRET_KEY="your-secret-key"
export EDGENEXT_ENDPOINT="https://api.edgenext.com"
export EDGENEXT_REGION="us-east-1"
terraform plan
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `access_key` - (Required) EdgeNext access key for authentication. It can also be sourced from the `EDGENEXT_ACCESS_KEY` environment variable.

* `secret_key` - (Required) EdgeNext secret key for authentication. It can also be sourced from the `EDGENEXT_SECRET_KEY` environment variable.

* `endpoint` - (Required) EdgeNext API endpoint address. It can also be sourced from the `EDGENEXT_ENDPOINT` environment variable.

* `region` - (Optional) EdgeNext region. It can also be sourced from the `EDGENEXT_REGION` environment variable.

## Resources and Data Sources

The EdgeNext provider supports the following resource types:

### Content Delivery Network (CDN)

#### Resources

* [`edgenext_cdn_domain`](resources/cdn_domain) - Manage CDN domain configuration
* [`edgenext_cdn_push`](resources/cdn_push) - Manage CDN cache push tasks
* [`edgenext_cdn_purge`](resources/cdn_purge) - Manage CDN cache purge tasks

#### Data Sources

* [`edgenext_cdn_domain`](data-sources/cdn_domain) - Query CDN domain configuration
* [`edgenext_cdn_domains`](data-sources/cdn_domains) - Query CDN domains
* [`edgenext_cdn_push`](data-sources/cdn_push) - Query CDN push task details
* [`edgenext_cdn_pushes`](data-sources/cdn_pushes) - Query CDN push tasks
* [`edgenext_cdn_purge`](data-sources/cdn_purge) - Query CDN purge task details
* [`edgenext_cdn_purges`](data-sources/cdn_purges) - Query CDN purge tasks

### SSL Certificate Management (SSL)

#### Resources

* [`edgenext_ssl_certificate`](resources/ssl_certificate) - Manage SSL certificates

#### Data Sources

* [`edgenext_ssl_certificate`](data-sources/ssl_certificate) - Query SSL certificate details
* [`edgenext_ssl_certificates`](data-sources/ssl_certificates) - Query SSL certificates

### Object Storage Service (OSS)

#### Resources

* [`edgenext_oss_bucket`](resources/oss_bucket) - Manage OSS buckets
* [`edgenext_oss_object`](resources/oss_object) - Manage OSS objects
* [`edgenext_oss_object_copy`](resources/oss_object_copy) - Manage OSS object copy

#### Data Sources

* [`edgenext_oss_buckets`](data-sources/oss_buckets) - Query OSS buckets
* [`edgenext_oss_object`](data-sources/oss_object) - Query OSS object details
* [`edgenext_oss_objects`](data-sources/oss_objects) - Query OSS objects


### Object Storage Service(OSS)


#### Resources


* [`edgenext_oss_bucket`](resources/oss_bucket) - Manage OSS buckets

* [`edgenext_oss_object`](resources/oss_object) - Manage OSS objects

* [`edgenext_oss_object_copy`](resources/oss_object_copy) - Manage OSS object copy




#### Data Sources


* [`edgenext_oss_buckets`](data-sources/oss_buckets) - Query OSS buckets

* [`edgenext_oss_object`](data-sources/oss_object) - Query OSS object details

* [`edgenext_oss_objects`](data-sources/oss_objects) - Query OSS objects




