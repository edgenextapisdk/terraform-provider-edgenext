Use this resource to create and manage an EdgeNext ELB listener on a load balancer.

Updatable in place: `name`, `description`, `default_tls_certificate_ref`. Other arguments are create-time only and force replacement when changed. Delete waits until the listener is removed from the API (configure `timeouts.delete` if Octavia is slow).

Example Usage

```hcl
resource "edgenext_elb_certificate" "site" {
  name        = "www-example-com"
  certificate = file("${path.module}/fullchain.pem")
  private_key = file("${path.module}/privkey.pem")
}

resource "edgenext_elb_listener" "https" {
  loadbalancer_id             = var.loadbalancer_id
  name                        = "https"
  description                 = "Public HTTPS"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  connection_limit            = -1
  insert_headers = {
    "X-Forwarded-For"   = "true"
    "X-Forwarded-Port"  = "true"
    "X-Forwarded-Proto" = "true"
  }
  default_tls_certificate_ref = edgenext_elb_certificate.site.certificate_ref
}
```

Import

Import format is `listener_id`.

```shell
terraform import edgenext_elb_listener.https 0c03c22d-b9b5-49c8-a5e1-aedd9bc6e6b9
```
