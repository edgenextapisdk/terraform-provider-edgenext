---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_routers"
sidebar_current: "docs-edgenext-datasource-ecs_routers"
description: |-
  Use this data source to query ECS routers.
---

# edgenext_ecs_routers

Use this data source to query ECS routers.

## Example Usage

```hcl
data "edgenext_ecs_routers" "example" {
  region = "tokyo-a"
  name   = "default-router"
  limit  = 10
}

output "first_router_id" {
  value = try(data.edgenext_ecs_routers.example.routers[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `limit` - (Optional, Int) Maximum number of routers to return.
* `name` - (Optional, String) The name to filter routers.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `routers` - A list of ECS routers.
  * `admin_state_up` - Whether router admin state is up.
  * `availability_zone_hints` - Availability zone hints.
  * `availability_zones` - Availability zones.
  * `created_at` - Creation time.
  * `description` - Description.
  * `external_gateway_info` - External gateway info.
    * `enable_snat` - Whether SNAT is enabled.
    * `external_fixed_ips` - External fixed IPs.
      * `ip_address` - External fixed IP address.
      * `subnet_id` - External fixed IP subnet ID.
    * `network_id` - External network ID.
    * `network_name` - External network name.
  * `flavor_id` - Flavor ID.
  * `id` - The ID of the router.
  * `name` - The name of the router.
  * `project_id` - Project ID.
  * `revision_number` - Revision number.
  * `routes` - Static routes.
    * `destination` - Route destination CIDR.
    * `nexthop` - Route next hop.
  * `status` - Router status.
  * `tags` - Tags.
  * `tenant_id` - Tenant ID.
  * `updated_at` - Last update time.
* `total` - Total number of matched routers.


