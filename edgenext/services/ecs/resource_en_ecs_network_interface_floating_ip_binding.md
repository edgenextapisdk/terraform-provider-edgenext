Use this resource to bind a floating IP to a network interface.

Example Usage

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

Import

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
