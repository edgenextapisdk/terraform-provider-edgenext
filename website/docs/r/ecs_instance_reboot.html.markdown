---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_instance_reboot"
sidebar_current: "docs-edgenext-resource-ecs_instance_reboot"
description: |-
  Use this resource to trigger reboot actions for an existing ECS instance.
---

# edgenext_ecs_instance_reboot

Use this resource to trigger reboot actions for an existing ECS instance.

## Example Usage

```hcl
data "edgenext_ecs_instances" "all" {
  limit = 1
}

resource "edgenext_ecs_instance_reboot" "example" {
  instance_id = data.edgenext_ecs_instances.all.instances[0].id
  trigger     = timestamp()
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) The instance ID.
* `reboot_type` - (Optional, String) Reboot action type. Currently only reboot_soft is supported.
* `trigger` - (Optional, String) Update this field to trigger reboot again.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_name` - Instance name.
* `status` - Current instance status from detail API.


