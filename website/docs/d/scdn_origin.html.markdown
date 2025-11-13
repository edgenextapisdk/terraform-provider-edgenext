---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin"
sidebar_current: "docs-edgenext-datasource-scdn_origin"
description: |-
  Use this data source to query details of a specific SCDN origin server.
---

# edgenext_scdn_origin

Use this data source to query details of a specific SCDN origin server.

## Example Usage

### Query origin by ID

```hcl
data "edgenext_scdn_origin" "example" {
  origin_id = 12345
  domain_id = 67890
}

output "origin_protocol" {
  value = data.edgenext_scdn_origin.example.protocol
}

output "origin_records" {
  value = data.edgenext_scdn_origin.example.records
}
```

### Query origin and save to file

```hcl
data "edgenext_scdn_origin" "example" {
  origin_id          = 12345
  domain_id          = 67890
  result_output_file = "origin.json"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int) The ID of the domain that owns the origin
* `origin_id` - (Required, Int) The ID of the origin to query
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the origin
* `listen_port` - The listening port of the origin server
* `load_balance` - The load balancing method
* `origin_protocol` - The origin protocol
* `origin_type` - The origin type
* `protocol` - The origin protocol
* `records` - The origin records
  * `port` - The port of the record
  * `priority` - The priority of the record
  * `value` - The value of the record
  * `view` - The view of the record


