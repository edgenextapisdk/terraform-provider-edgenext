---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_certificate"
sidebar_current: "docs-edgenext-resource-elb_certificate"
description: |-
  Use this resource to upload and manage EdgeNext ELB TLS certificates (Barbican containers). There is no update API; changing certificate or private key material requires replacement.
---

# edgenext_elb_certificate

Use this resource to upload and manage EdgeNext ELB TLS certificates (Barbican containers). There is no update API; changing certificate or private key material requires replacement.

## Example Usage

```hcl
resource "edgenext_elb_certificate" "site" {
  name        = "www-example-com"
  certificate = file("${path.module}/fullchain.pem")
  private_key = file("${path.module}/privkey.pem")
}

resource "edgenext_elb_listener" "https" {
  loadbalancer_id             = var.loadbalancer_id
  name                        = "https"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  default_tls_certificate_ref = edgenext_elb_certificate.site.certificate_ref
}
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


## Import

Import format is the Barbican container `certificate_id`.

```shell
terraform import edgenext_elb_certificate.site a037e728-a23f-4937-ae7d-a70e1fedf828
```

