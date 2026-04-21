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
  region            = "tokyo-a"
  security_group_id = "12f8f386-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  protocol          = "tcp"
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  description       = "allow ssh"
}
```

## Argument Reference

The following arguments are supported:

* `direction` - (Required, String, ForceNew) Traffic direction: ingress or egress.
* `ethertype` - (Required, String, ForceNew) IP version (e.g. IPv4, IPv6).
* `port_range_max` - (Required, Int, ForceNew) Maximum port number.
* `port_range_min` - (Required, Int, ForceNew) Minimum port number.
* `protocol` - (Required, String, ForceNew) Protocol name (e.g. tcp, udp, icmp).
* `region` - (Required, String, ForceNew) The region of the security group.
* `security_group_id` - (Required, String, ForceNew) The security group ID this rule belongs to.
* `description` - (Optional, String, ForceNew) Rule description.
* `remote_group_id` - (Optional, String, ForceNew) Remote security group ID. Leave empty when using remote_ip_prefix only.
* `remote_ip_prefix` - (Optional, String, ForceNew) Remote CIDR (e.g. 192.168.0.0/24).

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

Import format is `region/security_group_id/rule_id`.

```shell
terraform import edgenext_ecs_security_group_rule.example tokyo-a/12f8f386-xxxx-xxxx-xxxx-xxxxxxxxxxxx/df58bf0a-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `security_group_id` - (Required, ForceNew) Security group ID.
* `protocol` - (Required, ForceNew) Protocol, for example `tcp`, `udp`, `icmp`.
* `direction` - (Required, ForceNew) Rule direction, `ingress` or `egress`.
* `ethertype` - (Required, ForceNew) IP type, such as `IPv4` or `IPv6`.
* `port_range_min` - (Required, ForceNew) Minimum port.
* `port_range_max` - (Required, ForceNew) Maximum port.
* `remote_ip_prefix` - (Optional, ForceNew) Remote CIDR.
* `remote_group_id` - (Optional, ForceNew) Remote security group ID.
* `description` - (Optional, ForceNew) Rule description.

Attributes Reference

* `id` - Security group rule ID.
* `tenant_id`, `project_id`, `revision_number`, `created_at`, `updated_at`, `tags`

