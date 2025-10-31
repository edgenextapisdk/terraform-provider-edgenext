---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_rule_template_domains"
sidebar_current: "docs-edgenext-datasource-scdn_rule_template_domains"
description: |-
  Use this data source to query domains bound to a specific SCDN rule template.
---

# edgenext_scdn_rule_template_domains

Use this data source to query domains bound to a specific SCDN rule template.

## Example Usage

### Query template domains

```hcl
data "edgenext_scdn_rule_template_domains" "example" {
  id        = 12345
  app_type  = "network_speed"
  page      = 1
  page_size = 100
}

output "domain_count" {
  value = data.edgenext_scdn_rule_template_domains.example.total
}

output "domains" {
  value = data.edgenext_scdn_rule_template_domains.example.list
}
```

### Query with domain filter

```hcl
data "edgenext_scdn_rule_template_domains" "example" {
  id       = 12345
  app_type = "network_speed"
  domain   = "example.com"
}

output "filtered_domains" {
  value = data.edgenext_scdn_rule_template_domains.example.list
}
```

### Query and save to file

```hcl
data "edgenext_scdn_rule_template_domains" "example" {
  id                 = 12345
  app_type           = "network_speed"
  result_output_file = "template_domains.json"
}
```

## Argument Reference

The following arguments are supported:

* `app_type` - (Required, String) The application type (e.g., 'network_speed')
* `id` - (Required, Int) The rule template ID
* `domain` - (Optional, String) Filter by domain name
* `page_size` - (Optional, Int) Items per page
* `page` - (Optional, Int) Page number for pagination
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - List of domain information
  * `created_at` - Domain binding timestamp
  * `domain` - Domain name
  * `id` - Domain ID
* `total` - Total number of domains bound to template


