Use this resource to trigger reboot actions for an existing ECS instance.

Example Usage

```hcl
data "edgenext_ecs_instances" "all" {
  limit = 1
}

resource "edgenext_ecs_instance_reboot" "example" {
  instance_id = data.edgenext_ecs_instances.all.instances[0].id
  trigger     = timestamp()
}
```

Argument Reference

* `instance_id` - (Required) Target instance ID. Cannot be changed after creation.
* `reboot_type` - (Optional) Reboot action type. Default is `reboot_soft`.
* `trigger` - (Optional) Any value used to trigger reboot again on update.

Attributes Reference

* `id` - Uses `instance_id`.
* `status` - Current instance status from detail API.
* `instance_name` - Instance name.
