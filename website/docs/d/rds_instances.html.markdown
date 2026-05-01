---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_instances"
sidebar_current: "docs-edgenext-datasource-rds_instances"
description: |-
  Use this data source to list EdgeNext RDS instances.
---

# edgenext_rds_instances

Use this data source to list EdgeNext RDS instances.

## Example Usage

```hcl
data "edgenext_rds_instances" "mysql" {
  page_num       = 1
  page_size      = 1000
  datastore_type = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `datastore_type` - (Optional, String) Filter by datastore engine type (for example mysql). Omit when not filtering.
* `page_num` - (Optional, Int) Page number for the list request. The API request body uses the field name page_number.
* `page_size` - (Optional, Int) Page size for the list request.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - List of RDS instances returned by the API.
  * `access` - Network access settings.
    * `is_public` - Whether the instance is exposed on a public network.
  * `created` - Creation time (RFC3339).
  * `datastore` - Datastore engine and version.
    * `type` - Engine type.
    * `version_number` - Engine version number.
    * `version` - Engine version label.
  * `flavor` - Flavor specification.
    * `id` - Flavor ID.
    * `name` - Flavor name.
  * `hostname` - The instance hostname.
  * `id` - The instance ID.
  * `ip` - IP addresses attached to the instance.
    * `address` - IP address.
    * `network` - Network or subnet ID.
    * `type` - Address type (for example private).
    * `vpc_name` - VPC name.
  * `name` - The instance name.
  * `operating_status` - Operating status of the instance.
  * `region` - Region where the instance runs.
  * `root_enabled` - Whether root login is enabled.
  * `status` - The instance status.
  * `updated` - Last update time (RFC3339).
  * `volume` - Primary storage volume.
    * `size` - Volume size in GB.
* `total` - Total number of instances matching the query.


