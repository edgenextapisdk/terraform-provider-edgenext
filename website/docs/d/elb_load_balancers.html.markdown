---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_load_balancers"
sidebar_current: "docs-edgenext-datasource-elb_load_balancers"
description: |-
  Use this data source to list EdgeNext ELB load balancers.
---

# edgenext_elb_load_balancers

Use this data source to list EdgeNext ELB load balancers.

## Example Usage

```hcl
data "edgenext_elb_load_balancers" "all" {
  name  = ""
  limit = 20
}

output "load_balancer_ids" {
  value = [for lb in data.edgenext_elb_load_balancers.all.balancers : lb.id]
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional, Int) Maximum number of load balancers to return in one request.
* `name` - (Optional, String) Filter by load balancer name. Use empty string to omit the filter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `balancers` - Load balancers returned by the API.
  * `created_at` - Creation time as Unix timestamp (seconds).
  * `description` - Description.
  * `elastic_ip` - Elastic IP if bound.
  * `flavor_id` - Flavor ID.
  * `floating_ip` - Floating IP when present; JSON-encoded if the API returns an object, otherwise empty.
  * `iam_id` - IAM resource ID.
  * `id` - Load balancer ID.
  * `listeners` - Listener references.
    * `id` - Listener ID.
  * `mode` - Mode.
  * `name` - Load balancer name.
  * `operating_status` - Operating status (for example ONLINE, ERROR).
  * `provisioning_status` - Provisioning status (for example ACTIVE).
  * `spec` - Specification details.
    * `architecture` - Architecture (for example single).
    * `created_at` - Spec creation time as Unix timestamp.
    * `description` - Spec description.
    * `flavor_id` - Flavor ID in spec.
    * `id` - Spec ID.
    * `max_connections` - Maximum connections.
    * `new_connections_per_sec` - New connections per second.
    * `qps` - Queries per second.
    * `sort_order` - Sort order.
    * `spec_name` - Spec name.
    * `updated_at` - Spec update time as Unix timestamp.
  * `target_groups` - Target group references.
    * `id` - Target group ID.
  * `type` - Load balancer type (for example application, network).
  * `updated_at` - Last update time as Unix timestamp (seconds).
  * `vip_address` - VIP address.
  * `vip_network_id` - VIP network ID.
  * `vip_port_id` - VIP port ID.
  * `vip_subnet_id` - VIP subnet ID.
* `hasmore` - Whether more results exist beyond this page.
* `total` - Total number of load balancers matching the query (from the API when present).


