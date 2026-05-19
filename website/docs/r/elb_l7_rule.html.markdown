---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_l7_rule"
sidebar_current: "docs-edgenext-resource-elb_l7_rule"
description: |-
  Use this resource to manage a standalone L7 rule under an EdgeNext ELB L7 policy.
---

# edgenext_elb_l7_rule

Use this resource to manage a standalone L7 rule under an EdgeNext ELB L7 policy.

Updatable in place: `type`, `compare_type`, `value`. `l7policy_id` forces replacement. Read uses `l7rules/list` with an id filter. Delete waits until the rule is removed before create on replacement.

## Example Usage

```hcl
resource "edgenext_elb_l7_policy" "api_redirect" {
  listener_id  = edgenext_elb_listener.https.id
  name         = "api-redirect"
  action       = "REDIRECT_TO_URL"
  redirect_url = "https://api.example.com/"
}

resource "edgenext_elb_l7_rule" "api_prefix" {
  l7policy_id  = edgenext_elb_l7_policy.api_redirect.id
  type         = "PATH"
  compare_type = "STARTS_WITH"
  value        = "/api"
}
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


## Import

Import format is `l7policy_id/l7rule_id`. You may also use `listener_id/l7rule_id`; the provider resolves the policy on that listener.

```shell
terraform import edgenext_elb_l7_rule.api_prefix c6f3f52d-e331-487f-bb9f-c33c3f7f51ba/afea81e4-0a08-45e6-9d68-2ffefc10ba54
```

