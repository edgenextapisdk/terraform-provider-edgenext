---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_template"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_template"
description: |-
  Use this data source to query details of a specific SCDN security protection template.
---

# edgenext_scdn_security_protection_template

Use this data source to query details of a specific SCDN security protection template.

## Example Usage

### Query security protection template

```hcl
data "edgenext_scdn_security_protection_template" "example" {
  business_id = 12345
}

output "template_name" {
  value = data.edgenext_scdn_security_protection_template.example.name
}

output "bind_domain_count" {
  value = data.edgenext_scdn_security_protection_template.example.bind_domain_count
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_template" "example" {
  business_id        = 12345
  result_output_file = "template.json"
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int) Business ID (template ID)
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_domain_count` - Bind domain count
* `created_at` - Creation time
* `name` - Template name
* `remark` - Template remark
* `type` - Template type: global, only_domain, more_domain


