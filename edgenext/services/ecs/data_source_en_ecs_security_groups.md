Use this data source to query ECS security groups.

Example Usage

```hcl
data "edgenext_ecs_security_groups" "example" {
  name  = edgenext_ecs_security_group.example.name
  limit = 10
}

resource "edgenext_ecs_security_group" "example" {
  name = "example-sg"
}
```
