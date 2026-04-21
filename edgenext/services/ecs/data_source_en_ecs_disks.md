# en_ecs_disks

This data source provides a list of EdgeNext ECS disks in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_disks" "example" {
  name_regex = "example"
}

output "disk_id" {
  value = data.edgenext_ecs_disks.example.disks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of disk IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `disks` - A list of disks. Each element contains the following attributes:
  * `id` - The ID of the disk.
  * `name` - The name of the disk.
