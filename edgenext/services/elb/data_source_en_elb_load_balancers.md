Use this data source to list EdgeNext ELB load balancers.

Example Usage

```hcl
data "edgenext_elb_load_balancers" "all" {
  name  = ""
  limit = 20
}

output "load_balancer_ids" {
  value = [for lb in data.edgenext_elb_load_balancers.all.balancers : lb.id]
}
```
