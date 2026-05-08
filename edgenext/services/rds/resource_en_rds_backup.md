Use this resource to create and manage one manual EdgeNext RDS backup.

Example Usage

```hcl
resource "edgenext_rds_backup" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  name        = "daily-backup"
}
```

Import

Import format is `backup_id`.

```shell
terraform import edgenext_rds_backup.example backup-id-example
```
