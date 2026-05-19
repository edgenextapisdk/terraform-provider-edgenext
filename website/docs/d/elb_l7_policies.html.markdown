---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_l7_policies"
sidebar_current: "docs-edgenext-datasource-elb_l7_policies"
description: |-
  Use this data source to list EdgeNext ELB L7 policies on a listener.
---

# edgenext_elb_l7_policies

Use this data source to list EdgeNext ELB L7 policies on a listener.

## Example Usage

```hcl
data "edgenext_elb_l7_policies" "https" {
  listener_id = edgenext_elb_listener.https.id
  limit       = 100
  sort_key    = "position"
  sort_dir    = "asc"
}

output "l7_policy_ids" {
  value = [for p in data.edgenext_elb_l7_policies.https.l7_policies : p.id]
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required, String) Listener ID to list L7 policies for.
* `limit` - (Optional, Int) Maximum policies to return (0-1000).
* `sort_dir` - (Optional, String) Sort direction: asc, desc, or empty.
* `sort_key` - (Optional, String) Sort field key.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hasmore` - Whether more pages exist (from API hasmore).
* `l7_policies` - Policy list items returned by the API.
  * `action` - Policy action.
  * `created_at` - Creation time as Unix timestamp (seconds).
  * `description` - Policy description.
  * `id` - L7 policy ID.
  * `name` - Policy name.
  * `operating_status` - Operating status.
  * `position` - Policy position (priority).
  * `provisioning_status` - Provisioning status.
  * `rules_count` - Number of rules.
  * `updated_at` - Last update time as Unix timestamp (seconds).
* `total` - Total number of L7 policies matching the query (from the API when present).


