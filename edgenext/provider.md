---
layout: "edgenext"
page_title: "Provider: EdgeNext"
sidebar_current: "docs-edgenext-index"
description: |-
  The EdgeNext provider is used to interact with EdgeNext services.
---

# EdgeNext Provider

The EdgeNext provider is used to interact with the many resources supported by [EdgeNext](https://www.edgenext.com). The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to read about the available resources and data sources.

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

Resources List

Content Delivery Network (CDN)
Data Source
edgenext_cdn_domain
edgenext_cdn_domains
edgenext_cdn_push
edgenext_cdn_pushes
edgenext_cdn_purge
edgenext_cdn_purges

Resource
edgenext_cdn_domain
edgenext_cdn_push
edgenext_cdn_purge

SSL Certificate Management (SSL)
Data Source
edgenext_ssl_certificate
edgenext_ssl_certificates

Resource
edgenext_ssl_certificate

Object Storage Service (OSS)
Data Source
edgenext_oss_buckets
edgenext_oss_object
edgenext_oss_objects

Resource
edgenext_oss_bucket
edgenext_oss_object
edgenext_oss_object_copy
