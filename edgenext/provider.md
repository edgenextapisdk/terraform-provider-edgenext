The EdgeNext provider is used to interact with many resources supported
by [EdgeNext CDN](https://www.edgenext.com).
The provider needs to be configured with the proper API credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** This provider supports EdgeNext CDN v2.0+ API.

Example Usage

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

# Configure the EdgeNext Provider with timeout settings
provider "edgenext" {
  api_key     = var.api_key
  secret      = var.secret
  endpoint    = var.endpoint
  timeout     = 300
  retry_count = 3
}
```

Resources List

Content Delivery Network(CDN)
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

SSL Certificate Management(SSL)
Data Source
edgenext_ssl_certificate
edgenext_ssl_certificates

Resource
edgenext_ssl_certificate
