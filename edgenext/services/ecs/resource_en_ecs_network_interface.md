# en_ecs_network_interface

Provides an EdgeNext ECS network_interface resource. This allows you to manage network_interfaces within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_network_interface" "example" {
  name = "example-network_interface"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network_interface.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the network_interface.
