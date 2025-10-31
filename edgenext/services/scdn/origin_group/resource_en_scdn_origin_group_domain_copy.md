Provides a resource to copy an SCDN origin group configuration to a domain.

Example Usage

Copy origin group to domain

```hcl
resource "edgenext_scdn_origin_group_domain_copy" "example" {
  origin_group_id = 12345
  domain_id      = 67890
}
```

Import

SCDN origin group domain copies can be imported using the origin group ID and domain ID:

```shell
terraform import edgenext_scdn_origin_group_domain_copy.example 12345-67890
```

