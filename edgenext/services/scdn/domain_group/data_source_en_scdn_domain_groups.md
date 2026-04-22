# edgenext_scdn_domain_groups

Query SCDN Domain Groups.

Example Usage

Query domain groups by name

```hcl
data "edgenext_scdn_domain_groups" "example" {
  group_name = "my-domain-group"
}

output "groups" {
  value = data.edgenext_scdn_domain_groups.example.list
}
```

Argument Reference

The following arguments are supported:

* `group_name` - (Optional) Filter by group name.
* `domain` - (Optional) Filter by domain.
* `page` - (Optional) Page number. Default is 1.
* `per_page` - (Optional) Items per page. Default is 10.

Attributes Reference

The following attributes are exported:

* `list` - List of domain groups. Each group has the following attributes:
  * `id` - The ID of the domain group.
  * `group_name` - The name of the domain group.
  * `remark` - The remark of the domain group.
  * `created_at` - The creation time.
  * `updated_at` - The last update time.
* `total` - Total count of domain groups.