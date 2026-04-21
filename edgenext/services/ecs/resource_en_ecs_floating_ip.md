# en_ecs_floating_ip

Provides an EdgeNext ECS floating_ip resource. This allows you to manage floating_ips within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_floating_ip" "example" {
  name = "example-floating_ip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the floating_ip.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the floating_ip.
