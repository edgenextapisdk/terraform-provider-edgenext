---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_accounts"
sidebar_current: "docs-edgenext-datasource-rds_accounts"
description: |-
  Use this data source to list database users in one EdgeNext RDS instance.
---

# edgenext_rds_accounts

Use this data source to list database users in one EdgeNext RDS instance.

## Example Usage

```hcl
data "edgenext_rds_accounts" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) RDS instance ID to list database users for.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `users` - Database users returned by the API.
  * `databases` - Database names granted to this user when returned by the API; may be empty.
  * `host` - Host pattern the user is allowed to connect from (for example % or 127.0.0.1).
  * `user_name` - User name.


