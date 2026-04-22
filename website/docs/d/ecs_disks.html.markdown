---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_disks"
sidebar_current: "docs-edgenext-datasource-ecs_disks"
description: |-
  Use this data source to query ECS disks.
---

# edgenext_ecs_disks

Use this data source to query ECS disks.

## Example Usage

```hcl
data "edgenext_ecs_disks" "example" {
  region = "tokyo-a"
  name   = "example-disk"
}

output "first_disk_id" {
  value = try(data.edgenext_ecs_disks.example.disks[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `ids` - (Optional, List: [`String`]) A list of disk IDs.
* `name` - (Optional, String) The name to filter disks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `disks` - A list of ECS disks.
  * `id` - The ID of the disk.
  * `name` - The name of the disk.


