---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_log_download_template"
sidebar_current: "docs-edgenext-resource-scdn_log_download_template"
description: |-
  Provides a resource to create and manage SCDN log download templates.
---

# edgenext_scdn_log_download_template

Provides a resource to create and manage SCDN log download templates.

## Example Usage

### Create log download template

```hcl
resource "edgenext_scdn_log_download_template" "example" {
  template_name   = "my-template"
  group_name      = "my-group"
  group_id        = 1
  data_source     = "ng"
  download_fields = ["time", "domain", "url"]
}
```

### Create template with search terms

```hcl
resource "edgenext_scdn_log_download_template" "example" {
  template_name   = "my-template"
  group_name      = "my-group"
  group_id        = 1
  data_source     = "ng"
  download_fields = ["time", "domain", "url"]
  search_terms = [
    {
      key   = "domain"
      value = "example.com"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `data_source` - (Required, String) Data source: ng, cc, waf
* `download_fields` - (Required, List: [`String`]) Download fields
* `group_id` - (Required, Int) Group ID
* `group_name` - (Required, String) Group name
* `template_name` - (Required, String) Template name
* `domain_select_type` - (Optional, Int) Domain select type: 0-partial, 1-all, default: 0
* `search_terms` - (Optional, List) Search conditions
* `status` - (Optional, Int) Status: 1-enabled, 0-disabled, default: 1

The `search_terms` object supports the following:

* `key` - (Required, String) Search key
* `value` - (Required, String) Search value

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_at` - Creation timestamp
* `id` - The ID of the log download template
* `member_id` - The member ID
* `template_id` - The template ID
* `updated_at` - Last update timestamp


## Import

SCDN log download templates can be imported using the template ID:

```shell
terraform import edgenext_scdn_log_download_template.example 12345
```

