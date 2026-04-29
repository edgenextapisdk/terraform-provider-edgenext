Use this data source to query ECS instances by tag filters.

Example Usage

```hcl
data "edgenext_ecs_instance_tags" "example" {
  tag_key   = "env"
  tag_value = "dev"
  page_num  = 1
  page_size = 10
}

resource "edgenext_ecs_instance_tag" "binding" {
  instance_id   = data.edgenext_ecs_instance_tags.example.instance_tags[0].instance_id
  instance_name = data.edgenext_ecs_instance_tags.example.instance_tags[0].instance_name
  tag_ids       = [for t in data.edgenext_ecs_instance_tags.example.instance_tags[0].tags : t.id]
}
```
