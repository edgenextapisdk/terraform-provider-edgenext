---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_target_groups"
sidebar_current: "docs-edgenext-datasource-elb_target_groups"
description: |-
  Use this data source to query EdgeNext ELB target groups.
---

# edgenext_elb_target_groups

Use this data source to query EdgeNext ELB target groups.

## Example Usage

```hcl
data "edgenext_elb_target_groups" "example" {
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

* `loadbalancer_id` - (Required, String) Load balancer ID to list target groups for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `target_groups` - Target groups for the load balancer. Includes health_monitor when resolvable and per-target identifiers when present (targets are managed with edgenext_elb_target_group_attachment).
  * `created_at` - Creation time as Unix timestamp (seconds).
  * `description` - Target group description.
  * `health_monitor` - Resolved health monitor for this target group when available; at most one element.
    * `created_at` - Creation time as Unix timestamp (seconds).
    * `delay` - Seconds between health checks.
    * `expected_codes` - Expected HTTP status codes.
    * `http_method` - HTTP method for HTTP(S) checks.
    * `id` - Health monitor ID.
    * `max_retries` - Maximum retries.
    * `name` - Health monitor name.
    * `operating_status` - Operating status.
    * `provisioning_status` - Provisioning status.
    * `timeout` - Health check timeout in seconds.
    * `type` - Health monitor type.
    * `updated_at` - Last update time as Unix timestamp (seconds).
    * `url_path` - URL path for HTTP(S) checks.
  * `healthmonitor_id` - Health monitor id (same as health_monitor.0.id when present).
  * `id` - Target group ID.
  * `lb_algorithm` - Load balancing algorithm.
  * `listener_id` - First listener id from the target group.
  * `loadbalancer_id` - First load balancer id from the target group.
  * `member_ids` - Target identifiers when listed (manage targets in Terraform with edgenext_elb_target_group_attachment).
  * `name` - Target group name.
  * `operating_status` - Operating status.
  * `protocol` - Target group protocol.
  * `provisioning_status` - Provisioning status.
  * `updated_at` - Last update time as Unix timestamp (seconds).


