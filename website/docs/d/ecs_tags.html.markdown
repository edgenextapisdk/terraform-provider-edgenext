---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_tags"
sidebar_current: "docs-edgenext-datasource-ecs_tags"
description: |-
  Use this data source to query ECS tags.
---

# edgenext_ecs_tags

Use this data source to query ECS tags.

## Example Usage

```hcl
data "edgenext_ecs_tags" "example" {
  tag_key   = "env"
  tag_value = "dev"
  page_num  = 1
  page_size = 10
}

output "first_tag_id" {
  value = try(data.edgenext_ecs_tags.example.tags[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `page_num` - (Optional, Int) Page number for tag listing.
* `page_size` - (Optional, Int) Page size for tag listing.
* `tag_key` - (Optional, String) The tag key to filter tags.
* `tag_value` - (Optional, String) The tag value to filter tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tags` - A list of ECS tags.
  * `id` - The ID of the tag.
  * `resource_count` - The count of resources using this tag.
  * `tag_key` - The key of the tag.
  * `tag_value` - The value of the tag.
* `total` - Total number of tags.


