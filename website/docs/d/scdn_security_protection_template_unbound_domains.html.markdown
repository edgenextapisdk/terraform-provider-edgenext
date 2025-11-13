---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_template_unbound_domains"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_template_unbound_domains"
description: |-
  Use this data source to query domains not bound to any SCDN security protection template.
---

# edgenext_scdn_security_protection_template_unbound_domains

Use this data source to query domains not bound to any SCDN security protection template.

## Example Usage

### Query unbound domains

```hcl
data "edgenext_scdn_security_protection_template_unbound_domains" "example" {
  page      = 1
  page_size = 20
}

output "unbound_domain_count" {
  value = data.edgenext_scdn_security_protection_template_unbound_domains.example.total
}

output "unbound_domains" {
  value = data.edgenext_scdn_security_protection_template_unbound_domains.example.list
}
```

### Query with domain filter

```hcl
data "edgenext_scdn_security_protection_template_unbound_domains" "example" {
  domain = "example.com"
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_template_unbound_domains" "example" {
  result_output_file = "unbound_domains.json"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional, String) Domain filter
* `member_id` - (Optional, Int) Member ID
* `page_size` - (Optional, Int) Page size
* `page` - (Optional, Int) Page number
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domains` - Unbound domain list
  * `created_at` - Creation time
  * `domain` - Domain name
  * `id` - Domain ID
  * `remark` - Remark
  * `type` - Domain type
* `total` - Total number of unbound domains


