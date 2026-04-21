# en_ecs_floating_ips

This data source provides a list of EdgeNext ECS floating_ips in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_floating_ips" "example" {
  name_regex = "example"
}

output "floating_ip_id" {
  value = data.edgenext_ecs_floating_ips.example.floating_ips.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of floating_ip IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `floating_ips` - A list of floating_ips. Each element contains the following attributes:
  * `id` - The ID of the floating_ip.
  * `name` - The name of the floating_ip.
