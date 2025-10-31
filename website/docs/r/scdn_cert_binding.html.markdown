---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cert_binding"
sidebar_current: "docs-edgenext-resource-scdn_cert_binding"
description: |-
  Provides a resource to bind a certificate to an SCDN domain.
---

# edgenext_scdn_cert_binding

Provides a resource to bind a certificate to an SCDN domain.

## Example Usage

### Bind certificate to domain

```hcl
resource "edgenext_scdn_cert_binding" "example" {
  domain_id = 12345
  ca_id     = 67890
}
```

## Argument Reference

The following arguments are supported:

* `ca_id` - (Required, Int, ForceNew) The ID of the certificate to bind
* `domain_id` - (Required, Int, ForceNew) The ID of the domain to bind the certificate to

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `id` - The unique identifier for this certificate binding


## Import

SCDN certificate bindings can be imported using the domain ID and certificate ID:

```shell
terraform import edgenext_scdn_cert_binding.example 12345-67890
```

