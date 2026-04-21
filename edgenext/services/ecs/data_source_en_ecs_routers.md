# en_ecs_routers

This data source provides a list of EdgeNext ECS routers in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_routers" "example" {
  name_regex = "example"
}

output "router_id" {
  value = data.edgenext_ecs_routers.example.routers.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of router IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `routers` - A list of routers. Each element contains the following attributes:
  * `id` - The ID of the router.
  * `name` - The name of the router.
