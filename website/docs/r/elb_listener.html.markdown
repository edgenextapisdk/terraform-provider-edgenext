---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_listener"
sidebar_current: "docs-edgenext-resource-elb_listener"
description: |-
  Manages an EdgeNext ELB listener.
---

# edgenext_elb_listener

Manages an EdgeNext ELB listener.

## Example Usage

```hcl
# See examples/elb/main.tf
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
* `insert_headers` - (Optional, Map, ForceNew) Insert headers for the listener (string keys and values, for example X-Forwarded-For = true). Create only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `admin_state_up` - Administrative up/down state from the API.
* `created_at` - Creation time as Unix timestamp (seconds).
* `default_target_group_id` - Default target group ID when set.
* `operating_status` - Operating status.
* `provisioning_status` - Provisioning status.
* `updated_at` - Last update time as Unix timestamp (seconds).


