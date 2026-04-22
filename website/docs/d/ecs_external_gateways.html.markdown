---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_external_gateways"
sidebar_current: "docs-edgenext-datasource-ecs_external_gateways"
description: |-
  Use this data source to query ECS external gateway networks.
---

# edgenext_ecs_external_gateways

Use this data source to query ECS external gateway networks.

## Example Usage

```hcl
data "edgenext_ecs_external_gateways" "example" {
  region = "tokyo-a"
  limit  = 10
}

output "external_gateway_total" {
  value = data.edgenext_ecs_external_gateways.example.total
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `limit` - (Optional, Int) Maximum number of external gateways to return.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `external_gateways` - A list of external gateway networks.
  * `admin_state_up` - Whether admin state is up.
  * `availability_zone_hints` - Availability zone hints.
  * `availability_zones` - Availability zones.
  * `created_at` - Creation time.
  * `description` - Description.
  * `id` - Network ID.
  * `ipv4_address_scope` - IPv4 address scope.
  * `ipv6_address_scope` - IPv6 address scope.
  * `is_default` - Whether this is default network.
  * `mtu` - Network MTU.
  * `name` - Network name.
  * `port_security_enabled` - Whether port security is enabled.
  * `project_id` - Project ID.
  * `provider_network_type` - Provider network type.
  * `provider_physical_network` - Provider physical network.
  * `provider_segmentation_id` - Provider segmentation ID.
  * `qos_policy_id` - QoS policy ID.
  * `revision_number` - Revision number.
  * `router_external` - Whether this is an external gateway network.
  * `shared` - Whether network is shared.
  * `status` - Network status.
  * `subnets` - Subnet IDs.
  * `tags` - Tags.
  * `tenant_id` - Tenant ID.
  * `updated_at` - Last update time.
* `total` - Total number of matched external gateways.


