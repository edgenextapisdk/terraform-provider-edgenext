---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_log_download_template_status"
sidebar_current: "docs-edgenext-resource-scdn_log_download_template_status"
description: |-
  Provides a resource to manage the status (enable/disable) of SCDN log download templates.
---

# edgenext_scdn_log_download_template_status

Provides a resource to manage the status (enable/disable) of SCDN log download templates.

## Example Usage

### Enable log download template

```hcl
resource "edgenext_scdn_log_download_template_status" "example" {
  template_id = 12345
  status      = 1 # 1: enabled
}
```

### Disable log download template

```hcl
resource "edgenext_scdn_log_download_template_status" "example" {
  template_id = 12345
  status      = 0 # 0: disabled
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, Int) Status: 1-enabled, 0-disabled
* `template_id` - (Required, Int) Template ID

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `id` - The ID of the template status (same as template_id)


## Import

SCDN log download template status can be imported using the template ID:

```shell
terraform import edgenext_scdn_log_download_template_status.example 12345
```

