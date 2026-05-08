---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_backups"
sidebar_current: "docs-edgenext-datasource-rds_backups"
description: |-
  Use this data source to query EdgeNext RDS backups.
---

# edgenext_rds_backups

Use this data source to query EdgeNext RDS backups.

## Example Usage

```hcl
data "edgenext_rds_backups" "example" {
  page_num  = 1
  page_size = 1000
  backup_id = ""
  name      = ""
}
```

## Argument Reference

The following arguments are supported:

* `backup_id` - (Optional, String) Filter by backup ID. Maps to the API field id. Omit or leave empty to list without this filter.
* `name` - (Optional, String) Filter by backup name. Omit or leave empty to match the API empty-name filter.
* `page_num` - (Optional, Int) Page number for the list request. The API request body uses the field name page_number.
* `page_size` - (Optional, Int) Page size for the list request.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `backups` - Backups returned by the API.
  * `created` - Creation time (RFC3339).
  * `datastore` - Datastore engine and version.
    * `type` - Engine type.
    * `version` - Engine version.
  * `id` - Backup ID.
  * `instance_id` - RDS instance ID.
  * `instance_name` - RDS instance name.
  * `name` - Backup name.
  * `region` - Region.
  * `size` - Backup size in GB (as returned by the API).
  * `status` - Backup status.
  * `type` - Backup type (for example full).
  * `updated` - Last update time (RFC3339).
* `total` - Total number of backups matching the query.


