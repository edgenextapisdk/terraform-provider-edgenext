Use this resource to create and manage ECS security groups.

Example Usage

```hcl
resource "edgenext_ecs_security_group" "example" {
  region      = "tokyo-a"
  name        = "example-sg"
  description = "security group for app"
}
```

Import

Import format is `region/name`.

```shell
terraform import edgenext_ecs_security_group.example tokyo-a/example-sg
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Security group name.
* `description` - (Optional) Security group description.

Attributes Reference

* `id` - Security group ID.
