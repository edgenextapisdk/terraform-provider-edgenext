---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_security_group"
sidebar_current: "docs-edgenext-resource-ecs_security_group"
description: |-
  Use this resource to create and manage ECS security groups.
---

# edgenext_ecs_security_group

Use this resource to create and manage ECS security groups.

## Example Usage

```hcl
resource "edgenext_ecs_security_group" "example" {
  name        = "example-sg"
  description = "security group for app"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) name description
* `description` - (Optional, String) description description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `security_group_id`.

```shell
terraform import edgenext_ecs_security_group.example 2af2b1e5-344f-4184-9173-cf1b5d43bf7d
```

Argument Reference

* `name` - (Required) Security group name.
* `description` - (Optional) Security group description.

Attributes Reference

* `id` - Security group ID.

