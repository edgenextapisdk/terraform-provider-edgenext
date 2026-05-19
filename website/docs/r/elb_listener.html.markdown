---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_listener"
sidebar_current: "docs-edgenext-resource-elb_listener"
description: |-
  Use this resource to create and manage an EdgeNext ELB listener on a load balancer.
---

# edgenext_elb_listener

Use this resource to create and manage an EdgeNext ELB listener on a load balancer.

Updatable in place: `name`, `description`, `default_tls_certificate_ref`. Other arguments are create-time only and force replacement when changed. Delete waits until the listener is removed from the API (configure `timeouts.delete` if Octavia is slow).

## Example Usage

```hcl
resource "edgenext_elb_certificate" "site" {
  name        = "www-example-com"
  certificate = file("${path.module}/fullchain.pem")
  private_key = file("${path.module}/privkey.pem")
}

resource "edgenext_elb_listener" "https" {
  loadbalancer_id  = var.loadbalancer_id
  name             = "https"
  description      = "Public HTTPS"
  protocol         = "TERMINATED_HTTPS"
  protocol_port    = 443
  connection_limit = -1
  insert_headers = {
    "X-Forwarded-For"   = "true"
    "X-Forwarded-Port"  = "true"
    "X-Forwarded-Proto" = "true"
  }
  default_tls_certificate_ref = edgenext_elb_certificate.site.certificate_ref
}
```

## Argument Reference

The following arguments are supported:

* `loadbalancer_id` - (Required, String, ForceNew) Load balancer ID this listener belongs to.
* `name` - (Required, String) Listener name.
* `protocol_port` - (Required, Int, ForceNew) Protocol port.
* `protocol` - (Required, String, ForceNew) Listener protocol (for example TERMINATED_HTTPS, HTTP).
* `connection_limit` - (Optional, Int, ForceNew) Connection limit (-1 for unlimited). Create only.
* `default_tls_certificate_ref` - (Optional, String) Default TLS certificate (Barbican container) reference URL (for example TERMINATED_HTTPS). Updatable via listeners/update.
* `description` - (Optional, String) Listener description (updatable).
* `insert_headers` - (Optional, Map, ForceNew) Insert headers for the listener (string keys and values, for example X-Forwarded-For, X-Forwarded-Port, and X-Forwarded-Proto set to true). Create only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `admin_state_up` - Administrative up/down state from the API.
* `created_at` - Creation time as Unix timestamp (seconds).
* `default_target_group_id` - Default target group ID when set.
* `operating_status` - Operating status.
* `provisioning_status` - Provisioning status.
* `updated_at` - Last update time as Unix timestamp (seconds).


## Import

Import format is `listener_id`.

```shell
terraform import edgenext_elb_listener.https 0c03c22d-b9b5-49c8-a5e1-aedd9bc6e6b9
```

