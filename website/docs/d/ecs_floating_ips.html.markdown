---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_floating_ips"
sidebar_current: "docs-edgenext-datasource-ecs_floating_ips"
description: |-
  Use this data source to query ECS floating IPs.
---

# edgenext_ecs_floating_ips

Use this data source to query ECS floating IPs.

## Example Usage

```hcl
data "edgenext_ecs_floating_ips" "example" {
  region              = "tokyo-a"
  floating_ip_address = "148.222.161.86"
  limit               = 10
}

output "floating_ip_total" {
  value = data.edgenext_ecs_floating_ips.example.total
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `eid` - (Optional, String) The floating IP ID to filter.
* `floating_ip_address` - (Optional, String) The floating IP address to filter.
* `limit` - (Optional, Int) Maximum number of floating IPs to return.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `floating_ips` - A list of ECS floating_ips.
  * `bandwidth` - Bandwidth in Mbps.
  * `billing_model` - Billing model.
  * `charge_mode` - Charge mode.
  * `created_at` - Creation time.
  * `description` - The description.
  * `expiration_time` - Expiration time.
  * `fixed_ip_address` - The fixed IP address.
  * `floating_ip_address` - The floating IP address.
  * `floating_network_id` - The floating network ID.
  * `floating_network_name` - Floating network name.
  * `id` - The ID of the floating_ip.
  * `instance_name` - Instance name.
  * `port_forwardings` - Port forwarding entries.
  * `port_id` - The port ID.
  * `port_name` - Port name.
  * `project_id` - Project ID.
  * `qos_policy_id` - The QoS policy ID.
  * `revision_number` - Revision number.
  * `router_id` - The router ID.
  * `status` - The status.
  * `tags` - A list of tag strings.
  * `tenant_id` - The tenant ID.
  * `updated_at` - Last update time.
* `total` - Total number of matched floating IPs.


