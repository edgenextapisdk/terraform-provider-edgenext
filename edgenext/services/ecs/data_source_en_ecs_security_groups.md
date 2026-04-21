# en_ecs_security_groups

This data source provides a list of EdgeNext ECS security_groups in your current credentials.

## Example Usage

```hcl
data "edgenext_ecs_security_groups" "example" {
  name_regex = "example"
}

output "security_group_id" {
  value = data.edgenext_ecs_security_groups.example.security_groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by name.
* `ids` - (Optional) A list of security_group IDs.

## Attributes Reference

The following attributes are exported:

* `id` - The provider-assigned unique ID for this managed resource.
* `security_groups` - A list of security_groups. Each element contains the following attributes:
  * `id` - The ID of the security_group.
  * `name` - The name of the security_group.
