---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_origin_group"
sidebar_current: "docs-edgenext-datasource-scdn_origin_group"
description: |-
  Use this data source to query details of a specific SCDN origin group.
---

# edgenext_scdn_origin_group

Use this data source to query details of a specific SCDN origin group.

## Example Usage

### Query origin group by ID

```hcl
data "edgenext_scdn_origin_group" "example" {
  origin_group_id = 12345
}

output "origin_group_name" {
  value = data.edgenext_scdn_origin_group.example.name
}

output "origins" {
  value = data.edgenext_scdn_origin_group.example.origins
}
```

### Query and save to file

```hcl
data "edgenext_scdn_origin_group" "example" {
  origin_group_id    = 12345
  result_output_file = "origin_group.json"
}
```

## Argument Reference

The following arguments are supported:

* `origin_group_id` - (Required, Int) Origin group ID
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Creation time
* `member_id` - Member ID
* `name` - Origin group name
* `origins` - Origin list
  * `id` - Origin ID
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
* `remark` - Remark
* `updated_at` - Update time
* `username` - Username


