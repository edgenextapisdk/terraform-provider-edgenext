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

Argument Reference

* `name` - (Optional) Disk name filter; empty string lists all.
* `page_num` - (Optional) Page number, default 1.
* `page_size` - (Optional) Page size, default 10.

Attributes Reference

* `total` - Total count from API `data.total`.
* `disks` - List of disks for the page; see schema in `data_source_en_ecs_disks.go` for nested fields (`disk_type`, `attachment`, etc.).
