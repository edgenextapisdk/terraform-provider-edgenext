---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_instance_tags"
sidebar_current: "docs-edgenext-datasource-ecs_instance_tags"
description: |-
  Use this data source to query ECS instances by tag filters.
---

# edgenext_ecs_instance_tags

Use this data source to query ECS instances by tag filters.

## Example Usage

```hcl
data "edgenext_ecs_instance_tags" "example" {
  tag_key   = "env"
  tag_value = "dev"
  page_num  = 1
  page_size = 10
}

resource "edgenext_ecs_instance_tag" "binding" {
  instance_id   = data.edgenext_ecs_instance_tags.example.instance_tags[0].instance_id
  instance_name = data.edgenext_ecs_instance_tags.example.instance_tags[0].instance_name
  tag_ids       = [for t in data.edgenext_ecs_instance_tags.example.instance_tags[0].tags : t.id]
}
```

## Argument Reference

The following arguments are supported:

* `page_num` - (Optional, Int) Page number for instance tag listing.
* `page_size` - (Optional, Int) Page size for instance tag listing.
* `tag_id` - (Optional, Int) The tag ID to filter instances.
* `tag_key` - (Optional, String) The tag key to filter instances.
* `tag_value` - (Optional, String) The tag value to filter instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_tags` - A list of instances matched by tags.
  * `id` - The record ID.
  * `instance_id` - The instance ID.
  * `instance_name` - The instance name.
  * `instance_type` - The instance type, e.g. ECS.
  * `region` - The instance region.
  * `tag_count` - The number of tags on this instance.
  * `tags` - Detailed tag items for the instance.
    * `id` - Tag item ID.
    * `key` - Tag item key.
    * `value` - Tag item value.
* `total` - Total number of matched instances.


