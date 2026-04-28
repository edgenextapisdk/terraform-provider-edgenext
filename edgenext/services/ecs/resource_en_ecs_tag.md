Use this resource to create and manage ECS global tags.

Example Usage

```hcl
resource "edgenext_ecs_tag" "example" {
  tag_key   = "env"
  tag_value = "dev"
}
```

Import

Import format is `tag_id/tag_key/tag_value`.

```shell
terraform import edgenext_ecs_tag.example 52/env/dev
```

Argument Reference

* `tag_key` - (Required) Tag key. Cannot be changed after creation.
* `tag_value` - (Required) Tag value. Cannot be changed after creation.

Attributes Reference

* `id` - Tag ID.
