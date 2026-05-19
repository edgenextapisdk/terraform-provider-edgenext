---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_l7_policy"
sidebar_current: "docs-edgenext-resource-elb_l7_policy"
description: |-
  Use this resource to manage an EdgeNext ELB L7 policy on a listener (for example redirect to a URL or another target group).
---

# edgenext_elb_l7_policy

Use this resource to manage an EdgeNext ELB L7 policy on a listener (for example redirect to a URL or another target group).

Match rules are **not** nested in this resource; create them with `edgenext_elb_l7_rule`. `position` is computed from the API. Actions use `REDIRECT_TO_TARGET_GROUP` in Terraform (mapped to the Octavia API as needed).

## Example Usage

```hcl
resource "edgenext_elb_l7_policy" "redirect" {
  listener_id              = edgenext_elb_listener.https.id
  name                     = "redirect-legacy"
  description              = "Send legacy path to another pool"
  action                   = "REDIRECT_TO_TARGET_GROUP"
  redirect_target_group_id = edgenext_elb_target_group.legacy.id
}

resource "edgenext_elb_l7_rule" "legacy_path" {
  l7policy_id  = edgenext_elb_l7_policy.redirect.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/old"
}
```

## Argument Reference

The following arguments are supported:

* `action` - (Required, String) Policy action: REDIRECT_TO_URL or REDIRECT_TO_TARGET_GROUP.
* `listener_id` - (Required, String, ForceNew) Listener ID this L7 policy is attached to.
* `name` - (Required, String) Policy name.
* `description` - (Optional, String) Policy description.
* `redirect_target_group_id` - (Optional, String) Target group ID when action is REDIRECT_TO_TARGET_GROUP.
* `redirect_url` - (Optional, String) Redirect URL when action is REDIRECT_TO_URL.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `admin_state_up` - Administrative up/down state (computed).
* `created_at` - Creation time as Unix timestamp (seconds).
* `operating_status` - Operating status.
* `position` - Evaluation priority among L7 policies on the listener (computed).
* `project_id` - Project ID from the API.
* `provisioning_status` - Provisioning status.
* `redirect_target_group_name` - Name of the redirect target group when present in read results.
* `updated_at` - Last update time as Unix timestamp (seconds).


## Import

Import format is `l7policy_id`.

```shell
terraform import edgenext_elb_l7_policy.redirect 6232fc9e-a76e-4e53-bb17-d641ac21e91a
```

