Use this resource to create an EdgeNext ELB target group (pool) with a health monitor on a listener.

Backends are managed separately with `edgenext_elb_target_group_attachment`. Only `name` and `description` can be updated in place; `listener_id`, `protocol`, `lb_algorithm`, and `health_monitor` force replacement. Create and delete wait for Octavia provisioning and retry on transient immutable errors (see `timeouts`).

Example Usage

```hcl
resource "edgenext_elb_target_group" "app" {
  listener_id  = edgenext_elb_listener.https.id
  name         = "app"
  description  = "Application backends"
  protocol     = "HTTP"
  lb_algorithm = "ROUND_ROBIN"
  health_monitor {
    name           = "http-check"
    type           = "HTTP"
    max_retries    = 3
    delay          = 10
    timeout        = 5
    http_method    = "GET"
    url_path       = "/health"
    expected_codes = "200"
  }
}

resource "edgenext_elb_target_group_attachment" "app1" {
  target_group_id = edgenext_elb_target_group.app.id
  address         = "192.168.0.244"
  protocol_port   = 80
  name            = "app-1"
  weight          = 10
}
```
