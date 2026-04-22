---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_tag"
sidebar_current: "docs-edgenext-resource-ecs_tag"
description: |-
  Use this resource to create and manage ECS global tags.
---

# edgenext_ecs_tag

Use this resource to create and manage ECS global tags.

## Example Usage

```hcl
resource "edgenext_ecs_tag" "example" {
  key   = "env"
  value = "dev"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, String) Tag key. Cannot be changed after creation.
* `value` - (Required, String) Tag value. Cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `tag_id/key/value`.

```shell
terraform import edgenext_ecs_tag.example 52/env/dev
```

Argument Reference

* `key` - (Required) Tag key.
* `value` - (Required) Tag value.

Attributes Reference

* `id` - Tag ID.

