Use this resource to control the power state of an existing ECS instance.

Example Usage

```hcl
data "edgenext_ecs_instances" "all" {
  limit = 1
}

resource "edgenext_ecs_instance_power" "example" {
  instance_id   = data.edgenext_ecs_instances.all.instances[0].id
  desired_state = "ACTIVE"
}
```

Argument Reference

* `instance_id` - (Required) Target instance ID. Cannot be changed after creation.
* `desired_state` - (Required) Desired power state. Valid values: `ACTIVE`, `SHUTOFF`.

Attributes Reference

* `id` - Uses `instance_id`.
* `status` - Current instance status from detail API.
* `instance_name` - Instance name.
