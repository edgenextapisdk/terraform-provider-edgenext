Use this resource to create and manage one database in an EdgeNext RDS instance.

Example Usage

```hcl
resource "edgenext_rds_database" "example" {
  instance_id   = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  name          = "app_db"
  character_set = "utf8"
  collate       = "utf8_general_ci"
}
```

Import

Import format is `instance_id/database_name`.

```shell
terraform import edgenext_rds_database.example b4a406bc-2859-430d-8f95-9bc1e054d347/app_db
```
