Use this resource to create and manage ECS custom images.

Example Usage

```hcl
resource "edgenext_ecs_image" "example" {
  name        = "example-image"
  instance_id = data.edgenext_ecs_instances.example.instances[0].id
  description = "created from instance"
}

data "edgenext_ecs_instances" "example" {
  limit = 1
}
```

Import

Import format is `image_id`.

```shell
terraform import edgenext_ecs_image.example 7b6387c5-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `name` - (Required) Image name.
* `instance_id` - (Optional) Source instance ID.
* `description` - (Optional) Image description.

Attributes Reference

* `id` - Image ID.
* `os_distro` - OS distribution reported by API.
