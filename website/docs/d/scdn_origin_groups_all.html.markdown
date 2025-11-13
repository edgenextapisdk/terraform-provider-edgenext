---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin_groups_all"
sidebar_current: "docs-edgenext-datasource-scdn_origin_groups_all"
description: |-
  Use this data source to query all SCDN origin groups by protection status.
---

# edgenext_scdn_origin_groups_all

Use this data source to query all SCDN origin groups by protection status.

## Example Usage

### Query all origin groups for SCDN nodes

```hcl
data "edgenext_scdn_origin_groups_all" "example" {
  protect_status = "scdn"
}

output "origin_group_count" {
  value = data.edgenext_scdn_origin_groups_all.example.total
}

output "origin_groups" {
  value = data.edgenext_scdn_origin_groups_all.example.origin_groups
}
```

### Query all origin groups for exclusive nodes

```hcl
data "edgenext_scdn_origin_groups_all" "example" {
  protect_status = "exclusive"
}
```

### Query and save to file

```hcl
data "edgenext_scdn_origin_groups_all" "example" {
  protect_status     = "scdn"
  result_output_file = "origin_groups_all.json"
}
```

## Argument Reference

The following arguments are supported:

* `protect_status` - (Required, String) Protection status: scdn-shared nodes, exclusive-dedicated nodes
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origin_groups` - Origin group list
  * `id` - Origin group ID
  * `name` - Origin group name
  * `origins` - Origin list
    * `load_balance` - Load balance strategy: 0-ip_hash, 1-round_robin, 2-cookie
    * `origin_protocol` - Origin protocol: 0-http, 1-https, 2-follow
    * `origin_type` - Origin type: 0-IP, 1-domain
    * `protocol_ports` - Protocol port mapping
      * `listen_ports` - Listen port list
      * `protocol` - Protocol: 0-http, 1-https
    * `records` - Origin record list
      * `host` - Origin Host
      * `port` - Origin port
      * `priority` - Weight
      * `value` - Origin address
      * `view` - Origin type: primary-backup, backup-backup
* `total` - Total number of origin groups


