# en_ecs_key_pairs

This data source provides a list of EdgeNext ECS key_pairs in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_key_pairs" "example" {
  name_regex = "example"
}

output "key_pair_id" {
  value = data.edgenext_ecs_key_pairs.example.key_pairs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of key_pair IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `key_pairs` - A list of key_pairs. Each element contains the following attributes:
  * `id` - The ID of the key_pair.
  * `name` - The name of the key_pair.
