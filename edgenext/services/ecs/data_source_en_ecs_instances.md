# en_ecs_instances

This data source provides a list of EdgeNext ECS instances in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_instances" "example" {
  name_regex = "example"
}

output "instance_id" {
  value = data.edgenext_ecs_instances.example.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of instance IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - The ID of the instance.
  * `name` - The name of the instance.
