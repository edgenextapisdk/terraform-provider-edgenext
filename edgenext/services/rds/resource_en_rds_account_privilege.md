Use this resource to manage database grants for one user in an EdgeNext RDS instance.

Example Usage

```hcl
resource "edgenext_rds_account_privilege" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  user_name   = "app_user"
  host        = "%"
  databases   = ["app_db"]
}
```

Import

Import format is `instance_id/user_name/host`.

```shell
terraform import edgenext_rds_account_privilege.example b4a406bc-2859-430d-8f95-9bc1e054d347/app_user/%
```
