---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_certificate"
sidebar_current: "docs-edgenext-resource-scdn_certificate"
description: |-
  Provides a resource to create and manage SCDN certificates.
---

# edgenext_scdn_certificate

Provides a resource to create and manage SCDN certificates.

## Example Usage

### Create certificate with certificate and key

```hcl
resource "edgenext_scdn_certificate" "example" {
  ca_name = "my-certificate"
  ca_cert = file("certificate.pem")
  ca_key  = file("private_key.pem")
}
```

### Update existing certificate

```hcl
resource "edgenext_scdn_certificate" "example" {
  certificate_id = "12345"
  ca_name        = "my-certificate"
  ca_cert        = file("new_certificate.pem")
  ca_key         = file("new_private_key.pem")
}
```

## Argument Reference

The following arguments are supported:

* `ca_name` - (Required, String) The certificate name
* `ca_cert` - (Optional, String) The certificate public key (PEM format). Required for creation, optional for updates.
* `ca_key` - (Optional, String) The certificate private key (PEM format). Required for creation, optional for updates.
* `certificate_id` - (Optional, String) The certificate ID for updating an existing certificate. If provided, this will update the certificate instead of creating a new one.
* `product_flag` - (Optional, String) The product flag

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `apply_status` - The application status: 1-applying, 2-issued, 3-review failed, 4-uploaded
* `binded` - Whether the certificate is bound
* `ca_domain` - The domains in the certificate
* `ca_sn` - The certificate serial number
* `ca_type_domain` - The certificate domain type: 1-single domain, 2-multiple domains, 3-wildcard domain
* `ca_type` - The certificate type: 1-upload, 2-lets apply
* `created_at` - The creation timestamp
* `id` - The ID of the certificate
* `issuer_expiry_time_desc` - The certificate expiry time description
* `issuer_expiry_time` - The certificate expiry time
* `issuer_start_time` - The certificate start time
* `issuer` - The certificate issuer
* `member_id` - The member ID
* `renew_status` - The renewal status: 1-default, 2-renewing, 3-renewal failed, 4-renewal successful
* `updated_at` - The last update timestamp


## Import

SCDN certificates can be imported using the certificate ID:

```shell
terraform import edgenext_scdn_certificate.example 12345
```

