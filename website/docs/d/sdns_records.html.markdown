---
subcategory: "Security DNS (SDNS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_sdns_records"
sidebar_current: "docs-edgenext-datasource-sdns_records"
description: |-
  Use this data source to query a list of SDNS DNS records.
---

# edgenext_sdns_records

Use this data source to query a list of SDNS DNS records.

## Example Usage

### Query SDNS DNS records

```hcl
data "edgenext_sdns_records" "example" {
  domain_id = 12345
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int) Domain ID to list records for

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `records` - List of records in the domain


