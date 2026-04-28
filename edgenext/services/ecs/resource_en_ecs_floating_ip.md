Use this resource to create and manage ECS floating IPs.

Example Usage

```hcl
resource "edgenext_ecs_floating_ip" "example" {
  bandwidth = 10
}

data "edgenext_ecs_floating_ips" "example" {
  floating_ip_id = edgenext_ecs_floating_ip.example.id
  limit          = 1
}
```

Import

Import format is `floating_ip_id`.

```shell
terraform import edgenext_ecs_floating_ip.example c1eae862-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `bandwidth` - (Required) Floating IP bandwidth in Mbps.

Attributes Reference

* `id` - Floating IP ID.
* `ip_address` - Floating IP address.
