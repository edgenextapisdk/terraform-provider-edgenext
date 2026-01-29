---
subcategory: "Security DNS (SDNS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_sdns_record"
sidebar_current: "docs-edgenext-resource-sdns_record"
description: |-
  Provides a resource to create and manage SDNS DNS records.
---

# edgenext_sdns_record

Provides a resource to create and manage SDNS DNS records.

## Example Usage

### Create SDNS DNS record

```hcl
resource "edgenext_sdns_record" "example" {
  domain_id = 12345
  name      = "www"
  type      = "A"
  view      = "any"
  value     = "1.2.3.4"
  ttl       = 600
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, Int, ForceNew) The ID of the domain
* `name` - (Required, String) The name of the record (e.g., www)
* `type` - (Required, String) The type of the record (A, CNAME, etc.)
* `value` - (Required, String) The value of the record
* `view` - (Required, String) The view/line for the record
* `mx` - (Optional, Int) MX priority
* `remark` - (Optional, String) Remark for the record
* `status` - (Optional, Int) Status of the record (1 for enabled, 2 for paused)
* `ttl` - (Optional, Int) TTL in seconds

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SDNS DNS records can be imported using the record ID:

```shell
terraform import edgenext_sdns_record.example 67890
```

