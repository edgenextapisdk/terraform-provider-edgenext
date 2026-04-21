---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_resource_tag"
sidebar_current: "docs-edgenext-resource-ecs_resource_tag"
description: |-
  Use this resource to bind existing tag IDs to a specific ECS resource.
---

# edgenext_ecs_resource_tag

Use this resource to bind existing tag IDs to a specific ECS resource.

## Example Usage

```hcl
resource "edgenext_ecs_resource_tag" "example" {
  region        = "tokyo-a"
  resource_uuid = "55d747cd-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  resource_name = "example-instance"
  resource_type = 1
  tag_ids       = [52, 56, 57]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String, ForceNew) region description
* `resource_name` - (Required, String, ForceNew) The target resource name.
* `resource_type` - (Required, Int, ForceNew) The target resource type code.
* `resource_uuid` - (Required, String, ForceNew) The target resource UUID.
* `tag_ids` - (Required, List: [`Int`]) Tag IDs to bind to the target resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `region/resource_uuid/resource_name/resource_type`.

```shell
terraform import edgenext_ecs_resource_tag.example tokyo-a/55d747cd-xxxx-xxxx-xxxx-xxxxxxxxxxxx/example-instance/1
```

Argument Reference

* `region` - (Required) Region.
* `resource_uuid` - (Required, ForceNew) Target resource UUID.
* `resource_name` - (Required, ForceNew) Target resource name.
* `resource_type` - (Required, ForceNew) Target resource type code.
* `tag_ids` - (Required) Tag ID list to bind.

Attributes Reference

* `id` - Uses `resource_uuid`.

