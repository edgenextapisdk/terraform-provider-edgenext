Use this data source to list EdgeNext RDS instances.

Example Usage

```hcl
data "edgenext_rds_instances" "mysql" {
  page_num    = 1
  page_size   = 1000
  instance_id = ""
  name        = ""
}
```
