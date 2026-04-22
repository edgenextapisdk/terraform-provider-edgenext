---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_security_group_rules"
sidebar_current: "docs-edgenext-datasource-ecs_security_group_rules"
description: |-
  Use this data source to query rules of a specific ECS security group.
---

# edgenext_ecs_security_group_rules

Use this data source to query rules of a specific ECS security group.

## Example Usage

```hcl
data "edgenext_ecs_security_group_rules" "example" {
  region = "tokyo-a"
  id     = "de62db3d-c771-4d87-a233-344ef9574e76"
}

output "first_rule_id" {
  value = try(data.edgenext_ecs_security_group_rules.example.security_group_rules[0].id, null)
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) The security group ID.
* `region` - (Required, String) region description

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_group_rules` - A list of security group rules.
  * `created_at` - Creation time.
  * `description` - Rule description.
  * `direction` - Traffic direction (ingress/egress).
  * `ethertype` - IP version.
  * `id` - The rule ID.
  * `port_range_max` - Maximum port.
  * `port_range_min` - Minimum port.
  * `project_id` - Project ID.
  * `protocol` - Protocol name.
  * `remote_group_id` - Remote security group ID.
  * `remote_ip_prefix` - Remote IP prefix.
  * `revision_number` - Rule revision number.
  * `security_group_id` - The security group ID this rule belongs to.
  * `tags` - Rule tags.
  * `tenant_id` - The rule tenant ID.
  * `updated_at` - Last update time.


