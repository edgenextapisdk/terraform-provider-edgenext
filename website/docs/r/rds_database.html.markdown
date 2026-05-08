---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_database"
sidebar_current: "docs-edgenext-resource-rds_database"
description: |-
  Use this resource to create and manage one database in an EdgeNext RDS instance.
---

# edgenext_rds_database

Use this resource to create and manage one database in an EdgeNext RDS instance.

## Example Usage

```hcl
resource "edgenext_rds_database" "example" {
  instance_id   = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  name          = "app_db"
  character_set = "utf8"
  collate       = "utf8_general_ci"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) RDS instance ID.
* `name` - (Required, String, ForceNew) Database name.
* `character_set` - (Optional, String, ForceNew) MySQL character set for the new database (create only). Examples: utf8, utf8mb4, latin1. Default: utf8. Matching is case-insensitive; the create request sends the canonical lowercase name.
* `collate` - (Optional, String, ForceNew) MySQL collation for the selected character_set (create only). Examples: utf8_general_ci (for utf8), utf8mb4_general_ci (for utf8mb4). Default: utf8_general_ci. Matching is case-insensitive.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `instance_id/database_name`.

```shell
terraform import edgenext_rds_database.example b4a406bc-2859-430d-8f95-9bc1e054d347/app_db
```

