Use this resource to upload and manage EdgeNext ELB TLS certificates (Barbican containers). There is no update API; changing certificate or private key material requires replacement.

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
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  default_tls_certificate_ref = edgenext_elb_certificate.site.certificate_ref
}
```

Import

Import format is the Barbican container `certificate_id`.

```shell
terraform import edgenext_elb_certificate.site a037e728-a23f-4937-ae7d-a70e1fedf828
```
