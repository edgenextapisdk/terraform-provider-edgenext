Use this resource to create and manage ECS custom images.

Example Usage

```hcl
resource "edgenext_ecs_image" "example" {
  region      = "tokyo-a"
  name        = "example-image"
  instance_id = "80e47fca-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  description = "created from instance"
}
```

Import

Import format is `region/image_id`.

```shell
terraform import edgenext_ecs_image.example tokyo-a/7b6387c5-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Image name.
* `instance_id` - (Optional) Source instance ID.
* `description` - (Optional) Image description.

Attributes Reference

* `id` - Image ID.
* `os_distro` - OS distribution reported by API.
