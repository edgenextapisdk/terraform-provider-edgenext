# edgenext_scdn_domain_group_domains

Query domains in an SCDN Domain Group.

## Example Usage

```hcl
data "edgenext_scdn_domain_group_domains" "example" {
  group_id = 123
}

output "domains" {
  value = data.edgenext_scdn_domain_group_domains.example.list
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The ID of the domain group.
* `domain` - (Optional) Filter by domain name.
* `page` - (Optional) Page number. Default is 1.
* `per_page` - (Optional) Items per page. Default is 10.

## Attributes Reference

The following attributes are exported:

* `list` - List of domains. Each domain has the following attributes:
  * `domain_id` - The ID of the domain.
  * `domain` - The domain name.
* `total` - Total count of domains.
* `ports` - List of common ports.
