Use this resource to create and manage one database user in an EdgeNext RDS instance.

Example Usage

```hcl
resource "edgenext_rds_account" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  user_name   = "app_user"
  host        = "%"
  password    = "Strong@1234"
}
```

Import

Import format is `instance_id/user_name/host`.

```shell
terraform import edgenext_rds_account.example b4a406bc-2859-430d-8f95-9bc1e054d347/app_user/%
```
