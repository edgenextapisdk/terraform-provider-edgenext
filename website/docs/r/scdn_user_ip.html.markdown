---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_user_ip"
sidebar_current: "docs-edgenext-resource-scdn_user_ip"
description: |-
  # edgenext_scdn_user_ip
---

# edgenext_scdn_user_ip

# edgenext_scdn_user_ip

Provides a resource to manage SCDN User IP Lists. This resource allows you to create, update, and delete IP address lists that can be used in various SCDN configurations.

## Example Usage

### Create a user IP list

```hcl
resource "edgenext_scdn_user_ip" "example" {
  name      = "example-ip-list"
  remark    = "Managed by Terraform"
  file_path = "${path.module}/ip_list.txt"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the IP list
* `file_path` - (Optional, String) The path to the file containing IP list to upload
* `remark` - (Optional, String) The remark/description for the IP list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - Creation time
* `id` - The ID of the IP list
* `item_num` - Number of IPs in the list
* `updated_at` - Last update time


## Import

SCDN User IP Lists can be imported using the `id`, e.g.

```
$ terraform import edgenext_scdn_user_ip.example 123
```

