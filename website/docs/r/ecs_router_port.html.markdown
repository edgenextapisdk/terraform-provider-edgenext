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
resource "edgenext_ecs_router_port" "example" {
  region     = "tokyo-a"
  router_id  = "f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  network_id = "68451a78-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  subnet_id  = "b34fe463-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `network_id` - (Required, String, ForceNew) The network ID to attach.
* `region` - (Required, String, ForceNew) The region of the router.
* `router_id` - (Required, String, ForceNew) The router ID.
* `subnet_id` - (Required, String, ForceNew) The subnet ID to attach.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Port creation time.
* `ip_address` - Port IP address.
* `mac_address` - Port MAC address.
* `name` - Port name.
* `network_name` - Port network name.
* `port_id` - The created router port ID.
* `status` - Port status.


## Import

Import format is `region/router_id/port_id`.

```shell
terraform import edgenext_ecs_router_port.example tokyo-a/f9883769-xxxx-xxxx-xxxx-xxxxxxxxxxxx/74f3a422-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `router_id` - (Required, ForceNew) Router ID.
* `network_id` - (Required, ForceNew) Network ID.
* `subnet_id` - (Required, ForceNew) Subnet ID.

Attributes Reference

* `id` - Router port ID.
* `port_id` - Same as router port ID.
* `name`, `ip_address`, `mac_address`, `network_name`, `status`, `created_at`

