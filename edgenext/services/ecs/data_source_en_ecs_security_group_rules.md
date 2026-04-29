Use this data source to query rules of a specific ECS security group.

Example Usage

```hcl
data "edgenext_ecs_security_group_rules" "example" {
  id = data.edgenext_ecs_security_groups.example.security_groups[0].id
}

data "edgenext_ecs_security_groups" "example" {
  name  = "default"
  limit = 1
}
```
