Use this resource to manage a single ECS security group rule.

Example Usage

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

Import

Import format is `security_group_id/rule_id`.

```shell
terraform import edgenext_ecs_security_group_rule.example 12f8f386-xxxx-xxxx-xxxx-xxxxxxxxxxxx/df58bf0a-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```
