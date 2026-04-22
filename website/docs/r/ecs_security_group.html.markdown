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
  region      = "tokyo-a"
  name        = "example-sg"
  description = "security group for app"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) name description
* `region` - (Required, String, ForceNew) region description
* `description` - (Optional, String) description description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `region/name`.

```shell
terraform import edgenext_ecs_security_group.example tokyo-a/example-sg
```

Argument Reference

* `region` - (Required) Region.
* `name` - (Required) Security group name.
* `description` - (Optional) Security group description.

Attributes Reference

* `id` - Security group ID.

