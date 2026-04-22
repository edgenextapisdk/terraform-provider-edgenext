---
subcategory: "Security DNS (SDNS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_sdns_domains"
sidebar_current: "docs-edgenext-datasource-sdns_domains"
description: |-
  Use this data source to query a list of SDNS domains.
---

# edgenext_sdns_domains

Use this data source to query a list of SDNS domains.

## Example Usage

### Query SDNS domains

```hcl
data "edgenext_sdns_domains" "example" {
  domain = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional, String) Filter by domain name

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domains` - List of matched domains


