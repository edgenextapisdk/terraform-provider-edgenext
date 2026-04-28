---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_network_interface"
sidebar_current: "docs-edgenext-resource-ecs_network_interface"
description: |-
  Use this resource to create and manage ECS network interfaces (ENI / port).
---

# edgenext_ecs_network_interface

Use this resource to create and manage ECS network interfaces (ENI / port).

## Example Usage

```hcl
data "edgenext_ecs_vpcs" "all" {
  limit = 1
}

data "edgenext_ecs_vpc_subnets" "all" {
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}

data "edgenext_ecs_security_groups" "all" {
  limit = 1
}

resource "edgenext_ecs_network_interface" "example" {
  name                  = "example-eni"
  description           = "for application"
  vpc_id                = data.edgenext_ecs_vpcs.all.vpcs[0].id
  subnet_id             = data.edgenext_ecs_vpc_subnets.all.subnets[0].id
  port_security_enabled = true
  security_groups       = [data.edgenext_ecs_security_groups.all.security_groups[0].id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Network interface name.
* `subnet_id` - (Required, String) Subnet ID for the primary fixed IP. Cannot be changed after creation.
* `vpc_id` - (Required, String) VPC network ID. Cannot be changed after creation.
* `description` - (Optional, String) Network interface description.
* `port_security_enabled` - (Optional, Bool) Whether network interface security is enabled.
* `security_groups` - (Optional, List: [`String`]) Security group IDs to apply to the network interface.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `admin_state_up` - Administrative state of the network interface.
* `binding_vnic_type` - VNIC binding type.
* `created_at` - Creation time.
* `fixed_ips` - Fixed IP assignments.
  * `floating_ip` - Associated floating IP if present.
  * `ip_address` - Fixed IP address.
  * `subnet_id` - Subnet ID.
* `floating_ip_address` - Floating IP address bound to this network interface.
* `instance_id` - Attached instance ID to the network interface.
* `instance_name` - Resolved instance name when attached.
* `instance_owner` - Instance owner (e.g. compute:nova).
* `ipv4` - IPv4 addresses.
* `ipv6` - IPv6 addresses.
* `mac_address` - MAC address.
* `project_id` - Project ID.
* `qos_policy_id` - QoS policy ID.
* `revision_number` - Revision number.
* `status` - Network interface status.
* `tags` - Tags.
* `tenant_id` - Tenant ID.
* `updated_at` - Last update time.
* `vpc_name` - Resolved VPC name.


## Import

Import format is `port_id`.

```shell
terraform import edgenext_ecs_network_interface.example 29faf396-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `name` - (Required) Port name.
* `description` - (Optional) Port description.
* `vpc_id` - (Required) VPC ID. Cannot be changed after creation.
* `subnet_id` - (Required) Subnet ID. Cannot be changed after creation.
* `port_security_enabled` - (Optional) Port security switch.
* `security_groups` - (Optional) Security group IDs.

Attributes Reference

* `id` - Port ID.
* `instance_id`, `floating_ip_address`, `tenant_id`, `project_id`, `status`, `instance_owner`
* `fixed_ips` - Fixed IP list with `subnet_id`, `ip_address`, `floating_ip`.
* `qos_policy_id`, `tags`, `created_at`, `updated_at`, `revision_number`
* `mac_address`, `binding_vnic_type`, `instance_name`, `vpc_name`, `ipv4`, `ipv6`

