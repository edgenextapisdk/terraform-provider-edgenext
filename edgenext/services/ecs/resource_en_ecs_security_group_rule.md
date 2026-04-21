Use this resource to manage a single ECS security group rule.

Example Usage

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

Import

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
