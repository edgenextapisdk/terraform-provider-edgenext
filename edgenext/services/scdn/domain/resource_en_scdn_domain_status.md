Provides a resource to manage the status (enable/disable) of an SCDN domain.

Example Usage

Enable domain

```hcl
resource "edgenext_scdn_domain_status" "example" {
  domain_id = 12345
  enabled   = true
}
```

Disable domain

```hcl
resource "edgenext_scdn_domain_status" "example" {
  domain_id = 12345
  enabled   = false
}
```

Import

SCDN domain status can be imported using the domain ID:

```shell
terraform import edgenext_scdn_domain_status.example 12345
```

