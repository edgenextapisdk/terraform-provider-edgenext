# en_ecs_tag

Provides an EdgeNext ECS tag resource. This allows you to manage tags within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_tag" "example" {
  name = "example-tag"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the tag.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the tag.
