Provides a resource to create and manage SDNS domain groups.

Example Usage

Create SDNS domain group

```hcl
resource "edgenext_sdns_domain_group" "example" {
  group_name = "my-domain-group"
  remark     = "This is a test group"
  domain_ids = ["123", "456"]
}
```

Import

SDNS domain groups can be imported using the group ID:

```shell
terraform import edgenext_sdns_domain_group.example 67890
```
