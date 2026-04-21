# en_ecs_instance

Provides an EdgeNext ECS instance resource. This allows you to manage instances within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_instance" "example" {
  name = "example-instance"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance.
