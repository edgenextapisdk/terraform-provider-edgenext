---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_account"
sidebar_current: "docs-edgenext-resource-rds_account"
description: |-
  Use this resource to create and manage one database user in an EdgeNext RDS instance.
---

# edgenext_rds_account

Use this resource to create and manage one database user in an EdgeNext RDS instance.

## Example Usage

```hcl
resource "edgenext_rds_account" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  user_name   = "app_user"
  host        = "%"
  password    = "Strong@1234"
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) Client host for the database user: use % for any host, or a literal IPv4 address. Sent on update and delete.
* `instance_id` - (Required, String, ForceNew) RDS instance ID.
* `user_name` - (Required, String, ForceNew) User name.
* `password` - (Optional, String) User password. Required on create; omit after import unless rotating. Please enter 8-20 characters, must include all four: uppercase letters, lowercase letters, numbers, and special characters from ()~!@#$%^&*_-+=|{}[]:;'<>,.?/. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Import format is `instance_id/user_name/host`.

```shell
terraform import edgenext_rds_account.example b4a406bc-2859-430d-8f95-9bc1e054d347/app_user/%
```

