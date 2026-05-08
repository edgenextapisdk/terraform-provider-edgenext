---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_backup_policy"
sidebar_current: "docs-edgenext-resource-rds_backup_policy"
description: |-
  Use this resource to create and manage one EdgeNext RDS backup policy.
---

# edgenext_rds_backup_policy

Use this resource to create and manage one EdgeNext RDS backup policy.

## Example Usage

```hcl
resource "edgenext_rds_backup_policy" "example" {
  name             = "example-policy"
  cycle_type       = "weekly"
  weekdays         = [1, 7]
  schedule_times   = ["00:00:00", "01:00:00"]
  retention_type   = "days"
  retention_value  = 30
  stop_type        = "end_time"
  stop_value       = "2026-05-31 14:32:20"
  schedule_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `cycle_type` - (Required, String) Backup cycle type (for example weekly).
* `name` - (Required, String) Backup policy name.
* `retention_type` - (Required, String) Retention type. Allowed values: permanent, days, quantity.
* `retention_value` - (Required, Int) Retention value associated with retention_type.
* `schedule_times` - (Required, Set: [`String`]) Schedule time strings (for example 00:00:00).
* `stop_type` - (Required, String) Stop condition type. Allowed values: end_time, never, count.
* `weekdays` - (Required, Set: [`Int`]) Weekday numbers for scheduling (for example 1..7).
* `schedule_enabled` - (Optional, Bool) Whether to enable the schedule.
* `stop_value` - (Optional, String) Stop condition value, for example a datetime when stop_type is end_time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created` - Creation time.
* `hours` - Hour numbers in the schedule.
* `next_trigger_time_utc` - Next trigger time in UTC when returned by the API.
* `region_id` - Region ID.
* `schedule_description` - Human-readable schedule description.
* `status` - Policy status.


## Import

Import format is `policy_id`.

```shell
terraform import edgenext_rds_backup_policy.example backup-policy-id-example
```

