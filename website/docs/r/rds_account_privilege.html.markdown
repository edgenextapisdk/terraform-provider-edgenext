---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_account_privilege"
sidebar_current: "docs-edgenext-resource-rds_account_privilege"
description: |-
  Use this resource to manage database grants for one user in an EdgeNext RDS instance.
---

# edgenext_rds_account_privilege

Use this resource to manage database grants for one user in an EdgeNext RDS instance.

## Example Usage

```hcl
resource "edgenext_rds_account_privilege" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  user_name   = "app_user"
  host        = "%"
  databases   = ["app_db"]
}
```

## Argument Reference

The following arguments are supported:

* `databases` - (Required, Set: [`String`]) Database names granted to the user.
* `host` - (Required, String, ForceNew) Client host for the database user (same as edgenext_rds_account.host). Together with instance_id and user_name it identifies the user. Use % for any host.
* `instance_id` - (Required, String, ForceNew) RDS instance ID.
* `user_name` - (Required, String, ForceNew) Database user name to manage privileges for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `instance_id/user_name/host`.

```shell
terraform import edgenext_rds_account_privilege.example b4a406bc-2859-430d-8f95-9bc1e054d347/app_user/%
```

