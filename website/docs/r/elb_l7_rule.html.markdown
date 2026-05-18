---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_l7_rule"
sidebar_current: "docs-edgenext-resource-elb_l7_rule"
description: |-
  Manages an EdgeNext ELB l7 rule.
---

# edgenext_elb_l7_rule

Manages an EdgeNext ELB l7 rule.

## Example Usage

```hcl
# See examples/elb/main.tf
```

## Argument Reference

The following arguments are supported:

* `compare_type` - (Required, String) Compare type (for example EQUAL_TO, REGEX, STARTS_WITH, ENDS_WITH, CONTAINS).
* `l7policy_id` - (Required, String, ForceNew) L7 policy ID this rule belongs to.
* `type` - (Required, String) Rule type: HOSTNAME or PATH.
* `value` - (Required, String) Match value (1-255 characters).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `admin_state_up` - Administrative up/down state (computed).
* `created_at` - Creation time as Unix timestamp (seconds).
* `invert` - Whether the rule match is inverted (computed).
* `operating_status` - Operating status from the API.
* `project_id` - Project ID from the API.
* `provisioning_status` - Provisioning status from the API.
* `updated_at` - Last update time as Unix timestamp (seconds).


