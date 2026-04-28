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
  tag_key   = "env"
  tag_value = "dev"
}
```

## Argument Reference

The following arguments are supported:

* `tag_key` - (Required, String) Tag key. Cannot be changed after creation.
* `tag_value` - (Required, String) Tag value. Cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `tag_id/tag_key/tag_value`.

```shell
terraform import edgenext_ecs_tag.example 52/env/dev
```

Argument Reference

* `tag_key` - (Required) Tag key. Cannot be changed after creation.
* `tag_value` - (Required) Tag value. Cannot be changed after creation.

Attributes Reference

* `id` - Tag ID.

