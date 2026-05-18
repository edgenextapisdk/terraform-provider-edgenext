---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_target_group_attachment"
sidebar_current: "docs-edgenext-resource-elb_target_group_attachment"
description: |-
  Manages an EdgeNext ELB target group attachment.
---

# edgenext_elb_target_group_attachment

Manages an EdgeNext ELB target group attachment.

## Example Usage

```hcl
# See examples/elb/main.tf
```

## Argument Reference

The following arguments are supported:

* `address` - (Required, String, ForceNew) IP address of the target (backend).
* `protocol_port` - (Required, Int) Backend protocol port (updatable).
* `target_group_id` - (Required, String, ForceNew) Target group this target belongs to (same id as edgenext_elb_target_group).
* `name` - (Optional, String) Display name for the target (optional on create; updatable).
* `weight` - (Optional, Int) Target weight.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Creation time as Unix timestamp (seconds).
* `operating_status` - Operating status from the API.
* `provisioning_status` - Provisioning status from the API.
* `updated_at` - Last update time as Unix timestamp (seconds).


