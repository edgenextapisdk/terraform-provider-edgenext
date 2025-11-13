Provides a resource to bind domains to an SCDN origin group.

Example Usage

Bind domains to origin group

```hcl
resource "edgenext_scdn_origin_group_domain_bind" "example" {
  origin_group_id = 12345
  domain_ids      = [67890, 11111]
}
```

Bind by domain names

```hcl
resource "edgenext_scdn_origin_group_domain_bind" "example" {
  origin_group_id = 12345
  domains         = ["example.com", "www.example.com"]
}
```

Bind by domain group IDs

```hcl
resource "edgenext_scdn_origin_group_domain_bind" "example" {
  origin_group_id  = 12345
  domain_group_ids = [1, 2]
}
```

Import

SCDN origin group domain bindings can be imported using the origin group ID and domain IDs:

```shell
terraform import edgenext_scdn_origin_group_domain_bind.example 12345-67890,11111
```

