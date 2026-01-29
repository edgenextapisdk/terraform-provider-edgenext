---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domain_group_domains"
sidebar_current: "docs-edgenext-datasource-scdn_domain_group_domains"
description: |-
  # edgenext_scdn_domain_group_domains
---

# edgenext_scdn_domain_group_domains

# edgenext_scdn_domain_group_domains

Query domains in an SCDN Domain Group.

## Example Usage

### Query domains in a domain group

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

* `group_id` - (Required, Int) Group ID to list domains for
* `domain` - (Optional, String) Filter by domain
* `page` - (Optional, Int) Page number
* `per_page` - (Optional, Int) Items per page

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of domains
  * `domain_id` - Domain ID
  * `domain` - Domain name
* `ports` - Common ports
* `total` - Total count


