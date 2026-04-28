---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_network_interface_floating_ip_binding"
sidebar_current: "docs-edgenext-resource-ecs_network_interface_floating_ip_binding"
description: |-
  Use this resource to bind a floating IP to a network interface.
---

# edgenext_ecs_network_interface_floating_ip_binding

Use this resource to bind a floating IP to a network interface.

## Example Usage

```hcl
data "edgenext_ecs_network_interfaces" "all" {
  limit = 1
}

data "edgenext_ecs_floating_ips" "all" {
  limit = 1
}

resource "edgenext_ecs_network_interface_floating_ip_binding" "example" {
  network_interface_id = data.edgenext_ecs_network_interfaces.all.network_interfaces[0].id
  floating_ip_address  = data.edgenext_ecs_floating_ips.all.floating_ips[0].floating_ip_address
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip_address` - (Required, String, ForceNew) The floating IP address to bind. Changing this forces a new resource.
* `network_interface_id` - (Required, String, ForceNew) The network interface ID. Changing this forces a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `fixed_ip_address` - The fixed IP address used for floating IP binding.


## Import

Import format is `network_interface_id/floating_ip_address`.

```shell
terraform import edgenext_ecs_network_interface_floating_ip_binding.example 29faf396-xxxx-xxxx-xxxx-xxxxxxxxxxxx/156.246.18.218
```

Argument Reference

* `network_interface_id` - (Required) Network interface ID. Changing this forces a new resource.
* `floating_ip_address` - (Required) Floating IP address to bind. Changing this forces a new resource.

Attributes Reference

* `id` - The binding ID in format `network_interface_id/floating_ip_address`.
* `fixed_ip_address` - Fixed IP address used for this binding.

