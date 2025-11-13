---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_rule_templates"
sidebar_current: "docs-edgenext-datasource-scdn_rule_templates"
description: |-
  Use this data source to query a list of SCDN rule templates with optional filters.
---

# edgenext_scdn_rule_templates

Use this data source to query a list of SCDN rule templates with optional filters.

## Example Usage

### Query all rule templates

```hcl
data "edgenext_scdn_rule_templates" "all" {
  page      = 1
  page_size = 100
}

output "template_count" {
  value = data.edgenext_scdn_rule_templates.all.total
}

output "templates" {
  value = data.edgenext_scdn_rule_templates.all.list
}
```

### Query templates with filters

```hcl
data "edgenext_scdn_rule_templates" "filtered" {
  page      = 1
  page_size = 100
  name      = "my-template"
  app_type  = "network_speed"
}

output "filtered_templates" {
  value = data.edgenext_scdn_rule_templates.filtered.list
}
```

### Query and save to file

```hcl
data "edgenext_scdn_rule_templates" "all" {
  page               = 1
  page_size          = 100
  result_output_file = "templates.json"
}
```

## Argument Reference

The following arguments are supported:

* `app_type` - (Optional, String) Filter by application type
* `domain` - (Optional, String) Filter by associated domain
* `name` - (Optional, String) Filter by rule template name
* `page_size` - (Optional, Int) Items per page, max: 1000, default: 1000
* `page` - (Optional, Int) Page number for pagination, default: 1
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of rule templates
  * `app_type` - Application type
  * `bind_domains` - List of domains bound to this template
    * `created_at` - Domain binding timestamp
    * `domain_id` - Bound domain ID
    * `domain` - Bound domain name
  * `created_at` - Template creation timestamp
  * `description` - Rule template description
  * `id` - Rule template ID
  * `name` - Rule template name
* `total` - Total number of rule templates


