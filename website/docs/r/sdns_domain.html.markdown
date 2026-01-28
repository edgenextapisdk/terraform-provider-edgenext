---
subcategory: "Security DNS (SDNS)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_sdns_domain"
sidebar_current: "docs-edgenext-resource-sdns_domain"
description: |-
  Provides a resource to create and manage SDNS domains.
---

# edgenext_sdns_domain

Provides a resource to create and manage SDNS domains.

## Example Usage

### Create SDNS domain

```hcl
resource "edgenext_sdns_domain" "example" {
  domain = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String, ForceNew) The domain name to be added to DNS

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Status of the domain


## Import

SDNS domains can be imported using the domain ID:

```shell
terraform import edgenext_sdns_domain.example 12345
```

