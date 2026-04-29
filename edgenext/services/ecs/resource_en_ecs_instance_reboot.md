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
