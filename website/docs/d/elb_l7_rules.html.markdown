---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_l7_rules"
sidebar_current: "docs-edgenext-datasource-elb_l7_rules"
description: |-
  Use this data source to list EdgeNext ELB L7 rules for an L7 policy.
---

# edgenext_elb_l7_rules

Use this data source to list EdgeNext ELB L7 rules for an L7 policy.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String) L7 policy ID to list rules for.
* `filter_id` - (Optional, String) Optional filter: rule id.
* `filter_provisioning_status` - (Optional, String) Optional filter: provisioning_status.
* `filter_type` - (Optional, String) Optional filter: HOSTNAME or PATH.
* `limit` - (Optional, Int) Maximum rules to return (0-1000).
* `sort_dir` - (Optional, String) Sort direction: asc, desc, or empty.
* `sort_key` - (Optional, String) Sort field key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hasmore` - Whether more pages exist.
* `is_complete_task` - Whether async task is complete (from API is_complete_task).
* `rules` - Rules returned by the API.
  * `admin_state_up` - Administrative state.
  * `compare_type` - Compare type.
  * `created_at` - Creation time as Unix timestamp (seconds).
  * `id` - Rule ID.
  * `invert` - Invert match.
  * `operating_status` - Operating status.
  * `project_id` - Project ID.
  * `provisioning_status` - Provisioning status.
  * `type` - Rule type.
  * `updated_at` - Last update time as Unix timestamp (seconds).
  * `value` - Match value.
* `total` - Total number of L7 rules matching the query (from the API when present).


