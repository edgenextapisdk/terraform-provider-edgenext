---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_instance_power"
sidebar_current: "docs-edgenext-resource-ecs_instance_power"
description: |-
  Use this resource to control the power state of an existing ECS instance.
---

# edgenext_ecs_instance_power

Use this resource to control the power state of an existing ECS instance.

## Example Usage

```hcl
data "edgenext_ecs_instances" "all" {
  limit = 1
}

resource "edgenext_ecs_instance_power" "example" {
  instance_id   = data.edgenext_ecs_instances.all.instances[0].id
  desired_state = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

* `desired_state` - (Required, String) Desired instance power state. Valid values: ACTIVE, SHUTOFF.
* `instance_id` - (Required, String, ForceNew) The instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_name` - Instance name.
* `status` - Current instance status from detail API.


