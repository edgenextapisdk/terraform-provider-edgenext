Use this resource to create and manage ECS security groups.

Example Usage

```hcl
resource "edgenext_ecs_security_group" "example" {
  name        = "example-sg"
  description = "security group for app"
}
```

Import

Import format is `security_group_id`.

```shell
terraform import edgenext_ecs_security_group.example 2af2b1e5-344f-4184-9173-cf1b5d43bf7d
```

Argument Reference

* `name` - (Required) Security group name.
* `description` - (Optional) Security group description.

Attributes Reference

* `id` - Security group ID.
