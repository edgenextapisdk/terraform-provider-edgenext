---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origins"
sidebar_current: "docs-edgenext-datasource-scdn_origins"
description: |-
  Use this data source to query a list of SCDN origin servers for a domain.
---

# edgenext_scdn_origins

Use this data source to query a list of SCDN origin servers for a domain.

## Example Usage

### Query all origins for a domain

```hcl
data "edgenext_scdn_origins" "example" {
  domain_id = 12345
}

output "origin_count" {
  value = data.edgenext_scdn_origins.example.total
}

output "origin_details" {
  value = data.edgenext_scdn_origins.example.origins
}
```

### Query origins and save to file

```hcl
data "edgenext_scdn_origins" "example" {
  domain_id          = 12345
  result_output_file = "origins.json"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int) The ID of the domain to query origins for
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origins` - The list of origins
  * `domain_id` - The ID of the domain
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
* `total` - The total number of origins


