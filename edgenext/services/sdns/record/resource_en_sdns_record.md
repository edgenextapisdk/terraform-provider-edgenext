Provides a resource to create and manage SDNS DNS records.

Example Usage

Create SDNS DNS record

```hcl
resource "edgenext_sdns_record" "example" {
  domain_id = 12345
  name      = "www"
  type      = "A"
  view      = "any"
  value     = "1.2.3.4"
  ttl       = 600
}
```

Import

SDNS DNS records can be imported using the record ID:

```shell
terraform import edgenext_sdns_record.example 67890
```
