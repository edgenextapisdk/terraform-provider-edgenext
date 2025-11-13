---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_member_global_template"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_member_global_template"
description: |-
  Use this data source to query the member global SCDN security protection template.
---

# edgenext_scdn_security_protection_member_global_template

Use this data source to query the member global SCDN security protection template.

## Example Usage

### Query member global template

```hcl
data "edgenext_scdn_security_protection_member_global_template" "example" {
}

output "template" {
  value = data.edgenext_scdn_security_protection_member_global_template.example.template
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_member_global_template" "example" {
  result_output_file = "member_global_template.json"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_domain_count` - Bind domain count
* `template` - Global template information
  * `created_at` - Creation time
  * `id` - Template ID
  * `name` - Template name
  * `remark` - Template remark
  * `type` - Template type: global, only_domain, more_domain


