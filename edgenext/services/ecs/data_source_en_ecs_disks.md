Use this data source to query ECS disks via **GET** `/ecs/openapi/v2/volume/list`.

Example Usage

```hcl
data "edgenext_ecs_disks" "example" {
  name       = ""
  page_num   = 1
  page_size  = 10
}
```
