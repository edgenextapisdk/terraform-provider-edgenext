---
subcategory: "Relational Database Service (RDS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_rds_backup_policy_associate_instance"
sidebar_current: "docs-edgenext-resource-rds_backup_policy_associate_instance"
description: |-
  Use this resource to manage one instance association for an EdgeNext RDS backup policy.
---

# edgenext_rds_backup_policy_associate_instance

Use this resource to manage one instance association for an EdgeNext RDS backup policy.

## Example Usage

```hcl
resource "edgenext_rds_backup_policy_associate_instance" "example" {
  policy_id   = "backup-policy-d0ce79aeedfd76697daf6569"
  instance_id = "b4a406bc-2859-430d-8f95-9bc1e054d347"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `policy_id` - (Required, String, ForceNew) Backup policy ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_name` - Instance name returned from the API.


## Import

Import format is `policy_id/instance_id`.

```shell
terraform import edgenext_rds_backup_policy_associate_instance.example backup-policy-d0ce79aeedfd76697daf6569/b4a406bc-2859-430d-8f95-9bc1e054d347
```

