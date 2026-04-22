---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_security_groups"
sidebar_current: "docs-edgenext-datasource-ecs_security_groups"
description: |-
  Use this data source to query ECS security groups.
---

# edgenext_ecs_security_groups

Use this data source to query ECS security groups.

## Example Usage

```hcl
data "edgenext_ecs_security_groups" "example" {
  region = "tokyo-a"
  name   = ""
  limit  = 10
}

output "first_security_group_id" {
  value = try(data.edgenext_ecs_security_groups.example.security_groups[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required, String) region description
* `limit` - (Optional, Int) Maximum number of security_groups to return.
* `name` - (Optional, String) The name to filter security_groups.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - A list of ECS security_groups.
  * `created_at` - Creation time.
  * `description` - The security group description.
  * `id` - The ID of the security_group.
  * `name` - The name of the security_group.
  * `project_id` - Project ID.
  * `revision_number` - Revision number.
  * `security_group_rules` - Security group rules.
    * `created_at` - Rule creation time.
    * `description` - Rule description.
    * `direction` - Traffic direction.
    * `ethertype` - IP version type.
    * `id` - Rule ID.
    * `port_range_max` - Maximum port.
    * `port_range_min` - Minimum port.
    * `project_id` - Rule project ID.
    * `protocol` - Protocol name.
    * `remote_group_id` - Remote security group ID.
    * `remote_ip_prefix` - Remote IP CIDR.
    * `revision_number` - Rule revision number.
    * `security_group_id` - Security group ID.
    * `tags` - Rule tags.
    * `tenant_id` - Rule tenant ID.
    * `updated_at` - Rule update time.
  * `tags` - A list of tag strings.
  * `tenant_id` - The tenant ID.
  * `updated_at` - Last update time.
* `total` - Total number of matched security groups.


