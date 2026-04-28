Use this resource to create and manage ECS disks.

Example Usage

```hcl
resource "edgenext_ecs_disk" "example" {
  name        = "example-disk"
  volume_type = "SSD"
  size        = 50
}

data "edgenext_ecs_disks" "all" {
  name      = edgenext_ecs_disk.example.name
  page_num  = 1
  page_size = 10
}
```

Import

Import format is `disk_id`.

```shell
terraform import edgenext_ecs_disk.example 2c5c9f8d-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `name` - (Required) Disk name.
* `volume_type` - (Required) Volume type.
* `size` - (Required) Disk size in GiB.

Attributes Reference

* `id` - Disk ID.
