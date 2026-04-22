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
  region = "tokyo-a"
  name   = ""
  limit  = 10
}

output "first_port_id" {
  value = try(data.edgenext_ecs_network_interfaces.example.network_interfaces[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `limit` - (Optional, Int) Maximum number of ports to return.
* `name` - (Optional, String) Filter by port name (partial match per API behavior).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `network_interfaces` - A list of ECS network interfaces (Neutron ports).
  * `admin_state_up` - Administrative state of the port.
  * `binding_vnic_type` - VNIC binding type (from binding:vnic_type; may be filled from origin_data).
  * `created_at` - Creation time.
  * `description` - Description.
  * `device_id` - Attached device (e.g. instance) ID.
  * `device_owner` - Device owner string (e.g. compute:nova).
  * `fixed_ips` - Fixed IP assignments.
    * `floating_ip` - Associated floating IP if present.
    * `ip_address` - Fixed IP address.
    * `subnet_id` - Subnet ID.
  * `id` - The port ID.
  * `ipv4` - IPv4 addresses.
  * `ipv6` - IPv6 addresses.
  * `mac_address` - MAC address.
  * `name` - The port name.
  * `network_id` - The network (VPC) ID.
  * `network_name` - Resolved network name.
  * `port_security_enabled` - Whether port security is enabled.
  * `project_id` - Project ID.
  * `qos_policy_id` - QoS policy ID.
  * `revision_number` - Revision number (may be filled from origin_data when top-level is zero).
  * `security_groups` - Security group IDs (merged from origin_data when top-level is null).
  * `server_name` - Resolved server (instance) name.
  * `status` - Port status.
  * `tags` - Tags.
  * `tenant_id` - The tenant ID (may be filled from origin_data when top-level is empty).
  * `updated_at` - Last update time (may be filled from origin_data when top-level is empty).
* `total` - Total count from the API response (data.count).


