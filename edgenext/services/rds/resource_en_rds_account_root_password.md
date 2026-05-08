Use this resource to enable root access and set (or rotate) root password on an EdgeNext RDS instance.

Example Usage

```hcl
resource "edgenext_rds_account_root_password" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  password    = "Root@123456"
}
```
