---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_iota"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_iota"
description: |-
  Use this data source to query SCDN security protection IOTA (enum key-value pairs).
---

# edgenext_scdn_security_protection_iota

Use this data source to query SCDN security protection IOTA (enum key-value pairs).

## Example Usage

### Query IOTA

```hcl
data "edgenext_scdn_security_protection_iota" "example" {
}

output "iota" {
  value = data.edgenext_scdn_security_protection_iota.example.iota
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_iota" "example" {
  result_output_file = "iota.json"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `iota` - Enum key-value pairs


