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
resource "edgenext_ecs_network_interface" "example" {
  region                = "tokyo-a"
  name                  = "example-eni"
  description           = "for application"
  network_id            = "0e07db22-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id             = "50a0f20a-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  device_id             = "80e47fca-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  port_security_enabled = true
  security_groups       = ["aa2b7c0d-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]
  floating_ip_address   = "148.222.161.86"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Port name.
* `network_id` - (Required, String) VPC network ID. Cannot be changed after creation.
* `region` - (Required, String, ForceNew) The region of the port.
* `subnet_id` - (Required, String) Subnet ID for the primary fixed IP. Cannot be changed after creation.
* `description` - (Optional, String) Port description.
* `device_id` - (Optional, String) Attached server ID (instance ID).
* `floating_ip_address` - (Optional, String) Floating IP address bound to this port.
* `port_security_enabled` - (Optional, Bool) Whether port security is enabled.
* `security_groups` - (Optional, List: [`String`]) Security group IDs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `admin_state_up` - Administrative state of the port.
* `binding_vnic_type` - VNIC binding type.
* `created_at` - Creation time.
* `device_owner` - Device owner (e.g. compute:nova).
* `fixed_ips` - Fixed IP assignments.
  * `floating_ip` - Associated floating IP if any.
  * `ip_address` - Fixed IP address.
  * `subnet_id` - Subnet ID.
* `ipv4` - IPv4 addresses.
* `ipv6` - IPv6 addresses.
* `mac_address` - MAC address.
* `network_name` - Resolved network name.
* `project_id` - Project ID.
* `qos_policy_id` - QoS policy ID.
* `revision_number` - Revision number.
* `server_name` - Resolved server name when attached.
* `status` - Port status.
* `tags` - Tags.
* `tenant_id` - Tenant ID.
* `updated_at` - Last update time.


## Import

Import format is `region/port_id`.

```shell
terraform import edgenext_ecs_network_interface.example tokyo-a/29faf396-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Port name.
* `description` - (Optional) Port description.
* `network_id` - (Required) Network ID. Changing this value recreates the resource.
* `subnet_id` - (Required) Subnet ID. Changing this value recreates the resource.
* `device_id` - (Optional) Server ID to attach.
* `port_security_enabled` - (Optional) Port security switch.
* `security_groups` - (Optional) Security group IDs.
* `floating_ip_address` - (Optional) Floating IP address to bind.

Attributes Reference

* `id` - Importable ID in format `region/port_id`.
* `tenant_id`, `project_id`, `status`, `device_owner`
* `fixed_ips` - Fixed IP list with `subnet_id`, `ip_address`, `floating_ip`.
* `qos_policy_id`, `tags`, `created_at`, `updated_at`, `revision_number`
* `mac_address`, `binding_vnic_type`, `server_name`, `network_name`, `ipv4`, `ipv6`

