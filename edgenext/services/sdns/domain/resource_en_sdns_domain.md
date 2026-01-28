Provides a resource to create and manage SDNS domains.

Example Usage

Create SDNS domain

```hcl
resource "edgenext_sdns_domain" "example" {
  domain = "example.com"
}
```

Import

SDNS domains can be imported using the domain ID:

```shell
terraform import edgenext_sdns_domain.example 12345
```
