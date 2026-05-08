---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_backup_policies"
sidebar_current: "docs-edgenext-datasource-rds_backup_policies"
description: |-
  Use this data source to query EdgeNext RDS backup policies.
---

# edgenext_rds_backup_policies

Use this data source to query EdgeNext RDS backup policies.

## Example Usage

```hcl
data "edgenext_rds_backup_policies" "example" {
  page_num  = 1
  page_size = 10
  policy_id = ""
  name      = ""
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Filter by policy name. Use empty string to omit the filter.
* `page_num` - (Optional, Int) Page number for the list request. The API request body uses the field name page_number.
* `page_size` - (Optional, Int) Page size for the list request.
* `policy_id` - (Optional, String) Filter by policy ID. Maps to the API field id. Use empty string to omit the filter.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policies` - Backup policies returned by the API.
  * `created` - Creation time.
  * `cycle_type` - Backup cycle type (for example weekly).
  * `hours` - Hour numbers in the schedule.
  * `id` - Policy ID.
  * `name` - Policy name.
  * `next_trigger_time_utc` - Next trigger time in UTC when returned by the API.
  * `region_id` - Region ID.
  * `retention_type` - Retention type (for example days or quantity).
  * `retention_value` - Retention value from the API.
  * `schedule_description` - Human-readable schedule description.
  * `schedule_enabled` - Whether the schedule is enabled.
  * `schedule_times` - Scheduled time strings (for example HH:MM:SS).
  * `status` - Policy status.
  * `stop_type` - Stop condition type.
  * `stop_value` - Stop condition value (format depends on stop_type).
  * `weekdays` - Weekday numbers in the schedule.
* `total` - Total number of policies matching the query.


