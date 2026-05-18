---
subcategory: "Elastic IP (EIP)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_eip_floating_ips"
sidebar_current: "docs-edgenext-datasource-eip_floating_ips"
description: |-
  Use this data source to query EdgeNext floating IPs (EIPs).
---

# edgenext_eip_floating_ips

Use this data source to query EdgeNext floating IPs (EIPs).

## Example Usage

```hcl
data "edgenext_eip_floating_ips" "example" {
  limit = 50
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip_address` - (Optional, String) Floating IP address to filter.
* `floating_ip_id` - (Optional, String) Floating IP ID to filter.
* `limit` - (Optional, Int) Maximum number of floating IPs to return.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `floating_ips` - List of floating IPs.
  * `bandwidth` - Bandwidth in Mbps.
  * `billing_model` - Billing model.
  * `charge_mode` - Charge mode.
  * `created_at` - Creation time.
  * `description` - Description.
  * `expiration_time` - Expiration time.
  * `fixed_ip_address` - Fixed (private) IP when associated.
  * `floating_ip_address` - Public floating IP address.
  * `floating_network_id` - Floating network ID.
  * `floating_network_name` - Floating network name.
  * `id` - Floating IP ID.
  * `instance_name` - Associated instance name when present.
  * `network_interface_id` - Network interface (port) ID when associated.
  * `network_interface_name` - Network interface name.
  * `port_forwardings` - Port forwarding entries.
  * `project_id` - Project ID.
  * `qos_policy_id` - QoS policy ID.
  * `revision_number` - Revision number.
  * `router_id` - Router ID.
  * `status` - Status.
  * `tags` - Tag strings.
  * `tenant_id` - Tenant ID.
  * `updated_at` - Last update time.
* `total` - Total number of matched floating IPs.


