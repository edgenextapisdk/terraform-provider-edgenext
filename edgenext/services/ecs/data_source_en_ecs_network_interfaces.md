# en_ecs_network_interfaces

This data source provides a list of EdgeNext ECS network_interfaces in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_network_interfaces" "example" {
  name_regex = "example"
}

output "network_interface_id" {
  value = data.edgenext_ecs_network_interfaces.example.network_interfaces.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of network_interface IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `network_interfaces` - A list of network_interfaces. Each element contains the following attributes:
  * `id` - The ID of the network_interface.
  * `name` - The name of the network_interface.
