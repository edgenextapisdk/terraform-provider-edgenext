Use this data source to query ECS floating IPs.

Example Usage

```hcl
data "edgenext_ecs_floating_ips" "example" {
  floating_ip_id = ""
  limit          = 10
}
```
