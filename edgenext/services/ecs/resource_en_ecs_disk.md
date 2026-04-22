Use this resource to create and manage ECS disks.

Example Usage

```hcl
resource "edgenext_ecs_disk" "example" {
  region      = "tokyo-a"
  name        = "example-disk"
  volume_type = "SSD"
  size        = 50
}
```

Import

Import format is `region/disk_id`.

```shell
terraform import edgenext_ecs_disk.example tokyo-a/2c5c9f8d-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Disk name.
* `volume_type` - (Required) Volume type.
* `size` - (Required) Disk size in GiB.

Attributes Reference

* `id` - Disk ID.
