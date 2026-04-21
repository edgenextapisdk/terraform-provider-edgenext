# en_ecs_tags

This data source provides a list of EdgeNext ECS tags in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_tags" "example" {
  name_regex = "example"
}

output "tag_id" {
  value = data.edgenext_ecs_tags.example.tags.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of tag IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `tags` - A list of tags. Each element contains the following attributes:
  * `id` - The ID of the tag.
  * `name` - The name of the tag.
