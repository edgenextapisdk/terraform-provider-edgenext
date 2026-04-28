---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_instance_tag"
sidebar_current: "docs-edgenext-resource-ecs_instance_tag"
description: |-
  Use this resource to bind existing tag IDs to an ECS instance.
---

# edgenext_ecs_instance_tag

Use this resource to bind existing tag IDs to an ECS instance.

## Example Usage

```hcl
resource "edgenext_ecs_tag" "example" {
  for_each = {
    env  = "dev"
    team = "platform"
  }
  tag_key   = each.key
  tag_value = each.value
}

data "edgenext_ecs_instances" "all" {
  limit = 1
}

resource "edgenext_ecs_instance_tag" "example" {
  instance_id   = data.edgenext_ecs_instances.all.instances[0].id
  instance_name = data.edgenext_ecs_instances.all.instances[0].name
  tag_ids       = [for t in values(edgenext_ecs_tag.example) : tonumber(t.id)]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) The target instance ID. Cannot be changed after creation.
* `instance_name` - (Required, String) The target instance name. Cannot be changed after creation.
* `tag_ids` - (Required, Set: [`Int`]) Tag IDs to bind to the target instance. Element order is ignored when comparing updates.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_type` - The instance type returned by query API.
* `tag_count` - The number of tags on this instance.
* `tags` - Detailed tag items for the instance.
  * `id` - Tag item ID.
  * `key` - Tag item key.
  * `value` - Tag item value.


## Import

Import format is `instance_id`.

```shell
terraform import edgenext_ecs_instance_tag.example 0d4dd8b5-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `instance_id` - (Required) Target instance ID. Cannot be changed after creation.
* `instance_name` - (Required) Target instance name. Cannot be changed after creation.
* `tag_ids` - (Required) Tag ID list to bind.

Attributes Reference

* `id` - Uses `instance_id`.
* `instance_type` - Instance type returned by query API.
* `tag_count` - Number of tags on this instance.
* `tags` - Tag details with `id`, `key`, and `value`.

