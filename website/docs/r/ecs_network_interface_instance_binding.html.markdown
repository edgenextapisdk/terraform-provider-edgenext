---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_network_interface_instance_binding"
sidebar_current: "docs-edgenext-resource-ecs_network_interface_instance_binding"
description: |-
  Use this resource to bind an ECS instance to a network interface.
---

# edgenext_ecs_network_interface_instance_binding

Use this resource to bind an ECS instance to a network interface.

## Example Usage

```hcl
data "edgenext_ecs_network_interfaces" "all" {
  limit = 1
}

data "edgenext_ecs_instances" "all" {
  limit = 1
}

resource "edgenext_ecs_network_interface_instance_binding" "example" {
  network_interface_id = data.edgenext_ecs_network_interfaces.all.network_interfaces[0].id
  instance_id          = data.edgenext_ecs_instances.all.instances[0].id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The instance ID to bind. Changing this forces a new resource.
* `network_interface_id` - (Required, String, ForceNew) The network interface ID. Changing this forces a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `network_interface_id/instance_id`.

```shell
terraform import edgenext_ecs_network_interface_instance_binding.example 29faf396-xxxx-xxxx-xxxx-xxxxxxxxxxxx/0d4dd8b5-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `network_interface_id` - (Required) Network interface ID. Changing this forces a new resource.
* `instance_id` - (Required) Instance ID to bind. Changing this forces a new resource.

Attributes Reference

* `id` - The binding ID in format `network_interface_id/instance_id`.

