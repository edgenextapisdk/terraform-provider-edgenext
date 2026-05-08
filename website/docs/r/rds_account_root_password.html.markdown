---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_account_root_password"
sidebar_current: "docs-edgenext-resource-rds_account_root_password"
description: |-
  Use this resource to enable root access and set (or rotate) root password on an EdgeNext RDS instance.
---

# edgenext_rds_account_root_password

Use this resource to enable root access and set (or rotate) root password on an EdgeNext RDS instance.

## Example Usage

```hcl
resource "edgenext_rds_account_root_password" "example" {
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
  password    = "Root@123456"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) RDS instance ID.
* `password` - (Required, String) Root password. Please enter 8-20 characters, must include all four: uppercase letters, lowercase letters, numbers, and special characters from ()~!@#$%^&*_-+=|{}[]:;'<>,.?/. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



