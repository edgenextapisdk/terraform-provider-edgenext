# en_ecs_vpcs

This data source provides a list of EdgeNext ECS vpcs in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_vpcs" "example" {
  name_regex = "example"
}

output "vpc_id" {
  value = data.edgenext_ecs_vpcs.example.vpcs.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of vpc IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `vpcs` - A list of vpcs. Each element contains the following attributes:
  * `id` - The ID of the vpc.
  * `name` - The name of the vpc.
