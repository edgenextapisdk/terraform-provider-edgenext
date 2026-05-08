---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_backup_policy_associate_instances"
sidebar_current: "docs-edgenext-datasource-rds_backup_policy_associate_instances"
description: |-
  Use this data source to query instance associations for one EdgeNext RDS backup policy.
---

# edgenext_rds_backup_policy_associate_instances

Use this data source to query instance associations for one EdgeNext RDS backup policy.

## Example Usage

```hcl
data "edgenext_rds_backup_policy_associate_instances" "example" {
  policy_id = "backup-policy-d0ce79aeedfd76697daf6569"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String) Backup policy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - Associated instances returned by the API.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `instance_type` - Instance type.
* `total` - Total number of associated sources returned by the list API.


