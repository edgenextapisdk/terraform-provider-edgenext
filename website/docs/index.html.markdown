---
layout: "edgenext"
page_title: "Provider: EdgeNext"
sidebar_current: "docs-edgenext-index"
description: |-
  The EdgeNext provider is used to interact with EdgeNext CDN services.
---

# EdgeNext Provider

The EdgeNext provider is used to interact with many resources supported
by [EdgeNext CDN](https://www.edgenext.com).
The provider needs to be configured with the proper API credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** This provider supports EdgeNext CDN v2.0+ API.

## Example Usage

### Configure the EdgeNext Provider

```hcl
terraform {
  required_providers {
    edgenext = {
      source = "edgenextapisdk/edgenext"
    }
  }
}

# Configure the EdgeNext Provider
provider "edgenext" {
  api_key  = var.api_key
  secret   = var.secret
  endpoint = var.endpoint
}
```

### Configure with environment variables

```hcl
provider "edgenext" {
  # Configuration options
}
```

```shell
export EDGENEXT_API_KEY="your-api-key"
export EDGENEXT_SECRET="your-secret"  
export EDGENEXT_ENDPOINT="https://api.edgenext.com"
```

### Configure with timeout and retry settings

```hcl
provider "edgenext" {
  api_key     = var.api_key
  secret      = var.secret
  endpoint    = var.endpoint
  timeout     = 300
  retry_count = 3
}
```

## Authentication

The EdgeNext provider requires API credentials to authenticate with the EdgeNext API.

### API Credentials

You can provide your credentials via the following methods:

1. **Static credentials** (not recommended for production):

```hcl
provider "edgenext" {
  api_key  = "your-api-key"
  secret   = "your-secret"
  endpoint = "https://api.edgenext.com"
}
```

2. **Environment variables** (recommended):

```shell
export EDGENEXT_API_KEY="your-api-key"
export EDGENEXT_SECRET="your-secret"
export EDGENEXT_ENDPOINT="https://api.edgenext.com"
```

3. **Terraform variables**:

```hcl
variable "api_key" {
  description = "EdgeNext API Key"
  type        = string
  sensitive   = true
}

variable "secret" {
  description = "EdgeNext Secret"
  type        = string
  sensitive   = true
}

variable "endpoint" {
  description = "EdgeNext API Endpoint"
  type        = string
}

provider "edgenext" {
  api_key  = var.api_key
  secret   = var.secret
  endpoint = var.endpoint
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `api_key` - (Required) EdgeNext API key for authentication. This can also be specified with the `EDGENEXT_API_KEY` environment variable.

* `secret` - (Required) EdgeNext secret for authentication. This can also be specified with the `EDGENEXT_SECRET` environment variable.

* `endpoint` - (Required) EdgeNext API endpoint address. This can also be specified with the `EDGENEXT_ENDPOINT` environment variable.

* `timeout` - (Optional) API request timeout in seconds. Defaults to `300`.

* `retry_count` - (Optional) API request retry count. Defaults to `3`.

## Resources and Data Sources

The EdgeNext provider supports the following resource types:



### Content Delivery Network(CDN)


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





### SSL Certificate Management(SSL)


#### Resources


* [`edgenext_ssl_certificate`](resources/ssl_certificate) - Manage SSL certificates




#### Data Sources


* [`edgenext_ssl_certificate`](data-sources/ssl_certificate) - Query SSL certificate details

* [`edgenext_ssl_certificates`](data-sources/ssl_certificates) - Query SSL certificates




