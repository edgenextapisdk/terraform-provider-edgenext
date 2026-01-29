---
subcategory: "SCDN"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_user_ip"
sidebar_current: "docs-edgenext-resource-scdn-user-ip"
description: |-
  Provides a resource to manage SCDN User IP Lists.
---

# Resource: edgenext_scdn_user_ip

Provides a resource to manage SCDN User IP Lists. This resource allows you to create, update, and delete IP address lists that can be used in various SCDN configurations.

## Example Usage

```hcl
resource "edgenext_scdn_user_ip" "example" {
  name   = "example-ip-list"
  remark = "Managed by Terraform"
  file_path = "${path.module}/ip_list.txt"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the IP list.
* `remark` - (Optional) The remark or description for the IP list.
* `file_path` - (Optional) The path to the file containing IP list to upload.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IP list.
* `item_num` - The number of IP items in the list.
* `created_at` - The creation time of the list.
* `updated_at` - The last update time of the list.

## Import

SCDN User IP Lists can be imported using the `id`, e.g.

```
$ terraform import edgenext_scdn_user_ip.example 123
```
