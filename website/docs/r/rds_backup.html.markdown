---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_backup"
sidebar_current: "docs-edgenext-resource-rds_backup"
description: |-
  Use this resource to create and manage one manual EdgeNext RDS backup.
---

# edgenext_rds_backup

Use this resource to create and manage one manual EdgeNext RDS backup.

## Example Usage

```hcl
resource "edgenext_rds_backup" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  name        = "daily-backup"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) RDS instance ID to back up.
* `name` - (Required, String, ForceNew) Backup name passed to the create API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created` - Creation time (RFC3339).
* `datastore` - Datastore engine and version.
  * `type` - Engine type.
  * `version` - Engine version.
* `instance_name` - RDS instance name from the get API.
* `region` - Region from the get API.
* `size` - Backup size in GB from the get API.
* `status` - Backup status from the API.
* `type` - Backup type (for example full).
* `updated` - Last update time (RFC3339).


## Import

Import format is `backup_id`.

```shell
terraform import edgenext_rds_backup.example backup-id-example
```

