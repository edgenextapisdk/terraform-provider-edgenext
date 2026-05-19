Use this resource to attach one backend (target) to an EdgeNext ELB target group.

Updatable in place: `name`, `weight`, `protocol_port`. `target_group_id` and `address` force replacement. Delete polls until the target is gone so a subsequent create does not hit Octavia load balancer immutable errors.

Example Usage

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

Import

Import format is `target_group_id/target_id`.

```shell
terraform import edgenext_elb_target_group_attachment.backend b49ba19e-d737-4c6e-8067-62fe8f27448a/6292aaaa-388f-4a59-a379-83bd4f38246b
```
