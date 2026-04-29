Use this resource to bind existing tag IDs to an ECS instance.

Example Usage

```hcl
resource "edgenext_ecs_tag" "example" {
  for_each = {
    env  = "dev"
    team = "platform"
  }
  tag_key   = each.key
  tag_value = each.value
}

data "edgenext_ecs_instances" "all" {
  limit = 1
}

resource "edgenext_ecs_instance_tag" "example" {
  instance_id   = data.edgenext_ecs_instances.all.instances[0].id
  instance_name = data.edgenext_ecs_instances.all.instances[0].name
  tag_ids       = [for t in values(edgenext_ecs_tag.example) : tonumber(t.id)]
}
```

Import

Import format is `instance_id`.

```shell
terraform import edgenext_ecs_instance_tag.example 0d4dd8b5-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
