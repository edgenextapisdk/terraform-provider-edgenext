# en_ecs_security_group_rules

This data source provides security group rules from an EdgeNext ECS security group.

## Example Usage

```hcl
data "edgenext_ecs_security_group_rules" "example" {
  id = "de62db3d-c771-4d87-a233-344ef9574e76"
}

output "first_rule_id" {
  value = try(data.edgenext_ecs_security_group_rules.example.security_group_rules[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The security group ID.

## Attributes Reference

The following attributes are exported:

* `security_group_rules` - A list of security group rules. Each element contains the following attributes:
  * `id`
  * `tenant_id`
  * `security_group_id`
  * `ethertype`
  * `direction`
  * `protocol`
  * `port_range_min`
  * `port_range_max`
  * `remote_ip_prefix`
  * `remote_group_id`
  * `description`
  * `tags`
  * `created_at`
  * `updated_at`
  * `revision_number`
  * `project_id`

