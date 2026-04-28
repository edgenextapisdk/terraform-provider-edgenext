---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_router_port"
sidebar_current: "docs-edgenext-resource-ecs_router_port"
description: |-
  Use this resource to attach a subnet to an ECS router.
---

# edgenext_ecs_router_port

Use this resource to attach a subnet to an ECS router.

## Example Usage

```hcl
data "edgenext_ecs_external_gateways" "all" {
  limit = 1
}

resource "edgenext_ecs_router" "example" {
  name                = "example-router"
  external_network_id = data.edgenext_ecs_external_gateways.all.external_gateways[0].id
}

data "edgenext_ecs_vpcs" "all" {
  limit = 1
}

data "edgenext_ecs_vpc_subnets" "all" {
  vpc_id = data.edgenext_ecs_vpcs.all.vpcs[0].id
}

resource "edgenext_ecs_router_port" "example" {
  router_id = edgenext_ecs_router.example.id
  vpc_id    = data.edgenext_ecs_vpcs.all.vpcs[0].id
  subnet_id = data.edgenext_ecs_vpc_subnets.all.subnets[0].id
}
```

## Argument Reference

The following arguments are supported:

* `router_id` - (Required, String) The router ID. Cannot be changed after creation.
* `subnet_id` - (Required, String) The subnet ID to attach. Cannot be changed after creation.
* `vpc_id` - (Required, String) The VPC ID to attach. Cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Port creation time.
* `ip_address` - Port IP address.
* `mac_address` - Port MAC address.
* `name` - Port name.
* `port_id` - The created router port ID.
* `status` - Port status.
* `vpc_name` - VPC name.


## Import

Import format is `router_id/port_id`.

```shell
terraform import edgenext_ecs_router_port.example f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx/74f3a422-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `router_id` - (Required) Router ID. Cannot be changed after creation.
* `vpc_id` - (Required) VPC ID. Cannot be changed after creation.
* `subnet_id` - (Required) Subnet ID. Cannot be changed after creation.

Attributes Reference

* `id` - Router port ID.
* `port_id` - Same as router port ID.
* `name`, `ip_address`, `mac_address`, `vpc_name`, `status`, `created_at`

