# en_ecs_disk

Provides an EdgeNext ECS disk resource. This allows you to manage disks within your ECS environment.

## Example Usage

```hcl
resource "edgenext_ecs_disk" "example" {
  name = "example-disk"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the disk.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the disk.
