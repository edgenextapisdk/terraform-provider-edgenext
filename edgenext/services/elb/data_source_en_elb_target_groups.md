Use this data source to list EdgeNext ELB target groups for a load balancer, including health monitor details when available.

Example Usage

```hcl
data "edgenext_elb_target_groups" "lb" {
  loadbalancer_id = var.loadbalancer_id
}

output "target_group_ids" {
  value = [for tg in data.edgenext_elb_target_groups.lb.target_groups : tg.id]
}
```
