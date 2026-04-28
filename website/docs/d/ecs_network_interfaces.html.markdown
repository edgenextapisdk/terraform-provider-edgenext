---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_network_interfaces"
sidebar_current: "docs-edgenext-datasource-ecs_network_interfaces"
description: |-
  Use this data source to query ECS network interfaces (ports).
---

# edgenext_ecs_network_interfaces

Use this data source to query ECS network interfaces (ports).

## Example Usage

```hcl
data "edgenext_ecs_network_interfaces" "example" {
  network_interface_name = edgenext_ecs_network_interface.example.name
  limit                  = 10
}

resource "edgenext_ecs_network_interface" "example" {
  name      = "example-eni"
  vpc_id    = data.edgenext_ecs_vpcs.all.vpcs[0].id
  subnet_id = data.edgenext_ecs_vpc_subnets.all.subnets[0].id
}

data "edgenext_ecs_vpcs" "all" {
  limit = 1
}

data "edgenext_ecs_vpc_subnets" "all" {
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional, Int) Maximum number of ports to return.
* `network_interface_name` - (Optional, String) Filter by network interface name (partial match per API behavior).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_interfaces` - A list of ECS network interfaces (Neutron ports).
  * `admin_state_up` - Administrative state of the network interface.
  * `binding_vnic_type` - VNIC binding type (from binding:vnic_type; may be filled from origin_data).
  * `created_at` - Creation time.
  * `description` - Description.
  * `fixed_ips` - Fixed IP assignments.
    * `floating_ip` - Associated floating IP if present.
    * `ip_address` - Fixed IP address.
    * `subnet_id` - Subnet ID.
  * `id` - The network interface ID.
  * `instance_id` - Attached instance ID.
  * `instance_name` - Resolved instance name.
  * `instance_owner` - Instance owner (e.g. compute:nova).
  * `ipv4` - IPv4 addresses.
  * `ipv6` - IPv6 addresses.
  * `mac_address` - MAC address.
  * `name` - The network interface name.
  * `port_security_enabled` - Whether network interface security is enabled.
  * `project_id` - Project ID.
  * `qos_policy_id` - QoS policy ID.
  * `revision_number` - Revision number (may be filled from origin_data when top-level is zero).
  * `security_groups` - Security group IDs (merged from origin_data when top-level is null).
  * `status` - Network interface status.
  * `tags` - Tags.
  * `tenant_id` - The tenant ID (may be filled from origin_data when top-level is empty).
  * `updated_at` - Last update time (may be filled from origin_data when top-level is empty).
  * `vpc_id` - The VPC ID.
  * `vpc_name` - Resolved VPC name.
* `total` - Total count from the API response (data.count).


