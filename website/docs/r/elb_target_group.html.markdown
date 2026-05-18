---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_target_group"
sidebar_current: "docs-edgenext-resource-elb_target_group"
description: |-
  Manages an EdgeNext ELB target group.
---

# edgenext_elb_target_group

Manages an EdgeNext ELB target group.

## Example Usage

```hcl
# See examples/elb/main.tf
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
* `type` - (Required, String) Health monitor type (for example HTTP, HTTPS, TCP).
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


