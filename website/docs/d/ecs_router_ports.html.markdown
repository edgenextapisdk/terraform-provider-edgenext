---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_router_ports"
sidebar_current: "docs-edgenext-datasource-ecs_router_ports"
description: |-
  Use this data source to query ports attached to a specific ECS router.
---

# edgenext_ecs_router_ports

Use this data source to query ports attached to a specific ECS router.

## Example Usage

```hcl
data "edgenext_ecs_router_ports" "example" {
  router_id = data.edgenext_ecs_routers.example.routers[0].id
}

data "edgenext_ecs_routers" "example" {
  router_name = "default-router"
  limit       = 1
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required, String) The router ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ports` - A list of router ports.
  * `created_at` - Creation time.
  * `id` - Port ID.
  * `ip_address` - Port IP address.
  * `mac_address` - Port MAC address.
  * `name` - Port name.
  * `status` - Port status.
  * `vpc_name` - VPC name.
* `total` - Total number of router ports.


