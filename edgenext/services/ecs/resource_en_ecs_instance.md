Use this resource to create and manage ECS instances.

Example Usage

```hcl
resource "edgenext_ecs_instance" "example" {
  name            = "example-instance"
  flavor_ref      = "s1.small"
  image_ref       = data.edgenext_ecs_images.example.images[0].id
  admin_pass      = "SecurePass123!"
  bandwidth       = 5
  key_name        = edgenext_ecs_key_pair.example.name
  networks        = [data.edgenext_ecs_vpcs.all.vpcs[0].id]
  security_groups = [data.edgenext_ecs_security_groups.all.security_groups[0].id]
}

resource "edgenext_ecs_key_pair" "example" {
  name = "example-key"
}

data "edgenext_ecs_images" "example" {
  visibility = "public"
  page_size  = 1
}

data "edgenext_ecs_vpcs" "all" {
  limit = 1
}

data "edgenext_ecs_security_groups" "all" {
  limit = 1
}
```

Import

Import format is `instance_id`.

```shell
terraform import edgenext_ecs_instance.example 80e47fca-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `name` - (Required) Instance name.
* `flavor_ref` - (Required) Flavor ID or name.
* `image_ref` - (Required) Image ID or name.
* `admin_pass` - (Required) Initial admin password.
* `key_name` - (Optional) Key pair name.
* `project_id` - (Optional) Project ID.
* `bandwidth` - (Optional) Public bandwidth in Mbps.
* `networks` - (Optional) Network IDs.
* `security_groups` - (Optional) Security group IDs.

Attributes Reference

* `id` - Instance ID.
* `status` - Current instance status.
