---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_templates"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_templates"
description: |-
  Use this data source to query a list of SCDN security protection templates.
---

# edgenext_scdn_security_protection_templates

Use this data source to query a list of SCDN security protection templates.

## Example Usage

### Query all templates

```hcl
data "edgenext_scdn_security_protection_templates" "example" {
  tpl_type  = "global"
  page      = 1
  page_size = 20
}

output "template_count" {
  value = data.edgenext_scdn_security_protection_templates.example.total
}

output "templates" {
  value = data.edgenext_scdn_security_protection_templates.example.list
}
```

### Query with filters

```hcl
data "edgenext_scdn_security_protection_templates" "example" {
  tpl_type   = "only_domain"
  search_key = "my-template"
  page       = 1
  page_size  = 20
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_templates" "example" {
  tpl_type           = "global"
  result_output_file = "templates.json"
}
```

## Argument Reference

The following arguments are supported:

* `tpl_type` - (Required, String) Template type: global, only_domain, more_domain
* `page_size` - (Optional, Int) Page size
* `page` - (Optional, Int) Page number
* `result_output_file` - (Optional, String) Used to save results to a file
* `search_key` - (Optional, String) Search keyword
* `search_type` - (Optional, String) Search type

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `templates` - Template list
  * `created_at` - Creation time
  * `id` - Template ID
  * `name` - Template name
  * `remark` - Template remark
  * `type` - Template type
* `total` - Total number of templates


