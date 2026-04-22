Use this resource to create and manage ECS instances.

Example Usage

```hcl
resource "edgenext_ecs_instance" "example" {
  region          = "tokyo-a"
  name            = "example-instance"
  flavor_ref      = "s1.small"
  image_ref       = "debian-11"
  admin_pass      = "SecurePass123!"
  bandwidth       = 5
  networks        = ["0e07db22-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]
  security_groups = ["aa2b7c0d-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]
}
```

Import

Import format is `region/instance_id`.

```shell
terraform import edgenext_ecs_instance.example tokyo-a/80e47fca-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
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
