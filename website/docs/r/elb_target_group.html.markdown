---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_target_group"
sidebar_current: "docs-edgenext-resource-elb_target_group"
description: |-
  Use this resource to create an EdgeNext ELB target group (pool) with a health monitor on a listener.
---

# edgenext_elb_target_group

Use this resource to create an EdgeNext ELB target group (pool) with a health monitor on a listener.

Backends are managed separately with `edgenext_elb_target_group_attachment`. Only `name` and `description` can be updated in place; `listener_id`, `protocol`, `lb_algorithm`, and `health_monitor` force replacement. Create and delete wait for Octavia provisioning and retry on transient immutable errors (see `timeouts`).

## Example Usage

```hcl
resource "edgenext_elb_target_group" "app" {
  listener_id  = edgenext_elb_listener.https.id
  name         = "app"
  description  = "Application backends"
  protocol     = "HTTP"
  lb_algorithm = "ROUND_ROBIN"
  health_monitor {
    name           = "http-check"
    type           = "HTTP"
    max_retries    = 3
    delay          = 10
    timeout        = 5
    http_method    = "GET"
    url_path       = "/health"
    expected_codes = "200"
  }
}

resource "edgenext_elb_target_group_attachment" "app1" {
  target_group_id = edgenext_elb_target_group.app.id
  address         = "192.168.0.244"
  protocol_port   = 80
  name            = "app-1"
  weight          = 10
}
```

## Argument Reference

The following arguments are supported:

* `health_monitor` - (Required, List, ForceNew) Health check for this target group. Single block; changing it forces replacement.
* `lb_algorithm` - (Required, String, ForceNew) Load balancing algorithm (for example ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP).
* `listener_id` - (Required, String, ForceNew) Listener this target group is attached to.
* `name` - (Required, String) Target group name.
* `protocol` - (Required, String, ForceNew) Traffic protocol between the load balancer and backends (for example HTTP).
* `description` - (Optional, String) Target group description (updatable).

The `health_monitor` object supports the following:

* `delay` - (Required, Int) Seconds between health checks.
* `expected_codes` - (Required, String) Expected HTTP status codes (for example 200).
* `http_method` - (Required, String) HTTP method for HTTP(S) checks.
* `max_retries` - (Required, Int) Maximum retries before marking the target down.
* `name` - (Required, String) Health monitor name.
* `timeout` - (Required, Int) Health check timeout in seconds.
* `type` - (Required, String) Health monitor type (for example HTTP, HTTPS, TCP, PING, TLS-HELLO).
* `url_path` - (Required, String) URL path for HTTP(S) checks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Creation time as Unix timestamp (seconds).
* `healthmonitor_id` - Health monitor identifier for this target group (computed).
* `loadbalancer_id` - Load balancer this target group belongs to (computed).
* `operating_status` - Target group operating status.
* `provisioning_status` - Target group provisioning status.
* `target_ids` - Identifiers of targets (backends) attached to this target group (computed). Use edgenext_elb_target_group_attachment to add or remove targets.
* `updated_at` - Last update time as Unix timestamp (seconds).


