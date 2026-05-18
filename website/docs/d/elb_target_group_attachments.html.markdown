---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_target_group_attachments"
sidebar_current: "docs-edgenext-datasource-elb_target_group_attachments"
description: |-
  Use this data source to query EdgeNext ELB target group attachments.
---

# edgenext_elb_target_group_attachments

Use this data source to query EdgeNext ELB target group attachments.

## Example Usage

```hcl
data "edgenext_elb_target_group_attachments" "example" {
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

* `target_group_id` - (Required, String) Target group ID to list targets for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `targets` - Backend targets for the target group.
  * `address` - Target IP address.
  * `created_at` - Creation time as Unix timestamp (seconds).
  * `id` - Target ID.
  * `name` - Target name.
  * `operating_status` - Operating status.
  * `protocol_port` - Protocol port.
  * `provisioning_status` - Provisioning status.
  * `updated_at` - Last update time as Unix timestamp (seconds).
  * `weight` - Target weight.


