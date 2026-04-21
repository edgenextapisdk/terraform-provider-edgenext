# en_ecs_router

Provides an EdgeNext ECS router resource. This allows you to manage routers within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_router" "example" {
  name = "example-router"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the router.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the router.
