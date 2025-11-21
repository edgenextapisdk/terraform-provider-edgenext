---
subcategory: "Security CDN (SCDN)"
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
  * `load_balance` - The load balancing method. Valid values: 0 (IP hash), 1 (Round robin), 2 (Cookie)
  * `origin_protocol` - The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS), 2 (Follow)
  * `origin_type` - The origin type. Valid values: 0 (IP), 1 (Domain)
  * `protocol` - The origin protocol. Valid values: 0 (HTTP), 1 (HTTPS)
  * `records` - The origin records
    * `host` - The origin host, specifies the Host header when accessing the origin
    * `port` - The port of the record
    * `priority` - The priority of the record
    * `value` - The value of the record (IP address or domain)
    * `view` - The view of the record. Valid values: primary (primary line), backup (backup line)
* `total` - The total number of origins


