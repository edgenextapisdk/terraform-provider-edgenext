# en_ecs_image

Provides an EdgeNext ECS image resource. This allows you to manage images within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_image" "example" {
  name = "example-image"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the image.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the image.
