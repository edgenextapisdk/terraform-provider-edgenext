---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_listeners"
sidebar_current: "docs-edgenext-datasource-elb_listeners"
description: |-
  Use this data source to query EdgeNext ELB listeners.
---

# edgenext_elb_listeners

Use this data source to query EdgeNext ELB listeners.

## Example Usage

```hcl
data "edgenext_elb_listeners" "example" {
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

* `loadbalancer_id` - (Required, String) Load balancer ID to list listeners for.
* `limit` - (Optional, Int) Maximum number of listeners to return in one request.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hasmore` - Whether more results exist beyond this page.
* `listeners` - Listeners returned by the API.
  * `connection_limit` - Connection limit (-1 if unlimited).
  * `created_at` - Creation time as Unix timestamp (seconds).
  * `default_target_group_id` - Default target group ID when set.
  * `description` - Description.
  * `id` - Listener ID.
  * `insert_headers` - Insert headers configuration (string key to string value).
  * `name` - Listener name.
  * `operating_status` - Operating status.
  * `protocol_port` - Protocol port.
  * `protocol` - Listener protocol (for example TERMINATED_HTTPS, HTTP).
  * `provisioning_status` - Provisioning status.
* `total` - Total number of listeners matching the query (from the API when present).


