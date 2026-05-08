---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_databases"
sidebar_current: "docs-edgenext-datasource-rds_databases"
description: |-
  Use this data source to list databases in one EdgeNext RDS instance.
---

# edgenext_rds_databases

Use this data source to list databases in one EdgeNext RDS instance.

## Example Usage

```hcl
data "edgenext_rds_databases" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) RDS instance ID to list databases for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `databases` - Databases returned by the API for the instance.
  * `name` - Database name.


