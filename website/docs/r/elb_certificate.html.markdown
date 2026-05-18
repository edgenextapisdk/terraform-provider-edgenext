---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_certificate"
sidebar_current: "docs-edgenext-resource-elb_certificate"
description: |-
  Manages an EdgeNext ELB certificate.
---

# edgenext_elb_certificate

Manages an EdgeNext ELB certificate.

## Example Usage

```hcl
# See examples/elb/main.tf
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Required, String, ForceNew) PEM certificate body sent on create (chain supported).
* `name` - (Required, String, ForceNew) Certificate name.
* `private_key` - (Optional, String, ForceNew) PEM private key sent on create (required for create). Detail API usually omits the key; when omitted, the prior value is kept and must not force replacement.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cert_detail` - Extra certificate detail as JSON when the API returns an object; empty when null.
* `certificate_id` - Certificate (Barbican container) ID returned by the API.
* `certificate_ref` - Certificate (Barbican container) reference URL.
* `created` - Creation time as Unix timestamp (seconds).
* `expiration` - Expiration time as Unix timestamp (seconds).
* `private_key_passphrase` - Private key passphrase from detail when present.
* `status` - Certificate status.
* `type` - Container type (for example certificate).
* `updated` - Last update time as Unix timestamp (seconds).


