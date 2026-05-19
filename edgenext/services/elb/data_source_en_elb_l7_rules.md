Use this data source to list EdgeNext ELB L7 rules for an L7 policy.

Example Usage

```hcl
data "edgenext_elb_l7_rules" "redirect" {
  policy_id                  = edgenext_elb_l7_policy.redirect.id
  limit                      = 100
  sort_key                   = "created_at"
  sort_dir                   = "desc"
  filter_type                = "PATH"
  filter_provisioning_status = "ACTIVE"
}

output "rule_ids" {
  value = [for r in data.edgenext_elb_l7_rules.redirect.rules : r.id]
}
```
