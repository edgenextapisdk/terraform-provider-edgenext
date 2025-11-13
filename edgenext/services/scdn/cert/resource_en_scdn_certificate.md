Provides a resource to create and manage SCDN certificates.

Example Usage

Create certificate with certificate and key

```hcl
resource "edgenext_scdn_certificate" "example" {
  ca_name = "my-certificate"
  ca_cert = file("certificate.pem")
  ca_key  = file("private_key.pem")
}
```

Update existing certificate

```hcl
resource "edgenext_scdn_certificate" "example" {
  certificate_id = "12345"
  ca_name       = "my-certificate"
  ca_cert        = file("new_certificate.pem")
  ca_key         = file("new_private_key.pem")
}
```

Import

SCDN certificates can be imported using the certificate ID:

```shell
terraform import edgenext_scdn_certificate.example 12345
```

