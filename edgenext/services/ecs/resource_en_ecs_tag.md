Use this resource to create and manage ECS global tags.

Example Usage

```hcl
resource "edgenext_ecs_tag" "example" {
  key   = "env"
  value = "dev"
}
```

Import

Import format is `tag_id/key/value`.

```shell
terraform import edgenext_ecs_tag.example 52/env/dev
```

Argument Reference

* `key` - (Required) Tag key.
* `value` - (Required) Tag value.

Attributes Reference

* `id` - Tag ID.
