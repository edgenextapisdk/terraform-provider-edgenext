---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_certificate_export"
sidebar_current: "docs-edgenext-datasource-scdn_certificate_export"
description: |-
  Use this data source to export SCDN certificates.
---

# edgenext_scdn_certificate_export

Use this data source to export SCDN certificates.

## Example Usage

### Export certificate

```hcl
data "edgenext_scdn_certificate_export" "example" {
  id = "12345"
}

output "export_url" {
  value = data.edgenext_scdn_certificate_export.example.exports[0].real_url
}
```

### Export multiple certificates

```hcl
data "edgenext_scdn_certificate_export" "example" {
  id = "12345,67890"
}

output "exports" {
  value = data.edgenext_scdn_certificate_export.example.exports
}
```

### Query and save to file

```hcl
data "edgenext_scdn_certificate_export" "example" {
  id                 = "12345"
  result_output_file = "certificate_export.json"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) The certificate ID (can be a single ID or comma-separated IDs)
* `product_flag` - (Optional, String) The product flag
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `exports` - The list of exported certificate data
  * `hash` - The export hash
  * `key` - The export key
  * `real_url` - The real URL for downloading the exported certificate


