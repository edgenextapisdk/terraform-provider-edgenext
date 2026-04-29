Use this data source to query ECS floating IPs.

Example Usage

```hcl
data "edgenext_ecs_floating_ips" "example" {
  floating_ip_id = edgenext_ecs_floating_ip.example.id
  limit          = 10
}

resource "edgenext_ecs_floating_ip" "example" {
  bandwidth = 10
}
```
