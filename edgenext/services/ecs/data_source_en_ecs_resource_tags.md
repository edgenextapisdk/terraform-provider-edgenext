# en_ecs_resource_tags

This data source provides a list of EdgeNext ECS resources filtered by tag conditions.

## Example Usage

```hcl
data "edgenext_ecs_resource_tags" "example" {
  region    = "Singapore-a"
  tag_key   = "zyx_test"
  tag_value = "zyx_test1"
  page_num  = 1
  page_size = 10
}

output "first_resource_id" {
  value = try(data.edgenext_ecs_resource_tags.example.resource_tags[0].resource_id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) Region name, for example `Singapore-a`.
* `tag_key` - (Optional) Tag key to filter resources.
* `tag_value` - (Optional) Tag value to filter resources.
* `page_num` - (Optional) Page number, defaults to `1`.
* `page_size` - (Optional) Page size, defaults to `10`.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this data source query.
* `total` - Total number of matched resources.
* `resource_tags` - A list of matched resources. Each element contains:
  * `id` - Record ID.
  * `resource_id` - Resource ID.
  * `resource_name` - Resource name.
  * `product_type` - Product type, e.g. `ECS`.
  * `region` - Resource region.
  * `tag_count` - Number of tags on this resource.
  * `tags` - Detailed tag items on the resource:
    * `id` - Tag item ID.
    * `key` - Tag item key.
    * `value` - Tag item value.
