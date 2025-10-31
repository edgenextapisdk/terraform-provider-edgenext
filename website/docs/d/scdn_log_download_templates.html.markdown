---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_log_download_templates"
sidebar_current: "docs-edgenext-datasource-scdn_log_download_templates"
description: |-
  Use this data source to query SCDN log download templates.
---

# edgenext_scdn_log_download_templates

Use this data source to query SCDN log download templates.

## Example Usage

### Query all log download templates

```hcl
data "edgenext_scdn_log_download_templates" "all" {
  page     = 1
  per_page = 20
}

output "template_count" {
  value = data.edgenext_scdn_log_download_templates.all.total
}

output "templates" {
  value = data.edgenext_scdn_log_download_templates.all.templates
}
```

### Query with filters

```hcl
data "edgenext_scdn_log_download_templates" "filtered" {
  page     = 1
  per_page = 20
  search_terms = [
    {
      key   = "template_name"
      value = "my-template"
    }
  ]
}
```

### Query and save to file

```hcl
data "edgenext_scdn_log_download_templates" "all" {
  page               = 1
  per_page           = 20
  result_output_file = "log_download_templates.json"
}
```

## Argument Reference

The following arguments are supported:

* `data_source` - (Optional, String) Data source: ng, cc, waf
* `group_id` - (Optional, Int) Group ID
* `page` - (Optional, Int) Page number
* `per_page` - (Optional, Int) Items per page
* `result_output_file` - (Optional, String) Used to save results to a file
* `status` - (Optional, Int) Status: 1-enabled, 0-disabled
* `template_name` - (Optional, String) Template name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `templates` - Template list
  * `created_at` - Creation timestamp
  * `data_source` - Data source
  * `download_fields` - Download fields
  * `group_id` - Group ID
  * `group_name` - Group name
  * `member_id` - Member ID
  * `search_terms` - Search conditions
    * `key` - Search key
    * `value` - Search value
  * `status` - Status
  * `template_id` - Template ID
  * `template_name` - Template name
  * `updated_at` - Last update timestamp
* `total` - Total number of templates


