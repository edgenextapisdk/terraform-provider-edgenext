Provides a resource to switch the access mode of an SCDN domain.

Example Usage

Switch to CNAME access mode

```hcl
resource "edgenext_scdn_domain_access_mode" "example" {
  domain_id   = 12345
  access_mode = "cname"
}
```

Switch to NS access mode

```hcl
resource "edgenext_scdn_domain_access_mode" "example" {
  domain_id   = 12345
  access_mode = "ns"
}
```

Import

SCDN domain access mode can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain_access_mode.example 12345
```

