Use this data source to query ECS tags.

Example Usage

```hcl
data "edgenext_ecs_tags" "example" {
  tag_key   = edgenext_ecs_tag.example.tag_key
  tag_value = edgenext_ecs_tag.example.tag_value
  page_num  = 1
  page_size = 10
}

resource "edgenext_ecs_tag" "example" {
  tag_key   = "env"
  tag_value = "dev"
}
```
