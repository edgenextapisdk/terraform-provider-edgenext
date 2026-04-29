Use this data source to query ECS images.

Example Usage

```hcl
data "edgenext_ecs_images" "example" {
  visibility = "public"
  name       = "Debian"
  page_num   = 1
  page_size  = 10
}
```
