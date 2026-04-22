Use this resource to create and manage ECS floating IPs.

Example Usage

```hcl
resource "edgenext_ecs_floating_ip" "example" {
  region    = "tokyo-a"
  bandwidth = 10
}
```

Import

Import format is `region/floating_ip_id`.

```shell
terraform import edgenext_ecs_floating_ip.example tokyo-a/c1eae862-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Argument Reference

* `region` - (Required) Region.
* `bandwidth` - (Required) Floating IP bandwidth in Mbps.

Attributes Reference

* `id` - Floating IP ID.
* `ip_address` - Floating IP address.
