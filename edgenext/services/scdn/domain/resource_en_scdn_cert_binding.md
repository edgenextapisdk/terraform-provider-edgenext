Provides a resource to bind a certificate to an SCDN domain.

Example Usage

Bind certificate to domain

```hcl
resource "edgenext_scdn_cert_binding" "example" {
  domain_id = 12345
  ca_id     = 67890
}
```

Import

SCDN certificate bindings can be imported using the domain ID and certificate ID:

```shell
terraform import edgenext_scdn_cert_binding.example 12345-67890
```

