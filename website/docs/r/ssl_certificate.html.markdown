---
subcategory: "SSL Certificate Management (SSL)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ssl_certificate"
sidebar_current: "docs-edgenext-resource-ssl_certificate"
description: |-
  Provides a resource to create and manage SSL certificates.
---

# edgenext_ssl_certificate

Provides a resource to create and manage SSL certificates.

## Example Usage

### Basic SSL certificate upload

```hcl
resource "edgenext_ssl_certificate" "example" {
  name        = "example-com-cert"
  certificate = file("path/to/certificate.crt")
  key         = file("path/to/private.key")
}
```

### SSL certificate with certificate content

```hcl
resource "edgenext_ssl_certificate" "example" {
  name = "example-com-cert"

  certificate = <<-EOT
-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAJC1HiIAZAiIMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
...
-----END CERTIFICATE-----
EOT

  key = <<-EOT
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7S2IgnLVkjpQR
RIiTq86f1o3d6nOF4eU4h95tX3YfL8s6eSgfDNF2nDG8VQZ4m8Sv1JHbYDrDJ8Ac
...
-----END PRIVATE KEY-----
EOT
}
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Required, String) SSL certificate content
* `key` - (Required, String) SSL certificate private key content
* `name` - (Required, String) SSL certificate name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bind_domains` - List of bound domains
* `cert_expire_time` - Certificate end time
* `cert_id` - Certificate ID
* `cert_start_time` - Certificate start time


## Import

SSL certificates can be imported using the certificate ID:

```shell
terraform import edgenext_ssl_certificate.example ssl-cert-123456
```

