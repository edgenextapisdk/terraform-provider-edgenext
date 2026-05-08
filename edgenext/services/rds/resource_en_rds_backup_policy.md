Use this resource to create and manage one EdgeNext RDS backup policy.

Example Usage

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

Import

Import format is `policy_id`.

```shell
terraform import edgenext_rds_backup_policy.example backup-policy-id-example
```
