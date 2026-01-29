# edgenext_scdn_domain_group

Manages an SCDN Domain Group.

## Example Usage

```hcl
resource "edgenext_scdn_domain_group" "example" {
  group_name = "my-domain-group"
  remark     = "Production domains"
  domains    = ["example.com", "www.example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required) The name of the domain group.
* `remark` - (Optional) Remark or description for the domain group.
* `domain_ids` - (Optional) List of domain IDs to bind to the group. Conflicts with `domains`.
* `domains` - (Optional) List of domain names to bind to the group. Conflicts with `domain_ids`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the domain group.
* `created_at` - The creation time of the domain group.
* `updated_at` - The last update time of the domain group.

## Import

SCDN Domain Groups can be imported using the group ID, e.g.

```
$ terraform import edgenext_scdn_domain_group.example 123
```
