---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_rule_template"
sidebar_current: "docs-edgenext-datasource-scdn_rule_template"
description: |-
  Use this data source to query details of a specific SCDN rule template.
---

# edgenext_scdn_rule_template

Use this data source to query details of a specific SCDN rule template.

## Example Usage

### Query rule template by ID

```hcl
data "edgenext_scdn_rule_template" "example" {
  id       = "12345"
  app_type = "network_speed"
}

output "template_name" {
  value = data.edgenext_scdn_rule_template.example.name
}

output "bind_domains" {
  value = data.edgenext_scdn_rule_template.example.bind_domains
}
```

### Query and save to file

```hcl
data "edgenext_scdn_rule_template" "example" {
  id                 = "12345"
  app_type           = "network_speed"
  result_output_file = "template.json"
}
```

## Argument Reference

The following arguments are supported:

* `app_type` - (Required, String) The application type (e.g., 'network_speed')
* `id` - (Required, String) The rule template ID
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_domains` - List of domains bound to this template
  * `created_at` - Domain binding timestamp
  * `domain_id` - Bound domain ID
  * `domain` - Bound domain name
* `created_at` - The template creation timestamp
* `description` - The rule template description
* `name` - The rule template name


