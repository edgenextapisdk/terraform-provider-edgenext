---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_resource_tags"
sidebar_current: "docs-edgenext-datasource-ecs_resource_tags"
description: |-
  Use this data source to query ECS resources by tag filters.
---

# edgenext_ecs_resource_tags

Use this data source to query ECS resources by tag filters.

## Example Usage

```hcl
data "edgenext_ecs_resource_tags" "example" {
  region    = "tokyo-a"
  tag_key   = "env"
  tag_value = "dev"
  page_num  = 1
  page_size = 10
}

output "first_tagged_resource" {
  value = try(data.edgenext_ecs_resource_tags.example.resource_tags[0].resource_id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `page_num` - (Optional, Int) Page number for resource tag listing.
* `page_size` - (Optional, Int) Page size for resource tag listing.
* `tag_key` - (Optional, String) The tag key to filter resources.
* `tag_value` - (Optional, String) The tag value to filter resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_tags` - A list of resources matched by tags.
  * `id` - The record ID.
  * `product_type` - The product type, e.g. ECS.
  * `region` - The resource region.
  * `resource_id` - The resource ID.
  * `resource_name` - The resource name.
  * `tag_count` - The number of tags on this resource.
  * `tags` - Detailed tag items for the resource.
    * `id` - Tag item ID.
    * `key` - Tag item key.
    * `value` - Tag item value.
* `total` - Total number of matched resources.


