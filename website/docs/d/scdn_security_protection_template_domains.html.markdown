---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_template_domains"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_template_domains"
description: |-
  Use this data source to query domains bound to a specific SCDN security protection template.
---

# edgenext_scdn_security_protection_template_domains

Use this data source to query domains bound to a specific SCDN security protection template.

## Example Usage

### Query template domains

```hcl
data "edgenext_scdn_security_protection_template_domains" "example" {
  business_id = 12345
  page        = 1
  page_size   = 20
}

output "domain_count" {
  value = data.edgenext_scdn_security_protection_template_domains.example.total
}

output "domains" {
  value = data.edgenext_scdn_security_protection_template_domains.example.list
}
```

### Query with domain filter

```hcl
data "edgenext_scdn_security_protection_template_domains" "example" {
  business_id = 12345
  domain      = "example.com"
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_template_domains" "example" {
  business_id        = 12345
  result_output_file = "template_domains.json"
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int) Business ID (template ID)
* `domain` - (Optional, String) Domain filter
* `page_size` - (Optional, Int) Page size
* `page` - (Optional, Int) Page number
* `result_output_file` - (Optional, String) Used to save results to a file
* `tpl_type` - (Optional, String) Template type: global, only_domain, more_domain

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domains` - Domain list
  * `created_at` - Creation time
  * `domain` - Domain name
  * `id` - Domain ID
  * `remark` - Remark
  * `type` - Domain type
* `total` - Total number of domains


