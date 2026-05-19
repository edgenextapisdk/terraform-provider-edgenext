Use this data source to list backends (targets) on one EdgeNext ELB target group.

Example Usage

```hcl
data "edgenext_elb_target_group_attachments" "app" {
  target_group_id = edgenext_elb_target_group.app.id
}

output "backend_addresses" {
  value = [for t in data.edgenext_elb_target_group_attachments.app.targets : t.address]
}
```
