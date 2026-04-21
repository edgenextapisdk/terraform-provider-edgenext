# en_ecs_security_group

Provides an EdgeNext ECS security_group resource. This allows you to manage security_groups within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_security_group" "example" {
  name = "example-security_group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security_group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security_group.
