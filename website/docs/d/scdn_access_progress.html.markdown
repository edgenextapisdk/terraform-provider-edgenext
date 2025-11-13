---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_access_progress"
sidebar_current: "docs-edgenext-datasource-scdn_access_progress"
description: |-
  Use this data source to query available SCDN access progress status options.
---

# edgenext_scdn_access_progress

Use this data source to query available SCDN access progress status options.

## Example Usage

### Query access progress options

```hcl
data "edgenext_scdn_access_progress" "example" {
}

output "progress_options" {
  value = data.edgenext_scdn_access_progress.example.progress
}
```

### Query and save to file

```hcl
data "edgenext_scdn_access_progress" "example" {
  result_output_file = "access_progress.json"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `progress` - List of access progress status options
  * `key` - Progress key/identifier
  * `name` - Progress status name


