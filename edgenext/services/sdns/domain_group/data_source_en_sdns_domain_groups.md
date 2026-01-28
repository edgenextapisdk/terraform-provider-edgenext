Use this data source to query a list of SDNS domain groups.

Example Usage

Query SDNS domain groups

```hcl
data "edgenext_sdns_domain_groups" "example" {
  group_name = "my-domain-group"
}
```

Attributes Reference

The following attributes are exported:

* `groups` - List of matched groups
  * `id` - The ID of the domain group
  * `group_name` - The name of the group
  * `remark` - Remark for the group
