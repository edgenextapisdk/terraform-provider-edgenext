Use this resource to bind existing tag IDs to a specific ECS resource.

Example Usage

```hcl
resource "edgenext_ecs_resource_tag" "example" {
  region        = "tokyo-a"
  resource_uuid = "55d747cd-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  resource_name = "example-instance"
  resource_type = 1
  tag_ids       = [52, 56, 57]
}
```

Import

Import format is `region/resource_uuid/resource_name/resource_type`.

```shell
terraform import edgenext_ecs_resource_tag.example tokyo-a/55d747cd-xxxx-xxxx-xxxx-xxxxxxxxxxxx/example-instance/1
```

Argument Reference

* `region` - (Required) Region.
* `resource_uuid` - (Required) Target resource UUID. Cannot be changed after creation.
* `resource_name` - (Required) Target resource name. Cannot be changed after creation.
* `resource_type` - (Required) Target resource type code. Cannot be changed after creation.
* `tag_ids` - (Required) Tag ID list to bind.

Attributes Reference

* `id` - Uses `resource_uuid`.
