Use this data source to query ECS disks via **GET** `/ecs/openapi/v2/volume/list`.

Example Usage

```hcl
data "edgenext_ecs_disks" "example" {
  name       = edgenext_ecs_disk.example.name
  page_num   = 1
  page_size  = 10
}

resource "edgenext_ecs_disk" "example" {
  name        = "example-disk"
  volume_type = "SSD"
  size        = 50
}
```
