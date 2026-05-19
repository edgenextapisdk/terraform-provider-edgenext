---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_target_group_attachment"
sidebar_current: "docs-edgenext-resource-elb_target_group_attachment"
description: |-
  Use this resource to attach one backend (target) to an EdgeNext ELB target group.
---

# edgenext_elb_target_group_attachment

Use this resource to attach one backend (target) to an EdgeNext ELB target group.

Updatable in place: `name`, `weight`, `protocol_port`. `target_group_id` and `address` force replacement. Delete polls until the target is gone so a subsequent create does not hit Octavia load balancer immutable errors.

## Example Usage

```hcl
resource "edgenext_elb_target_group" "app" {
  listener_id  = edgenext_elb_listener.https.id
  name         = "app"
  protocol     = "HTTP"
  lb_algorithm = "ROUND_ROBIN"
  health_monitor {
    name           = "check"
    type           = "HTTP"
    max_retries    = 3
    delay          = 10
    timeout        = 5
    http_method    = "GET"
    url_path       = "/"
    expected_codes = "200"
  }
}

resource "edgenext_elb_target_group_attachment" "backend" {
  target_group_id = edgenext_elb_target_group.app.id
  address         = "192.168.0.244"
  protocol_port   = 8080
  name            = "app-server-1"
  weight          = 10
}
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


## Import

Import format is `target_group_id/target_id`.

```shell
terraform import edgenext_elb_target_group_attachment.backend b49ba19e-d737-4c6e-8067-62fe8f27448a/6292aaaa-388f-4a59-a379-83bd4f38246b
```

