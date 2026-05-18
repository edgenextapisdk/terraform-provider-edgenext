---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_l7_policy"
sidebar_current: "docs-edgenext-resource-elb_l7_policy"
description: |-
  Manages an EdgeNext ELB l7 policy.
---

# edgenext_elb_l7_policy

Manages an EdgeNext ELB l7 policy.

## Example Usage

```hcl
# See examples/elb/main.tf
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


