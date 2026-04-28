---
subcategory: "Elastic Compute Service (ECS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ecs_security_group_rule"
sidebar_current: "docs-edgenext-resource-ecs_security_group_rule"
description: |-
  Use this resource to manage a single ECS security group rule.
---

# edgenext_ecs_security_group_rule

Use this resource to manage a single ECS security group rule.

## Example Usage

```hcl
resource "edgenext_ecs_security_group_rule" "example" {
  security_group_id = edgenext_ecs_security_group.example.id
  protocol          = "tcp"
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  description       = "allow ssh"
}

resource "edgenext_ecs_security_group" "example" {
  name = "example-sg"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, String) Traffic direction: ingress or egress. Cannot be changed after creation.
* `ethertype` - (Required, String) IP version (e.g. IPv4, IPv6). Cannot be changed after creation.
* `port_range_max` - (Required, Int) Maximum port number. Cannot be changed after creation.
* `port_range_min` - (Required, Int) Minimum port number. Cannot be changed after creation.
* `protocol` - (Required, String) Protocol name (e.g. tcp, udp, icmp). Cannot be changed after creation.
* `security_group_id` - (Required, String) The security group ID this rule belongs to. Cannot be changed after creation.
* `description` - (Optional, String) Rule description. Cannot be changed after creation.
* `remote_group_id` - (Optional, String) Remote security group ID. Leave empty when using remote_ip_prefix only. Cannot be changed after creation.
* `remote_ip_prefix` - (Optional, String) Remote CIDR (e.g. 192.168.0.0/24). Cannot be changed after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Creation timestamp.
* `project_id` - Project ID of the rule.
* `revision_number` - Revision number.
* `tags` - Rule tags.
* `tenant_id` - Tenant ID of the rule.
* `updated_at` - Last update timestamp.


## Import

Import format is `security_group_id/rule_id`.

```shell
terraform import edgenext_ecs_security_group_rule.example 12f8f386-xxxx-xxxx-xxxx-xxxxxxxxxxxx/df58bf0a-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `security_group_id` - (Required) Security group ID. Cannot be changed after creation.
* `protocol` - (Required) Protocol, for example `tcp`, `udp`, `icmp`. Cannot be changed after creation.
* `direction` - (Required) Rule direction, `ingress` or `egress`. Cannot be changed after creation.
* `ethertype` - (Required) IP type, such as `IPv4` or `IPv6`. Cannot be changed after creation.
* `port_range_min` - (Required) Minimum port. Cannot be changed after creation.
* `port_range_max` - (Required) Maximum port. Cannot be changed after creation.
* `remote_ip_prefix` - (Optional) Remote CIDR. Cannot be changed after creation.
* `remote_group_id` - (Optional) Remote security group ID. Cannot be changed after creation.
* `description` - (Optional) Rule description. Cannot be changed after creation.

Attributes Reference

* `id` - Security group rule ID.
* `tenant_id`, `project_id`, `revision_number`, `created_at`, `updated_at`, `tags`

