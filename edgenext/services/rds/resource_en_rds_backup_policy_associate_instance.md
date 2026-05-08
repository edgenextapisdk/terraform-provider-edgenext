Use this resource to manage one instance association for an EdgeNext RDS backup policy.

Example Usage

```hcl
resource "edgenext_rds_backup_policy_associate_instance" "example" {
  policy_id   = "backup-policy-d0ce79aeedfd76697daf6569"
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
}
```

Import

Import format is `policy_id/instance_id`.

```shell
terraform import edgenext_rds_backup_policy_associate_instance.example backup-policy-d0ce79aeedfd76697daf6569/b4a406bc-2859-430d-8f95-9bc1e054d347
```
