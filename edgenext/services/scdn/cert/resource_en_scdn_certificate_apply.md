Provides a resource to apply for SCDN certificates for domains.

Example Usage

Apply certificate for single domain

```hcl
resource "edgenext_scdn_certificate_apply" "example" {
  domain = ["example.com"]
}
```

Apply certificate for multiple domains

```hcl
resource "edgenext_scdn_certificate_apply" "example" {
  domain = ["example.com", "www.example.com", "api.example.com"]
}
```

Import

SCDN certificate applications can be imported using the certificate application ID:

```shell
terraform import edgenext_scdn_certificate_apply.example 12345
```

