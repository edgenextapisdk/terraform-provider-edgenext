Use this data source to query instance associations for one EdgeNext RDS backup policy.

Example Usage

```hcl
data "edgenext_rds_backup_policy_associate_instances" "example" {
  policy_id = "backup-policy-d0ce79aeedfd76697daf6569"
}
```
